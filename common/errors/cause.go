package errors

// Cause returns the direct cause of an error.
//
// If the error doesn't implement "Cause()", nil is returned.
func Cause(err error) error {
	type causer interface {
		Cause() error
	}
	if err, ok := err.(causer); ok {
		return err.Cause()
	}
	return nil
}
