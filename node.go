package globe

import "context"

var globalContext context.Context

func StartNode() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	globalContext = ctx
	return ctx, cancel
}
