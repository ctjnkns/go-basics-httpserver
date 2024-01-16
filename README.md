# go-basics-httpserver

I find there's a lot of confusion about the differences between Handle, Handler, HandleFunc, HandlerFunc, etc.

This walkthrough takes a ground up approach to understanding what each function does by going over the various approaches to starting an http server in go and handling basic http routing.

Simply clone the repo and run: go run `./*\#-program-name*`

Run each program in order, and refer to the comments in the code for a description.



##A great summary I found descirbing the same things in a slightly different way:
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
