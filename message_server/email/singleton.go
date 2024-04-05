package email

import (
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

//
//-- SINGLETON: EClient
//

/*
Represents an SMTP client. This struct acts as a singleton wrapper on an
email client.
*/
type EClient struct {
	client *mail.SMTPClient
	config *EConfig
	// Guard mutex to ensure atomicity during connect/disconnect operations.
	mutex *sync.Mutex
}

// Holds the instance object for the global email client.
var instance *EClient

// Guard mutex to ensure that only one singleton object is created.
var once sync.Once

// Gets the currently active SMTP client instance.
func GetInstance() *EClient {
	once.Do(func() {
		instance = &EClient{}
		instance.mutex = &sync.Mutex{}
	})
	return instance
}

/*
Gets the underlying client instance that's used to interact with the
SMTP server. If the client is not currently connected, then this
object will be `nil`.
*/
func (m EClient) GetClient() *mail.SMTPClient {
	return m.client
}

/*
Gets the configuration used when the SMTP connection was established.
If the client is not currently connected, then this object will be `nil`.
*/
func (m EClient) GetConfig() *EConfig {
	return m.config
}

// Connects to the SMTP server specified in the given config object.
func (m *EClient) Connect(cfg *EConfig) (*mail.SMTPClient, error) {
	//Lock the mutex and defer its unlock
	m.mutex.Lock()
	defer m.mutex.Unlock()

	//Ensure there isn't already a connection open
	if m.client != nil {
		return m.client, fmt.Errorf("email: cannot establish a connection that is already open")
	}

	//Create a new SMTP client instance and set options
	smtps := mail.NewSMTPClient()
	smtps.Host = cfg.Host
	smtps.Port = cfg.Port
	smtps.Username = cfg.Username
	smtps.Password = cfg.Password
	smtps.Encryption = mail.Encryption(cfg.EncType)
	//smtps.KeepAlive = true //This must be true or the client will disconnect after sending one email

	//Set options for TLS if verification is not needed
	if !cfg.VerifyCert {
		smtps.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	//Connect to the SMTP server
	smtpc, err := smtps.Connect()
	m.client, m.config = smtpc, cfg

	//Return the client and any error that occurred
	return m.client, err
}

// Disconnects the client from the SMTP server and nullifies the instance.
func (m *EClient) Disconnect() error {
	//Lock the mutex and defer its unlock
	m.mutex.Lock()
	defer m.mutex.Unlock()

	//Disconnect from the SMTP server
	if m.client != nil {
		err := m.client.Close()
		m.client, m.config = nil, nil
		return err
	}
	return nil
}

/*
Pings the SMTP server to ensure the connection is ok. Returns the
ping time in microseconds.
*/
func (m EClient) Heartbeat() (int64, error) {
	//Ensure a connection actually exists
	if m.client == nil {
		return -1, fmt.Errorf("email: cannot perform a heartbeat; client is not currently connected to a server")
	}

	//Ping the server
	bm := time.Now()
	err := m.client.Noop()
	delta := time.Since(bm)

	//Return the ping time and any errors
	return delta.Microseconds(), err
}

/*
Sends an email using an email object. This method allows the email sending
routine to respawn the SMTP server connection if it goes down for whatever
reason, gracefully reconnecting in the process. See the following GitHub
issues for more information: https://github.com/xhit/go-simple-mail/issues/13,
https://github.com/xhit/go-simple-mail/issues/23
*/
func (m EClient) SendEmail(em *mail.Email) error {
	//Lock the mutex and defer its unlock
	m.mutex.Lock()
	defer m.mutex.Unlock()

	//Send a reset command to the server
	if err := m.client.Reset(); err != nil {
		return err
	}

	//Send the email using the given email object
	return em.Send(m.client)
}
