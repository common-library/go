package socket

import "github.com/common-library/go/socket/tcp"

// Deprecated: Use github.com/common-library/go/socket/tcp.Client instead.
// This type alias is provided for backward compatibility and will be removed in v2.0.0.
// To migrate, change:
//
//	import "github.com/common-library/go/socket"
//	client := &socket.Client{}
//
// To:
//
//	import "github.com/common-library/go/socket/tcp"
//	client := &tcp.Client{}
type Client = tcp.Client

// Deprecated: Use github.com/common-library/go/socket/tcp.Server instead.
// This type alias is provided for backward compatibility and will be removed in v2.0.0.
// To migrate, change:
//
//	import "github.com/common-library/go/socket"
//	server := &socket.Server{}
//
// To:
//
//	import "github.com/common-library/go/socket/tcp"
//	server := &tcp.Server{}
type Server = tcp.Server
