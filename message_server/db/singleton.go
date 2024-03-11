package db

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//-- Structs and singletons

/*
Represents a MongoDB client. This struct acts as a singleton wrapper on a
MongoDB client.
*/
type MClient struct {
	client *mongo.Client
	config *MConfig
}

/* Holds the instance object for the database client. */
var instance *MClient

// Guard mutex to ensure that only one singleton object is created
var once sync.Once

//-- Acquirance function

/* Gets the currently active MongoDB client instance. */
func GetInstance() *MClient {
	once.Do(func() {
		if instance == nil {
			instance = &MClient{}
		}
	})
	return instance
}

//-- Methods

/*
Gets the underlying client instance that's used to interact with the
MongoDB database. If the client is not currently connected, then this
object will be `nil`.
*/
func (m MClient) GetClient() *mongo.Client {
	return m.client
}

/*
Gets the configuration used when the database connection was established.
If the client is not currently connected, then this object will be `nil`.
*/
func (m MClient) GetConfig() *MConfig {
	return m.config
}

// Connects to the MongoDB server specified in the given config object.
func (m *MClient) Connect(cfg *MConfig) (*mongo.Client, error) {
	//Ensure there isn't already a connection open
	if m.client != nil {
		return m.client, fmt.Errorf("cannot establish a connection that is already open")
	}

	//Set client option
	clientOptions := options.Client().
		//ApplyURI(fmt.Sprintf("mongodb://%s:%d", cfg.MongoDB.Host, cfg.MongoDB.Port)).
		ApplyURI(fmt.Sprintf("mongodb://%s:%d", cfg.Host, cfg.Port))

	//Connect to the database
	client, err := mongo.Connect(context.TODO(), clientOptions)
	m.client, m.config = client, cfg

	//Return the client and any error that occurred
	return m.client, err
}

/* Disconnects the client from the database and nullifies the instance. */
func (m *MClient) Disconnect() error {
	if m.client != nil {
		err := m.client.Disconnect(context.Background())
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
		return -1, fmt.Errorf("cannot perform a heartbeat; client is not currently connected to a server")
	}

	//Ping the server
	bm := time.Now()
	err := m.client.Ping(context.Background(), nil)
	delta := time.Since(bm)

	//Return the ping time and any errors
	return delta.Microseconds(), err
}
