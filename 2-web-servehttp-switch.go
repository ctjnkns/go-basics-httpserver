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
	switch r.URL.Path {
	case "/foo":
		fmt.Fprintf(w, "foo: %s\n", db["foo"])
	case "/bar":
		fmt.Fprintf(w, "bar: %s\n", db["bar"])

	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "No switch case found for: %s\n", r.URL)

	}
}

func main() {
	db := database{
		"foo": 1,
		"bar": 2,
	}

	// "db" still implements ServerHTTP, so the structure has not changed.
	// however, db's ServerHTTP function now uses switch cases to match the url path and write different data depending on which url was visited.
	// this will lead to an enormous ServerHTTP function eventually.
	// How can we break these cases out into separate functions?
	http.ListenAndServe(":8081", db)
}
