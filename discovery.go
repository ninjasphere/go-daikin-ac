package daikin

import (
	"log"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/mostlygeek/arp"
)

var maxDatagramSize = 8192
var udpHost = "224.0.0.1:30050"
var udpClient = "224.0.0.1:30000"

func Discover(discoverInterval time.Duration) (chan AC, error) {

	seen := map[string]*wirelessAC{}

	found := make(chan AC)

	clientAddr, err := net.ResolveUDPAddr("udp", udpClient)
	if err != nil {
		log.Fatal(err)
	}
	l, err := net.ListenUDP("udp", clientAddr)
	if err != nil {
		return nil, err
	}

	hostAddr, err := net.ResolveUDPAddr("udp", udpHost)
	if err != nil {
		log.Fatal(err)
	}

	ping := func() error {
		_, err = l.WriteToUDP([]byte("DAIKIN_UDP/common/basic_info\n"), hostAddr)
		return err
	}

	arp := func() error {
		// We send out a broadcast pings first...

		ips, err := broadcastAddresses()
		if err != nil {
			return err
		}

		for _, ip := range ips {
			//log.Printf("Broadcast pinging %s", ip)
			cmd := exec.Command("ping", "-w", "5s", "-b", ip)
			err := cmd.Run()
			if err != nil {
				log.Printf("Failed to broadcast ping %s: %s", ip, err)
			}
		}

		for ip := range arp.Table() {
			mac := arp.Search(ip)

			if strings.HasPrefix(strings.ToUpper(mac), "90:B6:86") { // Murata Manufacturing Co., Ltd.
				log.Printf("Found a potential daikin AC: %s (mac: %s)", ip, mac)

				ac := NewWirelessAC(ip)
				info, err := ac.RefreshBasicInfo()

				if err == nil && info != nil && info.Ret == "OK" && info.Type == "aircon" {
					if existing, ok := seen[info.Id]; ok {
						existing.host = ip
					} else {
						seen[info.Id] = ac
						found <- ac
					}
				}

			} else {
				//log.Printf("Fail: %s %s", mac, ip)
			}

		}

		return nil
	}

	go func() {
		var t *time.Timer
		t = time.AfterFunc(0, func() {

			//log.Printf("Searching using ARP table...")
			if err := arp(); err != nil {
				log.Printf("Failed to discover daikin ACs using arp table: %s", err)
			}

			time.Sleep(time.Second * 5)

			//log.Printf("Searching using UDP...")
			if err := ping(); err != nil {
				log.Printf("Failed to discover daikin ACs: %s", err)
			}

			t.Reset(discoverInterval)
		})
	}()

	l.SetReadBuffer(maxDatagramSize)

	go func() {

		for {
			b := make([]byte, maxDatagramSize)
			n, src, err := l.ReadFromUDP(b)
			if err != nil {
				log.Printf("Daikin Discovery Error: ReadFromUDP failed: %s", err)
				continue
			}

			//log.Printf("Received response: %s", string(b[0:n]))

			info := &BasicInfo{}

			mapBytes(info, b[0:n])

			if existing, ok := seen[info.Id]; ok {
				existing.host = src.IP.String() + ":80"
			} else {
				ac := NewWirelessAC(src.IP.String() + ":80")
				ac.info = info

				seen[info.Id] = ac

				found <- ac
			}

		}
	}()

	if err := ping(); err != nil {
		return nil, err
	}

	return found, nil
}

func broadcastAddresses() ([]string, error) {

	ips := []string{}

	ifaces, err := net.Interfaces()
	if err != nil {
		return ips, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return ips, err
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

			ip[3] = 255 // Make it a broadcast address. ES: Is this ok?
			ips = append(ips, ip.String())
		}
	}
	return ips, nil
}
