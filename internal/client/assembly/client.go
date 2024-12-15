package assembly

import (
	"context"
	"fmt"
	"io"

	aai "github.com/AssemblyAI/assemblyai-go-sdk"

	"github.com/mrbelka12000/speak_freely/pkg/pointer"
)

type (
	// Client
	Client struct {
		client *aai.Client
	}

	FileData struct {
		Text          string
		AudioDuration float64
	}
)

// New
func New(apiKey string) *Client {
	return &Client{
		client: aai.NewClient(apiKey),
	}
}

func (c *Client) GetDataFromFile(ctx context.Context, file io.Reader, lang string) (out FileData, err error) {
	params := &aai.TranscriptOptionalParams{
		LanguageCode: aai.TranscriptLanguageCode(lang),
	}

	resp, err := c.client.Transcripts.TranscribeFromReader(ctx, file, params)
	if err != nil {
		return out, fmt.Errorf("transcribe from file: %w", err)
	}

	return FileData{
		Text:          pointer.Value(resp.Text),
		AudioDuration: pointer.Value(resp.AudioDuration),
	}, nil
}

func (c *Client) GetDataFromURL(ctx context.Context, url, languageCode string) (out FileData, err error) {
	params := aai.TranscriptOptionalParams{
		LanguageCode: aai.TranscriptLanguageCode(languageCode),
	}

	resp, err := c.client.Transcripts.TranscribeFromURL(ctx, url, &params)
	if err != nil {
		return out, err
	}

	return FileData{
		Text:          pointer.Value(resp.Text),
		AudioDuration: pointer.Value(resp.AudioDuration),
	}, nil
}
