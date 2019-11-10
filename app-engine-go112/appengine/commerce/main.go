// [START gae_go111_app]

// Sample helloworld is an App Engine app.
package main

// [START import]
import (
	"fmt"
	"log"
	"net/http"
	"os"

	"myproject/infra/mux"
)

// [END import]
// [START main_func]

var app = mux.CreateMux()

func init() {
	// curl -X POST -H 'Content-Type: application/json' -d '{"orderDetails": [{ "itemId": 901, "quantity": 2 }, { "itemId": 902, "quantity": 2 }, { "itemId": 903, "quantity": 2 }], "customer": { "name": "Alice", "phoneNumber": "XXX-XXXX-XXXX" }, "employee": { "id": 101 }}' localhost:8080
	// app.POST("/", controller.HandleMakeReservation)
}

func main() {
	http.Handle("/", app)

	// [START setting_port]
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
	// [END setting_port]
}

// [END main_func]

// [START indexHandler]

// indexHandler responds to requests with our greeting.

// [END indexHandler]
// [END gae_go111_app]
