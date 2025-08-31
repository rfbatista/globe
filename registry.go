package globe

import "errors"

type ProcessRegistryValue struct {
	localPids  map[ID]chan Message
	remotePids []SupervisorResolver
}

var registry ProcessRegistryValue

func init() {
	registry.localPids = make(map[ID]chan Message)
}

func (p *ProcessRegistryValue) Add(pid PID, c chan Message) error {
	p.localPids[pid.ID] = c
	return nil
}

func (p *ProcessRegistryValue) Remove(pid PID) error {
	return nil
}

func (p *ProcessRegistryValue) Get(pid PID) (chan Message, error) {
	c, ok := p.localPids[pid.ID]
	if !ok {
		return nil, errors.New("channel not found")
	}
	return c, nil
}
