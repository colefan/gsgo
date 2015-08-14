package orm

import (
	"fmt"
	"testing"

	"github.com/colefan/gsgo/db"
)

type UserInfo struct {
	OrmModule
	id       uint64
	username string
}

func (row *UserInfo) Init() {

}

func (row *UserInfo) SetId(id uint64) {
	row.id = id
	row.SaveValue("id", id)
}

func (row *UserInfo) GetId() uint64 {
	return row.id
}

func (row *UserInfo) GetUserName() string {
	return row.username
}

func (row *UserInfo) SetUserName(name string) {
	row.username = name
	row.SaveValue("username", name)
}

func (row *UserInfo) Equals(b OrmModuleInf) bool {
	brow, ok := b.(*UserInfo)
	if !ok {
		return false
	}
	return row.GetId() == brow.GetId()
}

type TableUserInfo struct {
	OrmTable
}

func (this *TableUserInfo) Init() {
	this.OrmTableDef = NewOrmTableDef()
	this.TableName = "user_info"
	this.PrimaryKey = append(this.PrimaryKey, NewOrmField("id", "uint64"))
	this.PrimaryKey = append(this.PrimaryKey, NewOrmField("username", "string"))

	this.FieldKeys = append(this.FieldKeys, "id")
	this.FieldKeys = append(this.FieldKeys, "username")
	this.Fields["id"] = NewOrmField("id", "uint64")
	this.Fields["username"] = NewOrmField("username", "string")
}

func (this *TableUserInfo) LoadData(strcond string) {

}

func (this *TableUserInfo) NewRow() *UserInfo {
	r := &UserInfo{}
	r.TableDef = this.OrmTableDef
	return r
}

func Test_userinfo(t *testing.T) {

	pool := db.NewDefaultDbPool()
	err := dbpool.InitPool("test", "mysql", "root:@tcp(192.168.16.58:3306)/test", 5)
	if err != nil {
		t.Fatal(err)
	}

	table := TableUserInfo{}
	table.Init()

	row1 := table.NewRow()
	row1.SetId(1)
	row1.SetUserName("yjx1")
	fmt.Println("add row1:", table.AddRow(row1))
	fmt.Println("repeat add row1:", table.AddRow(row1))
	fmt.Println("insert row1:", row1.InsertSql(nil))

	row2 := table.NewRow()
	row2.SetId(2)
	row2.SetUserName("yjx2")
	fmt.Println("add row2:", table.AddRow(row2))
	fmt.Println("insert row1:", row2.InsertSql(nil))

	row3 := table.NewRow()
	row3.SetId(1)
	b, row4 := table.FindRow(row3)
	fmt.Println("b,", b, "row3:", row4)
	row3.SetUserName("yjx21")
	fmt.Println("update", table.UpdateRow(row3))
	fmt.Println(table.OrmTable)
	fmt.Println("delete", table.DeleteRow(row2))

	fmt.Println(table.OrmTable)

}
