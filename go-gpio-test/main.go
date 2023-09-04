package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	milkvduo "github.com/pavelanni/gpiod-milkvduo"
	"github.com/warthog618/gpiod"
)

func main() {
	pin := "PWR_GPIO21"
	v := 0
	lineId, err := milkvduo.PinLineID(pin)
	if err != nil {
		log.Fatal(err)
	}

	l, err := gpiod.RequestLine(lineId.Chip, lineId.Offset, gpiod.AsOutput())
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := l.Reconfigure(gpiod.AsInput)
		if err != nil {
			log.Fatal(err)
		}
		l.Close()
	}()

	values := map[int]string{0: "inactive", 1: "active"}
	fmt.Printf("Set pin %s %s\n", pin, values[v])

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
			fmt.Printf("Set pin %s %s\n", pin, values[v])
		case <-quit:
			return
		}
	}
}
