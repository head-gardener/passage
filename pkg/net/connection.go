package net

import "net"

type Connection struct {
	tcp *net.TCPConn
}

func (conn *Connection) Dial(addr *net.TCPAddr) (err error) {
	if conn.tcp != nil {
		conn.tcp.Close()
	}

	tcp, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return
	}

	conn.tcp = tcp
	return
}

func (conn *Connection) EnsureOpen(addr *net.TCPAddr) (init bool, err error) {
	if conn.tcp != nil {
		return false, nil
	}

	return true, conn.Dial(addr)
}

func (conn *Connection) Read(b []byte) (n int, err error) {
	n, err = conn.tcp.Read(b)
	if err != nil {
		return
	}

	return
}

func (conn *Connection) Write(b []byte) (n int, err error) {
	n, err = conn.tcp.Write(b)
	// TODO: reopen on err?
	if err != nil {
		return
	}

	return
}
