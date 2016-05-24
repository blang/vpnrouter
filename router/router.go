package router

type Route struct {
	IP    string
	Table string
	Lease Host
}

type Router interface {
	Routes() ([]Route, error)
	SetRoute(ip string, table string) error
}

type VPNRouter struct {
	lp HostProvider
	rp RuleProvider
}

func NewVPNRouter(lp HostProvider, rp RuleProvider) *VPNRouter {
	return &VPNRouter{
		lp: lp,
		rp: rp,
	}
}

func (r *VPNRouter) Routes() ([]Route, error) {
	ls, err := r.lp.Hosts()
	if err != nil {
		return nil, err
	}
	rs, err := r.rp.Rules()
	if err != nil {
		return nil, err
	}
	rsMap := ruleMap(rs)

	var routes []Route
	var table string
	for _, l := range ls {
		table = "null"
		if rule, ok := rsMap[l.IP]; ok {
			table = rule.Table
		}
		routes = append(routes, Route{
			IP:    l.IP,
			Table: table,
			Lease: l,
		})
	}
	return routes, nil
}

func ruleMap(rs []Rule) map[string]Rule {
	m := make(map[string]Rule)
	for _, r := range rs {
		m[r.IP] = r
	}
	return m
}

func (r *VPNRouter) SetRoute(ip string, table string) error {
	return r.rp.Set(ip, table)
}
