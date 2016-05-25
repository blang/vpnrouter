package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/blang/vpnrouter/api"
	"github.com/blang/vpnrouter/router"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

var (
	//flagListen    = flag.String("listen", ":8080", "Listen addr")
	flagWebDir    = flag.String("web", "./web", "Path to static files")
	flagLeaseFile = flag.String("lease-file", "/var/lib/misc/dnsmasq.leases", "Lease file")
	flagARPFile   = flag.String("arp-file", "/proc/net/arp", "ARP file")
	flagNameFile  = flag.String("name-file", "./names.txt", "Static MAC to name mapping")
	flagDBFile    = flag.String("db-file", "./db.txt", "Database file")
	flagDevices   = flag.String("devices", "eth0,eth1", "Ethernet devices to get hosts from")
	flagAdminIPs  = flag.String("admin-ips", "127.0.0.1", "Admin IPs comma separated")
	flagTables    = flag.String("tables", "null=Gesperrt,defgw=KabelD", "Routing tables comma separated")
	flagDebug     = flag.Bool("debug", false, "Enable mock rules")
)

var (
	tables    []api.TableDef
	adminIPs  []string
	devices   []string
	leaseFile string
	nameFile  string
	arpFile   string
	dbFile    string
	listen    string
	webDir    string
)

func prepareFlags() {
	flag.Parse()

	// prepare tables
	tableParts := strings.Split(*flagTables, ",")
	if len(tableParts) == 0 {
		log.Fatal("No tables given")
	}
	for _, t := range tableParts {
		tp := strings.TrimSpace(t)
		nameTitle := strings.SplitN(tp, "=", 2)
		if len(nameTitle) != 2 {
			log.Printf("Ignore invalid table definition: %s", tp)
			continue
		}
		tables = append(tables, api.TableDef{
			Name: nameTitle[0],
			Text: nameTitle[1],
		})
	}

	ips := strings.Split(*flagAdminIPs, ",")
	for _, ip := range ips {
		ip = strings.TrimSpace(ip)
		adminIPs = append(adminIPs, ip)
	}

	devs := strings.Split(*flagDevices, ",")
	for _, dev := range devs {
		dev = strings.TrimSpace(dev)
		devices = append(devices, dev)
	}

	fi, err := os.Stat(*flagWebDir)
	if err != nil {
		log.Fatalf("Invalid webdir: %s", err)
	}
	if !fi.IsDir() {
		log.Fatalf("Webdir is not a directory")
	}
	webDir = *flagWebDir

	// check lease file
	f, err := os.Open(*flagLeaseFile)
	if err != nil {
		log.Fatalf("Error opening lease file: %s", err)
	}
	f.Close()
	leaseFile = *flagLeaseFile

	// check arp file
	f, err = os.Open(*flagARPFile)
	if err != nil {
		log.Fatalf("Error opening arp file: %s", err)
	}
	f.Close()
	arpFile = *flagARPFile

	// check name file
	f, err = os.Open(*flagNameFile)
	if err != nil {
		log.Fatalf("Error opening name file: %s", err)
	}
	f.Close()
	nameFile = *flagNameFile

	// no need to check dbFile for existence
	dbFile = *flagDBFile
}

func main() {
	prepareFlags()

	var ruleProv router.RuleProvider = router.NewIPRoute2RuleProvider()
	if *flagDebug {
		ruleProv = make(router.DummyRuleProvider)
	}

	// Add persistence layer
	persistence := router.NewRulePersistence(ruleProv, dbFile)
	persistence.Init()
	ruleProv = persistence

	dnsmasq := router.NewDNSMasqLeaseProvider(leaseFile)
	log.Printf("Devices: %s", devices)
	arp := router.NewARPProvider(devices, arpFile)
	staticNameProv := router.NewStaticNameProvider(nameFile)
	hostprov := router.HostMerger{arp, dnsmasq, staticNameProv}
	r := router.NewVPNRouter(hostprov, ruleProv)
	server := api.NewServer(r, api.NewIPAuth(adminIPs...), tables)
	apiMux := web.New()
	apiMux.Use(middleware.SubRouter)
	goji.Handle("/api/*", apiMux)
	apiMux.Get("/tables", server.GetTables)
	apiMux.Get("/routes", server.GetRoutes)
	apiMux.Post("/routes", server.SetRoute)

	goji.Get("/*", http.FileServer(http.Dir(webDir)))

	goji.Serve()
}
