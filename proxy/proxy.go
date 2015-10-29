package proxy

import (
	"io"
	"net"
	"time"

	"github.com/ninjasphere/go-ninja/logger"
)

var log = logger.GetLogger("DaikinProxy")

var maxDatagramSize = 8192

var udpHost = "224.0.0.1:30050"
var udpClient = "224.0.0.1:30000"

func Start(udpBroadcast string, incomingInterface string) error {

	// NOTE: We only support one AC atm
	var lastACSeen *net.TCPAddr

	clientAddr, err := net.ResolveUDPAddr("udp", udpClient)
	if err != nil {
		panic(err)
	}
	incomingResponseConn, err := net.ListenUDP("udp", clientAddr)
	if err != nil {
		panic(err)
	}

	broadcastAddr, err := net.ResolveUDPAddr("udp", udpBroadcast)
	if err != nil {
		panic(err)
	}

	send := func(data []byte) {
		log.Infof("Forwarding %d byte search to %s", len(data), udpHost)

		_, err = incomingResponseConn.WriteToUDP(data, broadcastAddr)
		if err != nil {
			log.Errorf("Failed to resend search: %s", err)
		}
	}

	fromAddr, err := net.ResolveUDPAddr("udp", udpHost)
	if err != nil {
		panic(err)
	}

	incomingIface, _ := net.InterfaceByName(incomingInterface)

	incomingSearchConn, err := net.ListenMulticastUDP("udp", incomingIface, fromAddr)
	if err != nil {
		panic(err)
	}

	incomingSearchConn.SetReadBuffer(maxDatagramSize)

	responses := make(chan []byte)

	go func() {
		for {
			b := make([]byte, maxDatagramSize)
			n, src, err := incomingSearchConn.ReadFromUDP(b)
			if err != nil {
				log.Fatalf("incomingSearchConn failed:%s", err)
			}

			if !isSelf(src) {
				log.Infof("Incoming search from %s", src.String())
				send(b[:n])

				select {
				case r := <-responses:
					log.Infof("Got search response. Sending back to %s", src.String())
					incomingSearchConn.WriteToUDP(r, src)
				case <-time.After(time.Second * 3):
					log.Infof("Timed out (5s). No response to search from %s", src.String())
				}
			}

		}
	}()

	go func() {
		for {
			b := make([]byte, maxDatagramSize)
			n, src, err := incomingResponseConn.ReadFromUDP(b)
			if err != nil {
				log.Fatalf("incomingResponseConn failed:%s", err)
			}

			log.Infof("Incoming response from %s", src.String())

			lastACSeen, _ = net.ResolveTCPAddr("tcp", src.IP.String()+":80")

			select {
			case responses <- b[:n]:
			default:
				log.Infof("No-one listened to search response!")
			}

		}
	}()

	go func() {

		listener, err := net.Listen("tcp", getIfaceIP(incomingInterface)+":80")
		if err != nil {
			log.Fatalf("Failed to setup listener: %v", err)
		}

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Fatalf("ERROR: failed to accept listener: %v", err)
			}
			log.Infof("New http connection from %s", conn.RemoteAddr().String())
			if lastACSeen != nil {
				log.Infof("Forwarding http connection to %s", lastACSeen.String())

				acConn, err := net.Dial("tcp", lastACSeen.String())
				if err != nil {
					log.Errorf("Failed to connect to AC. Closing incoming connection.")
					conn.Close()
				} else {
					go forward(conn, acConn)
				}

			}

		}
	}()

	return nil
}

func getIfaceIP(iface string) string {

	i, err := net.InterfaceByName(iface)
	if err != nil {
		time.Sleep(time.Second * 5)
		panic(err)
	}

	addrs, err := i.Addrs()
	if err != nil {
		time.Sleep(time.Second * 5)
		panic(err)
	}
	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		if ip == nil || ip.IsLoopback() {
			continue
		}
		ip = ip.To4()
		if ip == nil {
			continue // not an ipv4 address
		}
		return ip.String()
	}

	panic("No IP for " + iface)

}

func forward(from net.Conn, to net.Conn) {
	go func() {
		defer from.Close()
		defer to.Close()
		io.Copy(from, to)
		log.Infof("Forwarding conn from %s to %s closed.", from.RemoteAddr().String(), to.RemoteAddr().String())
	}()
	go func() {
		defer from.Close()
		defer to.Close()
		io.Copy(to, from)
	}()
}

func isSelf(addr *net.UDPAddr) bool {

	ifaces, _ := net.Interfaces()

	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		for _, a := range addrs {
			if a.(*net.IPNet).IP.String() == addr.IP.String() {
				return true
			}
		}
	}
	return false
}
