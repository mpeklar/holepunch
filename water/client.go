package water

import "net"
import "time"
import "fmt"

type ClientConn struct {
	conn             net.PacketConn
	server           net.UDPAddr
	maxWataPerSecond int
	count            int
	pack             chan net.UDPAddr
	closeConn        chan struct{}
}

func AttachClientConn(conn net.PacketConn, wataServer net.UDPAddr) ClientConn {
	retv := ClientConn{
		conn:             conn,
		server:           wataServer,
		maxWataPerSecond: 10,
		count:            1,
		pack:             make(chan net.UDPAddr),
		closeConn:        make(chan struct{}),
		//wataResponder:    func(addr net.UDPAddr, conn net.PacketConn) {},
	}
	go retv.heartbeatServer(5 * time.Second)
	go retv.serveRequests()
	return retv
}

func (self *ClientConn) serveRequests() {
	subtract := time.NewTicker(time.Millisecond * 100)
	defer subtract.Stop()
	var wataPacket [WATA_PACKET_SIZE]byte
	for {
		if self.count == self.maxWataPerSecond {
			<-subtract.C
			self.count--
		}
		select {
		case addr := <-self.pack:
			MakeWataPacket(addr, wataPacket[:])
			_, err := self.conn.WriteTo(wataPacket[:], &self.server)
			if err != nil {
				continue
			}
		case <-subtract.C:
			if self.count != 0 {
				self.count--
			}
		case <-self.closeConn:
			return
		}
	}
}

func (self *ClientConn) SetWatrNumber(count int) {
	self.count = count
}

func (self *ClientConn) heartbeatServer(interval time.Duration) {
	var ipport [10]byte
	MakeWatePacket(ipport[:])
	timer := time.NewTicker(interval)
	defer timer.Stop()
	for {
		select {
		case <-self.closeConn:
			return
		case <-timer.C:
			self.conn.WriteTo(ipport[:10], &self.server)
		}
	}
}

func (self *ClientConn) getInnerConn() net.PacketConn {
	return self.conn
}

func (self *ClientConn) WriteTo(b []byte, addr net.Addr) (n int, err error) {
	return self.conn.WriteTo(b, addr)
}

func (self *ClientConn) SetWriteDeadline(t time.Time) error {
	return self.conn.SetWriteDeadline(t)
}

func (self *ClientConn) SetReadDeadline(t time.Time) error {
	return self.conn.SetReadDeadline(t)
}

func (self *ClientConn) SetDeadline(t time.Time) error {
	err := self.SetWriteDeadline(t)
	if err != nil {
		return err
	}
	return self.SetReadDeadline(t)
}

func IsHolepunchPacket(b []byte) bool {
	return IsValidWath(b)
}

func ParseHolepunchPacket(b []byte) net.UDPAddr {
	if !IsHolepunchPacket(b) {
		panic("Invalid packet: You should always use IsHolepunchPacket")
	}
	addr, _ := ParseWath(b)
	return addr
}

func (self *ClientConn) ReadFrom(b []byte) (n int, addr net.Addr, err error) {
	for {
		n, addr, err := self.conn.ReadFrom(b)
		if err != nil {
			return n, addr, err
		}
		if IsValidWate(b[:n]) {
			fmt.Println("Received WATE!")
			continue
		}
		if !IsValidWatr(b[:n]) {
			return n, addr, err
		}
		req, err := ParseWatr(b[:n])
		if err != nil {
			panic(err)
		}
		fmt.Println("WATR recieved for", req.IP)
		if self.count == 0 {
			continue
		}
		//self.count -= 1
		n = MakeWathPacket(*addr.(*net.UDPAddr), b)
		self.WriteTo(b[:n], &req)
	}
}

func (self *ClientConn) RequestHolepunchFrom(addr net.UDPAddr) (e error) {
	fmt.Println("Address:", addr.IP)
	var packet [WATA_PACKET_SIZE]byte
	MakeWataPacket(addr, packet[:WATA_PACKET_SIZE])
	fmt.Println("Packet=", packet[:WATA_PACKET_SIZE])
	_, e = self.WriteTo(packet[:WATA_PACKET_SIZE], &self.server)
	fmt.Println("Holepunch sent, ip=", addr.IP, "port=", addr.Port)
	return
}

func (self *ClientConn) Close() error {
	close(self.closeConn)
	close(self.pack)
	return self.conn.Close()
}

func (self *ClientConn) LocalAddr() net.Addr {
	return self.conn.LocalAddr()
}

/*func (self *ClientConn) GetExteriorAddress(addr net.UDPAddr, timeout time.Duration) (net.UDPAddr, error) {
	var watePacket [WATE_PACKET_SIZE]byte
	time := time.NewTimer(timeout)
	defer time.Stop()
	MakeWatePacket(watePacket[:])
	_, err := self.conn.WriteTo(watePacket[:], &self.server)
	if err != nil {
		return net.UDPAddr{}, err
	}
	select {
	case <-time.C:
		return net.UDPAddr{}, errors.New("Echo timed out!")
	case addr := <-self.echo:
		return addr, nil
	}
}*/
