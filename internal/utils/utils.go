package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string, salt string) (string, error) {
	salted := password + salt
	hashedByte, err := bcrypt.GenerateFromPassword([]byte(salted), bcrypt.DefaultCost)
	hashed := string(hashedByte)

	if err != nil {
		return "", err
	}

	return hashed, nil
}

func ComparePassword(hashedPassword, password, salt string) bool {
	salted := password + salt
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(salted))
	return err == nil
}
