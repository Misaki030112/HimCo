package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// CrawlAlbum  deal with the request to CrawlAlbum
// Example  /album?id=203132&ConvertCount=......
func CrawlAlbum(w http.ResponseWriter, r *http.Request) {
	fmt.Println(1)
	log.Printf("request param %v \n", r.URL.Query())

	io.WriteString(w, "Hello, world!\n")
}
