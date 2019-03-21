package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rafaeldias/currency-converter/controllers"
	"github.com/rafaeldias/currency-converter/router"
	"github.com/rafaeldias/currency-converter/services/currency"
)

func getEnv(name string, def string) string {
	var env = os.Getenv(name)
	if env == "" {
		env = def
	}
	return env
}

func currencyLayerMiddleware(host, accessKey string) router.Handler {
	return func(hc router.HTTPContexter) {
		hc.Set("currency", currency.New(host, accessKey))
	}
}

// CORSHeaders sets headers to avoid cross-domain errors
func CORSHeaders() router.Handler {
	return func(hc router.HTTPContexter) {
		hc.Header().Set("Access-Control-Allow-Origin", "*")
		hc.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		hc.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
}

func noopHandler(hc router.HTTPContexter) {}

func main() {
	var appPort = getEnv("APP_PORT", "9000")
	var currHost = getEnv("CURRENCY_HOST", "service_host")
	var currAccessKey = getEnv("CURRENCY_ACCESSKEY", "service_access_key")

	var r = router.New()

	r.Use(currencyLayerMiddleware(currHost, currAccessKey))
	r.Use(CORSHeaders())

	r.Get("/currencies", controllers.List)
	r.Get("/currencies/:from/conversions/:to", controllers.ValidateConversion, controllers.Conversion)

	// To conform with CORS
	r.Options("/currencies", noopHandler)
	r.Options("/currencies/:from/conversions/:to", noopHandler)

	log.Printf("Listening to port %s\n", appPort)

	err := http.ListenAndServe(fmt.Sprintf(":%s", appPort), r)
	if err != nil {
		log.Fatalf("Failed starting the server: %s\n", err.Error())
		return
	}

}
