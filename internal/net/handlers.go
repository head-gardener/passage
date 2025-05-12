package net

import ()

type ConnectionHandler func(
	id int,
	done <-chan struct{},
	conn *Connection,
	st *State,
)

func HandleConnection(
	id int,
	done <-chan struct{},
	conn *Connection,
	st *State,
) {
	remote := conn.String()
	st.log.Debug("handler started", "remote", remote)
	bufs := make([][]byte, 1024)
	for i := range bufs {
		bufs[i] = make([]byte, 2000)
	}

	for {
		select {
		case <-done:
			st.log.Debug("closing connection handler", "remote", remote)
			return
		default:
		}
		n, total, err := conn.ReadWithTotal(bufs[0])
		if err != nil {
			st.log.Error("reading from peer", "err", err, "remote", remote)
			st.netw.Close(id, st)
			continue
		}
		st.log.Debug("peer read", "n", n, "remote", remote)
		st.updateLastSeen(id)
		if st.metrics != nil {
			st.metrics.PacketsReceived.WithLabelValues(remote).Inc()
			st.metrics.BytesReceived.WithLabelValues(remote).Add(float64(total))
		}

		n, err = st.tun.Write(bufs[0][:n])
		if err != nil {
			st.log.Error("writing to tun", "err", err, "remote", remote)
			continue
		}
		st.log.Debug("tun write", "n", n, "remote", remote)
		if st.metrics != nil {
			st.metrics.PacketsSentOs.Inc()
			st.metrics.BytesSentOs.Add(float64(n))
		}
	}
}
