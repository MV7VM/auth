package config

type ConfigModel struct {
	Server   ServerConfig   `yaml:"Server"`
	Postgres PostgresConfig `yaml:"Postgres"`
	//GRPC     struct {
	//	ContentClient Client `yaml:"client"`
	//} `yaml:"grpc"`
}

type Client struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"DBName"`
	SSLMode  string `yaml:"sslMode"`
	PgDriver string `yaml:"pgDriver"`
}

type ServerConfig struct {
	AppVersion string `yaml:"appVersion"`
	Host       string `yaml:"host" validate:"required"`
	Port       string `yaml:"port" validate:"required"`
}

type Rabbit struct {
	InvoiceConsumer  BaseConsumerConfig  `yaml:"invoiceConsumer"`
	InvoicePublisher BasePublisherConfig `yaml:"invoicePublisher"`
	EventPublisher   BasePublisherConfig `yaml:"eventPublisher"`
}

type BaseConsumerConfig struct {
	URL           string `yaml:"url"`
	QueueName     string `yaml:"queueName"`
	ExchangerName string `yaml:"exchangerName"`
	ExchangerType string `yaml:"exchangerType"`
	RoutingKey    string `yaml:"routingKey"`
	Durable       bool   `yaml:"durable"`
	AutoDelete    bool   `yaml:"autoDelete"`
	Interval      bool   `yaml:"interval"`
	NoWait        bool   `yaml:"noWait"`
	Tag           string `yaml:"tag"`
}
type BasePublisherConfig struct {
	URL           string `yaml:"url"`
	QueueName     string `yaml:"queueName"`
	ExchangerName string `yaml:"exchangerName"`
	ExchangerType string `yaml:"exchangerType"`
	RoutingKey    string `yaml:"routingKey"`
	Durable       bool   `yaml:"durable"`
	AutoDelete    bool   `yaml:"autoDelete"`
	Interval      bool   `yaml:"interval"`
	NoWait        bool   `yaml:"noWait"`
}
