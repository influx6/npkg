// +build go1.1,go1.2,go1.3,go1.4,go1.5,go1.6,go1.7

package nhttp

import (
	"context"
)

// Close closes the underline server.
// It will forcefully close the server.
func (s serverItem) Close(ctx context.Context) error {
	return s.listener.Close()
}
