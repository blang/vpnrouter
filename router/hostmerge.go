package router

// HostMerger merges the results from two HostProviders.
type HostMerger struct {
	First HostProvider

	// Backup Hostnames are used if duplicate found in prov2
	Backup HostProvider

	// Static names are used if found
	StaticName HostProvider
}

func (h HostMerger) Hosts() ([]Host, error) {
	hosts1, err := h.First.Hosts()
	if err != nil {
		return nil, err
	}
	hosts2, err := h.Backup.Hosts()
	if err != nil {
		return nil, err
	}
	hosts3, err := h.StaticName.Hosts()
	if err != nil {
		return nil, err
	}
	return mergeHosts(hosts1, hosts2, hosts3), nil
}

func mergeHosts(hosts1, hosts2, hosts3 []Host) []Host {
	var hosts []Host
	staticNameM := make(map[string]string) // MAC to Name
	for _, sn := range hosts3 {
		staticNameM[sn.MAC] = sn.Name
	}

	hm := make(map[string]Host)
	for _, h2 := range hosts2 {
		hm[h2.IP] = h2
	}
	for _, h1 := range hosts1 {
		if h2, ok := hm[h1.IP]; ok {
			h1.Name = h2.Name
			delete(hm, h1.IP)
		}
		if sname, ok := staticNameM[h1.MAC]; ok {
			h1.Name = sname
		}
		hosts = append(hosts, h1)
	}
	for _, h2 := range hm {
		if sname, ok := staticNameM[h2.MAC]; ok {
			h2.Name = sname
		}
		hosts = append(hosts, h2)
	}
	return hosts
}
