package kleos

import (
	"context"
	"sync"
)

// Stores context variable registrations.  See Register for more info.
var contexts contextFuncs

// ContextFunc can be registered to pull values from the context and add them to the supplied log
// message fields.  Won't be called if ctx is nil, and fields will not be nil either.
type ContextFunc func(ctx context.Context, fields Fields)

// Register a context function to pull variables from the context during logging (via WithContext).
// The context function should pull the desired variable out of a given context and add it to the
// Fields map.  For example:
//
// 	kleos.Register(func(ctx context.Context, fields kleos.Fields) {
//		requestID, ok := ctx.Value(CtxRequestID).(uint64)
//		if !ok {
//			return
//		}
//
//		fields["request"] = requestID
//	})
//
// The field will then be output with the rest of the fields.
func Register(fn ContextFunc) {
	contexts.Add(fn)
}

// Provides some synchronous update protections around registering and using the context functions.
type contextFuncs struct {
	fns []ContextFunc
	mutex sync.RWMutex
}

// Add a context function to the cache.
func (cf *contextFuncs) Add(fn ContextFunc) {
	cf.mutex.Lock()
	defer cf.mutex.Unlock()

	cf.fns = append(cf.fns, fn)
}

// Run the registered functions against the context and fields.  Updates the fields with context
// values when present.
func (cf *contextFuncs) Run(ctx context.Context, fields Fields) {
	if ctx == nil {
		return
	}

	cf.mutex.RLock()
	defer cf.mutex.RUnlock()

	for _, fn := range cf.fns {
		fn(ctx, fields)
	}
}
