// Package gemini provides communication with Gemini.
package gemini

import (
	"context"
)

// Question is to send a message or image.
//
// ex)
func Question(apiKey, text string, images []string) (string, error) {
	ctx := context.Background()

	client, err := getClient(ctx, apiKey)
	if err != nil {
		return "", err
	}
	defer client.Close()

	model := getModel(client, images)

	if parts, err := makeParts(text, images); err != nil {
		return "", err
	} else if response, err := model.GenerateContent(ctx, parts...); err != nil {
		return "", err
	} else {
		return responseToAnswer(response), nil
	}
}

// QuestionStream is to send a message or image and receive it as a stream.
//
// ex)
func QuestionStream(apiKey, text string, images []string) (chan answerForStream, error) {
	ctx := context.Background()

	client, err := getClient(ctx, apiKey)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	model := getModel(client, images)

	if parts, err := makeParts(text, images); err != nil {
		return nil, err
	} else {
		responseIterator := model.GenerateContentStream(ctx, parts...)

		return getStreamChannel(responseIterator), nil
	}
}
