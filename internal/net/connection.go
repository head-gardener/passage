package net

import (
	"fmt"
	"net"
	"sync"

	"github.com/flynn/noise"

	"github.com/head-gardener/passage/internal/handshake"
	"github.com/head-gardener/passage/pkg/bee2/belt"
)

const HeaderSize = len(belt.MAC{})
const SaltSize = len(belt.Key{})

// Connection is either:
// closed: tcp == nil
// open: tcp != nil
// transitioning: lock is locked
type Connection struct {
	tcp  net.Conn
	cx   *noise.CipherState
	cr   *noise.CipherState
	lock sync.Mutex
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

func handshakeInitiator(tcp net.Conn, pass []byte) (cx *noise.CipherState, cr *noise.CipherState, err error) {
	hs, err := handshake.Init(true, pass)
	if err != nil {
		return
	}

	msg, _, _, err := hs.WriteMessage(nil, nil)
	if err != nil {
		return
	}
	_, err = tcp.Write(msg)
	if err != nil {
		return
	}

	// `<- e` is the same length as `-> e`
	n, err := tcp.Read(msg)
	if err != nil {
		return
	}
	_, cx, cr, err = hs.ReadMessage(nil, msg[:n])
	return
}

func handshakeResponder(tcp net.Conn, pass []byte) (cx *noise.CipherState, cr *noise.CipherState, err error) {
	hs, err := handshake.Init(false, pass)
	if err != nil {
		return
	}
	msg := make([]byte, 4096)

	n, err := tcp.Read(msg)
	if err != nil {
		return
	}
	_, _, _, err = hs.ReadMessage(nil, msg[:n])
	if err != nil {
		return
	}

	msg, cx, cr, err = hs.WriteMessage(nil, nil)
	if err != nil {
		return
	}
	_, err = tcp.Write(msg)
	return
}

func (conn *Connection) Accept(tcp net.Conn, pass []byte) (err error) {
	conn.lock.Lock()
	defer conn.lock.Unlock()

	if conn.isOpenSimple() {
		return fmt.Errorf("already connected to %v", tcp.RemoteAddr())
	}

	cx, cr, err := handshakeResponder(tcp, pass)
	if err != nil {
		tcp.Close()
		return err
	}

	conn.tcp = tcp
	conn.cx = cx
	conn.cr = cr
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

	cx, cr, err := handshakeInitiator(tcp, pass)
	if err != nil {
		tcp.Close()
		return
	}

	conn.tcp = tcp
	conn.cx = cx
	conn.cr = cr
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

	conn.tcp = nil
	conn.cx = nil
	conn.cr = nil
	return
}

func (conn *Connection) ReadWithTotal(b []byte) (n int, total int, err error) {
	n, err = conn.tcp.Read(b)
	if err != nil {
		return
	}

	msg, err := conn.cr.Decrypt(b[:0], nil, b[:n])
	if err != nil {
		return
	}

	return len(msg), n, err
}

func (conn *Connection) Read(b []byte) (n int, err error) {
	n, _, err = conn.ReadWithTotal(b)
	return
}

func (conn *Connection) Write(b []byte) (n int, err error) {
	if full := len(b) + belt.MACSize; cap(b) < full {
		return 0, fmt.Errorf(
			"buffer capacity is too small for header: %d have, %d needed",
			cap(b),
			full,
		)
	}

	msg, err := conn.cr.Encrypt(b[:0], nil, b)
	if err != nil {
		return
	}

	n, err = conn.tcp.Write(msg)
	if err != nil {
		return
	}

	return len(msg), nil
}
