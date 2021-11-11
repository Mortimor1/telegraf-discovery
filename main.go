package main

import (
	"fmt"
	"github.com/tatsushid/go-fastping"
	"net"
	"time"
)

func main() {
	p := fastping.NewPinger()

	ips, err := Hosts("10.10.0.0/23")
	for _, ip := range ips {
		p.AddIP(ip)
	}

	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		fmt.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
	}
	p.OnIdle = func() {
		fmt.Println("finish")
	}
	err = p.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func Hosts(cidr string) ([]string, error) {
	ip, subnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(subnet.Mask); subnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	// remove network address and broadcast address
	return ips[1 : len(ips)-1], nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
