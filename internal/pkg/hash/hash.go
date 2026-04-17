package hash

import (
	"crypto/rand"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

const (
	otpChars   = "0123456789"
	tokenChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	codeChars  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func VerifyPassword(hashedPassword, plainPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}
	return err == nil, err
}

func GenerateOTP() (string, error) {
	return randomString(6, otpChars)
}

func GenerateSecureToken() (string, error) {
	return randomString(32, tokenChars)
}

func GenerateInviteCode() (string, error) {
	return randomString(6, codeChars)
}

func randomString(n int, chars string) (string, error) {
	b := make([]byte, n)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		b[i] = chars[num.Int64()]
	}
	return string(b), nil
}
