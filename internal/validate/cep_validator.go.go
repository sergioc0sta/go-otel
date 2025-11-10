package validate

import "strconv"

func CepValidator(cep string) bool {
	if len(cep) != 8 {
		return false
	}

	zone, errr := strconv.Atoi(cep[:5])
	suffix, err := strconv.Atoi(cep[5:])

	if err != nil || errr != nil {
		return false
	}

	if zone < 1000 || zone > 99999 {
		return false
	}

	if suffix < 0 || suffix > 999 {
		return false
	}

	return true
}
