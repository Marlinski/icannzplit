package util

import (
	"encoding/binary"
	"net"
)

// IToIP convert an Int into an IPAddress
func ITo2IP(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip
}
