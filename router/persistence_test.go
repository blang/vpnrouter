package router

import (
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
)

type mockRuleProvider struct {
	getFn func() ([]Rule, error)
	setFn func(ip, table string) error
}

func (m *mockRuleProvider) Rules() ([]Rule, error) {
	return m.getFn()
}

func (m *mockRuleProvider) Set(ip, table string) error {
	return m.setFn(ip, table)
}

func TestRulePersistence(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	f.Close()
	defer os.Remove(f.Name())

	rules := []Rule{
		{IP: "1", Table: "1"},
		{IP: "2", Table: "2"},
	}

	expRules := []Rule{
		{IP: "3", Table: "3"},
		{IP: "4", Table: "4"},
	}

	//save set rules
	setrules := []Rule{}

	mock := &mockRuleProvider{
		getFn: func() ([]Rule, error) {
			return rules, nil
		},
		setFn: func(ip, table string) error {
			setrules = append(setrules, Rule{
				IP:    ip,
				Table: table,
			})
			return nil
		},
	}
	rp := NewRulePersistence(mock, f.Name())
	if err := rp.Init(); err != nil {
		t.Fatalf("Error on init: %s", err)
	}
	if rs, err := rp.Rules(); err != nil || !reflect.DeepEqual(rs, rules) {
		t.Errorf("Wrong rules (error: %s):  %s", err, rs)
	}

	if err := rp.Set("3", "3"); err != nil {
		t.Errorf("Error on set: %s", err)
	}
	if err := rp.Set("4", "4"); err != nil {
		t.Errorf("Error on set: %s", err)
	}

	if !reflect.DeepEqual(setrules, expRules) {
		t.Errorf("Invalid tables saved: %s", setrules)
	}

	bs, err := ioutil.ReadFile(f.Name())
	if err != nil {
		t.Fatalf("Error reading file: %s", err)
	}

	str := string(bs)
	indexSecondLine := strings.Index(str, "\n") + 1
	if str[indexSecondLine:] != "3\t3\n4\t4\n" {
		t.Errorf("Invalid file contents: %s", str)
	}

	// Reset set rules
	setrules = []Rule{}
	// Create new RP
	rp = NewRulePersistence(mock, f.Name())
	if err := rp.Init(); err != nil {
		t.Errorf("Error on second init: %s", err)
	}

	if !reflect.DeepEqual(setrules, expRules) {
		t.Errorf("Invalid tables imported: %s", setrules)
	}

}
