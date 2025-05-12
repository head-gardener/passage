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

	for i := range st.conf.Peers {
		if st.metrics != nil {
			st.metrics.TunnelStatus.WithLabelValues(st.conf.Peers[i].Addr.String()).Set(0)
		}
	}

	for {
		var err error

		n, err := st.tun.Read(bufs[0])
		if err != nil {
			st.log.Error("error reading tun", "err", err)
			continue
		}
		st.log.Debug("tun read", "n", n)
		if st.metrics != nil {
			st.metrics.PacketsReceivedOs.Inc()
			st.metrics.BytesReceivedOs.Add(float64(n))
		}

		for i := range st.conf.Peers {
			init, err := st.netw.EnsureOpen(i, handler, st)
			if err != nil {
				st.log.Error("error connecting to peer", "err", err)
				if st.metrics != nil {
					st.metrics.PacketsFailed.WithLabelValues(st.conf.Peers[i].Addr.String()).Inc()
				}
				continue
			}
			if init {
				st.log.Info("dialed peer", "peer", st.conf.Peers[i])
			}

			sent, err := st.netw.Peers[i].conn.Write(bufs[0][:n])
			if err != nil {
				st.log.Error("error sending to peer", "err", err)
				st.netw.Close(i, st)
				if st.metrics != nil {
					st.metrics.PacketsFailed.WithLabelValues(st.conf.Peers[i].Addr.String()).Inc()
				}
				continue
			}
			st.log.Debug("peer write", "peer", st.conf.Peers[i])
			if st.metrics != nil {
				st.metrics.PacketsSent.WithLabelValues(st.conf.Peers[i].Addr.String()).Inc()
				st.metrics.BytesSent.WithLabelValues(st.conf.Peers[i].Addr.String()).Add(float64(sent))
			}
		}
	}
}
