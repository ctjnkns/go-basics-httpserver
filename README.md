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
> 
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

## 2-web-servehttp-switch 

```go
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
	http.ListenAndServe(":8081", db)
}
```

The struct "db" still implements ServerHTTP, so the structure has not changed.
However, db's ServerHTTP function now uses switch cases to match the url path and write different data depending on which url was visited.
This will lead to an enormous ServerHTTP function eventually.
How can we break these cases out into separate functions?

### http://localhost:8081/
```
No switch case found for: /
```

### http://localhost:8081/foo
```
foo: £1.00```

### [http://localhost:8081/](http://localhost:8081/bar)
```
bar: £2.00
```

## 3-web-handle-serverhttp

```go
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

	http.Handle("/", db)
	http.Handle("/foo", dbFoo)
	http.Handle("/bar", dbBar)

	http.ListenAndServe(":8081", nil)

}

```

We now use http.Handle to specify different ServeHTTP methods depending on which url was accessed.
Like http.ListenAndServe, http.Handle requires the http.Handler interface (something that implements ServeHTTP).
Unlike http.ListenAndServe, http.Handle allows you to specify a different handler depending on which string pattern is matched.
We pass in a different struct depending on which url is accessed; each one has it's own ServeHTTP method that satisfies the http.Handler interface.

> func http.Handle(pattern string, handler http.Handler)
> 
> type Handler interface {
>	ServeHTTP(ResponseWriter, *Request)
> }
> 
The down side is we had to create multiple duplicate structs that all implement ServeHttp method; now our structs are a mess!
It would be better if we could call various custom functions within the same struct instead of requiring the method to be named ServeHTTP...

### http://localhost:8081/
```
Welcome to the home page
```

### http://localhost:8081/foo
```
foo: £1.00
```

### http://localhost:8081/bar
```
bar: £2.00
```

##

##








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
