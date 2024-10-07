package helper

import "strconv"

func IsLuhn(text string) bool {
	sum := 0
	nDigits := len(text)
	parity := nDigits % 2
	for i := 0; i < nDigits; i++ {
		digit, err := strconv.Atoi(string(text[i]))
		if err != nil {
			return false
		}
		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}
	return (sum % 10) == 0
}
