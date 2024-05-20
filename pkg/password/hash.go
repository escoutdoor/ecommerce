package password

import "golang.org/x/crypto/bcrypt"

func HashPassword(pw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func ComparePasswords(hashed string, pw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pw))
	return err == nil
}
