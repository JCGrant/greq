greq
====
[![GoDoc](https://godoc.org/github.com/JCGrant/greq?status.svg)](https://godoc.org/github.com/JCGrant/greq)

Simple HTTP Requests for humans


Install
=======
``` sh
go get github.com/JCGrant/greq
```


Examples
========

Super simple API

``` go
res, err := greq.Get("my-favourite-books.com")
if err != nil {
  log.Fatalln(err)
}
var book Book
res.JSON(&book)
```

Can use functional options to configure the request. For example, adding a body to a POST request is as easy as wrapping up your Go object in `greq.Body()`.

``` go
res, err := greq.Post("my-favourite-books.com", greq.Body(NewBook{Title: "My great book"}))
if err != nil {
  log.Fatalln(err)
}
var book Book
res.JSON(&book)
```
