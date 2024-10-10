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
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
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
