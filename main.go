package main

import (
	"fmt"
)

// Didnt bother to check if strings for non playstation platforms are ok todo: check

func main() {

	s := NewScrapper()
	s2 := NewScrapper()
	newChannel := make(chan DiscountRecord, 1)
	bestChannel := make(chan DiscountRecord, 1)
	go s.FetchNewDiscounts(PS4, newChannel)
	go s2.FetchBestDiscounts(PS4, bestChannel)
	for {
		select {
		case n := <-newChannel:
			if n.Title != "" {
				fmt.Println(n.Title, " on new channel")
			}

		case b := <-bestChannel:
			if b.Title != "" {
				fmt.Println(b.Title, " on best channel")
			}
		}
	}
	fmt.Println("FINISH")

}
