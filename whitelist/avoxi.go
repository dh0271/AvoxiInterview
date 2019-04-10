package whitelist

import (
	"log"
	"net"
	"net/http"

	"avoxi/persistence"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type V1ErrorResponse struct {
	Message string `json:"errorMessage"`
}

type AvoxiV1Response struct {
	WhiteListed bool `json:"whitelisted"`
}

type Request struct {
	RequestingIP net.IP   `json:"requestingIp"`
	Whitelist    []string `json:"whitelist"`
}

type Record struct {
	Country struct {
		ISOCode string `maxminddb:"iso_code"`
	} `maxminddb:"country"`
}

func RoutesV1() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/whitelist", IsWhitelistedV1)

	return router
}

func IsWhitelistedV1(w http.ResponseWriter, r *http.Request) {
	log.Println("Received Whitlist request")

	var jsonBody Request

	err := render.DecodeJSON(r.Body, &jsonBody)

	// Verify that we have valid JSON and that our data was deserialized properly
	//  otherwise return an error response
	if err != nil {
		log.Printf("Error parsing request body %s\n", err)
		response := V1ErrorResponse{
			Message: "Error parsing request body",
		}
		render.JSON(w, r, response)
		return
	}

	log.Printf("Requesting IP: %s\n", jsonBody.RequestingIP)
	log.Printf("Whitelisted Countries: ")
	for i := 0; i < len(jsonBody.Whitelist); i++ {
		log.Printf("%s ", jsonBody.Whitelist[i])
	}

	var record Record

	// Lookup the given IP address and verify whether the database threw an error
	// 	or whether the lookup returned no results
	err = persistence.GetDB().Lookup(jsonBody.RequestingIP, &record)
	if err != nil || record == (Record{}) {
		log.Printf("IP Lookup Failed")
		response := V1ErrorResponse{
			Message: "Failed to lookup IP",
		}
		render.JSON(w, r, response)
		return
	}

	log.Printf("Country Code from IP %s", record.Country.ISOCode)

	// Now that we have successfully looked up the country code verify whether it's in the given whitelist
	whitelisted := false
	for i := 0; i < len(jsonBody.Whitelist); i++ {
		if record.Country.ISOCode == jsonBody.Whitelist[i] {
			whitelisted = true
		}
	}

	// Build and return the final response
	result := AvoxiV1Response{
		WhiteListed: whitelisted,
	}

	render.JSON(w, r, result)
}
