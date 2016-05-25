package router

import (
	"bufio"
	"os"
	"sort"
	"strings"
)

type Host struct {
	MAC  string
	IP   string
	Name string
}

type ByHostname []Host

func (a ByHostname) Len() int           { return len(a) }
func (a ByHostname) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByHostname) Less(i, j int) bool { return a[i].Name < a[j].Name }

type HostProvider interface {
	Hosts() ([]Host, error)
}

type DNSMasqLeaseProvider struct {
	leaseFile string
}

func NewDNSMasqLeaseProvider(leaseFile string) *DNSMasqLeaseProvider {
	return &DNSMasqLeaseProvider{
		leaseFile: leaseFile,
	}
}

func (p *DNSMasqLeaseProvider) Hosts() ([]Host, error) {
	f, err := os.Open(p.leaseFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var leases []Host

	hostMap := make(map[string][]Host)

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		parts := strings.Split(sc.Text(), " ")
		if len(parts) != 5 {
			continue
		}
		hostMap[parts[1]] = append(hostMap[parts[1]], Host{
			MAC:  parts[1],
			IP:   parts[2],
			Name: parts[3],
		})
	}

	// Save last x-last leases
	x := 2
	for _, host := range hostMap {
		for i := len(host) - 1; i >= len(host)-x && i >= 0; i-- {
			leases = append(leases, host[i])
		}
	}
	sort.Sort(ByHostname(leases))
	return leases, nil

}
