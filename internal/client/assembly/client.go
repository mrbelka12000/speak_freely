package assembly

import (
	"github.com/AssemblyAI/assemblyai-go-sdk"
)

// Assembly
type Assembly struct {
	client *assemblyai.Client
}

// NewAssembly
func NewAssembly(apiKey string) *Assembly {
	return &Assembly{
		client: assemblyai.NewClient(apiKey),
	}
}
