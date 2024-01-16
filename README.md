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
foo: £1.00
```

### http://localhost:8081/bar
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

## Recap 
Let's recap at this point:
http.Handler (with an r) is an interface (with an r) containing the ServerHTTP function

>type Handler interface {
>	ServeHTTP(ResponseWriter, *Request)
>}
>

Both http.ListenAndServe and http.Handle take http.Handler as an input parameter:

>func http.ListenAndServe(addr string, handler http.Handler) error

>func http.Handle(pattern string, handler http.Handler)

http.Handle let's us call different http.Handlers depending on the string that is matched.

Up until now, each type that we passed into http.ListenAndServer or http.Handle had to implement a function NAMED ServerHTTP to satisfy the Handler interface.

## 4-web-handle-handlerfunc

```go
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

	http.Handle("/foo", http.HandlerFunc(db.foo))
	http.Handle("/bar", http.HandlerFunc(db.bar))
	http.Handle("/", http.HandlerFunc(db.home))

	http.ListenAndServe(":8081", nil)

}
```

Now we'll see how we can name those functions anything we want
We're back to a single struct type (thankfully), and we now have three separate functions (home, foo, bar) that match the ServeHTTP signature.
We need to convince http.Handle to use our foo and bar functions instead of calling the default ServerHTTP function.
Enter http.HandlerFunc: an adapter that let's you register ordinary functions as an http.Handler that can be used by http.Handle

> type HandlerFunc func(ResponseWriter, *Request)

Now when "/foo" is matched, the db.foo function will be called INSTEAD of ServeHTTP
Finally, we can put as many methods in our struct as we like and call them depending on which url is accessed
This is a huge advantage over using switch statement or defining multiple structs as we tried before
Next, let's see how we can make this a little more concise

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

```go
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

	http.HandleFunc("/foo", db.foo)
	http.HandleFunc("/bar", db.bar)
	http.HandleFunc("/", db.home)

	http.ListenAndServe(":8081", nil)

}
```

The only change here is that we're using http.HandleFunc instead of http.Handle.
HandleFunc is just a convenience wrapper that that does the same thing we were doing before a little more concisely.
Instead of calling: http.Handle("/foo", http.HandlerFunc(db.foo)), http.HandleFunc expects us to be returning a custom function and does the http.HanderFunc for us.
This is a very common way of handling basic http routing and is probably what you will see in code examples.
But notice the ListenAndServe call: we're still passing in nil, which means the DefaultServeMux is used.
For better security, we want to create a local server mux and use that instead of the default one.

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
