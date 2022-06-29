package gosnmpHelper

import (
	"github.com/gosnmp/gosnmp"
	"strconv"
)

// Get PDU value as a uint32 value.  PDU value should be a numeric type, else 0 will be returned.
// Truncation due to signed/unsigned mismatch or numeric size are silently ignored, i.e. a
// 64-bit value will be truncated to a 32-bit value.
// Be warned that negative integer values are forced to unsigned.
// If PDU value is nil, 0 will be returned.
func GetAsUint32(pdu gosnmp.SnmpPDU) uint32 {
	if pdu.Value != nil {
		if pdu.Type == gosnmp.OctetString {
			var s string
			switch v := pdu.Value.(type) {
			case string:
				s = v
			case []byte:
				s = string(v)
			default:
				return 0
			}
			i, _ := strconv.ParseUint(s, 10, 32)
			return uint32(i)
		} else {
			switch v := pdu.Value.(type) {
			case uint8:
				return uint32(v)
			case uint16:
				return uint32(v)
			case uint32:
				return v
			case uint64:
				return uint32(v)
			case uint:
				return uint32(v)
			case int8:
				return uint32(v)
			case int16:
				return uint32(v)
			case int32:
				return uint32(v)
			case int64:
				return uint32(v)
			case int:
				return uint32(v)
			}
		}
	}
	return 0
}

// Get PDU value as a uint64 value.  PDU value should be a numeric type, else 0 will be returned.
// Be warned that negative integer values are forced to unsigned.
// If PDU value is nil, 0 will be returned.
func GetAsUint64(pdu gosnmp.SnmpPDU) uint64 {
	if pdu.Value != nil {
		if pdu.Type == gosnmp.OctetString {
			var s string
			switch v := pdu.Value.(type) {
			case string:
				s = v
			case []byte:
				s = string(v)
			default:
				return 0
			}
			i, _ := strconv.ParseUint(s, 10, 64)
			return i
		} else {
			switch v := pdu.Value.(type) {
			case uint8:
				return uint64(v)
			case uint16:
				return uint64(v)
			case uint32:
				return uint64(v)
			case uint64:
				return v
			case uint:
				return uint64(v)
			case int8:
				return uint64(v)
			case int16:
				return uint64(v)
			case int32:
				return uint64(v)
			case int64:
				return uint64(v)
			case int:
				return uint64(v)
			}
		}
	}
	return 0
}

// Get PDU value as a uint value.  PDU value should be a numeric type, else 0 will be returned.
// Be warned that negative integer values are forced to unsigned and truncation may occur on 32-bit architectures.
// If PDU value is nil, 0 will be returned.
func GetAsUint(pdu gosnmp.SnmpPDU) uint {
	if pdu.Value != nil {
		if pdu.Type == gosnmp.OctetString {
			var s string
			switch v := pdu.Value.(type) {
			case string:
				s = v
			case []byte:
				s = string(v)
			default:
				return 0
			}
			i, _ := strconv.ParseUint(s, 10, 64)
			return uint(i)
		} else {
			switch v := pdu.Value.(type) {
			case uint8:
				return uint(v)
			case uint16:
				return uint(v)
			case uint32:
				return uint(v)
			case uint64:
				return uint(v)
			case uint:
				return v
			case int8:
				return uint(v)
			case int16:
				return uint(v)
			case int32:
				return uint(v)
			case int64:
				return uint(v)
			case int:
				return uint(v)
			}
		}
	}
	return 0
}

// Get PDU value as an int32 value.  PDU value should be a numeric type, else 0 will be returned.
// Truncation due to signed/unsigned mismatch or numeric size are silently ignored, i.e. a
// 64-bit value will be truncated to a 32-bit value.
// Be warned that large unsigned values are forced to negative integers.
// If PDU value is nil, 0 will be returned.
func GetAsInt32(pdu gosnmp.SnmpPDU) int32 {
	if pdu.Value != nil {
		if pdu.Type == gosnmp.OctetString {
			var s string
			switch v := pdu.Value.(type) {
			case string:
				s = v
			case []byte:
				s = string(v)
			default:
				return 0
			}
			i, _ := strconv.ParseInt(s, 10, 32)
			return int32(i)
		} else {
			switch v := pdu.Value.(type) {
			case uint8:
				return int32(v)
			case uint16:
				return int32(v)
			case uint32:
				return int32(v)
			case uint64:
				return int32(v)
			case uint:
				return int32(v)
			case int8:
				return int32(v)
			case int16:
				return int32(v)
			case int32:
				return v
			case int64:
				return int32(v)
			case int:
				return int32(v)
			}
		}
	}
	return 0
}

// Get PDU value as a int64 value.  PDU value should be a numeric type, else 0 will be returned.
// Be warned that large unsigned values are forced to negative integers.
// If PDU value is nil, 0 will be returned.
func GetAsInt64(pdu gosnmp.SnmpPDU) int64 {
	if pdu.Value != nil {
		if pdu.Type == gosnmp.OctetString {
			var s string
			switch v := pdu.Value.(type) {
			case string:
				s = v
			case []byte:
				s = string(v)
			default:
				return 0
			}
			i, _ := strconv.ParseInt(s, 10, 64)
			return i
		} else {
			switch v := pdu.Value.(type) {
			case uint8:
				return int64(v)
			case uint16:
				return int64(v)
			case uint32:
				return int64(v)
			case uint64:
				return int64(v)
			case uint:
				return int64(v)
			case int8:
				return int64(v)
			case int16:
				return int64(v)
			case int32:
				return int64(v)
			case int64:
				return v
			case int:
				return int64(v)
			}
		}
	}
	return 0
}

// Get PDU value as an int value.  PDU value should be a numeric type, else 0 will be returned.
// Be warned that large unsigned values are forced to negative integers and truncation may occur
// on 32-bit architectures.
// If PDU value is nil, 0 will be returned.
func GetAsInt(pdu gosnmp.SnmpPDU) int {
	if pdu.Value != nil {
		if pdu.Type == gosnmp.OctetString {
			var s string
			switch v := pdu.Value.(type) {
			case string:
				s = v
			case []byte:
				s = string(v)
			default:
				return 0
			}
			i, _ := strconv.ParseInt(s, 10, 64)
			return int(i)
		} else {
			switch v := pdu.Value.(type) {
			case uint8:
				return int(v)
			case uint16:
				return int(v)
			case uint32:
				return int(v)
			case uint64:
				return int(v)
			case uint:
				return int(v)
			case int8:
				return int(v)
			case int16:
				return int(v)
			case int32:
				return int(v)
			case int64:
				return int(v)
			case int:
				return v
			}
		}
	}
	return 0
}

// Get PDU value as an float32 value.
// It is common for SNMP agents to return floating point values as strings since the ASN opaque float
// is not fully supported by SNMP systems.  If the PDU value is a string, an attempt will be made to
// convert back to float. Be warned that truncation may occur in multiple cases.
// If PDU value is nil, 0 will be returned.
func GetAsFloat32(pdu gosnmp.SnmpPDU) float32 {
	if pdu.Value != nil {
		if pdu.Type == gosnmp.OctetString {
			var s string
			switch v := pdu.Value.(type) {
			case string:
				s = v
			case []byte:
				s = string(v)
			default:
				return 0
			}
			f, _ := strconv.ParseFloat(s, 32)
			return float32(f)
		} else {
			switch v := pdu.Value.(type) {
			case uint8:
				return float32(v)
			case uint16:
				return float32(v)
			case uint32:
				return float32(v)
			case uint64:
				return float32(v)
			case uint:
				return float32(v)
			case int8:
				return float32(v)
			case int16:
				return float32(v)
			case int32:
				return float32(v)
			case int64:
				return float32(v)
			case int:
				return float32(v)
			case float32:
				return v
			case float64:
				return float32(v)
			}
		}
	}
	return 0
}

// Get PDU value as an float64 value.
// It is common for SNMP agents to return floating point values as strings since the ASN opaque float
// is not fully supported by SNMP systems.  If the PDU value is a string, an attempt will be made to
// convert back to float.)
// If PDU value is nil, 0 will be returned.
func GetAsFloat64(pdu gosnmp.SnmpPDU) float64 {
	if pdu.Value != nil {
		if pdu.Type == gosnmp.OctetString {
			var s string
			switch v := pdu.Value.(type) {
			case string:
				s = v
			case []byte:
				s = string(v)
			default:
				return 0
			}
			f, _ := strconv.ParseFloat(s, 64)
			return f
		} else {
			switch v := pdu.Value.(type) {
			case uint8:
				return float64(v)
			case uint16:
				return float64(v)
			case uint32:
				return float64(v)
			case uint64:
				return float64(v)
			case uint:
				return float64(v)
			case int8:
				return float64(v)
			case int16:
				return float64(v)
			case int32:
				return float64(v)
			case int64:
				return float64(v)
			case int:
				return float64(v)
			case float32:
				return float64(v)
			case float64:
				return v
			}
		}
	}
	return 0
}

// Get PDU value as a string.  An empty string will be returned for nil PDU values.
// Numeric values are converted to string format in base-10.
func GetAsString(pdu gosnmp.SnmpPDU) string {
	if pdu.Value != nil {
		switch v := pdu.Value.(type) {
		case uint8:
			return strconv.FormatUint(uint64(v), 10)
		case uint16:
			return strconv.FormatUint(uint64(v), 10)
		case uint32:
			return strconv.FormatUint(uint64(v), 10)
		case uint64:
			return strconv.FormatUint(v, 10)
		case uint:
			return strconv.FormatUint(uint64(v), 10)
		case int8:
			return strconv.FormatInt(int64(v), 10)
		case int16:
			return strconv.FormatInt(int64(v), 10)
		case int32:
			return strconv.FormatInt(int64(v), 10)
		case int64:
			return strconv.FormatInt(v, 10)
		case int:
			return strconv.FormatInt(int64(v), 10)
		case float32:
			return strconv.FormatFloat(float64(v), 'f', -1, 32)
		case float64:
			return strconv.FormatFloat(v, 'f', -1, 64)
		case []byte:
			return string(v)
		case string:
			return v
		}
	}
	return ""
}

// Get PDU value as a slice of bytes.  PDU nil values are returned as an empty slice.
// Any numeric values are first converted to strings, then returned as a byte slice.
func GetAsBytes(pdu gosnmp.SnmpPDU) []byte {
	if pdu.Value != nil {
		switch v := pdu.Value.(type) {
		case uint8, uint16, uint32, uint, uint64, int8, int16, int32, int, int64, float32, float64:
			return []byte(GetAsString(pdu))
		case string:
			return []byte(v)
		case []byte:
			return v
		}
	}
	return []byte{}
}
