package net

import ()

func Sender(
	handler ConnectionHandler,
	st *State,
) {
	bufs := make([][]byte, 1024)
	for i := range bufs {
		bufs[i] = make([]byte, 2000)
	}

	for {
		var err error

		n, err := st.tun.Read(bufs[0])
		if err != nil {
			st.log.Error("error reading tun", "err", err)
			continue
		}
		st.log.Debug("tun read", "n", n)

		for i := range st.conf.Peers {
			init, err := st.netw.EnsureOpen(i, handler, st)
			if err != nil {
				st.log.Error("error connecting to peer", "err", err)
				continue
			}
			if init {
				st.log.Info("dialed peer", "peer", st.conf.Peers[i])
			}

			_, err = st.netw.Peers[i].conn.Write(bufs[0][:n])
			if err != nil {
				st.log.Error("error sending to peer", "err", err)
				st.netw.Close(i, st.log)
				continue
			}
			st.log.Debug("peer write", "peer", st.conf.Peers[i])
		}
	}
}
