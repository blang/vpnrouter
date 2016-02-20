package router

import (
	"bufio"
	"os"
	"strings"
)

type Lease struct {
	MAC  string
	IP   string
	Name string
}

type LeaseProvider interface {
	Leases() ([]Lease, error)
}

type DNSMasqLeaseProvider struct {
	leaseFile string
}

func NewDNSMasqLeaseProvider(leaseFile string) *DNSMasqLeaseProvider {
	return &DNSMasqLeaseProvider{
		leaseFile: leaseFile,
	}
}

func (p *DNSMasqLeaseProvider) Leases() ([]Lease, error) {
	f, err := os.Open(p.leaseFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var leases []Lease
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		parts := strings.Split(sc.Text(), " ")
		if len(parts) != 5 {
			continue
		}
		leases = append(leases, Lease{
			MAC:  parts[1],
			IP:   parts[2],
			Name: parts[3],
		})
	}
	return leases, nil

}
