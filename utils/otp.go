package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/MegeKaplan/megebase-identity-service/messaging"
)

func GenerateOTP() (string, error) {
	otp := ""

	for range 6 {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		otp += fmt.Sprintf("%d", n.Int64())
	}

	return otp, nil
}

func SendOTP(messagingService messaging.MessagingService, channel string, to string, otp string) error {
	msg := messaging.MessageEvent{
		Service: "identity",
		Entity:  "otp",
		Action:  "sent",
		Channel: channel,
		To:      to,
		Data: map[string]interface{}{
			"otp": otp,
		},
	}

	if err := messagingService.PublishMessage(msg); err != nil {
		return err
	}

	return nil
}
