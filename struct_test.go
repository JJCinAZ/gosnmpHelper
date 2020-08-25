package gosnmpHelper

import (
	"fmt"
	snmp "github.com/soniah/gosnmp"
	"reflect"
	"testing"
	"time"
)

func TestMarshalPDUToStruct(t *testing.T) {
	var (
		err  error
		info struct {
			SysDesc      string            `oid:".1.3.6.1.2.1.1.1.0"`
			SysObjectId  string            `oid:".1.3.6.1.2.1.1.2.0"`
			SysUpTime    uint64            `oid:".1.3.6.1.2.1.1.3.0"`
			SysContact   string            `oid:".1.3.6.1.2.1.1.4.0"`
			SysName      string            `oid:".1.3.6.1.2.1.1.5.0"`
			IfDesc       map[string]string `oidx:"\\.1\\.3\\.6\\.1\\.2\\.1\\.2\\.2\\.1\\.2\\.(\\d+)"`
			IfOperStatus map[string]int    `oidx:\.1\.3\.6\.1\.2\.1\.2\.2\.1\.8\.(\d+)`
		}
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
	fmt.Printf("%+v\n", info)
}

type Test1 struct {
	SysDesc     string `oid:".1.3.6.1.2.1.1.1.0"`
	SysObjectId string `oid:".1.3.6.1.2.1.1.2.0"`
	SysUpTime   uint64 `oid:".1.3.6.1.2.1.1.3.0"`
	SysContact  string `oid:".1.3.6.1.2.1.1.4.0"`
	SysName     string `oid:".1.3.6.1.2.1.1.5.0"`
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

func TestGetOidsFromStructTags1(t *testing.T) {
	type args struct {
		source interface{}
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "Multiple", args: args{source: Test1{}}, want: []string{
			".1.3.6.1.2.1.1.1.0", ".1.3.6.1.2.1.1.2.0", ".1.3.6.1.2.1.1.3.0", ".1.3.6.1.2.1.1.4.0", ".1.3.6.1.2.1.1.5.0",
		}},
		{name: "Single", args: args{source: Test2{}}, want: []string{".1.3.6.1.2.1.1.1.0"}},
		{name: "None", args: args{source: Test3{}}, want: []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetOidsFromStructTags(tt.args.source); !reflect.DeepEqual(got, tt.want) {
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
	GetOidsFromStructTags(&x)
}

func TestGetOidsFromStructTags3(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("code did not panic")
		}
	}()
	var x int
	GetOidsFromStructTags(x)
}
