package assembly

import (
	"github.com/AssemblyAI/assemblyai-go-sdk"
)

type Assembly struct {
	client *assemblyai.Client
}

func NewAssembly(apiKey string) *Assembly {
	return &Assembly{
		client: assemblyai.NewClient(apiKey),
	}
}
