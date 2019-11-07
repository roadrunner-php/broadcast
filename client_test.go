package broadcast

import (
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spiral/roadrunner/service"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_Client_Consume_And_Produce(t *testing.T) {
	logger, _ := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	c := service.NewContainer(logger)
	c.Register(ID, &Service{})

	assert.NoError(t, c.Init(&testCfg{
		broadcast: `{}`,
	}))

	b, _ := c.Get(ID)
	br := b.(*Service)

	go func() { c.Serve() }()
	time.Sleep(time.Millisecond * 100)
	defer c.Stop()

	msg := make(chan *Message)

	client := br.NewClient(msg)
	assert.NoError(t, client.Subscribe("default"))

	assert.NoError(t, client.Publish(NewMessage("default", "hello")))

	assert.Equal(t, "\"hello\"", string((<-msg).Payload))
}

func Test_Client_Consume_And_Produce_On_Redis(t *testing.T) {
	logger, _ := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	c := service.NewContainer(logger)
	c.Register(ID, &Service{})

	assert.NoError(t, c.Init(&testCfg{
		broadcast: `{"redis":{"addr":"localhost:6379"}}`,
	}))

	b, _ := c.Get(ID)
	br := b.(*Service)

	go func() { c.Serve() }()
	time.Sleep(time.Millisecond * 100)
	defer c.Stop()

	msg := make(chan *Message)

	client := br.NewClient(msg)
	defer client.Close()

	assert.NoError(t, client.Subscribe("default"))

	assert.NoError(t, client.Publish(NewMessage("default", "hello")))

	assert.Equal(t, "\"hello\"", string((<-msg).Payload))
}

func Test_Client_Topics(t *testing.T) {
	logger, _ := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	c := service.NewContainer(logger)
	c.Register(ID, &Service{})

	assert.NoError(t, c.Init(&testCfg{
		broadcast: `{"redis":{"addr":"localhost:6379"}}`,
	}))

	b, _ := c.Get(ID)
	br := b.(*Service)

	go func() { c.Serve() }()
	time.Sleep(time.Millisecond * 100)
	defer c.Stop()

	msg := make(chan *Message)

	client := br.NewClient(msg)
	defer client.Close()

	assert.NoError(t, client.Subscribe("default"))

	assert.Equal(t, []string{"default"}, client.Topics())
	assert.NoError(t, client.Subscribe("default", "second"))
	assert.Equal(t, []string{"default", "second"}, client.Topics())

	client.Unsubscribe("default")
	assert.Equal(t, []string{"second"}, client.Topics())

	client.Unsubscribe("default")
	assert.Equal(t, []string{"second"}, client.Topics())

	client.Unsubscribe("second")
	assert.Equal(t, []string{}, client.Topics())
}
