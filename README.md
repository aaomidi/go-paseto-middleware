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

# Inspiration

This project was heavily inspired by [GO-JWT-Middleware](https://github.com/auth0/go-jwt-middleware).
