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

func (db database) home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the home page")
}

func (db database) foo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "foo: %s\n", db["foo"])
}

func (db database) bar(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "bar: %s\n", db["bar"])
}

func main() {
	db := database{
		"foo": 1,
		"bar": 2,
	}

	// now we create a mux using the NewServerMux function
	// instead of calling http.HandleFunc, we now set our handlers in our new mux using mux.HandleFunc (or whatever you may decide to call it)
	mux := http.NewServeMux()
	mux.HandleFunc("/foo", db.foo)
	mux.HandleFunc("/bar", db.bar)
	mux.HandleFunc("/", db.home)

	// we still call http.ListenAndServe, but instead of passing in nil, we pass it our own server mux so that it uses that instead of the default global one
	http.ListenAndServe(":8081", mux)

}
