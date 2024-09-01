package main

import (
	"fmt"
	"mine-game/internal/database"
	"mine-game/internal/router"
	"net/http"
)

func main() {
	// Connect to MongoDB
	mongodb, err := database.Connect()
	if err != nil {
		panic(err)
	}

	defer mongodb.Disconnect()

	routes := router.NewRouter(mongodb)
	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: routes,
	}

	fmt.Println("Server is running on port 8080")

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
