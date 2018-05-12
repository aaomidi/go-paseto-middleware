# GO Paseto Middleware

A middleware that will check that a Paseto token is sent in a request. It will then set the contents of the token into the context of the request.

This module allows you authenticate HTTP requests using Paseto.

## Installing

````bash
go get github.com/aaomidi/go-paseto-middleware
````

## Using it

This library is written for use with [o1egl's paseto library](https://github.com/o1egl/paseto).

You can use the `pasetomiddleware` with default `net/http` as follows:



# Inspiration

This project was heavily inspired by [GO-JWT-Middleware](https://github.com/auth0/go-jwt-middleware).
