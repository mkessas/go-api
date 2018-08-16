package main

import (
	"fmt"
)

const port = 8011

func main() {
	fmt.Printf("Starting server on port %d...\n", port)

	router()
}
