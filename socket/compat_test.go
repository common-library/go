package socket_test

import (
	"testing"

	"github.com/common-library/go/socket"
	"github.com/common-library/go/socket/tcp"
)

// Test backward compatibility - old code should still work
func TestBackwardCompatibility(t *testing.T) {
	t.Parallel()

	// Old way (using deprecated types)
	var oldServer *socket.Server
	var oldClient *socket.Client

	// Should be compatible with new types
	var newServer *tcp.Server
	var newClient *tcp.Client

	// Type assertion should work (using pointers to avoid copying sync primitives)
	oldServer = (*socket.Server)(newServer)
	oldClient = (*socket.Client)(newClient)

	newServer = (*tcp.Server)(oldServer)
	newClient = (*tcp.Client)(oldClient)

	// Verify assignments worked
	_ = oldServer
	_ = oldClient
	_ = newServer
	_ = newClient

	t.Log("Backward compatibility verified")
}
