package nats_email

type Config struct {
	EmailServerAddress  string `json:"emailServerHost"`
	EmailServerUserName string `json:"emailServerUserName"`
	EmailServerPassword string `json:"emailServerPassword"`

	ListenNatsSubject string `json:"listenNatsSubject"`
	OutputNatsSubject string `json:"outputNatsSubject"`
	PacketFormat      string `json:"packetFormat"`

	NatsAddress  string `json:"natsAddress"`
	NatsPoolSize int    `json:"natsPoolSize"`
}
