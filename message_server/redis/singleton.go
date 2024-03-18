package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

//
//-- SINGLETON: RClient
//

/*
Represents a Redis client. This struct acts as a singleton wrapper on a
Redis client.
*/
type RClient struct {
	client *redis.Client
	config *RConfig
	// Guard mutex to ensure atomicity during connect/disconnect operations.
	mutex *sync.Mutex
}

// Holds the instance object for the global database client.
var instance *RClient

// Guard mutex to ensure that only one singleton object is created.
var once sync.Once

// Gets the currently active Redis client instance.
func GetInstance() *RClient {
	once.Do(func() {
		instance = &RClient{}
		instance.mutex = &sync.Mutex{}
	})
	return instance
}

/*
Gets the underlying client instance that's used to interact with the
Redis database. If the client is not currently connected, then this
object will be `nil`.
*/
func (m RClient) GetClient() *redis.Client {
	return m.client
}

/*
Gets the configuration used when the database connection was established.
If the client is not currently connected, then this object will be `nil`.
*/
func (m RClient) GetConfig() *RConfig {
	return m.config
}

// Connects to the Redis server specified in the given config object.
func (m *RClient) Connect(cfg *RConfig) (*redis.Client, error) {
	//Lock the mutex and defer its unlock
	m.mutex.Lock()
	defer m.mutex.Unlock()

	//Ensure there isn't already a connection open
	if m.client != nil {
		return m.client, fmt.Errorf("cannot establish a connection that is already open")
	}

	//Set client options
	clientOptions := redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.DB,
	}

	//Connect to the database
	m.client = redis.NewClient(&clientOptions)
	m.config = cfg

	//Check if the connection was successful with a ping
	_, err := m.Heartbeat()

	//Return the client and any error that occurred
	return m.client, err
}

/*
Pings the Redis server to ensure the connection is ok. Returns the
ping time in microseconds.
*/
func (m RClient) Heartbeat() (delta int64, err error) {
	//Ensure a connection actually exists
	if m.client == nil {
		return -1, fmt.Errorf("cannot perform a heartbeat; client is not currently connected to a server")
	}

	//Ping the server
	bm := time.Now()
	var pong *redis.StatusCmd
	if pong = m.client.Ping(context.Background()); pong.String() != "ping: PONG" {
		err = fmt.Errorf("unexpected server response; %s", pong)
		return
	}
	delta = time.Since(bm).Microseconds()

	//Return the ping time and any errors
	return
}
