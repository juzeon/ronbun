package network

import (
	"bytes"
	"context"
	"errors"
	"github.com/sashabaranov/go-openai"
	"io"
	"log/slog"
	"ronbun/storage"
	"ronbun/util"
	"strings"
	"time"
)

func GetOpenAITranslation(text string) string {
	resp := util.AttemptMax(10, func() (string, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		stream, err := openaiClient.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
			Model: storage.Config.OpenAI.Model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: "system",
					Content: "You are a professional academic translator. " +
						"You will translate all the text user provided to Chinese. " +
						"You use techniques like reordering sentences, use native words and phrases to make the translation " +
						"feels natural to Chinese native speakers. " +
						"Output in MarkDown. Use LaTeX for formulas and wrap them with $. " +
						"Do not write explanations and output the translation directly.",
				},
				{
					Role:    "user",
					Content: text,
				},
			},
			Temperature: 0.4,
			Stream:      true,
		})
		if err != nil {
			slog.Warn("Create stream error", "err", err)
			return "", err
		}
		defer stream.Close()
		var out bytes.Buffer
		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				return out.String(), nil
			}
			if err != nil {
				slog.Error("Stream error", "err", err)
				return "", err
			}
			if len(response.Choices) == 0 {
				continue
			}
			out.WriteString(response.Choices[0].Delta.Content)
		}
	})
	content := strings.TrimSuffix(resp, "```")
	content = strings.TrimPrefix(content, "```markdown")
	content = strings.TrimSpace(content)
	return content
}
