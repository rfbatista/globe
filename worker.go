package globe

import (
	"log/slog"
	"time"
)

func SpawnWorker(h func(m *Context) error) PID {
	pid := NewPID()
	inboundchan := make(chan Message)
	m := NewMailbox()
	m.RegisterInboundChannel(inboundchan)
	actorchan := make(chan Message)
	a := worker{
		pid:     pid,
		handler: h,
		m:       m,
		c:       actorchan,
	}
	go func(p PID) {
		a.Start()
	}(pid)
	err := registry.Add(pid, inboundchan)
	if err != nil {
		slog.Error("failed to registry new actor")
	}
	return pid
}

type worker struct {
	pid     PID
	handler func(*Context) error
	m       Mailbox
	c       chan Message
}

func (a *worker) Start() {
	a.m.RegisterActorChannel(a.c)
	a.m.Start()
	a.Loop(a.c)
}

func (n *worker) Loop(ch chan Message) {
	timer := time.NewTimer(time.Second * 1)
	for {
		select {
		case m := <-ch:
			n.handler(NewContext(m, n.pid))
		case <-timer.C:
			continue
		}
	}
}
