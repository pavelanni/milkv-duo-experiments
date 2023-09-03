package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/warthog618/gpiod"
)

func main() {
	offset := 21
	v := 0
	c, err := gpiod.NewChip("gpiochip4")
	if err != nil {
		log.Fatal(err)
	}

	l, err := c.RequestLine(offset, gpiod.AsOutput())
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		l.Reconfigure(gpiod.AsInput)
		l.Close()
	}()

	values := map[int]string{0: "inactive", 1: "active"}
	fmt.Printf("Set pin %d %s\n", offset, values[v])

	// capture exit signals to ensure pin is reverted to input on exit.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	for {
		select {
		case <-time.After(2 * time.Second):
			v ^= 1
			err := l.SetValue(v)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Set pin %d %s\n", offset, values[v])
		case <-quit:
			return
		}
	}

}
