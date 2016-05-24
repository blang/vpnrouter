package router

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

var arpFixture = `IP address       HW type     Flags       HW address            Mask     Device
10.10.10.1      0x1         0x2         00:01:02:03:04:05     *        br0
10.10.11.1      0x1         0x2         00:01:02:03:04:06	  *        br1
10.10.12.1      0x1         0x0         00:01:02:03:04:07	  *        br2
`

func TestARPProvider(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	f.Close()
	defer os.Remove(f.Name())
	err = ioutil.WriteFile(f.Name(), []byte(arpFixture), 0666)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	p := NewARPProvider([]string{"br0", "br2"}, f.Name())
	l, err := p.Hosts()
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	ls := []Host{
		{MAC: "00:01:02:03:04:05", IP: "10.10.10.1"},
		{MAC: "00:01:02:03:04:07", IP: "10.10.12.1"},
	}
	if !reflect.DeepEqual(ls, l) {
		t.Errorf("Expected:\n%s\nGot:\n%s", ls, l)
	}
}
