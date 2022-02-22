package goatest

import (
	"io"
	"log"
	"testing"

	"github.com/shogo82148/goa-v1"
	"github.com/shogo82148/goa-v1/middleware"
)

// TInterface is an interface for Go's testing.T and testing.B.
//
// Deprecated: use testing.TB instead.
type TInterface = testing.TB

// ResponseSetterFunc func
type ResponseSetterFunc func(resp interface{})

// Encode implements a dummy encoder that returns the value being encoded
func (r ResponseSetterFunc) Encode(v interface{}) error {
	r(v)
	return nil
}

// Service provide a general goa.Service used for testing purposes
func Service(logBuf io.Writer, respSetter ResponseSetterFunc) *goa.Service {
	s := goa.New("test")
	logger := log.New(logBuf, "", log.Ltime)
	s.WithLogger(goa.NewLogger(logger))
	s.Use(middleware.LogRequest(true))
	s.Use(middleware.LogResponse())
	newEncoder := func(io.Writer) goa.Encoder {
		return respSetter
	}
	s.Decoder.Register(goa.NewJSONDecoder, "*/*")
	s.Encoder.Register(newEncoder, "*/*")
	return s
}
