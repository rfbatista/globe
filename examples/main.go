package main

import (
	"globe"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	_, cancel := globe.StartNode()
	defer cancel()
	p1 := globe.SpawnWorker(func(c *globe.Context) error {
		slog.Info("message received", "id", string(c.Body()), "my_pid", c.MyPID().ID)
		return nil
	})
	p2 := globe.SpawnWorker(func(c *globe.Context) error {
		slog.Info("message received", "id", string(c.Body()), "my_pid", c.MyPID().ID)
		return nil
	})
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-sigChan:
			slog.Info("received termination signal. Cancelling context...")
			return
		case <-ticker.C:
			globe.PublishMessage(p1, globe.NewActorMessage([]byte(p1.ID)))
			globe.PublishMessage(p2, globe.NewActorMessage([]byte(p2.ID)))
		}
	}
}
