package network

import (
	"context"
	"github.com/sashabaranov/go-openai"
	"log/slog"
	"ronbun/storage"
	"ronbun/util"
	"time"
)

func GetOpenAITranslation(text string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	resp := util.AttemptMax(10, func() (openai.ChatCompletionResponse, error) {
		r, err := openaiClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
			Model: storage.Config.OpenAI.Model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: "system",
					Content: "You are a professional academic translator. " +
						"You will translate all the text user provided to Chinese. " +
						"You use techniques like reordering sentences, use native words and phrases to make the translation " +
						"feels natural to Chinese native speakers. " +
						"Output in MarkDown. Do not write explanations and output the translation directly.",
				},
				{
					Role:    "user",
					Content: text,
				},
			},
			Temperature: 0.4,
			Stream:      false,
		})
		if err != nil {
			slog.Warn("Failed to request to translate, retrying...", "text", text)
			return openai.ChatCompletionResponse{}, err
		}
		if len(r.Choices) == 0 || r.Choices[0].Message.Content == "" {
			slog.Warn("Choices empty, retrying...", "response", r, "text", text)
			return openai.ChatCompletionResponse{}, err
		}
		return r, nil
	})
	return resp.Choices[0].Message.Content
}
