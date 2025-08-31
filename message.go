package globe

type Message interface {
	Reply(Message) error
	Body() []byte
}

type ActorMessage struct {
	rc   chan Message
	body []byte
}

func NewActorMessage(b []byte) Message {
	return ActorMessage{body: b}
}

func (a ActorMessage) Reply(m Message) error {
	return nil
}

func (a ActorMessage) Body() []byte {
	return a.body
}

type StopMessage struct {
	rc chan Message
}
