package helpers

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func Dump(i interface{}) {
	fmt.Println(ToJSON(i, "\t"))
}

func ToJSON(i interface{}, indent string) string {
	s, _ := json.MarshalIndent(i, "", indent)
	return string(s)
}

func StringReplacer(val string, replacer map[string]string) string {
	for k, v := range replacer {
		val = strings.Replace(val, fmt.Sprintf("{{%s}}", k), v, -1)
	}
	return val
}

func InArrayString(val string, haystack []string) bool {
	for _, v := range haystack {
		if val == v {
			return true
		}
	}
	return false
}
func InArrayInt(val int, haystack []int) bool {
	for _, v := range haystack {
		if val == v {
			return true
		}
	}
	return false
}

type Debug struct {
	Property   string
	Error      error
	Additional string
}

func (e Debug) String() string {
	return fmt.Sprintf("ERROR (%v): %v | %v", e.Property, e.Error, e.Additional)
}

func RemoveDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func RemoveDuplicateInt(intSlice []int) []int {
	allKeys := make(map[int]bool)
	list := []int{}
	for _, item := range intSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func GetTotalPage(total int64, limit int64) int64 {
	if total%limit == 0 {
		return total / limit
	}
	return (total / limit) + 1
}

func GenerateFormattedCode(code string, count int64, randomChar string) string {

	currentDate := time.Now()
	day := currentDate.Day()
	month := currentDate.Month()
	year := currentDate.Year() % 100

	formattedCode := fmt.Sprintf("%s/%02d%02d%d.%04d/%s", code, day, month, year, count, randomChar)
	return formattedCode
}

func IsValidHexColor(color string) bool {
	pattern := `^#([A-Fa-f0-9]{3}){1,2}$`
	match, _ := regexp.MatchString(pattern, color)
	return match
}

func IsValidEmail(email string) bool {
	if strings.Contains(email, "+") {
		return false
	}
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(pattern, email)
	return match
}

func IsValidURLWithoutProtocol(url string) bool {
	// Regex for a URL without protocol
	re := `^[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*(\.[a-zA-Z]{2,})((/[a-zA-Z0-9-._~:/?#[\]@!$&'()*+,;=]*)*)?$`
	matched, err := regexp.MatchString(re, url)
	if err != nil {
		fmt.Println("Error matching regex:", err)
		return false
	}
	return matched
}

func IsValidSubdomain(subdomain string) bool {
	// Regex for valid subdomain
	re := `^[a-zA-Z0-9]+(-[a-zA-Z0-9]+)*$`
	matched, err := regexp.MatchString(re, subdomain)
	if err != nil {
		fmt.Println("Error matching regex:", err)
		return false
	}
	return matched
}

func IsValidAlphanumeric(input string) bool {
	// Regex for alphanumeric validation
	re := `^[a-zA-Z0-9]+$`
	matched, err := regexp.MatchString(re, input)
	if err != nil {
		fmt.Println("Error matching regex:", err)
		return false
	}
	return matched
}

func IsValidNumeric(input string) bool {
	// Regex for alphanumeric validation
	re := `^[0-9]+$`
	matched, err := regexp.MatchString(re, input)
	if err != nil {
		fmt.Println("Error matching regex:", err)
		return false
	}
	return matched
}

func FormatFloat(format string, n float64) string {
	renderFloatPrecisionMultipliers := [...]float64{
		1,
		10,
		100,
		1000,
		10000,
		100000,
		1000000,
		10000000,
		100000000,
		1000000000,
	}

	renderFloatPrecisionRounders := [...]float64{
		0.5,
		0.05,
		0.005,
		0.0005,
		0.00005,
		0.000005,
		0.0000005,
		0.00000005,
		0.000000005,
		0.0000000005,
	}
	// Special cases:
	//   NaN = "NaN"
	//   +Inf = "+Infinity"
	//   -Inf = "-Infinity"
	if math.IsNaN(n) {
		return "NaN"
	}
	if n > math.MaxFloat64 {
		return "Infinity"
	}
	if n < -math.MaxFloat64 {
		return "-Infinity"
	}

	// default format
	precision := 2
	decimalStr := "."
	thousandStr := ","
	positiveStr := ""
	negativeStr := "-"

	if len(format) > 0 {
		format := []rune(format)

		// If there is an explicit format directive,
		// then default values are these:
		precision = 9
		thousandStr = ""

		// collect indices of meaningful formatting directives
		formatIndx := []int{}
		for i, char := range format {
			if char != '#' && char != '0' {
				formatIndx = append(formatIndx, i)
			}
		}

		if len(formatIndx) > 0 {
			// Directive at index 0:
			//   Must be a '+'
			//   Raise an error if not the case
			// index: 0123456789
			//        +0.000,000
			//        +000,000.0
			//        +0000.00
			//        +0000
			if formatIndx[0] == 0 {
				if format[formatIndx[0]] != '+' {
					panic("RenderFloat(): invalid positive sign directive")
				}
				positiveStr = "+"
				formatIndx = formatIndx[1:]
			}

			// Two directives:
			//   First is thousands separator
			//   Raise an error if not followed by 3-digit
			// 0123456789
			// 0.000,000
			// 000,000.00
			if len(formatIndx) == 2 {
				if (formatIndx[1] - formatIndx[0]) != 4 {
					panic("RenderFloat(): thousands separator directive must be followed by 3 digit-specifiers")
				}
				thousandStr = string(format[formatIndx[0]])
				formatIndx = formatIndx[1:]
			}

			// One directive:
			//   Directive is decimal separator
			//   The number of digit-specifier following the separator indicates wanted precision
			// 0123456789
			// 0.00
			// 000,0000
			if len(formatIndx) == 1 {
				decimalStr = string(format[formatIndx[0]])
				precision = len(format) - formatIndx[0] - 1
			}
		}
	}

	// generate sign part
	var signStr string
	if n >= 0.000000001 {
		signStr = positiveStr
	} else if n <= -0.000000001 {
		signStr = negativeStr
		n = -n
	} else {
		signStr = ""
		n = 0.0
	}

	// split number into integer and fractional parts
	intf, fracf := math.Modf(n + renderFloatPrecisionRounders[precision])

	// generate integer part string
	intStr := strconv.Itoa(int(intf))

	// add thousand separator if required
	if len(thousandStr) > 0 {
		for i := len(intStr); i > 3; {
			i -= 3
			intStr = intStr[:i] + thousandStr + intStr[i:]
		}
	}

	// no fractional part, we can leave now
	if precision == 0 {
		return signStr + intStr
	}

	// generate fractional part
	fracStr := strconv.Itoa(int(fracf * renderFloatPrecisionMultipliers[precision]))
	// may need padding
	if len(fracStr) < precision {
		fracStr = "000000000000000"[:precision-len(fracStr)] + fracStr
	}

	return signStr + intStr + decimalStr + fracStr
}
