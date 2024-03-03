// Package gemini provides communication with Gemini.
package gemini

import (
	"context"
	"errors"

	"github.com/google/generative-ai-go/genai"
)

// Chat is a struct that provides chat.
type Chat struct {
	ctx                   context.Context
	client                *genai.Client
	chatSession           *genai.ChatSession
	chatSessionWithImages *genai.ChatSession
}

// Start is start a chat.
//
// ex) err := chat.Start(API_KEY)
func (this *Chat) Start(apiKey string) error {
	this.ctx = context.Background()

	if client, err := getClient(this.ctx, apiKey); err != nil {
		return err
	} else {
		this.client = client

		this.chatSession = client.GenerativeModel("gemini-pro").StartChat()
		this.chatSessionWithImages = client.GenerativeModel("gemini-pro-vision").StartChat()

		return nil
	}
}

// Stop is stop a chat.
//
// ex) chat.Stop()
func (this *Chat) Stop() {
	this.chatSession = nil
	this.chatSessionWithImages = nil

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
	} else if response, err := this.getChatSession(images).SendMessage(this.ctx, parts...); err != nil {
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
		responseIterator := this.getChatSession(images).SendMessageStream(this.ctx, parts...)

		return getStreamChannel(responseIterator), nil
	}
}

// GetHistory returns the history for messages without images.
//
// ex) history := chat.GetHistory()
func (this *Chat) GetHistory() []chatHistory {
	return this.getHistory(this.chatSession)
}

// GetHistory returns the history of messages with images.
//
// ex) history := chat.GetHistoryWithImages()
func (this *Chat) GetHistoryWithImages() []chatHistory {
	return this.getHistory(this.chatSessionWithImages)
}

func (this *Chat) getHistory(chatSession *genai.ChatSession) []chatHistory {
	if chatSession == nil {
		return nil
	}

	histories := []chatHistory{}

	for _, history := range chatSession.History {
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

func (this *Chat) getChatSession(images []string) *genai.ChatSession {
	if len(images) == 0 {
		return this.chatSession
	} else {
		return this.chatSessionWithImages
	}
}
