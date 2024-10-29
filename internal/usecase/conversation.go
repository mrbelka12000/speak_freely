package usecase

import (
	"context"
	"fmt"

	"github.com/mrbelka12000/speak_freely/internal/client/ai"
	"github.com/mrbelka12000/speak_freely/internal/models"
)

func (uc *UseCase) Conversation(
	ctx context.Context,
	fileURL string,
	externalUserID int64,
) (string, error) {

	user, err := uc.UserGet(ctx, models.UserGet{
		ExternalID: fmt.Sprint(externalUserID),
	})
	if err != nil {
		return "", fmt.Errorf("get user: %w", err)
	}

	text, err := uc.transcriber.GetTextFromURL(ctx, fileURL, user.Language.ShortName)
	if err != nil {
		return "", fmt.Errorf("get text from message: %w", err)
	}

	dialog, err := uc.gen.Dialog(ctx, ai.DialogRequest{
		Text:     text,
		Language: user.Language.LongName,
	})
	if err != nil {
		return "", fmt.Errorf("get dialog: %w", err)
	}

	return dialog.Answer, nil
}
