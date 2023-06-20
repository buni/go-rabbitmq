package rabbitmq_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/wagslane/go-rabbitmq"
	"github.com/wagslane/go-rabbitmq/internal/dt"
)

func TestPublisherConfirmRaceCondition(t *testing.T) {
	sub, err := rabbitmq.NewConsumer(dt.RabbitMQ, func(d rabbitmq.Delivery) (action rabbitmq.Action) {
		log.Println(string(d.Body))
		return rabbitmq.Ack
	}, "queue", rabbitmq.WithConsumerOptionsRoutingKey("my_routing_key"), rabbitmq.WithConsumerOptionsExchangeName("events"), rabbitmq.WithConsumerOptionsExchangeDeclare)
	if err != nil {
		log.Fatal(err)
	}

	t.Cleanup(sub.Close)

	publisher, err := rabbitmq.NewPublisher(
		dt.RabbitMQ,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName("events"),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
		rabbitmq.WithPublisherOptionsConfirm,
	)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 100; i++ {
		time.Sleep(1 * time.Second)

		confirms, err := publisher.PublishWithDeferredConfirmWithContext(
			context.Background(),
			[]byte("hello, world"),
			[]string{"my_routing_key"},
			rabbitmq.WithPublishOptionsContentType("application/json"),
			rabbitmq.WithPublishOptionsMandatory,
			rabbitmq.WithPublishOptionsExchange("events"),
		)
		if err != nil {
			log.Println(err)
			continue
		}

		if len(confirms) == 0 || confirms[0] == nil {
			fmt.Println("message publishing not confirmed")
			continue
		}

		fmt.Println("message published")

		ok, err := confirms[0].WaitContext(context.Background())
		if err != nil {
			log.Println(err)
		}

		if ok {
			fmt.Println("message publishing confirmed")
		} else {
			fmt.Println("message publishing not confirmed")
		}
	}
}
