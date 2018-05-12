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

// Options is a struct for specifying configuration of the middleware
type Options struct {
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
}

type PasetoMiddleware struct {
	Options Options
}

func OnError(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, err.Error(), http.StatusUnauthorized)
}

func New(options ...Options) (*PasetoMiddleware, error) {
	var opts Options
	if len(options) == 0 {
		opts = Options{}
	} else {
		opts = options[0]
	}

	if opts.TokenProperty == "" {
		opts.TokenProperty = "token"
	}

	if opts.FooterProperty == "" {
		opts.FooterProperty = "paseto_footer"
	}

	if opts.ErrorHandler == nil {
		opts.ErrorHandler = OnError
	}

	if opts.Extractor == nil {
		return nil, errors.New("extractor not defined")
	}

	if opts.Decryptor == nil {
		return nil, errors.New("decryptor not defined")
	}

	return &PasetoMiddleware{Options: opts}, nil
}

func (p *PasetoMiddleware) logf(format string, args ...interface{}) {
	if p.Options.Debug {
		log.Printf(format, args)
	}
}

// Special implementation for Negroni, but could be used elsewhere.
func (p *PasetoMiddleware) HandlerWithNext(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	err := p.handlePaseto(w, r)

	// If there was an error, do not call next.
	if err == nil && next != nil {
		next(w, r)
	}
}

func (p *PasetoMiddleware) handlePaseto(w http.ResponseWriter, r *http.Request) error {
	pas, err := p.Options.Extractor(r)

	if err != nil {
		p.logf("Error extracting Paseto: %v", err)
	} else {
		p.logf("Token extracted: %s", pas)
	}

	if err != nil {
		p.Options.ErrorHandler(w, r, err)
		return fmt.Errorf("error extracting pas: %v", err)
	}

	if pas == "" {
		if p.Options.CredentialsOptional {
			p.logf("\tNo credentials found (CredentialsOptional=true)")
			return nil
		}

		errorMsg := "required auth paseto not found"
		p.Options.ErrorHandler(w, r, errors.New(errorMsg))
		p.logf("\tError: No credentials found (CredentialsOptional=false)")
		return fmt.Errorf(errorMsg)
	}

	var (
		footer string
		token  paseto.JSONToken
	)

	err = p.Options.Decryptor(pas, &token, &footer)

	if err != nil {
		p.logf("Error decrypting pas: %v", err)
		p.Options.ErrorHandler(w, r, err)
		return fmt.Errorf("error decrypting pas")
	} else {
		p.logf("Paseto decrypted: %s - %s", token, footer)
	}

	c := context.WithValue(r.Context(), p.Options.TokenProperty, token)
	c = context.WithValue(c, p.Options.FooterProperty, footer)
	newRequest := r.WithContext(c)

	*r = *newRequest
	return nil
}
