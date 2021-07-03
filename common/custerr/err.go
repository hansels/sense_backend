package custerr

import (
	"fmt"
)

type ErrChain struct {
	Message string
	Cause   error
	Fields  map[string]string
	Type    error
}

func (err ErrChain) Error() string {
	bcoz := ""
	fields := ""
	if err.Cause != nil {
		bcoz = fmt.Sprint(" because {", err.Cause.Error(), "}")
		if len(err.Fields) > 0 {
			fields = fmt.Sprintf(" with Fields {%+v}", err.Fields)
		}
	}
	return fmt.Sprint(err.Message, bcoz, fields)
}

func Type(err error) error {
	switch err.(type) {
	case ErrChain:
		return err.(ErrChain).Type
	}
	return nil
}

func (err ErrChain) SetField(key string, value string) ErrChain {
	if err.Fields == nil {
		err.Fields = map[string]string{}
	}
	err.Fields[key] = value
	return err
}
