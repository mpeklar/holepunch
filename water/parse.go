package water

import "net"
import "bytes"
import "errors"

const WATA_PACKET_SIZE = 10
const WATE_PACKET_SIZE = 10
const WATR_PACKET_SIZE = 10
const WAT_ZEROES_START = 10

func isFilledWithZeroes(slice []byte) bool {
	for i := range slice {
		if slice[i] != 0 {
			return false
		}
	}
	return true
}

func isSemiValid(packet []byte, size int) bool {
	if len(packet) != size {
		return false
	}
	if !bytes.Equal(packet[0:3], []byte("WAT")) {
		return false
	}
	//	ip := net.IPv4(packet[4], packet[5], packet[6], packet[7])
	//	if !ip.IsGlobalUnicast() || ip {
	//		return false
	//	}
	return true
}

func waterParse(packet []byte) (addr net.UDPAddr) {
	ip := net.IPv4(packet[4], packet[5], packet[6], packet[7])
	port := int(packet[8])<<8 | int(packet[9])
	addr.IP = ip
	addr.Port = port
	addr.Zone = ""
	return addr
}

func IsValidWatr(packet []byte) bool {
	if !isSemiValid(packet, WATR_PACKET_SIZE) {
		return false
	}
	return packet[3] == 'R'
}

func IsValidWata(packet []byte) bool {
	if !isSemiValid(packet, WATA_PACKET_SIZE) {
		return false
	}
	return packet[3] == 'A'
}

func IsValidWate(packet []byte) bool {
	if !isSemiValid(packet, WATE_PACKET_SIZE) {
		return false
	}
	return packet[3] == 'E'
}

func IsValidWath(packet []byte) bool {
	if !isSemiValid(packet, WATE_PACKET_SIZE) {
		return false
	}
	return packet[3] == 'H'
}

func ParseWatr(packet []byte) (addr net.UDPAddr, err error) {
	if !IsValidWatr(packet) {
		return addr, errors.New("Invalid WATR packet!")
	}
	return waterParse(packet), nil
}

func ParseWath(packet []byte) (addr net.UDPAddr, err error) {
	if !IsValidWath(packet) {
		return addr, errors.New("Invalid WATR packet!")
	}
	return waterParse(packet), nil
}

func ParseWata(packet []byte) (addr net.UDPAddr, err error) {
	if !IsValidWata(packet) {
		return addr, errors.New("Invalid WATA packet!")
	}
	return waterParse(packet), nil
}

// No parse functions for Wate
