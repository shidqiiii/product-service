package errmsg

type CostumError struct {
	Code   int
	Errors map[string][]string
	Msg    string
}

func (e *CostumError) Error() string {
	return e.Msg
}

func NewCostumErrors(errCode int, opts ...Option) *CostumError {
	err := &CostumError{
		Code:   errCode,
		Errors: make(map[string][]string),
		Msg:    "Your request has been failed to process",
	}

	for _, opt := range opts {
		opt(err)
	}

	return err
}

type Option func(*CostumError)

func WithMessage(msg string) Option {
	return func(err *CostumError) {
		err.Msg = msg
	}
}

func WithErrors(field string, msg string) Option {
	return func(err *CostumError) {
		err.Errors[field] = append(err.Errors[field], msg)
	}
}

func errorCustomHandler(err *CostumError) (int, *CostumError) {
	return err.Code, err
}
