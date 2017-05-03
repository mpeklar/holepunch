package water

import "net"
import "fmt"

func blit(source []byte, dest []byte) {
	for i := range source {
		dest[i] = source[i]
	}
}

func xorblit(source []byte, dest []byte) {
	for i := range source {
		dest[i] ^= source[i]
	}
}

func WriteUdpAddr(addr net.UDPAddr, dest []byte) (size int) {
	ip := addr.IP[len(addr.IP)-4:]
	for i := 0; i < 4; i++ {
		fmt.Println("Step", i, ":", ip[i])
		dest[i] = ip[i]
	}
	iplen := 4
	dest[iplen] ^= byte(addr.Port >> 8)
	dest[iplen+1] ^= byte(addr.Port & 0xFF)
	return iplen + 2
}

func XorUdpAddr(addr net.UDPAddr, dest []byte) (size int) {
	ip := addr.IP[len(addr.IP)-4:]
	for i := 0; i < 4; i++ {
		dest[i] ^= ip[i]
		//fmt.Println(addr.IP[i])
	}
	iplen := 4
	dest[iplen] ^= byte(addr.Port >> 8)
	dest[iplen+1] ^= byte(addr.Port & 0xFF)
	return iplen + 2
}

func MakeWataPacket(addr net.UDPAddr, dest []byte) (size int) {
	blit([]byte("WATA"), dest)
	return 4 + WriteUdpAddr(addr, dest[4:])
}

func MakeWathPacket(addr net.UDPAddr, dest []byte) (size int) {
	blit([]byte("WATH"), dest)
	return 4 + WriteUdpAddr(addr, dest[4:])
}

func MakeWatrPacket(addr net.UDPAddr, dest []byte) (size int) {
	blit([]byte("WATR"), dest)
	return 4 + WriteUdpAddr(addr, dest[4:])
}

func MakeWatePacket(dest []byte) (size int) {
	blit([]byte("WATE"), dest)
	for i := 5; i < 10; i++ {
		dest[i] = 0
	}
	return 8
}
