package router

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type StaticNameProvider struct {
	filename string
}

func NewStaticNameProvider(name string) *StaticNameProvider {
	return &StaticNameProvider{
		filename: name,
	}
}

func (s *StaticNameProvider) Hosts() ([]Host, error) {
	f, err := os.Open(s.filename)
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
		if len(parts) != 2 {
			continue
		}
		ls = append(ls, Host{
			MAC:  parts[0],
			Name: parts[1],
		})
	}
	return ls, nil
}
