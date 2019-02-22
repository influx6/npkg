package nhttp

import (
	"net/http"
	"strings"

	"github.com/dimfeld/httptreemux"

	"github.com/gorilla/mux"
)

// HandlerMW defines a function which wraps a provided http.handlerFunc
// which encapsulates the original for a underline operation.
type HandlerMW func(http.Handler, ...Middleware) http.Handler

// HandlerFuncMW defines a function which wraps a provided http.handlerFunc
// which encapsulates the original for a underline operation.
type HandlerFuncMW func(http.Handler, ...Middleware) http.HandlerFunc

// TreeMuxHandler defines a function type for the httptreemux.Handler type.
type TreeMuxHandler func(http.ResponseWriter, *http.Request, map[string]string)

// Middleware defines a function type which is used to create a chain
// of handlers for processing giving request.
type Middleware func(next http.Handler) http.Handler

// IdentityMW defines a http.Handler function that returns a the next http.Handler passed to it.
func IdentityMW(next http.Handler) http.Handler {
	return next
}

// MW combines multiple Middleware to return a single http.Handler.
func MW(mos ...Middleware) http.Handler {
	return CombineMoreMW(mos...)(IdentityHandler())
}

// CombineMW combines two middleware and returns a single http.Handler.
func CombineMW(mo, mi Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		handler := mo(mi(IdentityHandler()))

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler.ServeHTTP(w, r)

			if next != nil {
				next.ServeHTTP(w, r)
			}
		})
	}
}

// CombineMoreMW combines multiple Middleware to return a new Middleware.
func CombineMoreMW(mos ...Middleware) Middleware {
	var initial Middleware
	if len(mos) == 0 {
		initial = IdentityMW
	}

	if len(mos) == 1 {
		return mos[0]
	}

	for _, mw := range mos {
		if initial == nil {
			initial = mw
			continue
		}

		initial = CombineMW(initial, mw)
	}

	return initial
}

// ContextHandler defines a function type which accepts a function type.
type ContextHandler func(*NContext) error

// ErrorResponse defines a function which receives the possible error that
// occurs from a ContextHandler and applies necessary response as needed.
type ErrorResponse func(error, *NContext)

// ErrorHandler defines a function type which sets giving response to a Response object.
type ErrorHandler func(error, *NContext) error

// HandlerToContextHandler returns a new ContextHandler from a http.Handler.
func HandlerToContextHandler(handler http.Handler) ContextHandler {
	return func(context *NContext) error {
		handler.ServeHTTP(context.Response(), context.Request())
		return nil
	}
}

// Treemux returns httptreemux.Handler for use with a httptreemux router.
func Treemux(ops []Options, errHandler ErrorResponse, handler ContextHandler, before []Middleware, after []Middleware) httptreemux.HandlerFunc {
	beforeMW := MW(before...)
	afterMW := MW(after...)

	return func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		beforeMW.ServeHTTP(w, r)
		defer afterMW.ServeHTTP(w, r)

		ctx := NewContext(ops...)
		ctx.Reset(r, &Response{Writer: w})
		defer ctx.Reset(nil, nil)
		defer ctx.ClearFlashMessages()

		for key, val := range params {
			ctx.AddParam(key, val)
		}

		if err := handler(ctx); err != nil && errHandler != nil {
			errHandler(err, ctx)
			return
		}
	}
}

// HandlerWith defines a function which will return a http.Handler from a ErrorHandler,
// and a ContextHandler. If the middleware set is provided then it's executed
func HandlerWith(ops []Options, errHandler ErrorResponse, handle ContextHandler, before []Middleware, after []Middleware) http.Handler {
	beforeMW := MW(before...)
	afterMW := MW(after...)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		beforeMW.ServeHTTP(w, r)
		defer afterMW.ServeHTTP(w, r)

		var nctx = NewContext(ops...)
		if err := nctx.Reset(r, &Response{Writer: w}); err != nil {
			if errHandler != nil {
				errHandler(err, nctx)
				return
			}
		}

		if err := handle(nctx); err != nil {
			if errHandler != nil {
				errHandler(err, nctx)
			}
		}
	})
}

// IdentityHandler returns a non-op http.Handler
func IdentityHandler() http.Handler {
	return http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})
}

// NetworkAuthenticationNeeded implements a http.Handler which returns http.StatusNetworkAuthenticationRequired always.
func NetworkAuthenticationNeeded(ctx *NContext) error {
	ctx.Status(http.StatusNetworkAuthenticationRequired)
	return nil
}

// NoContentRequest implements a http.Handler which returns http.StatusNoContent always.
func NoContentRequest(ctx *NContext) error {
	ctx.Status(http.StatusNoContent)
	return nil
}

// OKRequest implements a http.Handler which returns http.StatusOK always.
func OKRequest(ctx *NContext) error {
	ctx.Status(http.StatusOK)
	return nil
}

// BadRequestWithError implements a http.Handler which returns http.StatusBagRequest always.
func BadRequestWithError(err error, ctx *NContext) error {
	if err != nil {
		if httperr, ok := err.(HTTPError); ok {
			http.Error(ctx.Response(), httperr.Error(), httperr.Code)
			return nil
		}
		http.Error(ctx.Response(), err.Error(), http.StatusBadRequest)
	}
	return nil
}

// BadRequest implements a http.Handler which returns http.StatusBagRequest always.
func BadRequest(ctx *NContext) error {
	ctx.Status(http.StatusBadRequest)
	return nil
}

// NotFound implements a http.Handler which returns http.StatusNotFound always.
func NotFound(ctx *NContext) error {
	ctx.Status(http.StatusNotFound)
	return nil
}

// StripPrefixMW returns a middleware which strips the URI of the request of
// the provided Prefix. All prefix must come in /prefix/ format.
func StripPrefixMW(prefix string) Middleware {
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqURL := r.URL.Path
			if !strings.HasPrefix(reqURL, "/") {
				reqURL = "/" + reqURL
			}

			r.URL.Path = strings.TrimPrefix(reqURL, prefix)
			if next != nil {
				next.ServeHTTP(w, r)
			}
		})
	}
}

// GorillaMuxVars retrieves the parameter lists from the underline
// variable map provided by the gorilla mux router and stores those
// into the context.
func GorillaMuxVars(ctx *NContext) error {
	for k, v := range mux.Vars(ctx.Request()) {
		ctx.AddParam(k, v)
	}
	return nil
}

// HTTPRedirect returns a http.Handler which always redirect to the given path.
func HTTPRedirect(to string, code int) ContextHandler {
	return func(ctx *NContext) error {
		return ctx.Redirect(code, to)
	}
}

// OnDone calls the next http.Handler after the condition handler returns without error.
func OnDone(condition ContextHandler, nexts ...ContextHandler) ContextHandler {
	if len(nexts) == 0 {
		return condition
	}

	return func(c *NContext) error {
		if err := condition(c); err != nil {
			return err
		}

		for _, next := range nexts {
			if err := next(c); err != nil {
				return err
			}
		}
		return nil
	}
}

// OnNoError calls the next http.Handler after the condition handler returns no error.
func OnNoError(condition ContextHandler, action ContextHandler) ContextHandler {
	return func(c *NContext) error {
		if err := condition(c); err != nil {
			return err
		}

		return action(c)
	}
}

// OnError calls the next ContextHandler after the condition handler returns an error.
func OnError(condition ContextHandler, errorAction ContextHandler) ContextHandler {
	return func(c *NContext) error {
		if err := condition(c); err != nil {
			return errorAction(c)
		}

		return nil
	}
}

// OnErrorAccess calls the next ErrorHandler after the condition handler returns an error.
func OnErrorAccess(condition ContextHandler, errorAction ErrorHandler) ContextHandler {
	return func(c *NContext) error {
		if err := condition(c); err != nil {
			return errorAction(err, c)
		}

		return nil
	}
}

// HTTPConditionFunc retusn a handler where a ContextHandler is used as a condition where if the handler
// returns an error then the errorAction is called else the noerrorAction gets called with
// context. This allows you create a binary switch where the final action is based on the
// success of the first. Generally if you wish to pass info around, use the context.Bag()
// to do so.
func HTTPConditionFunc(condition ContextHandler, noerrorAction, errorAction ContextHandler) ContextHandler {
	return func(ctx *NContext) error {
		if err := condition(ctx); err != nil {
			return errorAction(ctx)
		}
		return noerrorAction(ctx)
	}
}

// HTTPConditionErrorFunc returns a handler where a condition ContextHandler is called whoes result if with an error
// is passed to the errorAction for execution else using the noerrorAction. Differs from HTTPConditionFunc
// due to the assess to the error value.
func HTTPConditionErrorFunc(condition ContextHandler, noerrorAction ContextHandler, errorAction ErrorHandler) ContextHandler {
	return func(ctx *NContext) error {
		if err := condition(ctx); err != nil {
			return errorAction(err, ctx)
		}
		return noerrorAction(ctx)
	}
}

// ErrorsAsResponse returns a ContextHandler which will always write out any error that
// occurs as the response for a request if any occurs.
func ErrorsAsResponse(code int, next ContextHandler) ContextHandler {
	return func(ctx *NContext) error {
		if err := next(ctx); err != nil {
			if httperr, ok := err.(HTTPError); ok {
				http.Error(ctx.Response(), httperr.Error(), httperr.Code)
				return err
			}

			if code <= 0 {
				code = http.StatusBadRequest
			}

			http.Error(ctx.Response(), err.Error(), code)
			return err
		}
		return nil
	}
}

// HTTPConditionsFunc returns a ContextHandler where if an error occurs would match the returned
// error with a ContextHandler to be runned if the match is found.
func HTTPConditionsFunc(condition ContextHandler, noerrAction ContextHandler, errCons ...MatchableContextHandler) ContextHandler {
	return func(ctx *NContext) error {
		if err := condition(ctx); err != nil {
			for _, errcon := range errCons {
				if errcon.Match(err) {
					return errcon.Handle(ctx)
				}
			}
			return err
		}
		return noerrAction(ctx)
	}
}

// MatchableContextHandler defines a condition which matches expected error
// for performing giving action.
type MatchableContextHandler interface {
	Match(error) bool
	Handle(*NContext) error
}

// Matchable returns MatchableContextHandler using provided arguments.
func Matchable(err error, fn ContextHandler) MatchableContextHandler {
	return errorConditionImpl{
		Err: err,
		Fn:  fn,
	}
}

// errorConditionImpl defines a type which sets the error that occurs and the handler to be called
// for such an error.
type errorConditionImpl struct {
	Err error
	Fn  ContextHandler
}

// Handler calls the internal http.Handler with provided NContext returning error.
func (ec errorConditionImpl) Handle(ctx *NContext) error {
	return ec.Fn(ctx)
}

// Match validates the provided error matches expected error.
func (ec errorConditionImpl) Match(err error) bool {
	return ec.Err == err
}

// MatchableFunction returns MatchableContextHandler using provided arguments.
func MatchableFunction(err func(error) bool, fn ContextHandler) MatchableContextHandler {
	return fnErrorCondition{
		Err: err,
		Fn:  fn,
	}
}

// fnErrorCondition defines a type which sets the error that occurs and the handler to be called
// for such an error.
type fnErrorCondition struct {
	Fn  ContextHandler
	Err func(error) bool
}

// http.Handler calls the internal http.Handler with provided NContext returning error.
func (ec fnErrorCondition) Handle(ctx *NContext) error {
	return ec.Fn(ctx)
}

// Match validates the provided error matches expected error.
func (ec fnErrorCondition) Match(err error) bool {
	return ec.Err(err)
}
