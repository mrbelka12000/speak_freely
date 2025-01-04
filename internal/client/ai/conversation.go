package ai

import (
	"context"
	"encoding/json"
	"fmt"
)

type (
	DialogRequest struct {
		Text     string
		Language string
		Level    string

		Questions []string
		Answers   []string
	}

	DialogResponse struct {
		Answer string
	}
)

const (
	dialogPrompt = `
Using the following text as context: %s

Generate a response in %s according to the context of the previous text, questions, and answers.

Respond only with a raw, valid JSON object. Do not include any extra text, explanations, or formatting.

Response format example:
{
  "answer": "text"
}

Replace "text" with the actual answer.
Don't paste JSON in answer. Always give only raw text
Generate response according to level of student %s
`
)

func (c *Client) Dialog(ctx context.Context, req DialogRequest) (obj DialogResponse, err error) {
	var (
		out Out
		msg []Message
	)

	for i, question := range req.Questions {
		if i >= len(req.Answers) {
			break
		}
		msg = append(msg, Message{
			Role:    "user",
			Content: question,
		})

		msg = append(msg, Message{
			Role:    "assistant",
			Content: req.Answers[i],
		})
	}

	msg = append(msg, Message{
		Role:    "user",
		Content: fmt.Sprintf(dialogPrompt, req.Text, req.Language, unEmpty(req.Level, defaultLevel)),
	})

	err = c.do(ctx,
		In{
			Model:    c.gptModel,
			Messages: msg,
		},
		&out,
	)

	if err != nil {
		return obj, fmt.Errorf("generating dialog: %w", err)
	}

	if len(out.Choices) == 0 {
		return obj, fmt.Errorf("no choises found")
	}

	err = json.Unmarshal([]byte(out.Choices[0].Message.Content), &obj)
	if err != nil {
		return obj, fmt.Errorf("unmarshalling response: %w", err)
	}

	return obj, nil
}
