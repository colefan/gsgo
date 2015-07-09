package db

import (
	"fmt"
	"testing"
)

func TestDbPool(t *testing.T) {
	dbpool := NewDefaultDbPool()
	err := dbpool.InitPool("test", "mysql", "root:@tcp(192.168.16.58:3306)/test", 5)
	if err != nil {
		t.Fatal(err)
	}
	conn := dbpool.GetConnection()

	record, err := conn.Read("select id,name  from test01")
	if err != nil {
		t.Fatal(err)
	}
	for record.Next() {
		var id int
		var name string
		record.Scan(&id, &name)
		fmt.Println(id, name)

	}
	record.Close()
	err = conn.Delete("delete from test01 where id = 4")
	if err != nil {
		t.Fatal(err)
	}

	err = conn.Create("insert into test01(id,name) values(4,\"test04\")")
	if err != nil {
		t.Fatal(err)
	}

	err = conn.Update("update test01 set name=\"test04-up\" where id = 4")
	if err != nil {
		t.Fatal(err)
	}

	dbpool.ReleaseConnection(conn)
	dbpool.DestroyPool()
}
