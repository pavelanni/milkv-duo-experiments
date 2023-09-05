package main

import (
	"context"
	"fmt"
	"log"
	"time"

	milkvduo "github.com/pavelanni/gpiod-milkvduo"
	"github.com/warthog618/gpiod"
)

func main() {
	pinLed := "PWR_GPIO21"
	pinButton := "PWR_GPIO18"

	echan := make(chan gpiod.LineEvent, 6)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	eh := func(evt gpiod.LineEvent) {
		select {
		case echan <- evt: // the expected path
		default:
			// if you want the handler to block, rather than dropping
			// events when the channel fills then <- ctx.Done() instead
			// to ensure that the handler can't be left blocked
			fmt.Printf("event chan overflow - discarding event")
		}
	}

	lineLed, err := milkvduo.PinLineID(pinLed)
	if err != nil {
		log.Fatal(err)
	}
	lineButton, err := milkvduo.PinLineID(pinButton)
	if err != nil {
		log.Fatal(err)
	}

	lLed, err := gpiod.RequestLine(lineLed.Chip, lineLed.Offset, gpiod.AsOutput())
	if err != nil {
		log.Fatal(err)
	}
	lButton, err := gpiod.RequestLine(lineButton.Chip, lineButton.Offset, gpiod.WithBothEdges, gpiod.WithEventHandler(eh))
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := lLed.Reconfigure(gpiod.AsInput)
		if err != nil {
			log.Fatal(err)
		}
		lLed.Close()
		err = lButton.Reconfigure(gpiod.AsInput)
		if err != nil {
			log.Fatal(err)
		}
		lButton.Close()
	}()

	done := false
	for !done {
		select {
		// depending on the application other cases could deal with other channels
		case evt := <-echan:
			ledEvent(evt, lLed)
		case <-ctx.Done():
			fmt.Println("exiting...")
			lLed.Close()
			lButton.Close()
			done = true
		}
	}
	cancel()
}

func ledEvent(evt gpiod.LineEvent, l *gpiod.Line) {
	switch evt.Type {
	case gpiod.LineEventRisingEdge:
		ledDown(l)
	case gpiod.LineEventFallingEdge:
		ledUp(l)
	}
}

func ledUp(l *gpiod.Line) {
	l.SetValue(1)
}

func ledDown(l *gpiod.Line) {
	l.SetValue(0)
}
