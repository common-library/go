package long_polling_test

import (
	"math/rand"
	"net/http"
	"strconv"
	"testing"
	"time"

	long_polling "github.com/heaven-chp/common-library-go/long-polling"
)

func TestStart(t *testing.T) {
	const category = "category-1"
	const data = "data-1"
	const count = 10

	server := long_polling.Server{}
	address := ":" + strconv.Itoa(10000+rand.Intn(1000))
	dir := t.TempDir()

	start := func() {
		serverInfo := long_polling.ServerInfo{
			Address:                        address,
			Timeout:                        3600,
			SubscriptionURI:                "/subscription",
			HandlerToRunBeforeSubscription: func(w http.ResponseWriter, r *http.Request) bool { return true },
			PublishURI:                     "/publish",
			HandlerToRunBeforePublish:      func(w http.ResponseWriter, r *http.Request) bool { return true }}

		filePersistorInfo := long_polling.FilePersistorInfo{Use: true, FileName: dir + "/file-persistor.txt", WriteBufferSize: 250, WriteFlushPeriodSeconds: 1}

		err := server.Start(serverInfo, filePersistorInfo, func(err error) { panic(err) })
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Duration(200) * time.Millisecond)
	}

	stop := func() {
		err := server.Stop(10)
		if err != nil {
			t.Fatal(err)
		}
	}

	func() {
		start()
		defer stop()

		for i := 0; i < count; i++ {
			request := long_polling.PublishRequest{Category: category, Data: data}
			_, err := long_polling.Publish("http://"+address+"/publish", 10, nil, request, "", "")
			if err != nil {
				t.Fatal(err)
			}
		}
	}()

	func() {
		start()
		defer stop()

		request := long_polling.SubscriptionRequest{Category: category, Timeout: 300, SinceTime: 1}
		response, err := long_polling.Subscription("http://"+address+"/subscription", nil, request, "", "")
		if err != nil {
			t.Fatal(err)
		}

		if len(response.Events) != count {
			t.Fatalf("invalid count - (%d)(%d)", len(response.Events), count)
		}
	}()
}

func TestStop(t *testing.T) {
	server := long_polling.Server{}

	err := server.Stop(10)
	if err != nil {
		t.Fatal(err)
	}
}
