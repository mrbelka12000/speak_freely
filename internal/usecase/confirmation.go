package usecase

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/mrbelka12000/linguo_sphere_backend/internal/client/mail"
	"github.com/mrbelka12000/linguo_sphere_backend/internal/models"
)

// UserConfirm
func (uc *UseCase) UserConfirm(ctx context.Context, code string) error {
	id, ok := uc.cache.GetInt64(code)
	if !ok {
		return errors.New("can not find code in cache")
	}

	_, err := uc.UserUpdate(ctx, models.UserGet{
		ID: id,
	}, models.UserCU{
		Confirmed: true,
	})
	if err != nil {
		return fmt.Errorf("can not update user: %w", err)
	}

	uc.cache.Delete(code)
	uc.log.Info(fmt.Sprintf("user successfully confirmed: %d", id))

	return nil
}

func (uc *UseCase) sendConfirmationEmail(ctx context.Context, id int64) {
	log := uc.log.With("method", "sendConfirmationEmail")

	user, err := uc.UserGet(ctx, models.UserGet{ID: id})
	if err != nil {
		log.With("error", err).Error("can not get user")
		return
	}

	confirmationCode, err := generateRandomCode()
	if err != nil {
		log.With("error", err).Error("can not generate confirmation string")
		return
	}

	templateData := struct {
		ConfirmationURL string
	}{
		ConfirmationURL: uc.makeConfirmationURL(confirmationCode),
	}

	var b bytes.Buffer
	err = uc.emailConfirmTemplate.Execute(&b, templateData)
	if err != nil {
		log.With("error", err).Error("can not execute email confirmation template")
		return
	}

	err = uc.mailSender.Send(mail.Request{
		To:      user.Email,
		Body:    b.String(),
		Subject: "Confirmation Email",
	})
	if err != nil {
		log.With("error", err).Error("can not send email confirmation")
		return
	}

	err = uc.cache.Set(confirmationCode, user.ID, 1*time.Hour)
	if err != nil {
		log.With("error", err).Error("can not set user confirmation")
		return
	}

	log.Info("user confirmation email sent successfully")
}

func generateRandomCode() (string, error) {
	const digits = "0123456789"
	code := make([]byte, 8)
	for i := 0; i < 8; i++ {
		// Generate a random index for the digits string
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		code[i] = digits[n.Int64()]
	}
	return string(code), nil
}

func (uc *UseCase) makeConfirmationURL(code string) string {
	return fmt.Sprintf("%s/api/v1/confirm?code=%s", uc.publicURL, code)
}
