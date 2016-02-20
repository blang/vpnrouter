package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const fixture_rules = `
32745:  from 10.10.10.1 lookup vpn 
32747:  from 10.10.10.2 lookup defgw 
`

func TestParseRules(t *testing.T) {
	rs := parseRules(fixture_rules)
	assert := assert.New(t)
	assert.Equal(2, len(rs))
	assert.Equal(Rule{IP: "10.10.10.1", Table: "vpn"}, rs[0])
	assert.Equal(Rule{IP: "10.10.10.2", Table: "defgw"}, rs[1])
}
