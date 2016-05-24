package router

import (
	"bufio"
	"io"
	"os"
	"strings"
)

var arpFile = "/proc/net/arp"

type ARPProvider struct {
	devs map[string]struct{}
}

func NewARPProvider(devices []string) *ARPProvider {
	m := make(map[string]struct{})
	for _, dev := range devices {
		m[dev] = struct{}{}
	}
	return &ARPProvider{
		devs: m,
	}
}

func (p *ARPProvider) Leases() ([]Lease, error) {
	f, err := os.Open(arpFile)
	if err != nil {
		return nil, err
	}
	var ls []Lease
	br := bufio.NewReader(f)
	// skip first line
	br.ReadString('\n')
	for {
		line, err := br.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		parts := strings.Fields(line)
		if len(parts) != 6 {
			continue
		}
		if _, ok := p.devs[parts[5]]; !ok {
			continue
		}
		ls = append(ls, Lease{
			MAC: parts[3],
			IP:  parts[0],
		})
	}
	return ls, nil
}
