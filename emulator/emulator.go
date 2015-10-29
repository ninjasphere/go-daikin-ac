package emulator

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/ninjasphere/go-daikin-ac"
)

type emulatedWirelessAC struct {
	daikin.ControlState
	daikin.SensorState
}

func StartWirelessAC(port int) {
	ac := &emulatedWirelessAC{}
	ac.start(port)
}

var maxDatagramSize = 8192
var srvAddr = ":30050"
var srvAddr2 = "224.0.0.1:30000"

func (d *emulatedWirelessAC) start(port int) {

	go serveMulticastUDP(srvAddr, d.handleUDP)

	//go serveMulticastUDP(srvAddr2, d.handleUDP2)

	/*host, _ := os.Hostname()
	info := []string{"api=ninja"}
	service, _ := mdns.NewMDNSService(host, "_http._tcp", "", "", port, nil, info)

	// Create the mDNS server, defer shutdown
	server, _ := mdns.NewServer(&mdns.Config{Zone: service})
	defer server.Shutdown()*/

	//spew.Dump("Created mdns service", service)

	/*	http.HandleFunc("/aircon/get_control_info", d.getControlInfo)
		http.HandleFunc("/aircon/set_control_info", d.setControlInfo)
		http.HandleFunc("/aircon/get_sensor_info", d.getSensorInfo)*/

	go func() {
		for {
			time.Sleep(time.Second * 2)
			d.TargetTemperature += 0.5
			d.TargetHumidity++

			d.InsideTemperature += 0.3
			d.InsideHumidity++
			d.OutsideTemperature += 0.6
		}
	}()

	log.Printf("Starting emulated Daikin Wireless AC on port %d", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), d)
}

func (d *emulatedWirelessAC) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hah, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Printf("%s", err)
	}

	log.Printf("%s", hah)
	spew.Dump("incoming request", r.URL)

	spew.Dump("request", r)

	if r.URL.String() == "/aircon/get_control_info" {
		d.getControlInfo(w, r)
	} else if r.URL.String() == "/aircon/set_control_info" {
		d.setControlInfo(w, r)
	} else if r.URL.String() == "/aircon/get_sensor_info" {
		d.getSensorInfo(w, r)
	} else {
		spew.Dump("Unknown!")

		fmt.Fprint(w, "ret=OK")
	}
}

func (d *emulatedWirelessAC) handleUDP(conn *net.UDPConn, src *net.UDPAddr, n int, b []byte) {
	log.Println(n, "udp bytes read from", src)

	log.Println(hex.Dump(b[:n]))

	/*addr, err := net.ResolveUDPAddr("udp", "192.168.12.254:30000")
	if err != nil {
		log.Fatal(err)
	}*/

	i, err := conn.WriteToUDP([]byte("ret=OK,type=aircon,reg=th,dst=1,ver=2_2_5,pow=0,err=0,location=0,name=Ninja%20Blocks%20Wuz%20Here,icon=1,method=polling,port=30050,id=000005zr,pw=pyrbtetn,lpw_flag=0,adp_kind=2,pv=0,cpv=0,led=0,en_setzone=1,mac=FCC2DE460A5C,adp_mode=run\n"), src)
	spew.Dump("Wrote message", i, err)

}

func (d *emulatedWirelessAC) handleUDP2(src *net.UDPAddr, n int, b []byte) {
	log.Println(n, "udp bytes read from", src)
	log.Println(hex.Dump(b[:n]))
}

func serveMulticastUDP(a string, h func(*net.UDPConn, *net.UDPAddr, int, []byte)) {
	addr, err := net.ResolveUDPAddr("udp", a)
	if err != nil {
		log.Fatal(err)
	}

	//en0, _ := net.InterfaceByName("lo0")
	l, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}

	l.SetReadBuffer(maxDatagramSize)
	for {
		b := make([]byte, maxDatagramSize)
		n, src, err := l.ReadFromUDP(b)
		if err != nil {
			log.Fatal("ReadFromUDP failed:", err)
		}
		h(l, src, n, b)
	}
}

func (d *emulatedWirelessAC) getControlInfo(w http.ResponseWriter, r *http.Request) {
	out := "ret=OK"

	for k, v := range d.ControlState.GetWirelessValues() {
		out += fmt.Sprintf(",%s=%s", k, v[0])
	}

	//spew.Dump(out)

	fmt.Fprint(w, out)
}

func (d *emulatedWirelessAC) getModelInfo(w http.ResponseWriter, r *http.Request) {
	out := "ret=OK"

	//spew.Dump(out)

	fmt.Fprint(w, out)
}

func (d *emulatedWirelessAC) getSensorInfo(w http.ResponseWriter, r *http.Request) {
	out := "ret=OK"

	for k, v := range d.SensorState.GetWirelessValues() {
		out += fmt.Sprintf(",%s=%s", k, v[0])
	}

	fmt.Fprint(w, out)
}

func (d *emulatedWirelessAC) setControlInfo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	spew.Dump("CONTROL INFO", r.Form)

	d.ControlState.ParseWirelessValues(r.Form)

	fmt.Fprint(w, "ret=OK")
}
