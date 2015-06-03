package config

import (
	"os"
	"testing"
)

var inittestconfig = `

#comment two
appname = test01
listenport = 10002
mysqlport = 3600
PI = 3.1415976
runmode = "dev"
autorender = false
copyrequestbody = true
[demo]
key1="yjx"
key2 = "colefan"
CaseInsensitive = true
notelist = one;two;three`

func TestIniConfig(t *testing.T) {
	file, err := os.Create("testini.conf")
	if err != nil {
		t.Fatal("create file error")
	}
	_, err = file.WriteString(inittestconfig)
	if err != nil {
		file.Close()
		t.Fatal(err)
	}
	file.Close()
	defer os.Remove("testini.conf")
	conf, err := NewConfig("ini", "testini.conf")
	if err != nil {
		t.Fatal(err)
	}

	if conf.String("appname") != "test01" {
		t.Fatal("appname read error")
	}

	if p, err := conf.Int("listenport"); err != nil || p != 10002 {
		t.Error(p)
		t.Fatal(err)
	}

	if conf.String("demo::key1") != "yjx" {
		t.Fatal("read demo::key1 error")
	}

}
