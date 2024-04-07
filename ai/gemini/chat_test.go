package gemini_test

import (
	"testing"

	"github.com/common-library/go/ai/gemini"
	"github.com/common-library/go/ai/gemini/test-data"
)

func TestChat(t *testing.T) {
	return

	chat := gemini.Chat{}
	if err := chat.Start(test.API_KEY); err != nil {
		t.Fatal(err)
	}
	defer chat.Stop()

	if answer, err := chat.SendMessage("who are you?", nil); err != nil {
		t.Fatal(err)
	} else {
		t.Log(answer)
	}

	if answer, err := chat.SendMessage("let me know your opinion", []string{"test-data/image-01.webp"}); err != nil {
		t.Fatal(err)
	} else {
		t.Log(answer)
	}

	if channel, err := chat.SendMessageStream("please say something encouraging", nil); err != nil {
		t.Fatal(err)
	} else {
		for answer := range channel {
			if answer.Err != nil {
				t.Fatal(err)
			} else {
				t.Log(answer.Answer)
			}
		}
	}

	for _, history := range chat.GetHistory() {
		t.Log(history.Answer)
	}

	for _, history := range chat.GetHistoryWithImages() {
		t.Log(history.Answer)
	}
}
