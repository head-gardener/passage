package net

import (
	"fmt"
	"net"
)

func Listen(st *State) {
	for {
		err := Accept(st)
		if err != nil {
			st.log.Error("accepting connection", "err", err)
		}
	}
}

func Accept(st *State) (err error) {
	tcp, err := st.listener.Accept()
	st.log.Debug("received connection")
	if err != nil {
		return fmt.Errorf("accepting connection: %w", err)
	}

	remote, ok := tcp.RemoteAddr().(*net.TCPAddr)
	if !ok {
		tcp.Close()
		return fmt.Errorf("coercing remote addr to TCPAddr: %v", tcp.RemoteAddr())
	}
	remoteAddr := remote.IP

	found := false

	i := 0
	for ; i < len(st.conf.Peers); i++ {
		if st.conf.Peers[i].Addr.IP.Equal(remoteAddr) {
			found = true
			break
		}
	}

	if !found {
		tcp.Close()
		return fmt.Errorf("unexpected connection from %s", remoteAddr.String())
	}

	st.log.Info("initialized connection", "remote", remoteAddr.String())
	st.updateLastSeen(i)
	if st.netw.IsOpen(i) {
		tcp.Close()
		return fmt.Errorf("peer already established connection")
	}

	pass := st.conf.GetSecret()

	err = st.netw.Peers[i].conn.Accept(tcp, pass)
	if err != nil {
		tcp.Close()
		return err
	}

	st.netw.startConnectionHandler(i, HandleConnection, st)

	return
}
