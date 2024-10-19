package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type (
	// Client
	Client struct {
		hc       *http.Client
		token    string
		log      *slog.Logger
		gptModel string
	}

	clientOpt func(*Client)

	In struct {
		Model    string    `json:"model"`
		Messages []Message `json:"messages"`
	}

	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	Out struct {
		Id      string `json:"id"`
		Object  string `json:"object"`
		Created int    `json:"created"`
		Model   string `json:"model"`
		Choices []struct {
			Index   int `json:"index"`
			Message struct {
				Role    string      `json:"role"`
				Content string      `json:"content"`
				Refusal interface{} `json:"refusal"`
			} `json:"message"`
			Logprobs     interface{} `json:"logprobs"`
			FinishReason string      `json:"finish_reason"`
		} `json:"choices"`
		Usage struct {
			PromptTokens        int `json:"prompt_tokens"`
			CompletionTokens    int `json:"completion_tokens"`
			TotalTokens         int `json:"total_tokens"`
			PromptTokensDetails struct {
				CachedTokens int `json:"cached_tokens"`
			} `json:"prompt_tokens_details"`
			CompletionTokensDetails struct {
				ReasoningTokens int `json:"reasoning_tokens"`
			} `json:"completion_tokens_details"`
		} `json:"usage"`
		SystemFingerprint string `json:"system_fingerprint"`
	}
)

const (
	defaultTimeout = 30 * time.Second

	apiURL          = "https://api.openai.com"
	pathCompletions = "/v1/chat/completions"

	gptModel = "gpt-4o-mini"
)

// NewClient
func NewClient(token string, opts ...clientOpt) *Client {
	c := &Client{
		hc: &http.Client{
			Timeout: defaultTimeout,
		},
		token:    token,
		log:      slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		gptModel: gptModel,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// WithCustomTimeout set custom timeout
func WithCustomTimeout(t time.Duration) clientOpt {
	return func(c *Client) {
		c.hc.Timeout = t
	}
}

// WithLogger set custom logger
func WithLogger(l *slog.Logger) clientOpt {
	return func(c *Client) {
		c.log = l
	}
}

// WithGPTModel set custom gpt model
func WithGPTModel(m string) clientOpt {
	return func(c *Client) {
		c.gptModel = m
	}
}

func (c *Client) do(ctx context.Context, in, out any) (err error) {

	var (
		reqUrl   = apiURL + pathCompletions
		httpReq  *http.Request
		httpResp *http.Response
		reqBody  []byte
		respBody []byte
	)

	// Log request details.
	defer func() {
		log := c.log.With(
			"request_method", http.MethodPost,
			"request_url", reqUrl,
			"request_body", string(reqBody),
			"response_body", string(respBody),
		)

		if httpReq != nil {
			log.With("request_headers", httpReq.Header)
		}

		if httpResp != nil {
			log.With("response_headers", httpResp.Header)
			log.With("response_status", httpResp.Status)
		}

		//log.Info("execute request to ai")
	}()

	if in != nil {
		reqBody, _ = json.Marshal(in)
	}

	httpReq, err = http.NewRequest(http.MethodPost, apiURL+pathCompletions, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create http request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.token)

	httpResp, err = c.hc.Do(httpReq.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to do http request: %w", err)
	}
	defer httpResp.Body.Close()

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	err = json.Unmarshal(body, &out)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return nil
}
