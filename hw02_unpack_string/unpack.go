package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	// fmt.Printf("inp str: %s (% x)\n", s, s)

	var prevRune rune
	var resBuilder strings.Builder
	isEscape := false

	for _, r := range s {
		// fmt.Printf("%d: q=%q, v=%v\n", ri, r, r)

		// заэкранировать можно только цифру или слэш
		if isEscape && r != '\\' && !unicode.IsDigit(r) {
			return "", ErrInvalidString
		}

		// isDigit
		if unicode.IsDigit(r) && !isEscape {

			if prevRune == 0 {
				return "", ErrInvalidString
			}

			digit := int(r - '0')
			for i := 0; i < digit; i++ {
				resBuilder.WriteRune(prevRune)
			}
			prevRune = 0

		} else {

			if prevRune > 0 && !isEscape {
				resBuilder.WriteRune(prevRune)
			}

			if r == '\\' && !isEscape {
				isEscape = true
			} else {
				isEscape = false
				prevRune = r
			}

		}

	}
	// last rune
	if prevRune > 0 {
		resBuilder.WriteRune(prevRune)
	}

	resultStr := resBuilder.String()
	// fmt.Println(resultStr)
	return resultStr, nil
}
