package pasetomiddleware

import (
	"net/http"

	"github.com/o1egl/paseto"
)

// ErrorHandler receives an error and can use that to return an error to the HTTP Request
type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)

// TokenExtractor extracts the token from a request and returns that or an error
type TokenExtractor func(r *http.Request) (string, error)

// TokenDecryptor decrypts the encrypted paseto token and puts it in the token and footer pointers.
type TokenDecryptor func(pas string, token *paseto.JSONToken, footer *string) error

// Error handles any errors that occur in the middleware process
func Error(eh ErrorHandler) Option {
	return func(p *PasetoMiddleware) {
		p.ErrorHandler = eh
	}
}

// Extractor extracts the Paseto token from the request.
func Extractor(e TokenExtractor) Option {
	return func(p *PasetoMiddleware) {
		p.Extractor = e
	}
}

// Decryptor decrypts the Paseto token that was retrieved through Extractor
func Decryptor(d TokenDecryptor) Option {
	return func(p *PasetoMiddleware) {
		p.Decryptor = d
	}
}

// CredentialsOptional allows the user to not be authenticated on a path with this middleware.
func CredentialsOptional(o bool) Option {
	return func(p *PasetoMiddleware) {
		p.CredentialsOptional = o
	}
}

// TokenProperty defines where the unencrypted Paseto token should be stored in the context of the request.
func TokenProperty(tp string) Option {
	return func(p *PasetoMiddleware) {
		p.TokenProperty = tp
	}
}

// FooterProperty defines where the unencrypted Paseto footer should be stored in the context of the request.
func FooterProperty(fp string) Option {
	return func(p *PasetoMiddleware) {
		p.FooterProperty = fp
	}
}

// Debug defines if the decryption process should be printed to stdout
func Debug(d bool) Option {
	return func(p *PasetoMiddleware) {
		p.Debug = d
	}
}
