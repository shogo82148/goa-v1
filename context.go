package goa

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Keys used to store data in context.
var (
	reqKey            = &contextKey{"request"}
	respKey           = &contextKey{"response"}
	ctrlKey           = &contextKey{"controller"}
	actionKey         = &contextKey{"action"}
	logKey            = &contextKey{"logger"}
	errKey            = &contextKey{"error"}
	securityScopesKey = &contextKey{"security-scope"}
)

type (
	// RequestData provides access to the underlying HTTP request.
	RequestData struct {
		*http.Request

		// Payload returns the decoded request body.
		Payload interface{}
		// Params contains the raw values for the parameters defined in the design including
		// path parameters, query string parameters and header parameters.
		Params url.Values
	}

	// ResponseData provides access to the underlying HTTP response.
	ResponseData struct {
		http.ResponseWriter

		// The service used to encode the response.
		Service *Service
		// ErrorCode is the code of the error returned by the action if any.
		ErrorCode string
		// Status is the response HTTP status code.
		Status int
		// Length is the response body length.
		Length int
	}
)

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation.
type contextKey struct {
	name string
}

func (k *contextKey) String() string { return "goa-v1 context value " + k.name }

// NewContext builds a new goa request context.
// If ctx is nil then req.Context() is used.
func NewContext(ctx context.Context, rw http.ResponseWriter, req *http.Request, params url.Values) context.Context {
	if ctx == nil {
		ctx = req.Context()
	} else {
		// The parent of req.Context() should be ctx,
		// but actually they are not because of compatibility.
		// So, we emulates the context whose parent is ctx.
		ctx = mergeContext(ctx, req.Context())
	}
	request := &RequestData{Request: req, Params: params}
	response := &ResponseData{ResponseWriter: rw}
	ctx = context.WithValue(ctx, respKey, response)
	ctx = context.WithValue(ctx, reqKey, request)

	return ctx
}

type mergedContext struct {
	parent, child context.Context
	cancel        context.CancelFunc
}

func mergeContext(parent, child context.Context) context.Context {
	ctx := &mergedContext{
		parent: parent,
		child:  child,
	}
	if parent.Done() != nil {
		// propagate cancellation from the parent to the child.
		ctx.child, ctx.cancel = context.WithCancel(child)
		go ctx.watchCancel()
	}
	return ctx
}

func (ctx *mergedContext) watchCancel() {
	select {
	case <-ctx.parent.Done():
		ctx.cancel()
	case <-ctx.child.Done():
	}
}

func (ctx *mergedContext) Deadline() (deadline time.Time, ok bool) {
	parent, ok := ctx.parent.Deadline()
	if !ok {
		return ctx.child.Deadline()
	}
	child, ok := ctx.child.Deadline()
	if !ok {
		return parent, true
	}

	if parent.After(child) {
		return child, true
	}
	return parent, true
}

func (ctx *mergedContext) Done() <-chan struct{} {
	return ctx.child.Done()
}

func (ctx *mergedContext) Err() error {
	return ctx.child.Err()
}

func (ctx *mergedContext) Value(key interface{}) interface{} {
	if v := ctx.child.Value(key); v != nil {
		return v
	}
	return ctx.parent.Value(key)
}

// WithAction creates a context with the given action name.
func WithAction(ctx context.Context, action string) context.Context {
	return context.WithValue(ctx, actionKey, action)
}

// WithLogger sets the request context logger and returns the resulting new context.
func WithLogger(ctx context.Context, logger LogAdapter) context.Context {
	return context.WithValue(ctx, logKey, logger)
}

// WithLogContext instantiates a new logger by appending the given key/value pairs to the context
// logger and setting the resulting logger in the context.
func WithLogContext(ctx context.Context, keyvals ...interface{}) context.Context {
	logger := ContextLogger(ctx)
	if logger == nil {
		return ctx
	}
	nl := logger.New(keyvals...)
	return WithLogger(ctx, nl)
}

// WithError creates a context with the given error.
func WithError(ctx context.Context, err error) context.Context {
	return context.WithValue(ctx, errKey, err)
}

// ContextController extracts the controller name from the given context.
func ContextController(ctx context.Context) string {
	if c := ctx.Value(ctrlKey); c != nil {
		return c.(string)
	}
	return "<unknown>"
}

// ContextAction extracts the action name from the given context.
func ContextAction(ctx context.Context) string {
	if a := ctx.Value(actionKey); a != nil {
		return a.(string)
	}
	return "<unknown>"
}

// ContextRequest extracts the request data from the given context.
func ContextRequest(ctx context.Context) *RequestData {
	if r := ctx.Value(reqKey); r != nil {
		return r.(*RequestData)
	}
	return nil
}

// ContextResponse extracts the response data from the given context.
func ContextResponse(ctx context.Context) *ResponseData {
	if r := ctx.Value(respKey); r != nil {
		return r.(*ResponseData)
	}
	return nil
}

// ContextLogger extracts the logger from the given context.
func ContextLogger(ctx context.Context) LogAdapter {
	if v := ctx.Value(logKey); v != nil {
		return v.(LogAdapter)
	}
	return nil
}

// ContextError extracts the error from the given context.
func ContextError(ctx context.Context) error {
	if err := ctx.Value(errKey); err != nil {
		return err.(error)
	}
	return nil
}

// SwitchWriter overrides the underlying response writer. It returns the response
// writer that was previously set.
func (r *ResponseData) SwitchWriter(rw http.ResponseWriter) http.ResponseWriter {
	rwo := r.ResponseWriter
	r.ResponseWriter = rw
	return rwo
}

// Written returns true if the response was written, false otherwise.
func (r *ResponseData) Written() bool {
	return r.Status != 0
}

// WriteHeader records the response status code and calls the underlying writer.
func (r *ResponseData) WriteHeader(status int) {
	go IncrCounter([]string{"goa", "response", strconv.Itoa(status)}, 1.0)
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

// Write records the amount of data written and calls the underlying writer.
func (r *ResponseData) Write(b []byte) (int, error) {
	if !r.Written() {
		r.WriteHeader(http.StatusOK)
	}
	r.Length += len(b)
	return r.ResponseWriter.Write(b)
}
