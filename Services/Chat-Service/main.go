package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Server started at http://localhost:8082")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
