package main

import (
	"fmt"
	"github.com/mewa/krawlr"
	"os"
)

func main() {
	k := krawlr.New(1000)

	res, _ := k.Crawl(os.Args[1])
	for page, links := range res {
		fmt.Println(page)

		for link := range *links {
			fmt.Println("   -", link)
		}
	}
}
