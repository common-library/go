// Package gemini provides communication with Gemini.
package gemini

import (
	"context"
)

// Question is to send a message or image.
//
// ex) answer, err := gemini.Question("gemini-1.5-flash", API_KEY, "who are you?", nil)
func Question(model, apiKey, text string, images []string) (string, error) {
	ctx := context.Background()

	if client, model, err := getModel(ctx, model, apiKey); err != nil {
		return "", err
	} else {
		defer client.Close()

		if parts, err := makeParts(text, images); err != nil {
			return "", err
		} else if response, err := model.GenerateContent(ctx, parts...); err != nil {
			return "", err
		} else {
			return responseToAnswer(response), nil
		}
	}
}

// QuestionStream is to send a message or image and receive it as a stream.
//
// ex) channel, err := gemini.QuestionStream("gemini-1.5-flash", API_KEY, "please say something encouraging", nil)
func QuestionStream(model, apiKey, text string, images []string) (chan answerForStream, error) {
	ctx := context.Background()

	if client, model, err := getModel(ctx, model, apiKey); err != nil {
		return nil, err
	} else {
		defer client.Close()

		if parts, err := makeParts(text, images); err != nil {
			return nil, err
		} else {
			responseIterator := model.GenerateContentStream(ctx, parts...)

			return getStreamChannel(responseIterator), nil
		}
	}
}
