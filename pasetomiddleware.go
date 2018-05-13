package pasetomiddleware

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/o1egl/paseto"
)

// Option is a function for setting options within the PasetoMiddleware struct
type Option func(*PasetoMiddleware)

// TokenKey defines the key to access the decrypted paseto token
type TokenKey string

// FooterKey defines the key to access the decrypted paseto footer
type FooterKey string

// PasetoMiddleware struct for specifying all the configuration options for this middleware
type PasetoMiddleware struct {
	Extractor TokenExtractor
	Decryptor TokenDecryptor

	// The name of the property where the token will be stored
	// Default value: token
	TokenProperty TokenKey

	// The name of the property where the footer will be stored
	// Default value: paseto footer
	FooterProperty FooterKey

	ErrorHandler ErrorHandler

	CredentialsOptional bool

	Debug bool
}

// Next goes through the middleware and passes itself onto another http.Handler
func (p *PasetoMiddleware) Next(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := p.handlePaseto(w, r); err == nil && handler != nil {
			handler.ServeHTTP(w, r)
		}
	}
}

// NextFunc goes through the middleware and passes itself onto another http.HandlerFunc
func (p *PasetoMiddleware) NextFunc(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := p.handlePaseto(w, r); err == nil && handler != nil {
			handler.ServeHTTP(w, r)
		}
	}
}

// OnError is a default error handler
func OnError(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, err.Error(), http.StatusUnauthorized)
}

// New constructs a PasetoMiddleware structure with the supplied options
func New(options ...Option) (*PasetoMiddleware, error) {
	def := &PasetoMiddleware{
		TokenProperty:  "token",
		FooterProperty: "paseto_footer",
		ErrorHandler:   OnError,
	}

	for _, o := range options {
		o(def)
	}

	if def.Extractor == nil {
		return nil, errors.New("extractor not defined")
	}

	if def.Decryptor == nil {
		return nil, errors.New("decryptor not defined")
	}

	return def, nil
}

func (p *PasetoMiddleware) logf(format string, args ...interface{}) {
	if p.Debug {
		log.Printf(format, args)
	}
}

func (p *PasetoMiddleware) handlePaseto(w http.ResponseWriter, r *http.Request) error {
	pas, err := p.Extractor(r)

	if err != nil {
		p.logf("Error extracting Paseto: %v\n", err)
	} else {
		p.logf("Token extracted: %s\n", pas)
	}

	if err != nil {
		p.ErrorHandler(w, r, err)
		return fmt.Errorf("error extracting pas: %v", err)
	}

	if pas == "" {
		if p.CredentialsOptional {
			p.logf("\tNo credentials found (CredentialsOptional=true)\n")
			return nil
		}

		errorMsg := "required auth paseto not found"
		p.ErrorHandler(w, r, errors.New(errorMsg))
		p.logf("\tError: No credentials found (CredentialsOptional=false)\n")
		return fmt.Errorf(errorMsg)
	}

	var (
		footer string
		token  paseto.JSONToken
	)

	err = p.Decryptor(pas, &token, &footer)

	if err != nil {
		p.logf("Error decrypting pas: %v\n", err)
		p.ErrorHandler(w, r, err)
		return fmt.Errorf("error decrypting pas")
	}

	p.logf("Paseto decrypted: %s - %s\n", token, footer)

	c := context.WithValue(r.Context(), p.TokenProperty, &token)
	c = context.WithValue(c, p.FooterProperty, &footer)
	newRequest := r.WithContext(c)

	*r = *newRequest
	return nil
}
