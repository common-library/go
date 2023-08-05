package long_polling_test

import (
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	long_polling "github.com/heaven-chp/common-library-go/long-polling"
)

var address string

func subscription(t *testing.T, request long_polling.SubscriptionRequest, count int, data string) (int64, string) {
	response, err := long_polling.Subscription("http://"+address+"/subscription", nil, request, "", "")
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatalf("invalid status code - (%d)(%s)", response.StatusCode, http.StatusText(response.StatusCode))
	}

	if len(response.Events) != count {
		t.Fatalf("invalid count - (%d)(%d)", len(response.Events), count)
	}

	for _, event := range response.Events {
		if event.Category != request.Category {
			t.Fatalf("invalid category - (%s)(%s)", event.Category, request.Category)
		}

		if strings.HasPrefix(event.Data, data) == false {
			t.Fatalf("invalid data - (%s)(%s)", event.Data, data)
		}
	}

	return response.Events[len(response.Events)-1].Timestamp, response.Events[len(response.Events)-1].ID
}

func publish(t *testing.T, category, data string) {
	request := long_polling.PublishRequest{Category: category, Data: data}
	response, err := long_polling.Publish("http://"+address+"/publish", 10, nil, request, "", "")
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatalf("invalid status code - (%d)(%s)", response.StatusCode, http.StatusText(response.StatusCode))
	}

	if response.Body != `{"success": true}` {
		t.Fatalf("invalid body- (%s)", response.Body)
	}
}

func setUp(server *long_polling.Server) {
	address = ":" + strconv.Itoa(10000+rand.Intn(1000))

	serverInfo := long_polling.ServerInfo{
		Address:                        address,
		Timeout:                        3600,
		SubscriptionURI:                "/subscription",
		HandlerToRunBeforeSubscription: func(w http.ResponseWriter, r *http.Request) bool { return true },
		PublishURI:                     "/publish",
		HandlerToRunBeforePublish:      func(w http.ResponseWriter, r *http.Request) bool { return true }}

	filePersistorInfo := long_polling.FilePersistorInfo{Use: false}

	err := server.Start(serverInfo, filePersistorInfo, func(err error) { panic(err) })
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Duration(200) * time.Millisecond)
}

func tearDown(server *long_polling.Server) {
	err := server.Stop(1)
	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	server := long_polling.Server{}

	setUp(&server)

	code := m.Run()

	tearDown(&server)

	os.Exit(code)
}

func TestSubscription(t *testing.T) {
	wg := new(sync.WaitGroup)

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func(index int) {
			defer wg.Done()

			category := "category-1-" + strconv.Itoa(index)
			data := "data-1-" + strconv.Itoa(index)

			publish(t, category, data+"1")
			timestamp, id := subscription(t, long_polling.SubscriptionRequest{Category: category, Timeout: 300, SinceTime: 1}, 1, data)

			publish(t, category, data+"2")
			publish(t, category, data+"3")
			subscription(t, long_polling.SubscriptionRequest{Category: category, Timeout: 300, SinceTime: timestamp, LastID: id}, 2, data)
		}(i)
	}

	wg.Wait()
}

func TestPublish(t *testing.T) {
	wg := new(sync.WaitGroup)

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func(index int) {
			defer wg.Done()

			category := "category-2-" + strconv.Itoa(index)
			data := "data-2-" + strconv.Itoa(index)

			publish(t, category, data+"1")
			timestamp, id := subscription(t, long_polling.SubscriptionRequest{Category: category, Timeout: 300, SinceTime: 1}, 1, data)

			publish(t, category, data+"2")
			publish(t, category, data+"3")
			subscription(t, long_polling.SubscriptionRequest{Category: category, Timeout: 300, SinceTime: timestamp, LastID: id}, 2, data)
		}(i)
	}

	wg.Wait()
}
