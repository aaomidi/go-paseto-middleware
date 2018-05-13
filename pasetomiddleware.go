package pasetomiddleware

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/o1egl/paseto"
)

type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)

type TokenExtractor func(r *http.Request) (string, error)
type TokenDecryptor func(pas string, token *paseto.JSONToken, footer *string) error

// Option is a function for setting options within the PasetoMiddleware struct
type Option func(*PasetoMiddleware)

type PasetoMiddleware struct {
	Extractor TokenExtractor
	Decryptor TokenDecryptor

	// The name of the property where the token will be stored
	// Default value: token
	TokenProperty string

	// The name of the property where the footer will be stored
	// Default value: paseto footer
	FooterProperty string

	ErrorHandler ErrorHandler

	CredentialsOptional bool

	Debug bool

	log.Logger
}

func (p *PasetoMiddleware) Next(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := p.handlePaseto(w, r); err == nil && handler != nil {
			handler.ServeHTTP(w, r)
		}
	}
}

func (p *PasetoMiddleware) NextFunc(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := p.handlePaseto(w, r); err == nil && handler != nil {
			handler.ServeHTTP(w, r)
		}
	}
}

func OnError(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, err.Error(), http.StatusUnauthorized)
}

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
	} else {
		p.logf("Paseto decrypted: %s - %s\n", token, footer)
	}

	c := context.WithValue(r.Context(), p.TokenProperty, &token)
	c = context.WithValue(c, p.FooterProperty, &footer)
	newRequest := r.WithContext(c)

	*r = *newRequest
	return nil
}
