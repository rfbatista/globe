package globe

type Context struct {
	sender   PID
	receiver PID
	m        Message
}

func NewContext(m Message, receiver PID) *Context {
	return &Context{
		m:        m,
		receiver: receiver,
	}
}

func (c *Context) Send(p PID, m Message) {
	router.Route(p, m)
}

func (c *Context) MyPID() PID {
	return c.receiver
}

func (c *Context) Body() []byte {
	return c.m.Body()
}
