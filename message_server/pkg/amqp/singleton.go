package amqp

import (
	"errors"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

/*
Represents an AMQP client. This struct acts as a singleton wrapper on a message broker client.
*/
type AMQPClient struct {
	client *amqp.Connection
	config *AMQPConfig
	mutex  *sync.Mutex
}

// Holds the instance object for the global AMQP client.
var instance *AMQPClient

// Guard mutex to ensure that only one singleton object is created.
var once sync.Once

// Gets the currently active AMQP client instance.
func GetInstance() *AMQPClient {
	once.Do(func() {
		instance = &AMQPClient{
			mutex: &sync.Mutex{},
		}
	})
	return instance
}

/*
Gets the underlying client instance that's used to interact with the
AMQP server. If the client is not currently connected, then this
object will be `nil`.
*/
func (c *AMQPClient) GetConn() *amqp.Connection {
	return c.client
}

/*
Gets the configuration used when the AMQP connection was established.
If the client is not currently connected, then this object will be `nil`.
*/
func (c *AMQPClient) GetConfig() *AMQPConfig {
	return c.config
}

// Connects to the AMQP server specified in the given config object.
func (c *AMQPClient) Connect(cfg *AMQPConfig) (*amqp.Connection, error) {
	//Lock the mutex and defer its unlock
	c.mutex.Lock()
	defer c.mutex.Unlock()

	//Ensure there isn't already a connection open
	if c.client != nil {
		return nil, errors.New("amqp: cannot establish a connection that is already open")
	}

	//Connect to the AMQP server
	var err error
	c.client, err = amqp.Dial(cfg.ConnURL())
	if err != nil {
		return nil, err
	}

	//Assign the config and return the client
	c.config = cfg
	return c.client, nil
}

// Disconnects the client from the AMQP server and nullifies the instance.
func (c *AMQPClient) Disconnect() error {
	//Lock the mutex and defer its unlock
	c.mutex.Lock()
	defer c.mutex.Unlock()

	//Ensure a connection is open
	if c.client == nil {
		return errors.New("amqp: cannot sever a connection that is not open")
	}

	//Disconnect from the AMQP server and return any errors
	err := c.client.Close()
	c.client = nil
	instance = nil // Clear instance for future connections.
	return err
}

/*
Pings the AMQP server to ensure the connection is ok. Returns the
ping time in microseconds.
*/
func (c *AMQPClient) Heartbeat() (int64, error) {
	//Check if the connection is closed
	if c.client == nil || c.client.IsClosed() {
		return 0, errors.New("amqp: cannot perform a heartbeat; client is not currently connected to a server")
	}

	start := time.Now()
	//Attempt to create a new channel to check if the connection is alive
	channel, err := c.client.Channel()
	if err != nil {
		return 0, err
	}
	defer channel.Close() //Close the channel after checking

	return time.Since(start).Microseconds(), nil
}
