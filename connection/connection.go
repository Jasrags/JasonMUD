package connection

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	// "time"

	"github.com/satori/go.uuid"
	messagebus "github.com/vardius/message-bus"
	"go.uber.org/zap"

	// "github.com/logrusorgru/aurora"
	"github.com/jasrags/JasonMUD/event"
	// "github.com/jasrags/JasonMUD/message"
	// "github.com/jasrags/JasonMUD/player"
)

func New(logger *zap.Logger, conn net.Conn, eventBus messagebus.MessageBus) Connection {
	sessionID := uuid.NewV4()
	return &connection{
		sessionID: sessionID,
		conn:      conn,
		logger:    logger,
		eventBus:  eventBus,
	}
}

type Connection interface {
	SessionID() string
	Listen()
	Write(message string)
	Close()
}

type connection struct {
	sessionID uuid.UUID
	conn      net.Conn
	logger    *zap.Logger
	eventBus  messagebus.MessageBus
}

func (c *connection) SessionID() string {
	return c.sessionID.String()
}

func (c *connection) Listen() {
	logger := c.logger.With(zap.String("func", "connection.listen"), zap.String("session_id", c.SessionID()))
	logger.Debug("in:connection.listen")

	// Subscribe to events
	// c.eventBus.Subscribe(event.ServerShutdown, c.onServerShutdown)
	// c.eventBus.Subscribe(event.Tick250ms, c.onTick)
	// c.eventBus.Subscribe(event.Tick1s, c.onTick)
	// c.eventBus.Subscribe(event.Tick5s, c.onTick)
	// c.eventBus.Subscribe(event.Tick30s, c.onTick)
	// c.eventBus.Subscribe(event.Tick1m, c.onTick)

	c.eventBus.Publish(event.ConnectionOpened, c.SessionID())
	c.Write(">")

	rd := bufio.NewReader(c.conn)
	for {
		line, errReadString := rd.ReadString('\n')
		if errReadString != nil {
			if errReadString == io.EOF {
				c.logger.Info("connection closed")
			} else {
				c.logger.Error("connection buffer read error", zap.Error(errReadString))
			}
			c.Close()
			return
		}

		line = strings.TrimSpace(line)
		c.logger.Debug("line received", zap.String("line", line))
	}
}

func (c *connection) Write(message string) {
	logger := c.logger.With(zap.String("func", "connection.write"), zap.String("message", message))
	logger.Debug("in:connection.write")

	c.conn.Write([]byte(message))
}

func (c *connection) Close() {
	logger := c.logger.With(zap.String("func", "connection.close"))
	logger.Debug("in:connection.close")

	c.eventBus.Publish(event.ConnectionClosed, c.SessionID())
	c.conn.Close()
}

// Event handlers
func (c *connection) onServerShutdown() {
	c.Write("Server is shutting down\n")
	c.Close()
}

func (c *connection) onTick(topic string) {
	c.Write(fmt.Sprintf("tick[%v]\n", topic))
}

// const (
// 	StateWelcome int8 = iota
// 	StateLoginUsername
// 	StateLoginPassword
// 	StateLoginMenu
// 	StateCharacterCreation
// 	StatePlaying
// )

// type Connection interface {
// 	SessionID() string
// 	ConnectedAt() time.Time

// 	Write(message string)
// 	Listen()
// 	Close()
// }

// func New(logger *zap.Logger, eventBus messagebus.MessageBus, conn net.Conn, connectedAt time.Time,
// 	state int8) (Connection, error) {
// 	sessionID := uuid.NewV4()
// 	c := &connection{
// 		logger:      logger.With(zap.String("component", "connection"), zap.Stringer("session_id", sessionID)),
// 		eventBus:    eventBus,
// 		conn:        conn,
// 		sessionID:   sessionID,
// 		connectedAt: connectedAt,
// 		state:       state,
// 	}

// 	// Subscribe to events
// 	c.eventBus.Subscribe(event.MessageGlobal, c.onMessageGlobal)
// 	c.eventBus.Subscribe(fmt.Sprintf(event.MessageDirect, c.sessionID), c.onMessageDirect)

// 	return c, nil
// }

// type connection struct {
// 	logger      *zap.Logger
// 	eventBus    messagebus.MessageBus
// 	conn        net.Conn
// 	sessionID   uuid.UUID
// 	connectedAt time.Time
// 	state       int8
// 	player      player.Player
// }

// func (c *connection) SessionID() string {
// 	return c.sessionID.String()
// }

// func (c *connection) ConnectedAt() time.Time {
// 	return c.connectedAt
// }

// func (c *connection) Write(message string) {
// 	logger := c.logger.With(zap.String("func", "connection.write"), zap.String("message", message))
// 	logger.Debug("in:connection.write")

// 	c.conn.Write([]byte(message))
// }

// func (c *connection) WritePrompt() {
// 	c.Write("> ")
// }

// // func (c *connection) ReadLine(rd *bufio.Reader) (string, error) {
// // 	line, err := rd.ReadString('\n')
// // 	if err != nil {
// // 		c.logger.Error("unable to read line", zap.Error(err))
// // 		return "", err
// // 	}
// // 	return strings.TrimSpace(line), nil
// // }

// func GetCommands(line string) []string {
// 	return strings.Split(line, " ")
// }

// func GetWord(line string) string {
// 	cmds := GetCommands(line)
// 	return cmds[0]
// }

// // func (c *connection) ReadPassword() (string, error) {
// // 	password, err := terminal.ReadPassword(0)
// // 	if err != nil {
// // 		c.logger.Error("unable to read password", zap.Error(err))
// // 		return "", err
// // 	}
// // 	return string(password), nil
// // }

// func (c *connection) Listen() {
// 	logger := c.logger.With(zap.String("func", "connection.listen"))
// 	logger.Debug("in:connection.listen")

// 	c.Write("Welcome to whatever this is!\n")
// 	c.Write("Enter your username: ")
// 	c.state = StateLoginUsername

// 	rd := bufio.NewReader(c.conn)
// 	for {
// 		line, err := rd.ReadString('\n')
// 		if err != nil {
// 			c.Close()
// 			return
// 		}
// 		line = strings.TrimSpace(line)
// 		c.Write(fmt.Sprintf("%v\n", aurora.BgGray(aurora.Magenta(line))))

// 		// cmd = strings.TrimSpace(cmd)
// 		// logger.Debug("received cmd", zap.String("cmd", cmd))

// 		switch c.state {
// 		case StateLoginUsername:
// 			c.player.Username = GetWord(line)
// 			// c.Write("Enter your password: ")
// 			// c.state = StateLoginPassword
// 			c.state = StatePlaying
// 			c.Write("> ")
// 		// case StateLoginPassword:
// 		// 	c.player.password = GetWord(line)
// 		// 	c.state = StatePlaying
// 		// 	c.Write(fmt.Sprintf("%v %v\n",
// 		// 		aurora.BgGray(aurora.Magenta(c.player.Username)),
// 		// 		aurora.BgGray(aurora.Magenta(c.player.Password)),
// 		// 	))
// 		case StatePlaying:
// 			fallthrough
// 		default:
// 			c.Write("> ")
// 			// cmds, _ := c.ReadCommand(rd)
// 			// if err := c.ParseCommand(cmds); err != nil {
// 			// c.Write(err.Error())
// 			// }
// 		}
// 	}
// }

// func (c *connection) Close() {
// 	logger := c.logger.With(zap.String("func", "connection.close"))
// 	logger.Debug("in:connection.close")

// 	c.eventBus.Publish(event.ConnectionClosed, c.SessionID())

// 	c.conn.Close()
// }

// func (c *connection) ParseCommand(cmds []string) error {
// 	switch cmds[0] {
// 	case "tell":
// 		if len(cmds) < 3 {
// 			return fmt.Errorf("usage: tell <user> <message>")
// 		}
// 		to := cmds[1]
// 		body := strings.Join(cmds[2:], " ")
// 		c.SendDirectMessage(to, body)
// 	case "say":
// 		if len(cmds) < 2 {
// 			return fmt.Errorf("usage: say <message>")
// 		}
// 		body := strings.Join(cmds[1:], " ")
// 		c.SendGlobalMessage(body)
// 	}
// 	return nil
// }

// func (c *connection) SendDirectMessage(to, body string) {
// 	topic := fmt.Sprintf("%v:%v", event.MessageDirect, to)
// 	from := c.sessionID.String()
// 	c.eventBus.Publish(topic, message.Message{From: from, Body: body})
// }

// func (c *connection) SendGlobalMessage(body string) {
// 	c.eventBus.Publish(event.MessageGlobal, message.Message{Body: body})
// }

// // Event handlers
// func (c *connection) onTick() {
// 	logger := c.logger.With(zap.String("func", "connection.on_tick"))
// 	logger.Debug("in:connection.on_tick")
// }

// func (c *connection) onMessageGlobal(m message.Message) {
// 	logger := c.logger.With(zap.String("func", "connection.on_message_global"), zap.Any("message", m))
// 	logger.Debug("in:connection.on_message_global")

// 	sb := strings.Builder{}
// 	sb.WriteString(fmt.Sprintf("[%v] ", aurora.BgBlue(aurora.Red("GLOBAL"))))
// 	if m.From != "" {
// 		sb.WriteString(fmt.Sprintf("(%v) ", m.From))
// 	}
// 	sb.WriteString(fmt.Sprintf("%v\n", m.Body))
// 	c.Write(sb.String())
// 	c.WritePrompt()
// }

// func (c *connection) onMessageDirect(m message.Message) {
// 	logger := c.logger.With(zap.String("func", "connection.on_message_direct"), zap.Any("message", m))
// 	logger.Debug("in:connection.on_message_direct")

// 	sb := strings.Builder{}
// 	sb.WriteString(fmt.Sprintf("[%v] %v\n", aurora.BgGray(aurora.Black(m.From)), m.Body))
// 	c.Write(sb.String())
// 	c.WritePrompt()
// }
