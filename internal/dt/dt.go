package dt

import (
	"fmt"
	"log"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/wagslane/go-rabbitmq"
)

var RabbitMQ *rabbitmq.Conn

func SetupRabbitMQ() *dockertest.Resource {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalln("failed to create dt pool", err)
	}

	runOptions := &dockertest.RunOptions{
		Repository: "docker.io/bitnami/rabbitmq",
		Tag:        "latest",
		DNS:        []string{"rabbit"},
		Env: []string{
			"RABBITMQ_DEFAULT_USER=test",
			"RABBITMQ_DEFAULT_PASS=test",
		},
	}

	resource, err := pool.RunWithOptions(runOptions, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalln("could not start rabbitmq resource", err)
	}

	hostAndPort := resource.GetHostPort("5672/tcp")

	err = resource.Expire(600)
	if err != nil {
		log.Fatalln("could not set expire time for rabbitmq resource ", err)
	}

	url := fmt.Sprintf("amqp://test:test@%s/", hostAndPort)

	pool.MaxWait = 600 * time.Second

	err = pool.Retry(func() (err error) {
		RabbitMQ, err = rabbitmq.NewConn(url)
		if err != nil {
			log.Println("failed to ping rabbitmq, retrying ...", err)
			return fmt.Errorf("failed to ping rabbitmq: %w", err)
		}
		return nil
	})
	if err != nil {
		log.Fatalln("failed to connect to docker", err)
	}

	return resource
}
