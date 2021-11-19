package main

import (
	"github.com/Mortimor1/telegraf-discovery/internal/config"
	"github.com/Mortimor1/telegraf-discovery/internal/telegraf"
	"github.com/Mortimor1/telegraf-discovery/pkg/logging"
	"github.com/tatsushid/go-fastping"
	"net"
	"os/exec"
	"time"
)

func main() {
	// Init Logger
	logger := logging.GetLogger()

	// Load config
	cfg := config.GetConfig()

	for _, job := range cfg.Jobs {
		// Init Pinger
		p := fastping.NewPinger()
		p.MaxRTT = 10000
		ips, err := Hosts(job.Subnet)

		if err != nil {
			logger.Error(err)
		}

		for _, ip := range ips {
			err := p.AddIP(ip)
			if err != nil {
				logger.Fatal(err)
			}
		}

		pingHosts := make([]string, 0)

		p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
			pingHosts = append(pingHosts, addr.String())
		}
		p.OnIdle = func() {
			// Read Telegraf config
			t := telegraf.Telegraf{}
			t.LoadConfig(job.ConfigFile)
			// Check new ip address
			var exist bool
			save := false

			for _, ip := range pingHosts {
				for i, ping := range t.Inputs.Ping {
					exist = false
					for _, url := range ping.Urls {
						if ip == url {
							exist = true
							break
						}
					}
					if !exist {
						logger.Info("New ip address not found")
					}
					// Add new ip address
					if !exist {
						logger.Infof("Add ip address: %s", ip)
						ping.Urls = append(ping.Urls, ip)
						t.Inputs.Ping[i] = ping
						save = true
					}
				}
			}
			// Save config
			if save {
				t.SaveConfig(job.ConfigFile)

				// RESTART TELEGRAF
				logger.Infof("Restart docker container: %s", job.ContainerName)
				cmd := exec.Command("docker restart ", job.ContainerName)
				err := cmd.Run()
				if err != nil {
					logger.Fatal(err)
				}
			}
		}
		err = p.Run()
		if err != nil {
			logger.Fatal(err)
		}
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
