package auth

import "golang.org/x/crypto/bcrypt"

func HashedPassword(password string) string {
	// For simplicity, we're returning the password as-is.
	// In a real application, use a proper hashing algorithm like bcrypt.
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hash)
}
