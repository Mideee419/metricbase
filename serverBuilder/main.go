package serverBuilder

import (
	"github.com/msiebuhr/MetricBase"
	"github.com/msiebuhr/MetricBase/backends"
)

type MetricServer struct {
	frontends []MetricBase.Frontend
	backend   backends.Backend
	stopChan  chan bool
}

func NewMetricServer(f []MetricBase.Frontend, b backends.Backend) MetricServer {
	// Hook up backends
	for _, front := range f {
		front.SetBackend(b)
	}

	// Server construction
	return MetricServer{
		stopChan:  make(chan bool),
		frontends: f,
		backend:   b,
	}
}

func (m *MetricServer) Start() {
	// Start the backend
	go m.backend.Start()

	// Start all front-ends, now they can talk to something
	for i := range m.frontends {
		go m.frontends[i].Start()
	}

	// Wait for order to stop
	<-m.stopChan

	// Close up front-ends
	for i := range m.frontends {
		m.frontends[i].Stop()
	}

	m.backend.Stop()
}

func (m *MetricServer) Stop() {
	m.stopChan <- true
}
