package net

import (
	"crypto/rand"
	"fmt"
	"net"
	"sync"

	"github.com/head-gardener/passage/pkg/bee2/belt"
	"github.com/head-gardener/passage/pkg/crypto"
)

const HeaderSize = len(belt.MAC{})
const SaltSize = len(belt.Key{})

// Connection is either:
// closed: tcp == nil
// open: tcp != nil
// transitioning: lock is locked
type Connection struct {
	tcp    net.Conn
	cipher crypto.Cipher
	lock   sync.Mutex
}

func (conn *Connection) String() string {
	return conn.tcp.RemoteAddr().String()
}

func (conn *Connection) IsOpen() bool {
	conn.lock.Lock()
	defer conn.lock.Unlock()
	return conn.isOpenSimple()
}

func (conn *Connection) isOpenSimple() bool {
	return conn.tcp != nil
}

func handshakeInitiator(tcp net.Conn, pass []byte) (cipher crypto.Cipher, err error) {
	// TODO: sign salt with master key
	salt := make([]byte, SaltSize)
	rand.Read(salt)
	cipher, err = crypto.InitCHE(pass, salt)
	if err != nil {
		return
	}
	_, err = tcp.Write(salt)
	if err != nil {
		return
	}
	return
}

func handshakeResponder(tcp net.Conn, pass []byte) (cipher crypto.Cipher, err error) {
	salt := make([]byte, SaltSize)
	n, err := tcp.Read(salt)
	if err != nil {
		return
	}
	if n != 32 {
		return nil, fmt.Errorf("incorrect salt length %d", n)
	}
	cipher, err = crypto.InitCHE(pass, salt)
	if err != nil {
		return
	}
	return
}

func (conn *Connection) Accept(tcp net.Conn, pass []byte) (err error) {
	conn.lock.Lock()
	defer conn.lock.Unlock()

	if conn.isOpenSimple() {
		return fmt.Errorf("already connected to %v", tcp.RemoteAddr())
	}

	cipher, err := handshakeResponder(tcp, pass)
	if err != nil {
		return err
	}

	conn.tcp = tcp
	conn.cipher = cipher
	return
}

func (conn *Connection) Dial(addr *net.TCPAddr, pass []byte) (err error) {
	conn.lock.Lock()
	defer conn.lock.Unlock()

	if conn.isOpenSimple() {
		return fmt.Errorf("already connected to %v", addr)
	}

	tcp, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return
	}

	cipher, err := handshakeInitiator(tcp, pass)
	if err != nil {
		tcp.Close()
		return
	}

	conn.tcp = tcp
	conn.cipher = cipher
	return
}

func (conn *Connection) Close() (closed bool, err error) {
	conn.lock.Lock()
	defer conn.lock.Unlock()

	if !conn.isOpenSimple() {
		return false, fmt.Errorf("already closed")
	}

	err = conn.tcp.Close()
	if err != nil {
		return
	}
	closed = true

	conn.cipher.Finalize()

	conn.tcp = nil
	conn.cipher = nil
	return
}

func (conn *Connection) Read(b []byte) (n int, err error) {
	n, err = conn.tcp.Read(b)
	if err != nil {
		return
	}

	var mac crypto.MAC
	body := n - HeaderSize
	copy(mac[:], b[body:n])
	err = conn.cipher.Unwrap(b[:body], b[:body], nil, mac)
	if err != nil {
		return
	}

	return body, nil
}

func (conn *Connection) Write(b []byte) (n int, err error) {
	// TODO: zero pads and length in header
	full := len(b) + HeaderSize
	if cap(b) < full {
		return 0, fmt.Errorf(
			"buffer capacity is too small for header: %d have, %d needed",
			cap(b),
			full,
		)
	}

	err = conn.cipher.Wrap(b, b, nil, b[len(b):full])
	if err != nil {
		return
	}

	n, err = conn.tcp.Write(b[:full])
	if err != nil {
		return
	}

	return len(b), nil
}
