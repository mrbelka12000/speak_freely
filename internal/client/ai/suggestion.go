package ai

import "context"

type (
	SuggestionRequest struct {
		Text     string
		Topic    string
		Question string
	}

	SuggestionResponse struct {
		Accuracy    float64      `json:"accuracy"`
		Suggestions []Suggestion `json:"suggestions"`
	}

	Suggestion struct {
		From int    `json:"from"`
		To   int    `json:"to"`
		Text string `json:"text"`
	}
)

const (
	suggestionPrompt = `
Give me accuracy of this answer:
%s

Topic: %s
Question: %s

Please provide response format with json(without any extra text, without formatting, only raw json, without newlines)
Here is an example of a response:
{
	"accurancy": 0.91
	"suggestions" :[
		{
			"start" : from,
			"end" : to,
			"text" : text to replace
		}
	]
}
`
)

func (c *Client) GetSuggestions(ctx context.Context, req *SuggestionRequest) (empty SuggestionResponse, err error) {

	return
}
