# GO Paseto Middleware
[![License](http://img.shields.io/:license-mit-blue.svg)](LICENSE)
[![Build Status](http://img.shields.io/travis/o1egl/paseto.svg?style=flat-square)](https://travis-ci.org/aaomidi/go-paseto-middleware)
[![Travis](https://travis-ci.com/aaomidi/go-paseto-middleware.svg?branch=master&style=flat-square)](https://travis-ci.com/aaomidi/go-paseto-middleware)
[![Go Report Card](https://goreportcard.com/badge/github.com/aaomidi/go-paseto-middleware)](https://goreportcard.com/report/github.com/aaomidi/go-paseto-middleware)
[![GoDoc](https://godoc.org/github.com/aaomidi/go-paseto-middleware?status.svg)](https://godoc.org/github.com/aaomidi/go-pasteo-middleware)


A middleware that will check that a Paseto token is sent in a request. It will then set the contents of the token into the context of the request.

This module allows you authenticate HTTP requests using Paseto.

## Installing

````bash
go get github.com/aaomidi/go-paseto-middleware
````

## Using it

This library is written for use with [o1egl's paseto library](https://github.com/o1egl/paseto).

You can use the `pasetomiddleware` with default `net/http` as follows:

````golang
package auth

type Backend struct {
	secretKey string
}

func Auth() *Backend {
	if authBackendInstance == nil {
		authBackendInstance = &Backend{
			secretKey: "YELLOW SUBMARINE, BLACK WIZARDRY", // Obviously don't use this exact string.
		}
	}

	return authBackendInstance
}

func (backend *Backend) Middleware(optional bool) *pasetomiddleware.PasetoMiddleware {
	middleware, _ := pasetomiddleware.New(
		pasetomiddleware.Extractor(func(r *http.Request) (string, error) {
			cookie, err := r.Cookie("paseto")
			if err != nil {
				return "", err
			}
			return cookie.Value, nil
		}),

		pasetomiddleware.Decryptor(func(pas string, token *paseto.JSONToken, footer *string) error {
			v2 := paseto.NewV2()
			err := v2.Decrypt(pas, []byte(backend.secretKey), token, footer)
			return err
		}),
		pasetomiddleware.CredentialsOptional(optional),

		pasetomiddleware.Debug(optional),
	)

	return middleware
}
````

````golang
    package api

    router := mux.newRouter()

    router.Handle("/profile", auth.Auth().Middleware(false).NextFunc(profile)).Methods("GET")
    // You can also use .Next(http.Handler) to add another middleware.
````

To get the token from the request, you will do the following:


````golang
// This is initiated before.
var authedMiddleware *pasetomiddleware.PasetoMiddleware

func getUUIDFromRequest(r *http.Request) (uuid.UUID, error) {
	t, ok := r.Context().Value(authedMiddleware.TokenProperty).(*paseto.JSONToken)
	if !ok {
		return uuid.New(), errors.New("token not valid")
	}

    // I put the UUID string with the user key into the token Map
	id := t.Get("user")
	if id == "" {
		return uuid.New(), errors.New("token not valid")
	}

    // Parses the UUID
	uid, err := uuid.Parse(fmt.Sprint(id))
	return uid, err
}
````

# Inspiration

This project was heavily inspired by [GO-JWT-Middleware](https://github.com/auth0/go-jwt-middleware).
