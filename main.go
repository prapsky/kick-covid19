package main

import (
	"log"
	"net/http"

	"kick-covid19/controllers"
)

func main() {
	http.HandleFunc("/api/v1/covid19/", controllers.GetRealtimeData)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
