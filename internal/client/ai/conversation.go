package ai

import (
	"context"
	"encoding/json"
	"fmt"
)

type (
	DialogRequest struct {
		Text           string
		Language       string
		ConversationID string
	}

	DialogResponse struct {
		Answer string
	}
)

const (
	dialogPrompt = `
	Let's imagine that you are teacher, have a dialog with your student. 

Text to answer: %s

Please provide response format with json(without any extra text, without formatting, only raw json, without newlines)
Here is an example of a response:
{
"answer": text to continue dialog
}

Generate response in %s
Save information for conversation %s and answer with information according to this ID  
`
)

func (c *Client) Dialog(ctx context.Context, req DialogRequest) (obj DialogResponse, err error) {
	var out Out

	err = c.do(ctx,
		In{
			Model: c.gptModel,
			Messages: []Message{
				{
					Role:    "user",
					Content: fmt.Sprintf(dialogPrompt, req.Text, req.Language, req.ConversationID),
				},
			},
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
