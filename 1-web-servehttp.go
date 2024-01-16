package main

import (
	"fmt"
	"net/http"
)

type rupees float32

func (r rupees) String() string {
	return fmt.Sprintf("£%.2f", r)
}

// database is a pseudo database that maps strings to price
type database map[string]rupees

func (db database) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func main() {
	db := database{
		"foo": 1,
		"bar": 2,
	}

	// This is a simple implementation that does not involve any route parsing.
	// ListenAndServe starts Go's default http server on the specified port.
	// The custom db struct implements the ServeHTTP method, which satisfies the http.Handler interface required by ListenAndServe.
	/*
		func http.ListenAndServe(addr string, handler http.Handler) error

		type Handler interface {
			ServeHTTP(ResponseWriter, *Request)
		}
	*/
	// When a url on the specified port is accessed, the ServerHTTP method is called by default.
	// The same data is written to the http response writer no matter what url is visited.
	// How can we display different content depending on what url/page was visited?
	http.ListenAndServe(":8081", db)
}
