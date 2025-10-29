package minify

import "vimagination.zapto.org/javascript/internal"

const maxSafeInt = "9007199254740991"

func isIdentifier(str string) bool {
	for n, r := range str {
		if (n == 0 && !internal.IsIDStart(r)) || (n > 0 && !internal.IsIDContinue(r)) {
			return false
		}
	}

	return true
}

func isSimpleNumber(str string) bool {
	if str == "0" {
		return true
	} else if len(str) == 0 || len(str) > len(maxSafeInt) || str[0] < '1' || str[0] > '9' {
		return false
	}

	for _, r := range str[1:] {
		if r < '0' || r > '9' {
			return false
		}
	}

	return len(str) < len(maxSafeInt) || str <= maxSafeInt
}
