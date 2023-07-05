package messaging

import (
	"go.nandlabs.io/commons/errutils"
	"net/url"
	"sync"
)

var facade Messaging

//Producer interface is used to send message(s) to a specific provider
type Producer interface {
	// Send function sends an individual message to the url
	Send(*url.URL, Message, ...Option) error
	//SendBatch sends a batch of messages to the url
	SendBatch(*url.URL, []Message, ...Option) error
}

//Receiver interface provides the functions for receiving a message(s)
type Receiver interface {
	// Receive function performs on-demand receive of a single message.
	//This function may or may not wait for the messages to arrive. This is purely dependent on the implementation.
	Receive(*url.URL, ...Option) (Message, error)
	// ReceiveBatch function performs on-demand receive of a batch of messages .
	//This function may or may not wait for the messages to arrive. This is purely dependent on the implementation.
	ReceiveBatch(*url.URL, ...Option) ([]Message, error)
	// AddListener registers a listener for the message
	AddListener(*url.URL, func(msg Message), ...Option) error
}

type Provider interface {
	//Producer Interface included
	Producer
	//Receiver interface included
	Receiver
	// Schemes is array of URL schemes supported by this provider
	Schemes() []string
	// Setup method called
	Setup()
	//NewMessage function creates a new message that can be used by the clients. It expects the scheme to be provided
	NewMessage(string, ...Option) (Message, error)
}

type Messaging interface {
	Provider
	Register(Provider)
}

type Manager struct {
	knownProviders map[string]Provider
	mutex          sync.Mutex
}

func (m *Manager) getFor(scheme string) (provider Provider, err error) {
	var ok bool
	provider, ok = m.knownProviders[scheme]
	if !ok {
		err = errutils.FmtError("Unsupported scheme  %s ", scheme)
	}
	return
}

func (m *Manager) Send(u *url.URL, msg Message, options ...Option) (err error) {
	var provider Provider
	provider, err = m.getFor(u.Scheme)
	if err == nil {
		err = provider.Send(u, msg, options...)
	}
	return
}

func (m *Manager) Receive(u *url.URL, options ...Option) (msg Message, err error) {
	var provider Provider
	provider, err = m.getFor(u.Scheme)
	if err == nil {
		msg, err = provider.Receive(u, options...)
	}
	return
}

func (m *Manager) AddListener(u *url.URL, listener func(msg Message), options ...Option) (err error) {
	var provider Provider
	provider, err = m.getFor(u.Scheme)
	if err == nil {
		err = provider.AddListener(u, listener, options...)
	}
	return
}

func (m *Manager) ReceiveBatch(u *url.URL, options ...Option) (msgs []Message, err error) {
	var provider Provider
	provider, err = m.getFor(u.Scheme)
	if err == nil {
		msgs, err = provider.ReceiveBatch(u, options...)
	}
	return
}

func (m *Manager) SendBatch(u *url.URL, msgs []Message, options ...Option) (err error) {
	var provider Provider
	provider, err = m.getFor(u.Scheme)
	if err == nil {
		err = provider.SendBatch(u, msgs, options...)
	}
	return
}

func (m *Manager) Schemes() (schemes []string) {
	for k := range m.knownProviders {
		if k == "" {
			continue
		}
		schemes = append(schemes, k)
	}
	return
}

func (m *Manager) NewMessage(scheme string, options ...Option) (msg Message, err error) {
	var provider Provider
	provider, err = m.getFor(scheme)
	if err == nil {
		msg, err = provider.NewMessage(scheme, options...)
	}
	return
}

func (m *Manager) Setup() {
	//TODO The facade setup will register the local provider

}

func Get() Messaging {
	return facade
}

func (m *Manager) Register(provider Provider) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for _, s := range provider.Schemes() {
		if m.knownProviders == nil {
			m.knownProviders = make(map[string]Provider)
		}
		m.knownProviders[s] = provider
	}
}

func init() {
	facade = &Manager{
		knownProviders: make(map[string]Provider),
		mutex:          sync.Mutex{},
	}
	facade.Setup()
}
