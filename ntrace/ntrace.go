package ntrace

import "context"

// Trace defines a function
func Trace(ctx context.Context, traceName string, do func(context.Context, opentracing.Span)) {

}
