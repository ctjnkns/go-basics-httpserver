# go-basics-httpserver

I find there's a lot of confusion about the differences between Handle, Handler, HandleFunc, HandlerFunc, etc. in Go's net/http package.

This walkthrough takes a ground up approach to understanding what each function does by going over the various approaches to starting an http server in go and handling basic http routing.

Simply clone the repo and run: go run `./#-program-name`

Run each program in order, and refer to the comments in the code for a description.

## 1-web-servehttp


```go
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

	http.ListenAndServe(":8081", db)
}
```

This is a simple implementation that does not involve any route parsing.
ListenAndServe starts Go's default http server on the specified port.
The custom db struct implements the ServeHTTP method, which satisfies the http.Handler interface required by ListenAndServe.

> func http.ListenAndServe(addr string, handler http.Handler) error
> type Handler interface {
> 	ServeHTTP(ResponseWriter, *Request)
> }

When a url on the specified port is accessed, the ServerHTTP method is called by default.
The same data is written to the http response writer no matter what url is visited.
How can we display different content depending on what url/page was visited?

### http://localhost:8081/

```
bar: £2.00
foo: £1.00
```

## A great summary I found descirbing the same things in a slightly different way:
https://www.integralist.co.uk/posts/understanding-golangs-func-type/

Summary/Breakdown
Here is a useful summary for you…

http.Handler = interface
you support http.Handler if you have a ServeHTTP(w http.ResponseWriter, r *http.Request) method available.

http.Handle("/", <give me something that supports the http.Handler interface>)
e.g. an object with a ServeHTTP method.

http.HandleFunc("/", <give me any function with the same signature as ServeHTTP >)
e.g. a function that accepts the arguments (w http.ResponseWriter, r *http.Request).

http.HandlerFunc = func type used internally by http.HandleFunc
e.g. it adapts the given function to the http.HandlerFunc type, which has an associated ServeHTTP method (that is able to call your original incompatible function).
