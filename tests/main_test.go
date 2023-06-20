package rabbitmq_test

import (
	"os"
	"testing"

	"github.com/wagslane/go-rabbitmq/internal/dt"
)

func TestMain(m *testing.M) {
	res := dt.SetupRabbitMQ()
	code := m.Run()
	res.Close()
	os.Exit(code)
}
