package gosnmpHelper

import (
	"fmt"
	"regexp"
	"strings"
)

var regexpMacSeparated = regexp.MustCompile(`^([[:xdigit:]]{2})[:\-]?([[:xdigit:]]{2})[:\-]?([[:xdigit:]]{2})[:\-]?([[:xdigit:]]{2})[:\-]?([[:xdigit:]]{2})[:\-]?([[:xdigit:]]{2})$`)
var regegpMacSep2 = regexp.MustCompile(`^([[:xdigit:]]{1,2})[:\-]([[:xdigit:]]{1,2})[:\-]([[:xdigit:]]{1,2})[:\-]([[:xdigit:]]{1,2})[:\-]([[:xdigit:]]{1,2})[:\-]([[:xdigit:]]{1,2})$`)

// Normalizes a MAC address in one of the following formats:
//    xx-xx-xx-xx-xx-xx
//    xx:xx:xx:xx:xx:xx
//    xxxxxx-xxxxxx
// To (all lower case):
//    xxxxxxxxxxxx
// Also handles a case where the MAC address is a 6-byte value
func NormalizeMac(input string) (string, error) {
	if len(input) == 6 {
		input = fmt.Sprintf("%x", []byte(input))
	}
	if m := regexpMacSeparated.FindStringSubmatch(strings.TrimSpace(input)); m != nil && len(m) == 7 {
		return strings.ToLower(strings.Join(m[1:], "")), nil
	}
	if m := regegpMacSep2.FindStringSubmatch(strings.TrimSpace(input)); m != nil && len(m) == 7 {
		// Must have been some mix of single and double digit bytes all with separators
		// Pad all to two digits
		for i := 1; i < 7; i++ {
			if len(m[i]) == 1 {
				m[i] = "0" + m[i]
			}
		}
		return strings.ToLower(strings.Join(m[1:], "")), nil
	}
	return "", fmt.Errorf("invalid format for MAC address '%-1.20s'", input)
}
