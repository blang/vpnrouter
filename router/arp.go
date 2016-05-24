package router

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type ARPProvider struct {
	devs    map[string]struct{}
	arpFile string
}

func NewARPProvider(devices []string, arpFile string) *ARPProvider {
	m := make(map[string]struct{})
	for _, dev := range devices {
		m[dev] = struct{}{}
	}
	return &ARPProvider{
		devs:    m,
		arpFile: arpFile,
	}
}

func (p *ARPProvider) Hosts() ([]Host, error) {
	f, err := os.Open(p.arpFile)
	if err != nil {
		return nil, err
	}
	var ls []Host
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
		ls = append(ls, Host{
			MAC: parts[3],
			IP:  parts[0],
		})
	}
	return ls, nil
}
