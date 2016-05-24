package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mock struct {
	rules  []Rule
	leases []Host
}

func (m mock) Hosts() ([]Host, error) {
	return m.leases, nil
}

func (m mock) Rules() ([]Rule, error) {
	return m.rules, nil
}

func (m mock) Set(ip string, table string) error {
	for i, r := range m.rules {
		if r.IP == ip {
			m.rules[i].Table = table
		}
	}
	return nil
}

func TestRoutes(t *testing.T) {
	m := mock{
		rules: []Rule{
			{IP: "127.0.0.1", Table: "table1"},
			{IP: "127.0.0.2", Table: "table2"},
		},
		leases: []Host{
			{MAC: "abc", IP: "127.0.0.1", Name: "pc1"},
			{MAC: "def", IP: "127.0.0.2", Name: "pc2"},
		},
	}
	var r Router = NewVPNRouter(m, m)
	rs, err := r.Routes()
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	assert := assert.New(t)
	assert.Equal(2, len(rs))
	assert.Equal(Route{
		IP:    "127.0.0.1",
		Table: "table1",
		Lease: m.leases[0],
	}, rs[0])
	assert.Equal(Route{
		IP:    "127.0.0.2",
		Table: "table2",
		Lease: m.leases[1],
	}, rs[1])
}

func TestSetRoute(t *testing.T) {
	m := mock{
		rules: []Rule{
			{IP: "127.0.0.1", Table: "table1"},
		},
		leases: []Host{
			{MAC: "abc", IP: "127.0.0.1", Name: "pc1"},
		},
	}
	var r Router = NewVPNRouter(m, m)
	err := r.SetRoute("127.0.0.1", "table3")
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	rs, err := r.Routes()
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	assert := assert.New(t)
	assert.Equal(1, len(rs))
	assert.Equal(Route{
		IP:    "127.0.0.1",
		Table: "table3",
		Lease: m.leases[0],
	}, rs[0])

}
