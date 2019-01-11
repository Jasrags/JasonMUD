package message

import "time"

const (
	TypeGlobal = "global"
	TypeDirect = "direct"
)

type Messages []Message

// type Message interface {
// Body() string
// }

type Message struct {
	From   string
	Body   string
	SentAt time.Time
}

// func New(t, from, body string) Message {
// 	m := message{
// 		t, from
// 		body: body,
// 	}

// 	return &m
// }

// func (m *message) Body() string {
// 	return m.body
// }
