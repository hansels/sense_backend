package errors

// Formattable represents a formattable error.
type Formattable interface {
	error
	Message() string
	Cause() error
}

// Error formats an error on a single line.
func Error(err Formattable) string {
	msg := err.Message()
	cause := err.Cause().Error()
	if msg == "" {
		return cause
	}
	return msg + ": " + cause
}
