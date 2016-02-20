package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/blang/vpnrouter/router"
	"github.com/stretchr/testify/assert"
)

type mockRouter struct {
	routesFn   func() ([]router.Route, error)
	setRouteFn func(ip, table string) error
}

func (r mockRouter) Routes() ([]router.Route, error) {
	return r.routesFn()
}

func (r mockRouter) SetRoute(ip, table string) error {
	return r.setRouteFn(ip, table)
}

func TestRoutes(t *testing.T) {
	assert := assert.New(t)
	mock := mockRouter{
		routesFn: func() ([]router.Route, error) {
			return []router.Route{
				{IP: "127.0.0.1", Table: "table1", Lease: router.Lease{MAC: "abc", IP: "127.0.0.1", Name: "name"}},
			}, nil
		},
	}

	server := Server{
		router: mock,
	}
	req, err := http.NewRequest("GET", "http://127.0.0.1", nil)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	req.RemoteAddr = "127.0.0.2:6000"
	w := httptest.NewRecorder()
	server.GetRoutes(w, req)

	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("application/vnd.api+json", w.HeaderMap.Get("Content-Type"))
	const expected = `{"request-ip":"127.0.0.2","data":[{"ip":"127.0.0.1","table":"table1","hostname":"name","mac":"abc"}]}`
	assert.Equal(expected, strings.TrimSpace(w.Body.String()), "Invalid response")
	// t.Logf("%s", strings.TrimSpace(w.Body.String()))
}

func TestSetRoute(t *testing.T) {
	assert := assert.New(t)

	mock_routes := []router.Route{
		{IP: "127.0.0.1", Table: "table1", Lease: router.Lease{MAC: "abc", IP: "127.0.0.1", Name: "name"}},
	}
	mock := mockRouter{
		routesFn: func() ([]router.Route, error) {
			return mock_routes, nil
		},
		setRouteFn: func(ip, table string) error {
			if ip == mock_routes[0].IP {
				mock_routes[0].Table = table
				return nil
			}
			return errors.New("No Route")
		},
	}

	server := Server{
		router: mock,
		auth:   NewTokenAuth(),
	}
	const reqStr = `{"data":{"ip":"127.0.0.1","table":"table2"}}`
	req, err := http.NewRequest("POST", "http://127.0.0.1", strings.NewReader(reqStr))
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	req.RemoteAddr = "127.0.0.2:6000"
	w := httptest.NewRecorder()
	server.SetRoute(w, req)

	assert.Equal(http.StatusUnauthorized, w.Code, "Invalid status code")

}

func TestSetRouteAuthorizedSameIP(t *testing.T) {
	assert := assert.New(t)

	mock_routes := []router.Route{
		{IP: "127.0.0.1", Table: "table1", Lease: router.Lease{MAC: "abc", IP: "127.0.0.1", Name: "name"}},
	}
	mock := mockRouter{
		routesFn: func() ([]router.Route, error) {
			return mock_routes, nil
		},
		setRouteFn: func(ip, table string) error {
			if ip == mock_routes[0].IP {
				mock_routes[0].Table = table
				return nil
			}
			return errors.New("No Route")
		},
	}

	server := Server{
		router: mock,
		auth:   NewTokenAuth(""),
	}
	const reqStr = `{"data":{"ip":"127.0.0.1","table":"table2"}}`
	req, err := http.NewRequest("POST", "http://127.0.0.1", strings.NewReader(reqStr))
	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	req.RemoteAddr = "127.0.0.1:6000"
	w := httptest.NewRecorder()
	server.SetRoute(w, req)

	assert.Equal(http.StatusOK, w.Code, "Invalid status code")
	assert.Equal("table2", mock_routes[0].Table)
	assert.Equal(`{"data":{"ip":"127.0.0.1","table":"table2","hostname":"name","mac":"abc"}}`, strings.TrimSpace(w.Body.String()))
}

func TestSetRouteAuthorized(t *testing.T) {
	assert := assert.New(t)

	mock_routes := []router.Route{
		{IP: "127.0.0.1", Table: "table1", Lease: router.Lease{MAC: "abc", IP: "127.0.0.1", Name: "name"}},
	}
	mock := mockRouter{
		routesFn: func() ([]router.Route, error) {
			return mock_routes, nil
		},
		setRouteFn: func(ip, table string) error {
			if ip == mock_routes[0].IP {
				mock_routes[0].Table = table
				return nil
			}
			return errors.New("No Route")
		},
	}

	server := Server{
		router: mock,
		auth:   NewTokenAuth("token"),
	}
	const reqStr = `{"data":{"ip":"127.0.0.1","table":"table2"}}`
	req, err := http.NewRequest("POST", "http://127.0.0.1", strings.NewReader(reqStr))
	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	//Different IP
	req.RemoteAddr = "127.0.0.5:6000"
	req.Header.Set("Authorization", authHelper("token"))
	w := httptest.NewRecorder()
	server.SetRoute(w, req)

	assert.Equal(http.StatusOK, w.Code, "Invalid status code")
	assert.Equal("table2", mock_routes[0].Table)
	assert.Equal(`{"data":{"ip":"127.0.0.1","table":"table2","hostname":"name","mac":"abc"}}`, strings.TrimSpace(w.Body.String()))
}

func TestParseIP(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("127.0.0.1", parseIP("127.0.0.1:1000"))
	assert.Equal("192.168.0.1", parseIP("192.168.0.1:4567"))
}

func TestRoutesError(t *testing.T) {
	assert := assert.New(t)
	mock := mockRouter{
		routesFn: func() ([]router.Route, error) {
			return nil, errors.New("Error Message")
		},
	}

	server := Server{
		router: mock,
	}
	req, err := http.NewRequest("GET", "http://127.0.0.1", nil)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	req.RemoteAddr = "127.0.0.2:6000"
	w := httptest.NewRecorder()
	server.GetRoutes(w, req)

	assert.Equal(http.StatusInternalServerError, w.Code)
	assert.Equal("application/vnd.api+json", w.HeaderMap.Get("Content-Type"))
	const expected = `{"errors":[{"code":"500","title":"Unable to fetch routes"}]}`
	assert.Equal(expected, strings.TrimSpace(w.Body.String()), "Invalid response")
}
