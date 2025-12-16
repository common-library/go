// Package eslock provides a global lock for elasticsearch client initialization
// to prevent data races in the underlying elastictransport library.
package eslock

import "sync"

var (
	// InitMu protects concurrent calls to elasticsearch.NewClient across all versions
	// to avoid data races in github.com/elastic/elastic-transport-go/v8/elastictransport
	/*
		==================
		WARNING: DATA RACE
		Write at 0x000002326490 by goroutine 123:
		  github.com/elastic/elastic-transport-go/v8/elastictransport.New()
		      /home/chp/go/pkg/mod/github.com/elastic/elastic-transport-go/v8@v8.7.0/elastictransport/elastictransport.go:257 +0x1424
		  github.com/elastic/go-elasticsearch/v9.newTransport()
		      /home/chp/go/pkg/mod/github.com/elastic/go-elasticsearch/v9@v9.1.0/elasticsearch.go:337 +0x897
		  github.com/elastic/go-elasticsearch/v9.NewClient()
		      /home/chp/go/pkg/mod/github.com/elastic/go-elasticsearch/v9@v9.1.0/elasticsearch.go:203 +0x57
	*/
	InitMu sync.Mutex
)
