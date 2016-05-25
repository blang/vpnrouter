package router

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

type RulePersistence struct {
	base RuleProvider
	file string
	db   map[string]string
	mu   *sync.Mutex
}

func NewRulePersistence(base RuleProvider, file string) *RulePersistence {
	return &RulePersistence{
		base: base,
		file: file,
		db:   make(map[string]string),
		mu:   &sync.Mutex{},
	}
}

type persRule struct {
	IP    string
	Table string
}

// Init applies all saved rules
func (r *RulePersistence) Init() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	err := r.readFromFile()
	if err != nil {
		return err
	}
	r.applyRulesInDB()
	return nil
}

func (r *RulePersistence) readFromFile() error {
	f, err := os.Open(r.file)
	if err != nil {
		return err
	}
	br := bufio.NewReader(f)
	// skip first line
	br.ReadString('\n')
	for {
		line, err := br.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			continue
		}
		ip := strings.TrimSpace(parts[0])
		table := strings.TrimSpace(parts[1])
		r.db[ip] = table
	}
	return nil
}

func (r *RulePersistence) applyRulesInDB() {
	for ip, table := range r.db {
		r.base.Set(ip, table)
	}
}

func (r *RulePersistence) saveRulesToDB() error {
	var buf bytes.Buffer
	buf.WriteString("IP\tTable\n")
	for ip, table := range r.db {
		buf.WriteString(ip)
		buf.WriteString("\t")
		buf.WriteString(table)
		buf.WriteString("\n")
	}
	return ioutil.WriteFile(r.file, buf.Bytes(), 0644)
}

// Wrap base
func (r *RulePersistence) Rules() ([]Rule, error) {
	return r.base.Rules()
}

func (r *RulePersistence) Set(ip string, table string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.db[ip] = table
	r.saveRulesToDB()
	return r.base.Set(ip, table)
}
