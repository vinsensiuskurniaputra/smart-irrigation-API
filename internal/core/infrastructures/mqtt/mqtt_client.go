package mqtt

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/core/infrastructures/config"
)

type Client struct {
	cli mqtt.Client
}

func NewClient(cfg *config.Config) *Client {
	broker := fmt.Sprintf("tcp://%s:%s", cfg.MQTT.Broker, cfg.MQTT.Port)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	if cfg.MQTT.ClientID != "" {
		opts.SetClientID(cfg.MQTT.ClientID)
	} else {
		opts.SetClientID(fmt.Sprintf("smart-irrigation-api-%d", time.Now().UnixNano()))
	}
	if cfg.MQTT.Username != "" {
		opts.SetUsername(cfg.MQTT.Username)
		opts.SetPassword(cfg.MQTT.Password)
	}
	opts.AutoReconnect = true
	opts.ConnectRetry = true
	opts.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})
	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Printf("MQTT connection lost: %v", err)
	}
	opts.OnConnect = func(c mqtt.Client) {
		log.Printf("Connected to MQTT broker %s", broker)
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Printf("Failed to connect MQTT: %v", token.Error())
	}
	return &Client{cli: client}
}

func (c *Client) Subscribe(topic string, cb mqtt.MessageHandler) error {
	if token := c.cli.Subscribe(topic, 0, cb); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	log.Printf("Subscribed to MQTT topic: %s", topic)
	return nil
}

func (c *Client) IsConnected() bool { return c.cli.IsConnected() }
