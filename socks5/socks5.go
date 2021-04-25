package socks5

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/Howard0o0/shadowsocks-mini/util"
	"github.com/pkg/errors"
)

// SOCKS request commands as defined in RFC 1928 section 4.
const (
	CmdConnect      = 1
	CmdBind         = 2
	CmdUDPAssociate = 3
)

// SOCKS address types as defined in RFC 1928 section 5.
const (
	AtypIPv4       = 1
	AtypDomainName = 3
	AtypIPv6       = 4
)

// MaxAddrLen is the maximum size of SOCKS address in bytes.
const (
	MaxAddrLen = 1 + 1 + 255 + 2
	PortLen    = 2
)

func Auth(rw io.ReadWriter) error {

	// Read RFC 1928 for request and reply structure and sizes.
	buf := make([]byte, MaxAddrLen)
	// read VER, NMETHODS, METHODS
	if _, err := io.ReadFull(rw, buf[:2]); err != nil {
		return err
	}
	nmethods := buf[1]
	if _, err := io.ReadFull(rw, buf[:nmethods]); err != nil {
		return err
	}
	// write VER METHOD
	if _, err := rw.Write([]byte{5, 0}); err != nil {
		return err
	}

	return nil
}

func ReadAddr(rw io.ReadWriter) ([]byte, error) {
	buf := make([]byte, MaxAddrLen)
	erw := util.ErrReadWrite{RW: rw}

	// read [VER CMD RSV ATYP DST.ADDR DST.PORT]
	erw.ReadFull(buf[:4])
	if vers, cmd := buf[0], buf[1]; vers != 5 || cmd != 1 {
		return nil, fmt.Errorf("invalid version/cmd : %d,%d ", vers, cmd)
	}

	atyp := int(buf[3])
	addr := bytes.NewBuffer([]byte{byte(atyp)})
	switch atyp {
	case AtypIPv4:
		erw.ReadFull(buf[:net.IPv4len+PortLen])
		addr.Write(buf[:net.IPv4len+PortLen])
	case AtypDomainName:
		erw.ReadFull(buf[1:2])
		len := int(buf[1])
		addr.WriteByte(buf[1])
		erw.ReadFull(buf[:len+PortLen])
		addr.Write(buf[:len+PortLen])
	case AtypIPv6:
		erw.ReadFull(buf[:net.IPv6len+PortLen])
		addr.Write(buf[:net.IPv6len+PortLen])
	default:
		return nil, fmt.Errorf("invalid atyp : %d", atyp)
	}
	erw.Write(BuildSocks5ConnectOkRsp())

	return addr.Bytes(), erw.Err

}

func ParseAddr(buf []byte) (string, error) {

	if len(buf) < 3 {
		return "", errors.New("len illegal")
	}
	atyp := buf[0]
	buf = buf[1:]
	ipAddr := ""

	switch atyp {
	case AtypIPv4:
		if len(buf) < net.IPv4len+2 {
			return "", errors.New("proxy ipv4 addr len illegal")
		}
		ipAddr = fmt.Sprintf("%d.%d.%d.%d", buf[0], buf[1], buf[2], buf[3])
		buf = buf[4:]
	case AtypDomainName:
		hostnameLen := int(buf[0])
		buf = buf[1:]
		if len(buf) < hostnameLen+2 {
			return "", errors.New("proxy hostname addr len illegal")
		}
		ipAddr = string(buf[:hostnameLen])
		buf = buf[hostnameLen:]
	case AtypIPv6:
		ipAddr = net.IP(buf[0:net.IPv6len]).String()
		buf = buf[net.IPv6len:]
	default:
		return "", fmt.Errorf("invalid atyp : %x ", atyp)
	}

	port := binary.BigEndian.Uint16(buf[:2])
	if len(buf[2:]) > 0 {
		return "", errors.New("proxy addr message too long")
	}

	return fmt.Sprintf("%s:%d", ipAddr, port), nil
}

func BuildSocks5AuthOkRsp() []byte {
	return []byte{0x05, 0x00}
}

func BuildSocks5ConnectOkRsp() []byte {
	return []byte{0x05, 0x0, 0x00, 0x01, 0, 0, 0, 0, 0, 0}

}
