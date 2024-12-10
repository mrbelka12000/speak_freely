package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/mrbelka12000/speak_freely/internal/client/ai"
	"github.com/mrbelka12000/speak_freely/internal/models"
)

func (uc *UseCase) Conversation(
	ctx context.Context,
	fileURL string,
	externalUserID int64,
) (string, error) {

	conversationID, ok := uc.cache.Get(fmt.Sprintf("conversation_%d", externalUserID))
	if !ok {
		conversationID = uuid.New().String()
		err := uc.cache.Set(fmt.Sprintf("conversation_%d", externalUserID), conversationID, 1*time.Hour)
		if err != nil {
			uc.log.With("error", err).Error("can not save conversation")
		}
	}

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
		Text:           text,
		Language:       user.Language.LongName,
		ConversationID: conversationID,
	})
	if err != nil {
		return "", fmt.Errorf("get dialog: %w", err)
	}

	return dialog.Answer, nil
}
