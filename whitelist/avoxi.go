package whitelist

import (
	"log"
	"net"
	"net/http"

	"avoxi/persistence"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type v1ErrorResponse struct {
	Message string `json:"errorMessage"`
}

type v1Response struct {
	WhiteListed bool `json:"whitelisted"`
}

type v1Request struct {
	RequestingIP net.IP   `json:"requestingIp"`
	Whitelist    []string `json:"whitelist"`
}

type v1Record struct {
	Country struct {
		ISOCode string `maxminddb:"iso_code"`
	} `maxminddb:"country"`
}

// Used by main to generate V1 whitelist route
func RoutesV1() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/whitelist", IsWhitelistedV1)

	return router
}

// IsWhitelistedV1 is the main function that handles the POST request on our endpoint
func IsWhitelistedV1(w http.ResponseWriter, r *http.Request) {
	log.Println("Received Whitlist request")

	var jsonBody v1Request

	err := render.DecodeJSON(r.Body, &jsonBody)

	// Verify that we have valid JSON and that our data was deserialized properly
	//  otherwise return an error response
	if err != nil {
		log.Printf("Error parsing request body %s\n", err)
		response := v1ErrorResponse{
			Message: "Error parsing request body",
		}

		render.Status(r, 422)
		render.JSON(w, r, response)
		return
	}

	logRequestInfo(jsonBody.RequestingIP, jsonBody.Whitelist)

	var record v1Record

	// Lookup the given IP address and verify whether the database threw an error
	err = persistence.GetDB().Lookup(jsonBody.RequestingIP, &record)
	if err != nil {
		log.Printf("Database Error: %s", err)
		response := v1ErrorResponse{
			Message: "Internal Server Error",
		}
		render.Status(r, 500)
		render.JSON(w, r, response)
		return
	}

	// 	Check whether the lookup did not find a matching result
	// 		treating this case as not whitelisted
	if record == (v1Record{}) {
		log.Printf("IP Lookup Failed %s", err)
		response := v1Response{
			WhiteListed: false,
		}
		render.JSON(w, r, response)
		return
	}
	log.Printf("Country Code from IP %s", record.Country.ISOCode)

	// Now that we have successfully looked up the country code verify whether it's in the given whitelist
	whitelisted := checkCountry(record.Country.ISOCode, jsonBody.Whitelist)

	// Build and return the final response
	result := v1Response{
		WhiteListed: whitelisted,
	}

	render.JSON(w, r, result)
}

// Helper function that checks the country in whitelist
func checkCountry(iso string, whitelist []string) bool {
	whitelisted := false
	for i := 0; i < len(whitelist); i++ {
		if iso == whitelist[i] {
			whitelisted = true
			break
		}
	}

	return whitelisted
}

// Helper function that logs request info
func logRequestInfo(ip net.IP, whitelist []string) {
	log.Printf("Requesting IP: %s\n", ip)
	log.Printf("Whitelisted Countries: ")
	for i := 0; i < len(whitelist); i++ {
		log.Printf("%s ", whitelist[i])
	}
}
