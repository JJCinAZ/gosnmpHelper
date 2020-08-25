# Library to assist with GoSNMP operations

## Value assignment utilities

This set of functions can be used to safely pull values from PDUs
in the format you desire.  Here are some examples:

---
	gosnmp.Default.Target = "192.168.1.1"
	err := gosnmp.Default.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer gosnmp.Default.Conn.Close()

	oids := []string{".1.3.6.1.2.1.2.1.0"}
	result, err := gosnmp.Default.Get(oids)
    if err != nil {
        fmt.Println(gosnmpHelper.GetAsUint(result.Variables[0]))
    }
---

## Struct Tag Helpers

It's often tedious to parse returned PDUs into struct member values.  The struct tags
in Go provide a way to streamline this somewhat.  Let's use the following struct:

---
    type BasicInfo struct {
        ID          int
        IPAddress   string 
        SysDesc     string `oid:".1.3.6.1.2.1.1.1.0"`
        SysObjectId string `oid:".1.3.6.1.2.1.1.2.0"`
        SysUpTime   uint64 `oid:".1.3.6.1.2.1.1.3.0"`
        SysContact  string `oid:".1.3.6.1.2.1.1.4.0"`
        SysName     string `oid:".1.3.6.1.2.1.1.5.0"`
    }
---

This structure has tags indicating which OIDs are used to get the values.
If we want to do an SNMP Get, we only need to use the GetOidsFromStructTags function as follows:

---
    info := BasicInfo{}
    result, err = gosnmp.Default.Get(gosnmpHelper.GetOidsFromStructTags(&info))
---

To extract the values, you can use the MarshalPDUsToStruct function.  It utilizes the struct tags
to transcribe the returned SNMP values into the struct member fields. 

---
    gosnmpHelper.MarshalPDUsToStruct(result.Variables, &info)
---
