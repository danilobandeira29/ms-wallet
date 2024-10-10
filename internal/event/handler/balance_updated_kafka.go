package handler

import (
	"fmt"
	"github.com/danilobandeira29/ms-wallet/pkg/events"
	"github.com/danilobandeira29/ms-wallet/pkg/kafka"
	"log"
	"sync"
)

type BalanceUpdatedKafkaHandler struct {
	Kafka *kafka.Producer
}

func NewBalanceUpdatedKafkaHandler(k *kafka.Producer) *BalanceUpdatedKafkaHandler {
	return &BalanceUpdatedKafkaHandler{
		Kafka: k,
	}
}

func (h *BalanceUpdatedKafkaHandler) Handle(message events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	err := h.Kafka.Publish(message, nil, "balances")
	if err != nil {
		log.Printf("error when trying to publish in kafka %v\n", err)
		return
	}
	fmt.Printf("calling kafka for balance updated with %v\n", message.GetPayload())
}
