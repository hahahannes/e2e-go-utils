package mqtt

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"errors"
	"os"
	"time"
	"strconv"
	"net/url"	
	"fmt"
	"crypto/tls"
	"github.com/hahahannes/e2e-go-utils/lib/streaming"
)

type TopicConfig = map[string]byte

type MQTTClient struct {
	Client      MQTT.Client
	Retained    *bool
	OnConnectHandler func(MQTT.Client)
	MsgChannel chan streaming.Message
	Host string
	Port string
	TopicConfig TopicConfig
	SubscribeInitial bool
}

func NewMQTTClient(host, port string, topicConfig TopicConfig, msgChannel chan streaming.Message, subscribeInitial bool) *MQTTClient {
	return &MQTTClient{
		Host: host,
		Port: port,
		TopicConfig: topicConfig,
		MsgChannel: msgChannel,
		OnConnectHandler: func(MQTT.Client) {},
		SubscribeInitial: subscribeInitial,
	}
}

func (client *MQTTClient) ConnectMQTTBroker(username, password *string) error {
	//MQTT.DEBUG = log.New(os.Stdout, "", 0)

	hostname, _ := os.Hostname()

	server := "tcp://"+client.Host+":"+client.Port
	retained := false
	client.Retained = &retained
	clientId := hostname+strconv.Itoa(time.Now().Second())

	connOpts := MQTT.NewClientOptions().
		AddBroker(server).
		SetClientID(clientId).
		SetCleanSession(true).
		SetConnectionLostHandler(func(c MQTT.Client, err error) {
			fmt.Println("Connection Lost to " + server)
		}).
		SetConnectionAttemptHandler(func(broker *url.URL, tlsCfg *tls.Config) *tls.Config {
			fmt.Println("Attempt to connect to " + server)
			return tlsCfg
		}).
		SetReconnectingHandler(func(mqttClient MQTT.Client, opt *MQTT.ClientOptions) {
			fmt.Println("Try to reconnect to " + server)
		}).
		// debug
		SetMaxReconnectInterval(5 * time.Second).
		SetPingTimeout(10 * time.Second).
		SetKeepAlive(10 * time.Second).
		SetAutoReconnect(true)

	if username != nil && *username != "" {
		connOpts.SetUsername(*username)
		if *password != "" {
			connOpts.SetPassword(*password)
		}
	}
	// TODO insecure skip ?
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)

	connOpts.OnConnect = func(c MQTT.Client) {
		fmt.Println("Connected to " + server)
	}
	
	client.Client = MQTT.NewClient(connOpts)

	loopCounter := 0
	for {
		fmt.Printf(fmt.Sprintf("Try to connect to: %s [%d/240]", server, loopCounter))

		if loopCounter > 240 {
			return errors.New("Could not connect with broker")
		}

		if token := client.Client.Connect(); token.Wait() && token.Error() != nil {
			fmt.Errorf("Could not connect to %s : %s\n", server, token.Error())
			time.Sleep(5 * time.Second)
		} else {
			fmt.Printf("Connected to %s\n", server)
			if client.SubscribeInitial {
				return client.InitialSubscribe()
			}
			return nil
		}
		loopCounter += 1
	}
}

func (client *MQTTClient) InitialSubscribe() error {
	if(len(client.TopicConfig) == 0) {
		return errors.New("No topics configured")
	}

	fmt.Printf("Subscribed to topics: %v", client.TopicConfig)
	if token := client.Client.SubscribeMultiple(client.TopicConfig, client.OnMessageReceived); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (client *MQTTClient) CloseConnection() {
	client.Client.Disconnect(250)
	time.Sleep(1 * time.Second)
}

func (client *MQTTClient) Publish(topic string, message string, qos int) error {
	if !client.Client.IsConnected() {
		fmt.Errorf("WARNING: mqtt client not connected")
		return errors.New("mqtt client not connected")
	}

	token := client.Client.Publish(topic, byte(qos), *client.Retained, message)
	if token.Wait() && token.Error() != nil {
		fmt.Errorf("Error on Publish: ", token.Error())
		return token.Error()
	}
	return nil
}

func (client *MQTTClient) Subscribe(topic string, qos int) error {
	token := client.Client.Subscribe(topic, 2, client.OnMessageReceived)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (client *MQTTClient) Unsubscribe(topic string) error {
	token := client.Client.Unsubscribe(topic)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (client *MQTTClient) OnMessageReceived(mqttClient MQTT.Client, message MQTT.Message) {
	client.MsgChannel <- streaming.Message{
		Topic: message.Topic(),
		Value: string(message.Payload()),
	}
}
