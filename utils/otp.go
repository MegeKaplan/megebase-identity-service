package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateOTP() (string, error) {
	otp := ""

	for range 6 {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		otp += fmt.Sprintf("%d", n.Int64())
	}

	return otp, nil
}

func SendOTP(otp string, to string) error {
	fmt.Printf("Sending OTP %s to %s\n", otp, to)
	return nil
}
