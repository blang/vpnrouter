package router

import (
	"os/exec"
	"strings"
	"sync"
)

type Rule struct {
	IP    string
	Table string
}

type RuleProvider interface {
	Rules() ([]Rule, error)
	Set(ip string, table string) error
}

type IPRoute2RuleProvider struct {
	sync.Mutex
}

func NewIPRoute2RuleProvider() *IPRoute2RuleProvider {
	return &IPRoute2RuleProvider{}
}

func (p *IPRoute2RuleProvider) Rules() ([]Rule, error) {
	b, err := exec.Command("ip", "rule", "show").Output()
	if err != nil {
		return nil, err
	}
	return parseRules(string(b)), nil
}

func parseRules(s string) []Rule {
	lines := strings.Split(s, "\n")
	var rules []Rule
	for _, line := range lines {
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 5 {
			continue
		}
		rules = append(rules, Rule{
			IP:    parts[2],
			Table: parts[4],
		})
	}
	return rules
}

func findByIP(rules []Rule, ip string) (Rule, bool) {
	for _, r := range rules {
		if r.IP == ip {
			return r, true
		}
	}
	return Rule{}, false
}

func (p *IPRoute2RuleProvider) Set(ip string, table string) error {
	oldRules, err := p.Rules()
	if err != nil {
		return err
	}
	p.Lock()
	defer p.Unlock()
	oldRule, found := findByIP(oldRules, ip)
	// Old Rule exists, delete
	if found {
		err = p.delRoute(ip, oldRule.Table)
		if err != nil {
			return err
		}
	}
	err = p.addRoute(ip, table)
	if err != nil {
		return err
	}
	return nil
}

func (p *IPRoute2RuleProvider) delRoute(ip string, table string) error {
	return exec.Command("ip", "rule", "del", "from", ip, "table", table).Run()
}

func (p *IPRoute2RuleProvider) addRoute(ip string, table string) error {
	return exec.Command("ip", "rule", "add", "from", ip, "table", table).Run()
}

type DummyRuleProvider map[string]string

func (p DummyRuleProvider) Rules() ([]Rule, error) {
	var rules []Rule
	for k, v := range p {
		rules = append(rules, Rule{
			IP:    k,
			Table: v,
		})
	}
	return rules, nil
}

func (p DummyRuleProvider) Set(ip string, table string) error {
	p[ip] = table
	return nil
}
