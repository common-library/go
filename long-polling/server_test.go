package long_polling_test

import (
	"math/rand/v2"
	"net/http"
	"strconv"
	"testing"
	"time"

	long_polling "github.com/common-library/go/long-polling"
)

func TestStart(t *testing.T) {
	t.Parallel()

	const category = "category-1"
	const data = "data-1"
	const count = 10

	server := long_polling.Server{}
	address := ":" + strconv.Itoa(10000+rand.IntN(1000))
	dir := t.TempDir()

	start := func() {
		serverInfo := long_polling.ServerInfo{
			Address:                        address,
			TimeoutSeconds:                 3600,
			SubscriptionURI:                "/subscription",
			HandlerToRunBeforeSubscription: func(w http.ResponseWriter, r *http.Request) bool { return true },
			PublishURI:                     "/publish",
			HandlerToRunBeforePublish:      func(w http.ResponseWriter, r *http.Request) bool { return true }}

		filePersistorInfo := long_polling.FilePersistorInfo{Use: true, FileName: dir + "/file-persistor.txt", WriteBufferSize: 250, WriteFlushPeriodSeconds: 1}

		if err := server.Start(serverInfo, filePersistorInfo, func(err error) { panic(err) }); err != nil {
			t.Fatal(err)
		}
		time.Sleep(200 * time.Millisecond)
	}

	stop := func() {
		if err := server.Stop(100 * time.Second); err != nil {
			t.Fatal(err)
		}
	}

	func() {
		start()
		defer stop()

		for i := 0; i < count; i++ {
			request := long_polling.PublishRequest{Category: category, Data: data}
			if _, err := long_polling.Publish("http://"+address+"/publish", 10*time.Second, nil, request, "", "", nil); err != nil {
				t.Fatal(err)
			}
		}
	}()

	func() {
		start()
		defer stop()

		request := long_polling.SubscriptionRequest{Category: category, TimeoutSeconds: 300, SinceTime: 1}
		if response, err := long_polling.Subscription("http://"+address+"/subscription", nil, request, "", "", nil); err != nil {
			t.Fatal(err)
		} else if len(response.Events) != count {
			t.Fatal(len(response.Events), count)
		}
	}()
}

func TestStop(t *testing.T) {
	t.Parallel()

	server := long_polling.Server{}

	if err := server.Stop(100 * time.Second); err != nil {
		t.Fatal(err)
	}
}
