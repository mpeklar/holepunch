package main

import "fmt"
import "github.com/mpeklar/water"
import "net"

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":10000")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	UDPConn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	socket := water.AttachServerConn(UDPConn)

	defer socket.Close()

	buf := make([]byte, 10000)
	for {
		n, _, _ := socket.ReadFrom(buf)
		fmt.Println("Packet recieved!", string(buf[:n]))
	}

	water.ParseWatr(make([]byte, 10))
	fmt.Printf("Hello world!")
}
