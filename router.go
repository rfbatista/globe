package globe

type Router struct{}

var router Router

func (r *Router) Route(p PID, m Message) {
	c, _ := registry.Get(p)
	c <- m
}

func PublishMessage(p PID, m Message) {
	router.Route(p, m)
}
