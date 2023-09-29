package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func CountDecimalPlaces(decimal float64) (int, error) {
	// Convert the float to a string in scientific notation
	strNum := fmt.Sprintf("%e", decimal)
	// Split the string using "e" as the delimiter
	parts := strings.Split(strNum, "e")
	// If there is no exponent (e.g., num is a whole number), return 0
	if len(parts) == 1 {
		return 0, nil
	}
	// Extract the exponent part and convert it to an integer
	exponent := strings.TrimSpace(parts[1])
	exp, parseIntError := strconv.Atoi(exponent)
	if parseIntError != nil {
		return 0, parseIntError
	}
	// Calculate the number of decimal places
	decimalPlaces := -exp
	return decimalPlaces, nil
}
