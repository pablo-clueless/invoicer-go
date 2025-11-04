package lib

import (
	"fmt"
	"strings"
	"time"
)

func ToKebabCase(str string) string {
	return strings.ReplaceAll(strings.ToLower(str), " ", "-")
}

func GenerateTransactionReference() string {
	timestamp := time.Now().Unix()
	randomStr := fmt.Sprintf("%d", timestamp)
	return fmt.Sprintf("%d%s", timestamp, randomStr)
}
