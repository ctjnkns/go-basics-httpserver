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

	// let's recap at this point:
	// http.Handler (with an r) is an interface containing the ServerHTTP function
	/*
		type Handler interface {
			ServeHTTP(ResponseWriter, *Request)
		}
	*/

	// both http.ListenAndServe and http.Handle take http.Handler as an input parameter:
	/*
		func http.ListenAndServe(addr string, handler http.Handler) error
		func http.Handle(pattern string, handler http.Handler)
	*/
	// http.Handle let's us call different http.Handlers depending on the string that is matched.

	// up until now, each type that we passed into http.ListenAndServer or http.Handle had to implement a function NAMED ServerHTTP to satisfy the Handler interface.

	// now we'll see how we can name those functions anything we want
	// we're back to a single struct type (thankfully), and we now have three separate functions (home, foo, bar) that match the ServeHTTP signature.
	// we need to convince http.Handle to use our foo and bar functions instead of calling the default ServerHTTP function.
	// enter http.HandlerFunc: an adapter that let's you register ordinary functions as an http.Handler that can be used by http.Handle
	//
	/*
		type HandlerFunc func(ResponseWriter, *Request)
	*/
	// now when "/foo" is matched, the db.foo function will be called INSTEAD of ServeHTTP
	// finally, we can put as many methods in our struct as wel like and call them depending on which url is accessed
	// this is a huge advantage over using switch statement or defining multiple structs as we tried before

	http.Handle("/foo", http.HandlerFunc(db.foo))
	http.Handle("/bar", http.HandlerFunc(db.bar))
	http.Handle("/", http.HandlerFunc(db.home))

	http.ListenAndServe(":8081", nil)

}
