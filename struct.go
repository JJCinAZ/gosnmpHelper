package gosnmpHelper

import (
	"fmt"
	"github.com/soniah/gosnmp"
	"reflect"
	"regexp"
	"strings"
)

func GetOidsFromStructTags(source interface{}) []string {
	var destType reflect.Type

	if reflect.TypeOf(source).Kind() == reflect.Ptr {
		destType = reflect.TypeOf(source).Elem()
		if destType.Kind() != reflect.Struct {
			panic("source must be a struct or a pointer to a struct")
		}
	} else if reflect.TypeOf(source).Kind() == reflect.Struct {
		destType = reflect.TypeOf(source)
	} else {
		panic("source must be a struct or a pointer to a struct")
	}
	numFields := destType.NumField()
	result := make([]string, 0, numFields)
	for i := 0; i < numFields; i++ {
		tag := destType.Field(i).Tag
		if oid := tag.Get("oid"); len(oid) > 0 {
			result = append(result, oid)
		}
	}
	return result
}

func MarshalPDUsToStruct(pdus []gosnmp.SnmpPDU, dest interface{}) {
	for _, pdu := range pdus {
		MarshalPDUToStruct(pdu, dest)
	}
}

/*

 */
func MarshalPDUToStruct(pdu gosnmp.SnmpPDU, dest interface{}) {
	var (
		oidMatches []string
	)
	if reflect.TypeOf(dest).Kind() != reflect.Ptr {
		panic(fmt.Errorf("dest must be a pointer to a struct"))
	}
	destType := reflect.TypeOf(dest).Elem()
	// Walk through fields of struct and check for oid tag
	for i := 0; i < destType.NumField(); i++ {
		processValue := false
		tag := destType.Field(i).Tag
		if oid := tag.Get("oid"); len(oid) > 0 && oid == pdu.Name {
			processValue = true
		} else if strings.HasPrefix(string(tag), "oidx:") {
			if oid = tag.Get("oidx"); len(oid) > 0 {
				processValue, oidMatches = oidxMatch(oid, pdu.Name)
			} else if oid = string(tag)[5:]; len(oid) > 0 {
				processValue, oidMatches = oidxMatch(oid, pdu.Name)
			}
		}
		if processValue {
			v := reflect.ValueOf(dest).Elem().Field(i)
			switch v.Kind() {
			case reflect.Map:
				assignToMap(oidMatches[1], pdu, v)
			case reflect.Uint, reflect.Uint32, reflect.Uint64:
				v.SetUint(GetAsUint64(pdu))
			case reflect.Int, reflect.Int32, reflect.Int64:
				v.SetInt(GetAsInt64(pdu))
			case reflect.Float32, reflect.Float64:
				v.SetFloat(GetAsFloat64(pdu))
			case reflect.String:
				v.SetString(GetAsString(pdu))
			case reflect.Slice:
				// Only []byte is supported here
				if v.Type().Elem().Kind() == reflect.Uint8 {
					v.SetBytes(GetAsBytes(pdu))
				} else {
					panic(fmt.Errorf("[]byte are the only slice types supported"))
				}
			case reflect.Struct:
				MarshalPDUToStruct(pdu, &v)
			}
			break
		}
	}
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
