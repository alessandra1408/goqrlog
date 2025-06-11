package util

import (
	"fmt"
	"strings"
)

func GetMaskedToken(token string) string {
	token = strings.TrimPrefix(token, "Bearer ")

	tokenLen := len(token)
	if tokenLen <= 10 {
		return strings.Repeat("*", 5)
	}

	return fmt.Sprintf("%s%s", strings.Repeat("*", 5), token[tokenLen-5:])
}
