package router

import "github.com/stretchr/testify/assert"
import "io/ioutil"
import "testing"

const fixture_leases = `
0 00:11:22:33:44:55 192.168.0.1 pc1 ff:ed:10:bd:b8:00:02:00:00:ab:11:04:2f:88:c8:b3:5e:6f:b7
0 00:11:22:33:44:66 192.168.0.2 pc2 ff:13:69:93:b7:00:01:00:01:1c:d3:ee:ce:00:27:13:69:93:b7
`

func writeTempFile(s string) (string, error) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return "", err
	}
	defer f.Close()
	err = ioutil.WriteFile(f.Name(), []byte(s), 0666)
	if err != nil {
		return "", err
	}
	return f.Name(), nil
}

func TestDNSMasqLeases(t *testing.T) {
	assert := assert.New(t)
	name, err := writeTempFile(fixture_leases)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	var p HostProvider = NewDNSMasqLeaseProvider(name)
	ls, err := p.Hosts()
	if err != nil {
		t.Fatalf("Error getting leases: %s", err)
	}
	assert.Equal(2, len(ls), "Need 2 leases")
	assert.Equal(Host{
		IP:   "192.168.0.1",
		MAC:  "00:11:22:33:44:55",
		Name: "pc1",
	}, ls[0], "Invalid lease")
	assert.Equal(Host{
		IP:   "192.168.0.2",
		MAC:  "00:11:22:33:44:66",
		Name: "pc2",
	}, ls[1], "Invalid lease")
}
