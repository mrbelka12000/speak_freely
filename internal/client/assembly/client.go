package assembly

import (
	"context"
	"fmt"
	"io"

	aai "github.com/AssemblyAI/assemblyai-go-sdk"
)

// Client
type Client struct {
	client *aai.Client
}

// New
func New(apiKey string) *Client {
	return &Client{
		client: aai.NewClient(apiKey),
	}
}

func (c *Client) GetTextFromFile(ctx context.Context, file io.Reader, lang string) (string, error) {
	params := &aai.TranscriptOptionalParams{
		LanguageCode: aai.TranscriptLanguageCode(lang),
	}

	resp, err := c.client.Transcripts.TranscribeFromReader(ctx, file, params)
	if err != nil {
		return "", fmt.Errorf("transcribe from file: %w", err)
	}

	return aai.ToString(resp.Text), nil
}

func (c *Client) GetTextFromURL(ctx context.Context, url, languageCode string) (string, error) {
	params := aai.TranscriptOptionalParams{
		LanguageCode: aai.TranscriptLanguageCode(languageCode),
	}

	transcript, err := c.client.Transcripts.TranscribeFromURL(ctx, url, &params)
	if err != nil {
		return "", err
	}

	return aai.ToString(transcript.Text), nil
}
