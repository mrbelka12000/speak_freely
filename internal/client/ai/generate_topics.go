package ai

import (
	"context"
	"encoding/json"
	"fmt"
)

type (
	GenerateThemeRequest struct {
		Level    string
		Language string
		Topic    string
	}
	GenerateThemeResponse struct {
		Question string `json:"question"`
	}
)

const (
	topicsPrompt = `Generate random question to discuss and learn foreign language
please provide response format with json(without any extra text, without formatting, only raw json, without newlines)
Here is an example of a response:
{
		"question":"What are the best things about the place where you grew up?"
}
The level of difficulty should be %v
Provide for all languages that I need
Question should be in the same language as "lang"
Topic of question is %v
Language of question is %v
`
)

func (c *Client) GenerateTheme(ctx context.Context, request GenerateThemeRequest) (empty GenerateThemeResponse, err error) {
	var out Out

	err = c.do(ctx, In{
		Model: c.gptModel,
		Messages: []Message{
			{
				Role:    "user",
				Content: fmt.Sprintf(topicsPrompt, request.Level, request.Topic, request.Language),
			},
		},
	},
		&out,
	)
	if err != nil {
		return empty, fmt.Errorf("generating topics: %w", err)
	}

	if len(out.Choices) == 0 {
		return empty, fmt.Errorf("no choices found")
	}

	fmt.Println(out.Choices[0].Message.Content)
	var result GenerateThemeResponse
	err = json.Unmarshal([]byte(out.Choices[0].Message.Content), &result)
	if err != nil {
		return empty, fmt.Errorf("unmarshaling choice content: %w", err)
	}

	return result, nil
}
