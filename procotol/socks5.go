package procotol

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/pkg/errors"

	"github.com/Howard0o0/shadowsocks-mini/tinylog"
)

func Socks5Auth(readBuf []byte) bool {
	// 	+----+----------+----------+
	// 	|VER | NMETHODS | METHODS  |
	// 	+----+----------+----------+
	// 	| 1  |    1     | 1 to 255 |
	// 	+----+----------+----------+

	if len(readBuf) < 3 || readBuf[0] != 0x05 {
		return false
	}
	nmethods := readBuf[1]
	if len(readBuf) != (2 + int(nmethods)) {
		tinylog.LogError("auth message len ilegal\n")
		return false
	}

	return true
}

func Socks5ReadProxyAddr(buf []byte) (string, error) {

	if len(buf) < 7 {
		return "", errors.New("invalid proxy addr len")
	}

	vers, cmd, atyp := int(buf[0]), int(buf[1]), int(buf[3])
	if vers != 5 || cmd != 1 {
		return "", fmt.Errorf("invalid version/cmd : %d,%d ", vers, cmd)
	}
	return parseAddr(buf[4:], atyp)

}

func parseAddr(buf []byte, atyp int) (string, error) {

	ipAddr := ""

	switch atyp {
	case 1:
		// x.x.x.x
		if len(buf) < net.IPv4len+2 {
			return "", errors.New("proxy ipv4 addr len illegal")
		}
		ipAddr = fmt.Sprintf("%d.%d.%d.%d", buf[0], buf[1], buf[2], buf[3])
		buf = buf[4:]
	case 3:
		// hostname
		hostnameLen := int(buf[0])
		buf = buf[1:]
		if len(buf) < hostnameLen+2 {
			return "", errors.New("proxy hostname addr len illegal")
		}
		ipAddr = string(buf[:hostnameLen])
		buf = buf[hostnameLen:]
	case 4:
		// ipv6
		ipAddr = net.IP(buf[0:net.IPv6len]).String()
		buf = buf[net.IPv6len:]
	default:
		return "", fmt.Errorf("invalid atyp : %x ", atyp)
	}

	port := binary.BigEndian.Uint16(buf[:2])
	if len(buf[2:]) > 0 {
		err := errors.New("proxy addr message too long")
		tinylog.LogError("%s\n", err)
		return "", err
	}

	return fmt.Sprintf("%s:%d", ipAddr, port), nil
}

func BuildSocks5AuthOkRsp() []byte {
	return []byte{0x05, 0x00}
}

func BuildSocks5ConnectOkRsp() []byte {
	return []byte{0x05, 0x0, 0x00, 0x01, 0, 0, 0, 0, 0, 0}

}
