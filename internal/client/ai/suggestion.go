package ai

import (
	"context"
	"encoding/json"
	"fmt"
)

type (
	SuggestionRequest struct {
		Text     string
		Topic    string
		Question string
		Language string
		Level    string
	}

	SuggestionResponse struct {
		Accuracy float64  `json:"accuracy"`
		Text     string   `json:"text"`
		Hints    []string `json:"hints"`
	}
)

const (
	suggestionPrompt = `
Give me accuracy of this answer and fix grammar errors:
%s

Topic: %s
Question: %s

Please provide response format with json(without any extra text, without formatting, only raw json, without newlines)
Here is an example of a response:
{
	"accuracy": 0.59,
	"text": new text with corrections,
	"hints:" [
	provide some hints 
]
}

Do not provide hints like this "Add a comma after 'Also'."
Generate response in %s
Generate response according to level of student %s
`
)

func (c *Client) GetSuggestions(ctx context.Context, req SuggestionRequest) (obj SuggestionResponse, err error) {
	var out Out

	err = c.do(
		ctx,
		In{
			Model: c.gptModel,
			Messages: []Message{
				{
					Role:    "user",
					Content: fmt.Sprintf(suggestionPrompt, req.Text, req.Topic, req.Question, req.Language, unEmpty(req.Level, defaultLevel)),
				},
			},
		},
		&out,
	)
	if err != nil {
		return obj, fmt.Errorf("generating suggestions: %w", err)
	}

	if len(out.Choices) == 0 {
		return obj, fmt.Errorf("no choices found")
	}

	err = json.Unmarshal([]byte(out.Choices[0].Message.Content), &obj)
	if err != nil {
		return obj, fmt.Errorf("unmarshalling response: %w", err)
	}

	return obj, nil
}
