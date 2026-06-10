package main

import (
	"fmt"
	"log"
)

func main() {
	// Starting point
	entries, err := ReadDirectory(".")

	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		fmt.Println(entry.Name)
	}
}
