package utils

import (
	"errors"
	"strings"
)

func CleanUpHex(s string) string {
	s = strings.Replace(strings.TrimPrefix(s, "0x"), " ", "", -1)

	return strings.ToLower(s)
}

func ValidateAccount(accountAddress string) (string, error) {
	accountAddress = CleanUpHex(accountAddress)
	// check account length
	if len(accountAddress) != 40 {
		return "", errors.New("invalid account address")
	}
	return accountAddress, nil
}

func AppendNotEmpty(slice []string, str string) []string {
	if str != "" {
		return append(slice, str)
	}

	return slice
}
