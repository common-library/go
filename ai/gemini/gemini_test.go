package gemini_test

import (
	"testing"

	"github.com/common-library/go/ai/gemini"
	"github.com/common-library/go/ai/gemini/test-data"
)

func TestQuestion(t *testing.T) {
	t.Parallel()

	if len(test.API_KEY) == 0 {
		return
	}

	if answer, err := gemini.Question(test.MODEL, test.API_KEY, "who are you?", nil); err != nil {
		t.Fatal(err)
	} else {
		t.Log(answer)
	}

	if answer, err := gemini.Question(test.MODEL, test.API_KEY, "let me know your opinion", []string{"test-data/image-01.webp"}); err != nil {
		t.Fatal(err)
	} else {
		t.Log(answer)
	}
}

func TestQuestionStream(t *testing.T) {
	t.Parallel()

	if len(test.API_KEY) == 0 {
		return
	}

	if channel, err := gemini.QuestionStream(test.MODEL, test.API_KEY, "please say something encouraging", nil); err != nil {
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
}
