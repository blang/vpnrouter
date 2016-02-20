package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/blang/vpnrouter/router"
)

func NewServer(router router.Router, auth AuthProvider, tables []TableDef) *Server {
	return &Server{
		router: router,
		auth:   auth,
		tables: tables,
	}
}

type Server struct {
	router router.Router
	auth   AuthProvider
	tables []TableDef
}

type TableDef struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

type routesResp struct {
	IP       string `json:"ip"`
	Table    string `json:"table"`
	Hostname string `json:"hostname"`
	MAC      string `json:"mac"`
}

type ByHostname []routesResp

func (a ByHostname) Len() int           { return len(a) }
func (a ByHostname) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByHostname) Less(i, j int) bool { return a[i].Hostname < a[j].Hostname }

func parseIP(addr string) string {
	parts := strings.SplitN(addr, ":", 2)
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}

func routeToRespRoute(r router.Route) routesResp {
	return routesResp{
		IP:       r.IP,
		Table:    r.Table,
		Hostname: r.Lease.Name,
		MAC:      r.Lease.MAC,
	}
}

func (s *Server) GetRoutes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	ip := parseIP(r.RemoteAddr)

	rs, err := s.router.Routes()
	if err != nil {
		log.Printf("GetRoutes/Error: %s", err)
		sendError(w, http.StatusInternalServerError, "500", "Unable to fetch routes")
		return
	}

	resps := make([]routesResp, 0, len(rs))
	for _, r := range rs {
		resps = append(resps, routeToRespRoute(r))
	}
	sort.Sort(ByHostname(resps))

	t := struct {
		IP   string       `json:"request-ip"`
		Data []routesResp `json:"data"`
	}{
		IP:   ip,
		Data: resps,
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(t)
}

func (s *Server) GetTables(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	resp := struct {
		Data []TableDef `json:"data"`
	}{
		Data: s.tables,
	}
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		sendError(w, http.StatusBadRequest, "400", "Unable to process request")
	}
}

type setReq struct {
	Data struct {
		IP    string
		Table string
	} `json:"data"`
}

func (s *Server) SetRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	ip := parseIP(r.RemoteAddr)
	dec := json.NewDecoder(r.Body)
	var req setReq
	err := dec.Decode(&req)
	if err != nil {
		sendError(w, http.StatusBadRequest, "400", "Unable to process request")
		return
	}
	defer r.Body.Close()

	changeReq := req.Data
	// Auth if ip does not match
	if changeReq.IP != ip {
		authd := s.auth.Auth(r)
		if !authd {
			sendError(w, http.StatusUnauthorized, "401", "Invalid authorization")
			return
		}
	}
	err = s.router.SetRoute(changeReq.IP, changeReq.Table)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "500", "Could not process request")
		return
	}
	rs, err := s.router.Routes()
	if err != nil {
		sendError(w, http.StatusInternalServerError, "500", "Could not get routes")
		return
	}
	route, found := routeByIP(rs, changeReq.IP)
	if !found {
		sendError(w, http.StatusNotFound, "404", "Route not found")
		return
	}

	sendRoute(w, route)
}

func routeByIP(rs []router.Route, ip string) (router.Route, bool) {
	for _, r := range rs {
		if r.IP == ip {
			return r, true
		}
	}
	return router.Route{}, false
}

func sendRoute(w http.ResponseWriter, route router.Route) {
	t := struct {
		Data routesResp `json:"data"`
	}{
		Data: routeToRespRoute(route),
	}
	err := json.NewEncoder(w).Encode(&t)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "500", "Could not get routes")
		return
	}
}
