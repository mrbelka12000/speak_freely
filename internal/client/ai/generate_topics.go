package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
)

type (
	GenerateTopicsRequest struct {
		Level string
		//Preferences []string // TODO after start of project
	}
	GenerateTopicsResponse struct {
		Lang     string `json:"lang"`
		Topic    string `json:"topic"`
		Question string `json:"question"`
	}
)

const (
	topicsPrompt = `Generate random topics to discuss and learn foreign language
please provide response format with json(without any extra text, without formatting, only raw json, without newlines)
Here is an example of a response:
[
   {
		"lang": "en",
		"topic":"Describe your hometown",
		"question":"What are the best things about the place where you grew up?"
   }
]
"lang" should be replaced with "ru" / "en" / "fr" / "es" / "pt" / "ko" / "de" / "it" / "ja" / "tr"
The level of difficulty should be %v
Provide for all languages that I need
Topic and question should be in the same language as "lang"
Topic is %v
`
)

// preferences TODO make dynamic and personalized
var preferences = []string{
	"Travelling", "Family", "Books", "Films", "Science", "Education", "Friends", "Social Media",
	"Work", "Cooking", "Personal Information", "Daily Routine", "Weather", "Food and Drink",
	"Hobbies", "Health", "Fitness", "Plans", "Cultural Differences",
}

func (c *Client) GenerateTopics(ctx context.Context, request GenerateTopicsRequest) ([]GenerateTopicsResponse, error) {
	var out Out

	err := c.do(ctx, In{
		Model: c.gptModel,
		Messages: []Message{
			{
				Role:    "user",
				Content: fmt.Sprintf(topicsPrompt, request.Level, getRandomTopic(preferences)),
			},
		},
	},
		&out,
	)
	if err != nil {
		return nil, fmt.Errorf("generating topics: %w", err)
	}

	if len(out.Choices) == 0 {
		return nil, fmt.Errorf("no choices found")
	}

	var result []GenerateTopicsResponse
	err = json.Unmarshal([]byte(out.Choices[0].Message.Content), &result)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling choice content: %w", err)
	}

	return result, nil
}

func getRandomTopic(preferences []string) string {
	return preferences[rand.Intn(len(preferences))]
}
