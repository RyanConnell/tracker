package main

import (
	"fmt"
	"tracker/trackable/show"
)

func main() {
	fmt.Printf("Starting...\n")
	err := show.ScrapeAll()
	if err != nil {
		fmt.Printf("Error encountered; %v", err)
	}
}
