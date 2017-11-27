package nats_email

import (
	"encoding/json"
	"fmt"
	"github.com/akaumov/nats-email/js"
	"github.com/akaumov/nats-pool"
	"github.com/nats-io/go-nats"
	"log"
	"net/smtp"
	"os"
	"os/signal"
	"syscall"
)

type NatsEmail struct {
	config   *Config
	natsPool *nats_pool.Pool
}

func New(config *Config) *NatsEmail {
	return &NatsEmail{
		config: config,
	}
}

func (e *NatsEmail) handleRequest(natsClient *nats.Conn, msg *nats.Msg) {

	var request js.RequestSendEmail

	err := json.Unmarshal(msg.Data, &request)
	if err != nil {
		return
	}

	auth := smtp.PlainAuth(
		"",
		e.config.EmailServerUserName,
		e.config.EmailServerPassword,
		e.config.EmailServerAddress)

	err = smtp.SendMail(
		e.config.EmailServerAddress,
		auth,
		request.From,
		request.To,
		[]byte(request.Body))

	if err != nil {
		response, _ := json.Marshal(js.ResponseSendEmail{
			Result: "",
			Error:  fmt.Sprintf("can't send email: %v", err),
		})

		natsClient.PublishRequest(e.config.OutputNatsSubject, msg.Reply, response)
		return
	}

	response, _ := json.Marshal(js.ResponseSendEmail{
		Result: "ok",
	})

	natsClient.PublishRequest(e.config.OutputNatsSubject, msg.Reply, response)
}

func (e *NatsEmail) startListenBus(stopSignal chan bool) error {

	natsClient, err := e.natsPool.Get()
	if err != nil {
		return err
	}

	subscription, err := natsClient.Subscribe(e.config.ListenNatsSubject, func(msg *nats.Msg) {
		e.handleRequest(natsClient, msg)
	})

	if err != nil {
		return err
	}

	<-stopSignal
	subscription.Unsubscribe()
	return nil
}

func getOsSignalWatcher() chan os.Signal {

	stopChannel := make(chan os.Signal)
	signal.Notify(stopChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)

	return stopChannel
}

func (e *NatsEmail) Start() {

	stopSignal := getOsSignalWatcher()

	natsPool, err := nats_pool.New(e.config.NatsAddress, e.config.NatsPoolSize)
	if err != nil {
		log.Panicf("can't connect to nats: %v", err)
	}

	e.natsPool = natsPool
	defer func() { natsPool.Empty() }()

	stopListenerSignal := make(chan bool)
	e.startListenBus(stopListenerSignal)

	go func() {
		<-stopSignal
		e.Stop()
	}()
}

func (e *NatsEmail) Stop() {
	e.natsPool.Empty()
	log.Println("natspool: empty")
}
