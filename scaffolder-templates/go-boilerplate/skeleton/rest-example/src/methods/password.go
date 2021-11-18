package methods

import "strings"

// ValidatePassword is used to check password validation
func ValidatePassword(pass string) int {
	// checking length
	if len(pass) < 8 {
		return 1
	}
	// checking password contains any special characters
	if !strings.ContainsAny(pass, "!@#$%^&*-?") {
		return 1
	}
	// checking password contains any uppercase aphabets
	if !strings.ContainsAny(pass, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return 1
	}
	// checking password contains any lowercase aphabets
	if !strings.ContainsAny(pass, "abcdefghijklmnopqrstuvwxyz") {
		return 1
	}
	// checking password contains any numeric
	if !strings.ContainsAny(pass, "1234567890") {
		return 1
	}
	return 0
}

// HashForNewPassword is a function for encoding password
func HashForNewPassword(pass string) string {

	// random string is generated as a key
	key := RandomString(5)
	// passhash generated
	passwordHash := key + "." + Sign(key, pass)
	// return hashed password
	return passwordHash
}

// CheckHashForPassword is a method for matching two passwords
func CheckHashForPassword(passwordHash, password string) bool {

	//spliting hash password on basis of .
	passHashParts := strings.Split(passwordHash, ".")
	// if there are less then or more then two parts the hashed password is wrong
	if len(passHashParts) != 2 {
		return false
	}

	// calculate hashpassword using new password
	reCalculatedHash := passHashParts[0] + "." + Sign(passHashParts[0], password)
	// chech old hash with new hash
	if reCalculatedHash == passwordHash {
		return true
	}

	return false
}
