package main

import (
	"context"
	"database/sql"
	"fmt"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/danilobandeira29/ms-wallet/internal/database"
	"github.com/danilobandeira29/ms-wallet/internal/event"
	"github.com/danilobandeira29/ms-wallet/internal/event/handler"
	"github.com/danilobandeira29/ms-wallet/internal/usecase/createaccount"
	"github.com/danilobandeira29/ms-wallet/internal/usecase/createclient"
	"github.com/danilobandeira29/ms-wallet/internal/usecase/createtransction"
	"github.com/danilobandeira29/ms-wallet/internal/web"
	"github.com/danilobandeira29/ms-wallet/internal/web/webserver"
	"github.com/danilobandeira29/ms-wallet/pkg/events"
	"github.com/danilobandeira29/ms-wallet/pkg/kafka"
	"github.com/danilobandeira29/ms-wallet/pkg/uow"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

const createClientsTable = `
	create table if not exists clients(
		id varchar(255) primary key, 
		name varchar(255) not null, 
		email varchar(255) not null, 
		created_at date not null, 
		updated_at date
	)
`
const createAccountsTable = `
	create table if not exists accounts(
	    id varchar(255) primary key, 
	    balance int not null, 
	    client_id varchar(255),
	    created_at date not null,
	    updated_at date,
	    foreign key (client_id) references clients(id) 
	)`
const createTransactionsTable = `
	create table if not exists transactions(
	    id varchar(255) primary key,
	    account_from_id varchar(255), 
	    account_to_id varchar(255), 
	    amount int not null, 
	    created_at date,
	    foreign key (account_from_id) references accounts(id),
	    foreign key (account_to_id) references accounts(id)
	)`
const insertIntoClients = `
	insert into clients(id, name, email, created_at) values 
		('abdd1759-b88b-46c3-bb55-28651ba3ab59', 'John Doe', 'john@doe.com', now()),
		('b36d0ec8-0ade-420e-9d3d-6f61e4634ca7', 'Janne Doe', 'janne@doe.com', now()),
		('63b85fe6-081e-484e-83c7-4e11daafc80c', 'Jannet Doe', 'jannet@doe.com', now()),
		('748b39ec-4505-4e74-b235-9971af2e2254', 'Jonhny Doe', 'jonhny@doe.com', now())
	on duplicate key update id = id;
`
const insertIntoAccounts = `
	insert into accounts(id, balance, client_id, created_at) values 
		('0502ee10-1536-435e-99c0-123c265e96c3', 1000, 'abdd1759-b88b-46c3-bb55-28651ba3ab59', now()),
		('a42582b5-9a81-4a34-a7f9-573ed825b189', 1000, 'b36d0ec8-0ade-420e-9d3d-6f61e4634ca7', now()),
		('23950fec-1a45-4656-bb78-6f025ff9f709', 1000, '63b85fe6-081e-484e-83c7-4e11daafc80c', now()),
		('ce99e4db-646f-4b7e-bb9b-1e5243a92248', 1000, '748b39ec-4505-4e74-b235-9971af2e2254', now())
	on duplicate key update id = id;
`

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(mysql:3306)/wallet?parseTime=true", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD")))
	if err != nil {
		log.Fatalln("error when trying to connect to the database: ", err)
	}
	defer func(d *sql.DB) {
		err = d.Close()
		if err != nil {
			log.Printf("error when trying to close database's connection %v\n", err)
		}
	}(db)
	for _, q := range []string{createClientsTable, createAccountsTable, createTransactionsTable, insertIntoClients, insertIntoAccounts} {
		_, err := db.Exec(q)
		if err != nil {
			log.Fatalf("not possible to create table %v\n", err)
		}
	}
	accountDb := database.NewAccountDB(db)
	clientDb := database.NewClientDB(db)
	createClientUseCase := createclient.NewCreateClientUseCase(clientDb)
	createAccountUseCase, _ := createaccount.NewCreateAccountUseCase(accountDb, clientDb)
	eventDispatcher := events.NewEventDispatcher()
	configMap := ckafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"group.id":          "wallet",
	}
	fmt.Println("connecting to kafka")
	kafkaAdmin, err := ckafka.NewAdminClient(&configMap)
	if err != nil {
		log.Fatalf("cannot create ckafka admin client: %v\n", err)
	}
	_, err = kafkaAdmin.CreateTopics(context.Background(), []ckafka.TopicSpecification{
		{
			Topic:         "balances",
			NumPartitions: 1,
		},
		{
			Topic:         "transactions",
			NumPartitions: 1,
		},
	})
	if err != nil {
		log.Fatalf("cannot create ckafka's topics: %v\n", err)
	}
	kafkaProduct := kafka.NewKafkaProducer(&configMap)
	err = eventDispatcher.Register("TransactionCreated", handler.NewTransactionCreatedKafkaHandler(kafkaProduct))
	if err != nil {
		log.Printf("error when trying to TransactionCreated handler %v\n", err)
		return
	}
	err = eventDispatcher.Register("BalanceUpdated", handler.NewBalanceUpdatedKafkaHandler(kafkaProduct))
	if err != nil {
		log.Printf("error when trying to BalanceUpdated handler %v\n", err)
		return
	}
	transactionCreatedEvent := event.NewTransactionCreated()
	balanceUpdatedEvent := event.NewBalanceUpdated()
	ctx := context.Background()
	unitOfWork := uow.NewUow(ctx, db)
	unitOfWork.Register("AccountDB", func(tx *sql.Tx) interface{} {
		return database.NewAccountDB(db)
	})
	unitOfWork.Register("TransactionDB", func(tx *sql.Tx) interface{} {
		return database.NewTransactionDB(db)
	})
	createTransactionUseCase, _ := createtransction.NewCreateTransactionUseCase(unitOfWork, eventDispatcher, transactionCreatedEvent, balanceUpdatedEvent)
	webServer := webserver.NewWebServer(":8080")
	transactionHandler := web.NewTransactionHandler(*createTransactionUseCase)
	accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	webServer.AddHandler("/transactions", transactionHandler.CreateClient)
	webServer.AddHandler("/accounts", accountHandler.CreateClient)
	webServer.AddHandler("/clients", clientHandler.CreateClient)
	fmt.Println("webserver started at port :8080")
	webServer.Start()
}
