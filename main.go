package main

import (
	"fmt"
	"log"
	"os"

	"github.com/karalabe/hid"
	"github.com/reiver/go-telnet"
)

// Pedal constants
const (
	PedalRight  byte = 0x4
	PedalMiddle byte = 0x2
	PedalLeft   byte = 0x1
)

type event struct {
	down  bool
	pedal byte
}

var (
	pedalLeftUp     = event{down: false, pedal: PedalLeft}
	pedalLeftDown   = event{down: true, pedal: PedalLeft}
	pedalMiddleUp   = event{down: false, pedal: PedalMiddle}
	pedalMiddleDown = event{down: true, pedal: PedalMiddle}
	pedalRightUp    = event{down: false, pedal: PedalRight}
	pedalRightDown  = event{down: true, pedal: PedalRight}
)

func createPedalEvent(prevBuf, buf []byte, pedal byte) event {
	if buf[0]&pedal == pedal && prevBuf[0]&pedal != pedal {
		return event{down: true, pedal: pedal}
	}

	if buf[0]&pedal != pedal && prevBuf[0]&pedal == pedal {
		return event{down: false, pedal: pedal}
	}

	return event{}
}

func eventLoop(eventChannel chan event) {
	deviceInfos := hid.Enumerate(0x5f3, 0x00ff)
	if len(deviceInfos) != 1 {
		log.Fatalf("Got wrong number of devices: %d", len(deviceInfos))
	}

	deviceInfo := deviceInfos[0]
	device, err := deviceInfo.Open()
	if err != nil {
		log.Fatalf("Oh no! %v", err)
	}

	buf := []byte{0}
	prevBuf := []byte{0}
	for {
		_, err := device.Read(buf)
		if err != nil {
			log.Fatalf("Oh no! %v", err)
		}

		if e := createPedalEvent(prevBuf, buf, PedalLeft); e.pedal == PedalLeft {
			eventChannel <- e
		}
		if e := createPedalEvent(prevBuf, buf, PedalMiddle); e.pedal == PedalMiddle {
			eventChannel <- e
		}
		if e := createPedalEvent(prevBuf, buf, PedalRight); e.pedal == PedalRight {
			eventChannel <- e
		}

		prevBuf[0] = buf[0]
	}
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Requries one argument: the song to be played")
	}
	myMP3 := os.Args[1]

	eventChannel := make(chan event)

	conn, err := telnet.DialTo("localhost:9001")
	if err != nil {
		log.Fatalf("are you running vlc with: vlc -I rc --rc-host localhost:9001")
	}
	conn.Write([]byte("add " + myMP3 + "\n"))

	rate := 1.0

	go eventLoop(eventChannel)
	for {
		e := <-eventChannel
		// do things!

		if e == pedalLeftDown {
			fmt.Println("Pedal Left down")
			switch {
			case rate == 1.0:
				rate = 0.5
			case rate == 0.5:
				rate = 3.0
			case rate == 3.0:
				rate = 1.0
			}
			conn.Write([]byte(fmt.Sprintf("rate %f\n", rate)))
			fmt.Println(rate)
		}
		if e == pedalLeftUp {
			fmt.Println("Pedal Left up")
			// conn.Write([]byte("rate 1.0\n"))
		}

		if e == pedalMiddleDown {
			fmt.Println("Pedal Middle down")
			conn.Write([]byte("play\n"))
		}
		if e == pedalMiddleUp {
			fmt.Println("Pedal Middle up")
			conn.Write([]byte("pause\n"))
		}

		if e == pedalRightDown {
			fmt.Println("Pedal Right down")
			conn.Write([]byte("rewind\nrewind\n"))
		}
		if e == pedalRightUp {
			fmt.Println("Pedal Right up")
		}
	}
}
