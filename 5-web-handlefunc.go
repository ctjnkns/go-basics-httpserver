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

	// the only change here is that we're using http.HandleFunc instead of http.Handle
	// HandleFunc is just a convenience wrapper that that does the same thing we were doing before a little more concisely
	// instead of calling: http.Handle("/foo", http.HandlerFunc(db.foo)), http.HandleFunc expects us to be returning a custom function and does the http.HanderFunc for us.
	// this is a very common way of handling basic http routing and is probably what you will see in code examples
	// but notice the ListenAndServe call: we're still passing in nil, which means the DefaultServeMux is used
	// for better security, we want to create a local server mux and use that instead of the default one

	http.HandleFunc("/foo", db.foo)
	http.HandleFunc("/bar", db.bar)
	http.HandleFunc("/", db.home)

	http.ListenAndServe(":8081", nil)

}
