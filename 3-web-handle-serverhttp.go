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

type databaseFoo map[string]rupees

type databaseBar map[string]rupees

func (db database) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the home page")
}

func (db databaseFoo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (db databaseBar) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func main() {

	db := database{}

	dbFoo := databaseFoo{
		"foo": 1,
	}

	dbBar := databaseBar{
		"bar": 1,
	}

	// we now use http.Handle to specify different ServeHTTP methods depending on which url was accessed.
	// like http.ListenAndServe, http.Handle requires the http.Handler interface (something that implements ServeHTTP).
	// unlike http.ListenAndServe, http.Handle allows you to specify a different handler depending on which string pattern is matched.
	// we pass in a different struct depending on which url is accessed; each one has it's own ServeHTTP method that satisfies the http.Handler interface
	/*
		func http.Handle(pattern string, handler http.Handler)

		type Handler interface {
			ServeHTTP(ResponseWriter, *Request)
		}
	*/
	// the down side is we had to create multiple duplicate structs that all implement ServeHttp method; now our structs are a mess!
	// it would be better if we could call various custom methods within the same type struct instead of requiring the method to be named ServeHTTP...

	http.Handle("/", db)
	http.Handle("/foo", dbFoo)
	http.Handle("/bar", dbBar)

	http.ListenAndServe(":8081", nil)

}
