package server

import (
	"net"
	"os"
	"os/signal"
	"time"

	// "github.com/satori/go.uuid"
	messagebus "github.com/vardius/message-bus"
	"go.uber.org/zap"

	"github.com/jasrags/JasonMUD/connection"
	"github.com/jasrags/JasonMUD/event"
)

type Server interface {
	Start(host, port, network string) error
	AddConnection(conn net.Conn) (connection.Connection, error)
}

type server struct {
	logger         *zap.Logger
	eventBus       messagebus.MessageBus
	connectionList map[string]connection.Connection
}

func New(logger *zap.Logger) Server {
	return &server{
		logger:         logger.With(zap.String("component", "server")),
		eventBus:       messagebus.New(),
		connectionList: map[string]connection.Connection{},
	}

}

func (s *server) Start(host, port, network string) error {
	logger := s.logger.With(zap.String("func", "server.start"))
	logger.Debug("in:server.start")

	logger.Info("starting server")
	l, err := net.Listen(network, net.JoinHostPort(host, port))
	if err != nil {
		logger.Error("unable to start server", zap.Error(err), zap.String("host", host),
			zap.String("port", port), zap.String("network", network))

		return err
	}
	defer l.Close()

	// Subscribe to events
	logger.Info("setting up event listeners")
	s.eventBus.Subscribe(event.ConnectionOpened, s.onConnectionOpened)
	s.eventBus.Subscribe(event.ConnectionClosed, s.onConnectionClosed)

	// logger.Info("starting tickers")
	s.StartTickers()

	// listen for ctrl+c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			logger.Info("caught ctrl+c, shutting down")
			s.eventBus.Publish(event.ServerShutdown)
			time.Sleep(time.Second)
			l.Close()
			os.Exit(0)
			break
		}
	}()

	logger.Debug("server started, listening for connections")
	for {
		conn, err := l.Accept()
		if err != nil {
			logger.Error("unable to accept connection", zap.Error(err))
		}
		s.AddConnection(conn)
	}
}

func (s *server) StartTickers() {
	go s.runTicker(time.Tick(250*time.Millisecond), event.Tick250ms)
	go s.runTicker(time.Tick(time.Second), event.Tick1s)
	go s.runTicker(time.Tick(5*time.Second), event.Tick5s)
	go s.runTicker(time.Tick(30*time.Second), event.Tick30s)
	go s.runTicker(time.Tick(time.Minute), event.Tick1m)
}

func (s *server) runTicker(tick <-chan time.Time, topic string) {
	for range tick {
		s.eventBus.Publish(topic, topic)
	}
}

func (s *server) AddConnection(conn net.Conn) (connection.Connection, error) {
	logger := s.logger.With(zap.String("func", "server.add_connection"))
	logger.Debug("in:server.add_connection")

	c := connection.New(logger, conn, s.eventBus)
	s.connectionList[c.SessionID()] = c
	go c.Listen()

	// s.eventBus.Publish(fmt.Sprintf(event.MessageDirect, c.SessionID()), "Admin", "Welcome!")
	// s.eventBus.Publish(event.MessageGlobal, message.Message{
	// 	Body: fmt.Sprintf("Welcome %v", c.SessionID()),
	// })

	return c, nil
}

// Event handlers
func (s *server) onConnectionOpened(id string) {
	logger := s.logger.With(zap.String("func", "server.on_connection_opened"), zap.String("id", id))
	logger.Debug("in:server.on_connection_opened")

	s.logger.Info("connection count", zap.Int("count", len(s.connectionList)))
}

func (s *server) onConnectionClosed(id string) {
	logger := s.logger.With(zap.String("func", "server.on_connection_closed"), zap.String("id", id))
	logger.Debug("in:server.on_connection_closed")

	delete(s.connectionList, id)

	s.logger.Info("connection count", zap.Int("count", len(s.connectionList)))
}
