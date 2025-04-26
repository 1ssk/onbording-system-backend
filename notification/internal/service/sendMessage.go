package service

import (
	"fmt"
	"math"
	"time"

	"go.uber.org/zap"
	"gopkg.in/gomail.v2"

	"notification/internal/config"
)

const (
	maxRetries      = 3
	basicRetryPause = 1
)

func (s *Service) SendMessage(cfg *config.Config, logger *zap.Logger, to string, subject string, message string) error {
	msg := gomail.NewMessage()

	msg.SetHeader("From", cfg.SendMail.SenderEmail)
	if to == "" {
		return fmt.Errorf("empty address")
	}
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", message)

	dialer := gomail.NewDialer(
		"smtp.mail.ru",
		587,
		cfg.SendMail.SenderEmail,
		cfg.SendMail.SenderPassword,
	)

	err := dialer.DialAndSend(msg)
	if err != nil {
		for i := 0; i < maxRetries; i++ {
			pause := math.Pow(2, float64(i)) + basicRetryPause
			time.Sleep(time.Duration(pause) * time.Second)

			err = dialer.DialAndSend(msg)
			if err == nil {
				logger.Info("Successfully sent message")
				return nil
			}
		}

		logger.Error(cannotSendMessage(to, err))
		return fmt.Errorf(cannotSendMessage(to, err))
	}

	logger.Info("Successfully sent message")
	return nil
}

func cannotSendMessage(to string, err error) string {
	return fmt.Sprintf("cannot send message to %s : %v", to, err)
}
