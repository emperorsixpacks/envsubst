package envsubt

type typeError struct {
	s string
}

func (e typeError) Error() string {
	return e.s
}

type configError struct {
	s string
}

func (e configError) Error() string {
	return e.s
}

func NewTypeError(s string) typeError {
	return typeError{s}
}

func NewConfigError(s string) configError {
	return configError{s}
}
