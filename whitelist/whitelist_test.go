package whitelist

import (
	"avoxi/persistence"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func GenerateRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.RedirectSlashes,
		middleware.Recoverer,
	)

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/api", RoutesV1())
	})

	return router
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	GenerateRoutes().ServeHTTP(rr, req)

	return rr
}

func TestValidRequestWhitelistTrue(t *testing.T) {
	persistence.Init("../GeoLite2-Country/GeoLite2-Country.mmdb")

	defer persistence.Close()
	router := GenerateRoutes()

	httptest.NewServer(router)

	payload := []byte(`{"requestingIp": "8.8.8.8","whitelist": ["US", "FR"]}`)
	req, _ := http.NewRequest("POST", "/v1/api/whitelist", bytes.NewBuffer(payload))

	r := executeRequest(req)

	if r.Code != 200 {
		t.Fatal("Didn't receive 200 status, received" + string(r.Code))
	}

	var jsonBody v1Response

	err := render.DecodeJSON(r.Body, &jsonBody)

	if err != nil {
		t.Fatal("Received invalid JSON response")
	}

	if !jsonBody.WhiteListed {
		t.Error("Incorrect result recieved, expected true received false")
	}
}

func TestValidRequestWhitelistFalse(t *testing.T) {
	persistence.Init("../GeoLite2-Country/GeoLite2-Country.mmdb")

	defer persistence.Close()
	router := GenerateRoutes()

	httptest.NewServer(router)

	payload := []byte(`{"requestingIp": "8.8.8.8","whitelist": ["FR"]}`)
	req, _ := http.NewRequest("POST", "/v1/api/whitelist", bytes.NewBuffer(payload))

	r := executeRequest(req)

	if r.Code != 200 {
		t.Fatal("Didn't receive 200 status, received" + string(r.Code))
	}

	var jsonBody v1Response

	err := render.DecodeJSON(r.Body, &jsonBody)

	if err != nil {
		t.Fatal("Received invalid JSON response")
	}

	if jsonBody.WhiteListed {
		t.Error("Incorrect result recieved, expected false received true")
	}
}

func TestValidRequestNonexistentIP(t *testing.T) {
	persistence.Init("../GeoLite2-Country/GeoLite2-Country.mmdb")

	defer persistence.Close()
	router := GenerateRoutes()

	httptest.NewServer(router)

	payload := []byte(`{"requestingIp": "0.240.141.100","whitelist": ["FR"]}`)
	req, _ := http.NewRequest("POST", "/v1/api/whitelist", bytes.NewBuffer(payload))

	r := executeRequest(req)

	if r.Code != 200 {
		t.Fatal("Didn't receive 200 status, received" + string(r.Code))
	}

	var jsonBody v1Response

	err := render.DecodeJSON(r.Body, &jsonBody)

	if err != nil {
		t.Fatal("Received invalid JSON response")
	}

	if jsonBody.WhiteListed {
		t.Error("Incorrect result recieved, expected false received received true")
	}
}

func TestInvalidJsonWhitelist(t *testing.T) {
	persistence.Init("../GeoLite2-Country/GeoLite2-Country.mmdb")

	defer persistence.Close()
	router := GenerateRoutes()

	httptest.NewServer(router)

	payload := []byte(`{"requestingIp":::: "8.8.8.8","whitelist": ["FR"]}`)
	req, _ := http.NewRequest("POST", "/v1/api/whitelist", bytes.NewBuffer(payload))

	r := executeRequest(req)

	if r.Code != 422 {
		t.Fatal("Didn't receive 422 status, received" + string(r.Code))
	}

	var jsonBody v1ErrorResponse

	err := render.DecodeJSON(r.Body, &jsonBody)

	if err != nil {
		t.Fatal("Received invalid JSON response")
	}

	if jsonBody.Message != "Error parsing request body" {
		t.Error("Incorrect error message recieved")
	}
}
