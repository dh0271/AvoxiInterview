package main

import (
	"avoxi/persistence"
	"avoxi/whitelist"

	"log"
	"net/http"
	"os"

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
		r.Mount("/api", whitelist.RoutesV1())
	})

	return router
}

func main() {
	persistence.Init("GeoLite2-Country/GeoLite2-Country.mmdb")

	defer persistence.Close()

	router := GenerateRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// define our function to be used to list all routes and associated HTTP verbs
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(router, walkFunc); err != nil {
		log.Panicf("Error: %s\n", err.Error())
	}

	log.Printf("Starting webservice on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
