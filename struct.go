package gosnmpHelper

import (
	"fmt"
	"github.com/soniah/gosnmp"
	"reflect"
	"regexp"
	"strings"
)

/*
Given a struct with oid tags, this function returns a slice of strings of all the OID values which
would need to be queried to fill in the struct and is suitable for passing into gosnmp.Get().
For example:

	var s struct {
        SysDesc     string `oid:".1.3.6.1.2.1.1.1.0"`
        SysObjectId string `oid:".1.3.6.1.2.1.1.2.0"`
        SysUpTime   uint64 `oid:".1.3.6.1.2.1.1.3.0"`
    }
	a := GetOidsFromStructTags(s, true)

The returned slice looks like:
    a = []string{".1.3.6.1.2.1.1.1.0", ".1.3.6.1.2.1.1.2.0", ".1.3.6.1.2.1.1.3.0"}

Nested structures are also parsed if the getNested parameter is true, but pointers to nested structures cannot be nil,
else they will be skipped.

For example:
	type InterfaceInfo struct {
		InterfaceCount int `oid:".1.3.6.1.2.1.2.1.0"`
	}

	type Test4 struct {
		SysDesc     string `oid:".1.3.6.1.2.1.1.1.0"`
		SysObjectId string `oid:".1.3.6.1.2.1.1.2.0"`
		Intfs       *InterfaceInfo
	}
	a := GetOidsFromStructTags(Test4{}, true)

The returned slice looks like the following because the Intfs pointer was nil:
    a = []string{".1.3.6.1.2.1.1.1.0", ".1.3.6.1.2.1.1.2.0"}

Had the call been like:
	a := GetOidsFromStructTags(Test4{ Intfs: new(InterfaceInfo) }, true)

Then the returned slice would look like:
    a = []string{".1.3.6.1.2.1.1.1.0", ".1.3.6.1.2.1.1.2.0", ".1.3.6.1.2.1.2.1.0"}
*/
func GetOidsFromStructTags(source interface{}, getNested bool) []string {
	if source == nil {
		return []string{}
	}
	isnil := false
	srcT := reflect.TypeOf(source)
	if srcT.Kind() == reflect.Ptr {
		srcT = srcT.Elem()
		isnil = reflect.ValueOf(source).IsNil()
	}
	srcV := reflect.Indirect(reflect.ValueOf(source))
	numfields := srcT.NumField()
	result := make([]string, 0, numfields)
	for i := 0; i < numfields; i++ {
		fInfo := srcT.Field(i)
		if oid := fInfo.Tag.Get("oid"); len(oid) > 0 {
			result = append(result, oid)
		}
		if !isnil && getNested {
			field := srcV.Field(i)
			switch field.Kind() {
			case reflect.Ptr:
				if field.Elem().Kind() == reflect.Struct {
					result = append(result, GetOidsFromStructTags(field.Interface(), true)...)
				}
			case reflect.Struct:
				result = append(result, GetOidsFromStructTags(field.Interface(), true)...)
			}
		}
	}
	return result
}

/*
Helper function which processes a slice of PDUs into a struct.
See help text on MarshalPDUToStruct() for details.
*/
func MarshalPDUsToStruct(pdus []gosnmp.SnmpPDU, dest interface{}) {
	for _, pdu := range pdus {
		MarshalPDUToStruct(pdu, dest)
	}
}

/*
Given a struct with oid tags, this function will attempt to copy the value from the supplied PDU to the
matching field in the struct.  For example, given this struct:

    var s struct {
        SysDesc     string `oid:".1.3.6.1.2.1.1.1.0"`
        SysObjectId string `oid:".1.3.6.1.2.1.1.2.0"`
        SysUpTime   uint64 `oid:".1.3.6.1.2.1.1.3.0"`
        SysContact  string `oid:".1.3.6.1.2.1.1.4.0"`
        SysName     string `oid:".1.3.6.1.2.1.1.5.0"`
    }
	var pdu gosnmp.SnmpPDU {
		Name: ".1.3.6.1.2.1.1.4.0",
		Type: gosnmp.OctetString,
		Value: []uint8{'F','O','O'},
	}
	MarshalPDUToStruct(pdu, &s)

The s.SysContact now equals "FOO".  If the OID is not found, False will be returned, else True will be returned.

Nested structs are supported, so the following will also work:

	var s struct {
		ID string
		IP net.Addr
		Data struct {
			SysDesc    string `oid:".1.3.6.1.2.1.1.1.0"`
			SysContact string `oid:".1.3.6.1.2.1.1.4.0"`
		}
		InterfaceCount int `oid:".1.3.6.1.2.1.2.1.0"`
	}

If a struct has a duplicate tag (same OID in two different tags), only the first will be used.

As is idiomatic, the struct fields must be exported (start with upper case) to be elidgble for use.

All fields in the passed struct should be values and not pointers, with the exception of nested structs.
For example, the following SysInfo1 and SysInfo2 are allowed:

	type SysIntfs struct {
		IfDesc       map[string]string `oidx:"\\.1\\.3\\.6\\.1\\.2\\.1\\.2\\.2\\.1\\.2\\.(\\d+)"`
		IfOperStatus map[string]int    `oidx:\.1\.3\.6\.1\.2\.1\.2\.2\.1\.8\.(\d+)`
	}
	type SysInfo1 struct {
		SysDesc     string `oid:".1.3.6.1.2.1.1.1.0"`
		Intfs       SysIntfs
	}
	type SysInfo2 struct {
		SysDesc     string `oid:".1.3.6.1.2.1.1.1.0"`
		Intfs       *SysIntfs
	}

The following is not allowed and no OID match will be made for the field SysName:

	type SysInfo1 struct {
		SysDesc     string `oid:".1.3.6.1.2.1.1.1.0"`
		SysUpTime   uint64 `oid:".1.3.6.1.2.1.1.3.0"`
		SysContact  string `oid:".1.3.6.1.2.1.1.4.0"`
		SysName     *string `oid:".1.3.6.1.2.1.1.5.0"`
	}

Struct tags can be a simple string or a regular expressing.  In the simple string case:
	ex: `oid:".1.3.6.1.2.1.1.1.0"`
The OID value matching the tag will cause the PDU value to be copied into the struct member.
In the case of a struct member which is of type map[string]<type>, the map key value can be
derrived from a regular expression with a capture group.  For example:

		IfDesc map[string]string `oidx:"\\.1\\.3\\.6\\.1\\.2\\.1\\.2\\.2\\.1\\.2\\.(\\d+)"`

In this example, a returned value with OID .1.3.6.1.2.1.2.2.1.2.6 will get matched and will cause the value
to be inserted into the map with key="6".  Because of the way Go processes the internal double quoted strings
within a struct tag, the internal regx string must be doubly escaped; hence all the "\\" in the previous example.
An alternate tag format is allowed where no internal double quotes are used as follows:

		IfDesc map[string]string `oidx:\.1\.3\.6\.1\.2\.1\.2\.2\.1\.2\.(\d+)`
*/
func MarshalPDUToStruct(pdu gosnmp.SnmpPDU, dest interface{}) bool {
	if dest == nil {
		return false
	}
	if reflect.TypeOf(dest).Kind() != reflect.Ptr {
		panic(fmt.Errorf("dest must be a pointer to a struct"))
	}
	if reflect.TypeOf(dest).Elem().Kind() != reflect.Struct {
		panic(fmt.Errorf("dest must be a pointer to a struct"))
	}
	destType := reflect.TypeOf(dest).Elem()
	destV := reflect.ValueOf(dest)
	// Walk through fields of struct and check for oid tag
	for i := 0; i < destType.NumField(); i++ {
		tag := destType.Field(i).Tag
		v := destV.Elem().Field(i)
		switch v.Kind() {
		case reflect.Map:
			if m, found := processOidTag(tag, pdu.Name); found {
				assignToMap(m[1], pdu, v)
				return true
			}
		case reflect.Uint, reflect.Uint32, reflect.Uint64:
			if _, found := processOidTag(tag, pdu.Name); found {
				v.SetUint(GetAsUint64(pdu))
				return true
			}
		case reflect.Int, reflect.Int32, reflect.Int64:
			if _, found := processOidTag(tag, pdu.Name); found {
				v.SetInt(GetAsInt64(pdu))
				return true
			}
		case reflect.Float32, reflect.Float64:
			if _, found := processOidTag(tag, pdu.Name); found {
				v.SetFloat(GetAsFloat64(pdu))
				return true
			}
		case reflect.String:
			if _, found := processOidTag(tag, pdu.Name); found {
				v.SetString(GetAsString(pdu))
				return true
			}
		case reflect.Slice:
			if _, found := processOidTag(tag, pdu.Name); found {
				// Only []byte is supported here
				if v.Type().Elem().Kind() == reflect.Uint8 {
					v.SetBytes(GetAsBytes(pdu))
				} else {
					panic(fmt.Errorf("[]byte are the only slice types supported"))
				}
				return true
			}
		case reflect.Struct:
			found := MarshalPDUToStruct(pdu, v.Addr().Interface())
			return found
		case reflect.Ptr:
			// We only deal with pointers to structs here
			if v.Type().Elem().Kind() == reflect.Struct {
				if v.IsNil() {
					v.Set(reflect.New(v.Type().Elem()))
				}
				found := MarshalPDUToStruct(pdu, v.Interface())
				return found
			}
		}
	}
	return false
}

func processOidTag(tag reflect.StructTag, pduName string) ([]string, bool) {
	var oidMatches []string
	processValue := false
	if oid := tag.Get("oid"); len(oid) > 0 && oid == pduName {
		processValue = true
	} else if strings.HasPrefix(string(tag), "oidx:") {
		if oid = tag.Get("oidx"); len(oid) > 0 {
			processValue, oidMatches = oidxMatch(oid, pduName)
		} else if oid = string(tag)[5:]; len(oid) > 0 {
			processValue, oidMatches = oidxMatch(oid, pduName)
		}
	}
	return oidMatches, processValue
}

// oidxMatch will take a RegX pattern and match it against the PDU Name value
// Any captured values will be returned
func oidxMatch(oidPattern string, pduName string) (bool, []string) {
	var (
		rx  *regexp.Regexp
		err error
	)
	if rx, err = regexp.Compile(oidPattern); err != nil {
		return false, nil
	}
	vars := rx.FindStringSubmatch(pduName)
	if vars == nil {
		return false, nil
	}
	return true, vars
}

func assignToMap(key string, pdu gosnmp.SnmpPDU, v reflect.Value) {
	t := v.Type()
	if v.IsNil() {
		v.Set(reflect.MakeMap(t))
	}
	et := t.Elem()
	v.SetMapIndex(reflect.ValueOf(key), getAsValue(pdu, et.Kind()))
}

func getAsValue(pdu gosnmp.SnmpPDU, destKind reflect.Kind) reflect.Value {
	switch destKind {
	case reflect.Uint:
		return reflect.ValueOf(GetAsUint(pdu))
	case reflect.Uint32:
		return reflect.ValueOf(GetAsUint32(pdu))
	case reflect.Uint64:
		return reflect.ValueOf(GetAsUint64(pdu))
	case reflect.Int:
		return reflect.ValueOf(GetAsInt(pdu))
	case reflect.Int32:
		return reflect.ValueOf(GetAsInt32(pdu))
	case reflect.Int64:
		return reflect.ValueOf(GetAsInt64(pdu))
	case reflect.Float32:
		return reflect.ValueOf(GetAsFloat32(pdu))
	case reflect.Float64:
		return reflect.ValueOf(GetAsFloat64(pdu))
	case reflect.String:
		return reflect.ValueOf(GetAsString(pdu))
	}
	panic(fmt.Errorf("unsupported destination type, must be Int, Uint, Float, or String"))
}
