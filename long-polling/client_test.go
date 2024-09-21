package long_polling_test

import (
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	long_polling "github.com/common-library/go/long-polling"
	"github.com/google/uuid"
)

var address string

func subscription(t *testing.T, request long_polling.SubscriptionRequest, count int, data string) (int64, string) {
	response, err := long_polling.Subscription("http://"+address+"/subscription", nil, request, "", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatal(response.StatusCode, http.StatusText(response.StatusCode))
	}

	if len(response.Events) != count {
		t.Fatal(len(response.Events), count)
	}

	for _, event := range response.Events {
		if event.Category != request.Category {
			t.Fatal(event.Category, request.Category)
		}

		if strings.HasPrefix(event.Data, data) == false {
			t.Fatal(event.Data, data)
		}
	}

	return response.Events[len(response.Events)-1].Timestamp, response.Events[len(response.Events)-1].ID
}

func publish(t *testing.T, category, data string) {
	request := long_polling.PublishRequest{Category: category, Data: data}
	if response, err := long_polling.Publish("http://"+address+"/publish", 10, nil, request, "", "", nil); err != nil {
		t.Fatal(err)
	} else if response.StatusCode != http.StatusOK {
		t.Fatal(response.StatusCode, http.StatusText(response.StatusCode))
	} else if response.Body != `{"success": true}` {
		t.Fatal(response.Body)
	}
}

func setUp(server *long_polling.Server) {
	address = ":" + strconv.Itoa(10000+rand.IntN(1000))

	serverInfo := long_polling.ServerInfo{
		Address:                        address,
		TimeoutSeconds:                 3600,
		SubscriptionURI:                "/subscription",
		HandlerToRunBeforeSubscription: func(w http.ResponseWriter, r *http.Request) bool { return true },
		PublishURI:                     "/publish",
		HandlerToRunBeforePublish:      func(w http.ResponseWriter, r *http.Request) bool { return true }}

	filePersistorInfo := long_polling.FilePersistorInfo{Use: false}

	if err := server.Start(serverInfo, filePersistorInfo, func(err error) { panic(err) }); err != nil {
		panic(err)
	}

	time.Sleep(200 * time.Millisecond)
}

func tearDown(server *long_polling.Server) {
	if err := server.Stop(1 * time.Second); err != nil {
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
	t.Parallel()

	category := "category-" + uuid.New().String()
	data := "data-" + uuid.New().String()

	publish(t, category, data+"1")
	timestamp, id := subscription(t, long_polling.SubscriptionRequest{Category: category, TimeoutSeconds: 300, SinceTime: 1}, 1, data)

	time.Sleep(100 * time.Millisecond)

	publish(t, category, data+"2")
	publish(t, category, data+"3")
	subscription(t, long_polling.SubscriptionRequest{Category: category, TimeoutSeconds: 300, SinceTime: timestamp, LastID: id}, 2, data)
}

func TestPublish(t *testing.T) {
	TestSubscription(t)
}
