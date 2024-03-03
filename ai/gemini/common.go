// Package gemini provides communication with Gemini.
package gemini

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type answerForStream struct {
	Answer string
	Err    error
}

type chatHistory struct {
	Role   string
	Answer string
}

func responseToAnswer(response *genai.GenerateContentResponse) string {
	if response == nil {
		return ""
	}

	answer := ""

	for _, candidate := range response.Candidates {
		answer += contentToAnswer(candidate.Content)
	}

	return answer
}

func contentToAnswer(content *genai.Content) string {
	if content == nil {
		return ""
	}

	answer := ""

	for _, part := range content.Parts {
		builder := strings.Builder{}
		fmt.Fprintf(&builder, "%v", part)
		answer += builder.String()
	}

	return answer
}

func makeParts(text string, images []string) ([]genai.Part, error) {
	parts := []genai.Part{}
	if len(text) != 0 {
		parts = append(parts, genai.Text(text))
	}
	for _, image := range images {
		if data, err := os.ReadFile(image); err != nil {
			return nil, err
		} else if extension := strings.ToLower(filepath.Ext(image)); len(extension) > 0 && extension[0] == '.' {
			parts = append(parts, genai.ImageData(extension[1:], data))
		} else {
			parts = append(parts, genai.ImageData(extension, data))
		}
	}

	return parts, nil
}

func getClient(ctx context.Context, apiKey string) (*genai.Client, error) {
	return genai.NewClient(ctx, option.WithAPIKey(apiKey))
}

func getModel(client *genai.Client, images []string) *genai.GenerativeModel {
	if len(images) == 0 {
		return client.GenerativeModel("gemini-pro")
	} else {
		return client.GenerativeModel("gemini-pro-vision")
	}
}

func getStreamChannel(responseIterator *genai.GenerateContentResponseIterator) chan answerForStream {
	channel := make(chan answerForStream)

	go func() {
		for {
			if response, err := responseIterator.Next(); err == iterator.Done {
				close(channel)
				break
			} else if err != nil {
				channel <- answerForStream{Answer: "", Err: err}
				close(channel)
				break
			} else {
				channel <- answerForStream{Answer: responseToAnswer(response), Err: nil}
			}
		}
	}()

	return channel
}
