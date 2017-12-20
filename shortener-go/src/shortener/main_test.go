package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/guregu/kami"

	"github.com/alicebob/miniredis"
)

func TestRegisterAndRedirect(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	urlForRedirect := "https://google.com"

	redisURL = s.Addr()

	data := url.Values{}
	data.Set("url", urlForRedirect)
	req, err := http.NewRequest("POST", "/api/register/", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	kami.Reset()
	setupRoutes()

	ctx := req.Context()
	req = req.WithContext(ctx)

	kami.Handler().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	alias := rr.Body.String()

	redReq, err := http.NewRequest("GET", "/api/r/"+alias, nil)
	rr = httptest.NewRecorder()
	kami.Handler().ServeHTTP(rr, redReq)

	if status := rr.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusFound)
	}
	if location := rr.Header().Get("location"); location != urlForRedirect {
		t.Errorf("handler returned wrong redirect: got %v want %v",
			location, urlForRedirect)
	}
}
