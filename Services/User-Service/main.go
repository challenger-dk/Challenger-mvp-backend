package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Server started at http://localhost:8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
