package db

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/qiniu/qmgo"
)

//
//-- SINGLETON: MClient
//

/*
Represents a qmgo client. This struct acts as a singleton wrapper on a
qmgo client.
*/
type MClient struct {
	client *qmgo.Client
	config *MConfig
	// Guard mutex to ensure atomicity during connect/disconnect operations.
	mutex *sync.Mutex
}

// Holds the instance object for the global database client.
var instance *MClient

// Guard mutex to ensure that only one singleton object is created.
var once sync.Once

// Gets the currently active qmgo client instance.
func GetInstance() *MClient {
	once.Do(func() {
		instance = &MClient{}
		instance.mutex = &sync.Mutex{}
	})
	return instance
}

/*
Gets the underlying client instance that's used to interact with the
MongoDB database. If the client is not currently connected, then this
object will be `nil`.
*/
func (m MClient) GetClient() *qmgo.Client {
	return m.client
}

/*
Gets the configuration used when the database connection was established.
If the client is not currently connected, then this object will be `nil`.
*/
func (m MClient) GetConfig() *MConfig {
	return m.config
}

// Returns whether there is an active connection to the database.
func (m MClient) IsConnected() bool {
	return m.client != nil
}

// Connects to the MongoDB server specified in the given config object.
func (m *MClient) Connect(cfg *MConfig) (*qmgo.Client, error) {
	//Lock the mutex and defer its unlock
	m.mutex.Lock()
	defer m.mutex.Unlock()

	//Ensure there isn't already a connection open
	if m.client != nil {
		return m.client, fmt.Errorf("mongodb: cannot establish a connection that is already open")
	}

	//Set client options
	clientOptions := qmgo.Config{
		//Uri: fmt.Sprintf("mongodb://%s:%d", cfg.MongoDB.Host, cfg.MongoDB.Port),
		Uri:              cfg.ConnStr,
		ConnectTimeoutMS: func(i int64) *int64 { return &i }(cfg.Timeout * 1000), //This is expected to be a pointer; see: https://stackoverflow.com/a/30716481
	}

	//Connect to the database
	client, err := qmgo.NewClient(context.Background(), &clientOptions)
	m.client, m.config = client, cfg

	//Return the client and any error that occurred
	return m.client, err
}

// Disconnects the client from the database and nullifies the instance.
func (m *MClient) Disconnect() error {
	//Lock the mutex and defer its unlock
	m.mutex.Lock()
	defer m.mutex.Unlock()

	//Disconnect from the db
	if m.client != nil {
		err := m.client.Close(context.Background())
		m.client, m.config = nil, nil
		return err
	}
	return nil
}

/*
Pings the MongoDB server to ensure the connection is ok. Returns the
ping time in microseconds.
*/
func (m MClient) Heartbeat() (int64, error) {
	//Ensure a connection actually exists
	if m.client == nil {
		return -1, fmt.Errorf("mongodb: cannot perform a heartbeat; client is not currently connected to a server")
	}

	//Ping the server
	bm := time.Now()
	err := m.client.Ping(m.config.Timeout) //Treated in seconds
	delta := time.Since(bm)

	//Return the ping time and any errors
	return delta.Microseconds(), err
}
