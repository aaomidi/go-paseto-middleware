package pasetomiddleware

func Error(eh ErrorHandler) Option {
	return func(p *PasetoMiddleware) {
		p.ErrorHandler = eh
	}
}

func Extractor(e TokenExtractor) Option {
	return func(p *PasetoMiddleware) {
		p.Extractor = e
	}
}

func Decryptor(d TokenDecryptor) Option {
	return func(p *PasetoMiddleware) {
		p.Decryptor = d
	}
}

func CredentialsOptional(o bool) Option {
	return func(p *PasetoMiddleware) {
		p.CredentialsOptional = o
	}
}
func TokenProperty(tp string) Option {
	return func(p *PasetoMiddleware) {
		p.TokenProperty = tp
	}
}
func FooterProperty(fp string) Option {
	return func(p *PasetoMiddleware) {
		p.FooterProperty = fp
	}
}

func Debug(d bool) Option {
	return func(p *PasetoMiddleware) {
		p.Debug = d
	}
}
