package event

const (
	ServerInit     = "server:init"
	ServerStart    = "server:start"
	ServerShutdown = "server:shutdown"

	Tick250ms = "tick:250ms"
	Tick1s    = "tick:1s"
	Tick5s    = "tick:5s"
	Tick30s   = "tick:30s"
	Tick1m    = "tick:1m"

	// MessageGlobal = "message:global"
	// MessageDirect = "message:%v"

	ConnectionOpened = "connection:opened"
	ConnectionClosed = "connection:closed"
)
