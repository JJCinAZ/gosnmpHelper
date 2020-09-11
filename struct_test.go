package gosnmpHelper

import (
	"github.com/davecgh/go-spew/spew"
	snmp "github.com/soniah/gosnmp"
	"reflect"
	"testing"
	"time"
)

type SysInfo1 struct {
	SysDesc     string `oid:".1.3.6.1.2.1.1.1.0"`
	SysObjectId string `oid:".1.3.6.1.2.1.1.2.0"`
	SysUpTime   uint64 `oid:".1.3.6.1.2.1.1.3.0"`
	SysContact  string `oid:".1.3.6.1.2.1.1.4.0"`
	SysName     string `oid:".1.3.6.1.2.1.1.5.0"`
	Intfs       SysIntfs
}

type SysInfo2 struct {
	SysDesc     string `oid:".1.3.6.1.2.1.1.1.0"`
	SysObjectId string `oid:".1.3.6.1.2.1.1.2.0"`
	SysUpTime   uint64 `oid:".1.3.6.1.2.1.1.3.0"`
	SysContact  string `oid:".1.3.6.1.2.1.1.4.0"`
	SysName     string `oid:".1.3.6.1.2.1.1.5.0"`
	Intfs       *SysIntfs
}

type SysIntfs struct {
	IfDesc       map[string]string `oidx:"\\.1\\.3\\.6\\.1\\.2\\.1\\.2\\.2\\.1\\.2\\.(\\d+)"`
	IfOperStatus map[string]int    `oidx:\.1\.3\.6\.1\.2\.1\.2\.2\.1\.8\.(\d+)`
}

func TestMarshalPDUToStruct1(t *testing.T) {
	var (
		err error
	)
	params := &snmp.GoSNMP{
		Target:    "192.168.91.1",
		Port:      161,
		Community: "public",
		Version:   snmp.Version2c,
		Timeout:   time.Duration(5) * time.Second,
		Retries:   3,
		Logger:    nil,
	}
	err = params.Connect()
	if err != nil {
		t.Errorf("Connect() err: %v", err)
	}
	defer params.Conn.Close()

	info := SysInfo2{}
	err = params.BulkWalk(".1.3.6.1.2.1.2.2.1",
		func(pdu snmp.SnmpPDU) error {
			MarshalPDUToStruct(pdu, &info)
			return nil
		})
	if err != nil {
		t.Errorf("bulkwalk failure %s", err)
	}
	err = params.BulkWalk(".1.3.6.1.2.1.1",
		func(pdu snmp.SnmpPDU) error {
			MarshalPDUToStruct(pdu, &info)
			return nil
		})
	if err != nil {
		t.Errorf("bulkwalk failure %s", err)
	}
	spew.Dump(info)
}

func TestMarshalPDUToStruct2(t *testing.T) {
	var (
		err error
	)
	params := &snmp.GoSNMP{
		Target:    "192.168.91.1",
		Port:      161,
		Community: "public",
		Version:   snmp.Version2c,
		Timeout:   time.Duration(5) * time.Second,
		Retries:   3,
		Logger:    nil,
	}
	err = params.Connect()
	if err != nil {
		t.Errorf("Connect() err: %v", err)
	}
	defer params.Conn.Close()

	info := SysInfo2{
		Intfs: new(SysIntfs),
	}
	err = params.BulkWalk(".1.3.6.1.2.1.2.2.1",
		func(pdu snmp.SnmpPDU) error {
			MarshalPDUToStruct(pdu, &info)
			return nil
		})
	if err != nil {
		t.Errorf("bulkwalk failure %s", err)
	}
	err = params.BulkWalk(".1.3.6.1.2.1.1",
		func(pdu snmp.SnmpPDU) error {
			MarshalPDUToStruct(pdu, &info)
			return nil
		})
	if err != nil {
		t.Errorf("bulkwalk failure %s", err)
	}
	spew.Dump(info)
}

type Test1 struct {
	SysDesc     string `oid:".1.3.6.1.2.1.1.1.0"`
	SysObjectId string `oid:".1.3.6.1.2.1.1.2.0"`
	SysUpTime   uint64 `oid:".1.3.6.1.2.1.1.3.0"`
	SysContact  string `oid:".1.3.6.1.2.1.1.4.0"`
	SysName     string `oid:".1.3.6.1.2.1.1.5.0"`
	Nested      struct {
		InterfaceCount int `oid:".1.3.6.1.2.1.2.1.0"`
	}
}

type Test2 struct {
	SysDesc   string `oid:".1.3.6.1.2.1.1.1.0"`
	OtherData string
}

type Test3 struct {
	SysDesc   string `json:"foo"`
	OtherData string
	Bar       int `json:"bar"`
}

type Test4a struct {
	InterfaceCount int `oid:".1.3.6.1.2.1.2.1.0"`
}

type Test4 struct {
	SysDesc     string `oid:".1.3.6.1.2.1.1.1.0"`
	SysObjectId string `oid:".1.3.6.1.2.1.1.2.0"`
	Nested      *Test4a
}

func TestGetOidsFromStructTags1(t *testing.T) {
	type args struct {
		source    interface{}
		getNested bool
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "Multiple", args: args{Test1{}, false}, want: []string{
			".1.3.6.1.2.1.1.1.0", ".1.3.6.1.2.1.1.2.0", ".1.3.6.1.2.1.1.3.0", ".1.3.6.1.2.1.1.4.0", ".1.3.6.1.2.1.1.5.0",
		}},
		{name: "Single", args: args{Test2{}, false}, want: []string{".1.3.6.1.2.1.1.1.0"}},
		{name: "None", args: args{source: Test3{}}, want: []string{}},
		{name: "Pointer to struct", args: args{source: &Test2{}}, want: []string{".1.3.6.1.2.1.1.1.0"}},
		{name: "Single-Nested", args: args{Test2{}, true}, want: []string{".1.3.6.1.2.1.1.1.0"}},
		{name: "Multiple-Nested", args: args{Test1{}, true}, want: []string{
			".1.3.6.1.2.1.1.1.0", ".1.3.6.1.2.1.1.2.0", ".1.3.6.1.2.1.1.3.0", ".1.3.6.1.2.1.1.4.0", ".1.3.6.1.2.1.1.5.0",
			".1.3.6.1.2.1.2.1.0",
		}},
		{name: "Multiple-NestedPtr", args: args{Test4{Nested: new(Test4a)}, true}, want: []string{
			".1.3.6.1.2.1.1.1.0", ".1.3.6.1.2.1.1.2.0", ".1.3.6.1.2.1.2.1.0",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetOidsFromStructTags(tt.args.source, tt.args.getNested); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOidsFromStructTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetOidsFromStructTags2(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("code did not panic")
		}
	}()
	var x int
	GetOidsFromStructTags(&x, false)
}

func TestGetOidsFromStructTags3(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("code did not panic")
		}
	}()
	var x int
	GetOidsFromStructTags(x, false)
}
