package teeceepeego

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/alexmherrmann/amh_go_utils"
)

/*
This file defines commands that will cause the server
to do specific things when it sees them
*/

type BadCommand struct {
	message string
}

func (b BadCommand) Error() string {
	return b.message
}

func BadCommandErr(message string) error {
	return BadCommand{message}
}

type IoErr struct {
	because error
	message string
}

func (e *IoErr) Unwrap() error {
	return e.because
}

func (e *IoErr) Error() string {
	return e.message
}

//ioErrWrapper is a helper to wrap up a lower level error in something we can work with
func ioErrWrapper(becauseof error, message string, formats ...any) error {
	return &IoErr{
		message: fmt.Sprintf(message, formats...),
		because: becauseof,
	}
}

//DealResult is just an empty struct we return from Deal that will contain a possible error
type DealResult struct {
	// conn net.Conn
	Err error
}

func escape(input string) string {
	return strings.ReplaceAll(input, "\n", "\\n")
}

//Deal will read bytes from the connection and execute the command given
func Deal(deal net.Conn) <-chan DealResult {
	command := bytes.Buffer{}
	tmp := make([]byte, 1024)
	commchan := make(chan DealResult)

	go func() {
		for {
			n, err := deal.Read(tmp)
			log.Printf("Read %4d bytes: \"%s\" with err?: %t", n, escape(string(tmp[:n])), err != nil)
			if err != nil {
				log.Println("EOF!")
				if err != io.EOF {
					log.Println("read error:", err)
				}
				break
			}
			command.Write(tmp[:n])

			if n > 1 {
				// Last char
				if tmp[n-1] == '\n' {
					break
				}
			}

			if n == 0 {
				<-time.After(time.Second * 1)
			}

		}
		if strings.HasPrefix(command.String(), "drip") {
			args := strings.Split(strings.TrimSpace(command.String()), " ")[1:]

			bytes, bytesErr := strconv.Atoi(args[0])
			times, timesErr := strconv.Atoi(args[1])
			delay, delayErr := strconv.Atoi(args[2])

			if bytesErr == nil && timesErr == nil && delayErr == nil {
				// Push it back down
				commchan <- <-Drip(deal, bytes, times, delay)
			} else {
				msg := fmt.Sprintf("Invalid command \"%s\"", escape(command.String()))
				commchan <- DealResult{BadCommandErr(msg)}
			}
		}
	}()

	return commchan
}

//Drip will drip x bytes y times with a delay of z milliseconds between
func Drip(dripTo net.Conn, bytes, times, delay int) <-chan DealResult {
	// toWrite := amh_go_utils.CreateRandomString(bytes * times)

	// Don't actually so anything right now
	commchan := make(chan DealResult)

	go func() {
		numruns := 0
		for numruns < times {
			numruns++

			// Wait the delay
			<-time.After(time.Second * time.Duration(delay))

			tosend := amh_go_utils.CreateRandomString(bytes)
			// Total number of tosend bytes written
			total := 0
			for total < len(tosend) {
				// Write whatever we can
				written, writeErr := dripTo.Write([]byte(tosend[total:]))
				if writeErr != nil {
					commchan <- DealResult{ioErrWrapper(writeErr, "Writing drip %d", numruns)}
					return
				}
				total += written

			}
		}
		commchan <- DealResult{}
	}()
	return commchan
}
