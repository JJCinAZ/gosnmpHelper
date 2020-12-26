package gosnmpHelper

import "strings"

// Return the last portion of an OID
// For example GetOIDIndex(".1.2.4.5.2.3.6.7") return "7"
func GetOIDIndex(oid string) string {
	parts := strings.Split(oid, ".")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}
