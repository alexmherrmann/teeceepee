package teeceepeego

import (
	"log"
	"net"
)

func DoListen() {

	// TODO AH: Make the ip address and port a param
	listen, listenErr := net.Listen("tcp", "127.0.0.1:4246")

	if listenErr != nil {
		panic("Couldn't bind to the port: " + listenErr.Error())
	}

	log.Print("Bound and listening!")
	for {
		if conn, connErr := listen.Accept(); connErr == nil {
			go func() {
				log.Println("Got a new connection, passing it to Deal")
				charn := Deal(conn)

				result := <-charn
				log.Println("Done Dealing!")
				if result.Err != nil {
					log.Printf("Had a problem: %s\n", result.Err.Error())
				}
				conn.Close()
			}()
		}

	}
}
