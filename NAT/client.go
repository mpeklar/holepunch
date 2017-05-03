package main

import "fmt"
import "github.com/mpeklar/water"
import "net"
import "os"
import "flag"

//import "io"

var sendchan = make(chan *net.UDPAddr)

func receiver(socket net.PacketConn) {
	var b [60000]byte
	var destaddr *net.UDPAddr
	var started bool = false
	for {
		n, addr, err := socket.ReadFrom(b[:])
		if err != nil {
			fmt.Println("Read failed!")
			return
		}

		if !started {
			started = true
			sendchan <- addr.(*net.UDPAddr)
			destaddr = addr.(*net.UDPAddr)
		}
		if addr != destaddr {
			fmt.Println("Packet ignored")
			continue
		}
		fmt.Println("Packet received:", string(b[:n]), ":size=", n)

	}
}

func main() {

	holepunchServer := flag.String("w", "127.0.0.1:10000", "Use a WATA holepunch server")
	listen := flag.Bool("l", false, "Listen on port")

	flag.Parse()

	if flag.NArg() < 0 || flag.NArg() > 2 {
		fmt.Println("Wrong number of arguments specified")
	}
	var ip, port string
	if flag.NArg() == 2 {
		ip = flag.Arg(0)
		port = flag.Arg(1)
	} else if flag.NArg() == 1 && *listen {
		ip = ""
		port = flag.Arg(0)
	} else {
		fmt.Println("Wrong number of arguments specified")
		return
	}
	addrstr := ip + ":" + port
	fmt.Println(addrstr)

	addr, err := net.ResolveUDPAddr("udp", addrstr)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	var UDPConn *net.UDPConn

	if *listen {
		UDPConn, err = net.ListenUDP("udp", addr)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
	} else {
		dd, _ := net.ResolveUDPAddr("udp", "")
		UDPConn, err = net.ListenUDP("udp", dd)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
	}

	if *holepunchServer == "" {
		return
	}

	addr2, _ := net.ResolveUDPAddr("udp", *holepunchServer)
	socket := water.AttachClientConn(UDPConn, *addr2)

	defer socket.Close()

	go receiver(&socket)

	var destaddr *net.UDPAddr

	if !*listen {
		destaddr, _ = net.ResolveUDPAddr("udp", addrstr)
	}

	for {
		var b [1024]byte
		if destaddr == nil {
			destaddr = <-sendchan
		}
		n, _ := os.Stdin.Read(b[:])
		fmt.Println("Written!")
		if !*listen {
			_ = socket.RequestHolepunchFrom(*destaddr)
		}

		_, err := socket.WriteTo(b[:n], destaddr)
		if err != nil {
			fmt.Println("Cannot write:", err)
		}

	}
}
