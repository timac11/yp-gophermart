package util

import (
	"regexp"
)

var digitsExpression *regexp.Regexp = regexp.MustCompile(`^[0-9]+$`)

func CheckOrderNum(number string) bool {
	if !digitsExpression.MatchString(number) {
		return false
	}

	sum := 0

	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0')

		if i%2 == 1 {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
	}

	return sum%10 == 0
}
