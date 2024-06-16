// Package gemini provides communication with Gemini.
package gemini

import (
	"context"
	"errors"

	"github.com/google/generative-ai-go/genai"
)

// Chat is a struct that provides chat.
type Chat struct {
	ctx         context.Context
	client      *genai.Client
	chatSession *genai.ChatSession
}

// Start is start a chat.
//
// ex) err := chat.Start("gemini-1.5-flash", API_KEY)
func (this *Chat) Start(model, apiKey string) error {
	this.ctx = context.Background()

	if client, model, err := getModel(this.ctx, model, apiKey); err != nil {
		return err
	} else {
		this.client = client
		this.chatSession = model.StartChat()

		return nil
	}
}

// Stop is stop a chat.
//
// ex) chat.Stop()
func (this *Chat) Stop() {
	this.chatSession = nil

	if this.client != nil {
		this.client.Close()
		this.client = nil
	}
}

// SendMessage is to send a message or image.
//
// ex 1) answer, err := chat.SendMessage("who are you?", nil)
// ex 2) answer, err := chat.SendMessage("let me know your opinion", []string{"test/image-01.webp"})
func (this *Chat) SendMessage(text string, images []string) (string, error) {
	if parts, err := this.getParts(text, images); err != nil {
		return "", err
	} else if response, err := this.chatSession.SendMessage(this.ctx, parts...); err != nil {
		return "", err
	} else {
		return responseToAnswer(response), nil
	}
}

// SendMessageStream is to send a message or image and receive it as a stream.
//
// ex) channel, err := chat.SendMessageStream("please say something encouraging", nil)
func (this *Chat) SendMessageStream(text string, images []string) (chan answerForStream, error) {
	if parts, err := this.getParts(text, images); err != nil {
		return nil, err
	} else {
		responseIterator := this.chatSession.SendMessageStream(this.ctx, parts...)

		return getStreamChannel(responseIterator), nil
	}
}

// GetHistory returns the history for messages without images.
//
// ex) history := chat.GetHistory()
func (this *Chat) GetHistory() []chatHistory {
	if this.chatSession == nil {
		return nil
	}

	histories := []chatHistory{}

	for _, history := range this.chatSession.History {
		histories = append(histories, chatHistory{Role: history.Role, Answer: contentToAnswer(history)})

	}

	return histories
}

func (this *Chat) getParts(text string, images []string) ([]genai.Part, error) {
	if this.client == nil || this.chatSession == nil {
		return nil, errors.New("Please call the Start method first.")
	}

	return makeParts(text, images)
}
