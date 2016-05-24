package router

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

var staticNameFixture = `MAC 	Name
00:01:02:03:04:05     Host1
00:01:02:03:04:06	  Host2
`

func TestStaticNameProvider(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	f.Close()
	defer os.Remove(f.Name())
	err = ioutil.WriteFile(f.Name(), []byte(staticNameFixture), 0666)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	p := NewStaticNameProvider(f.Name())
	l, err := p.Hosts()
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	ls := []Host{
		{MAC: "00:01:02:03:04:05", Name: "Host1"},
		{MAC: "00:01:02:03:04:06", Name: "Host2"},
	}
	if !reflect.DeepEqual(ls, l) {
		t.Errorf("Expected:\n%s\nGot:\n%s", ls, l)
	}
}
