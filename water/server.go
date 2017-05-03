package water

import "net"
import "time"
import "fmt"

type ServerConn struct {
	conn            net.PacketConn
	maxWatPerSecond int
	count           int
	closeConn       chan struct{}
}

func AttachServerConn(conn net.PacketConn) ServerConn {
	retv := ServerConn{
		conn:            conn,
		maxWatPerSecond: 10,
		count:           0,
		closeConn:       make(chan struct{}),
	}
	return retv
}

func (self *ServerConn) getInnerConn() net.PacketConn {
	return self.conn
}

func (self *ServerConn) WriteTo(b []byte, addr net.Addr) (n int, err error) {
	return self.conn.WriteTo(b, addr)
}

func (self *ServerConn) SetWriteDeadline(t time.Time) error {
	return self.conn.SetWriteDeadline(t)
}

func (self *ServerConn) SetReadDeadline(t time.Time) error {
	return self.conn.SetReadDeadline(t)
}

func (self *ServerConn) SetDeadline(t time.Time) error {
	err := self.SetWriteDeadline(t)
	if err != nil {
		return err
	}
	return self.SetReadDeadline(t)
}

func (self *ServerConn) ReadFrom(b []byte) (n int, addr net.Addr, err error) {
	for {
		n, addr, err := self.conn.ReadFrom(b)
		if err != nil {
			return n, addr, err
		}
		if IsValidWata(b[:n]) {
			req, err := ParseWata(b[:n])
			if err != nil {
				fmt.Println("Parse error!")
				continue
			}

			var watrPacket [WATR_PACKET_SIZE]byte
			MakeWatrPacket(*addr.(*net.UDPAddr), watrPacket[:])
			self.WriteTo(watrPacket[:], &req)
			fmt.Println("WATR packet sent to:", req.IP, "With content", watrPacket)
			continue
		}
		if IsValidWate(b[:n]) {
			var watePacket = b[0:WATE_PACKET_SIZE]
			XorUdpAddr(*addr.(*net.UDPAddr), watePacket[4:])
			self.WriteTo(watePacket[:], addr)
			continue
		}
		return n, addr, err
	}
}

func (self *ServerConn) Close() error {
	close(self.closeConn)
	return self.conn.Close()
}

func (self *ServerConn) LocalAddr() net.Addr {
	return self.conn.LocalAddr()
}
