package main

import (
	"fmt"
	"github.com/jjcinaz/gosnmpHelper"
	"github.com/soniah/gosnmp"
	"log"
)

type BasicInfo struct {
	SysDesc     string `oid:".1.3.6.1.2.1.1.1.0"`
	SysObjectId string `oid:".1.3.6.1.2.1.1.2.0"`
	SysUpTime   uint64 `oid:".1.3.6.1.2.1.1.3.0"`
	SysContact  string `oid:".1.3.6.1.2.1.1.4.0"`
	SysName     string `oid:".1.3.6.1.2.1.1.5.0"`
}

func main() {
	var (
		info   BasicInfo
		result *gosnmp.SnmpPacket
	)
	gosnmp.Default.Target = "192.168.91.1"
	err := gosnmp.Default.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer gosnmp.Default.Conn.Close()
	result, err = gosnmp.Default.Get(gosnmpHelper.GetOidsFromStructTags(&info))
	if err != nil {
		log.Fatalf("snmp failure %s", err)
	}
	gosnmpHelper.MarshalPDUsToStruct(result.Variables, &info)
	fmt.Printf("%#v\n", info)
}
