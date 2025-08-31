package globe

import "github.com/hashicorp/memberlist"

type SupervisorResolver interface{}

type MemberlistResolver struct {
	m memberlist.Memberlist
}

func New() {
}
