package bstest
import (
	"github.com/colefan/gsgo/orm"
)
type Row_role_info struct {
	orm.OrmModule
	role_id	uint64	//
	user_id	int64	//
}
func (row *Row_role_info) Set_role_id(v uint64) {
	row.role_id = v
	row.SaveValue("role_id", v)
}

func (row *Row_role_info) Get_role_id() uint64 {
	return	row.role_id
}

func (row *Row_role_info) Set_user_id(v int64) {
	row.user_id = v
	row.SaveValue("user_id", v)
}

func (row *Row_role_info) Get_user_id() int64 {
	return	row.user_id
}

func (row *Row_role_info) Equals(b orm.OrmModuleInf) bool {
	brow, ok := b.(*Row_role_info)
	if !ok {
		return false
	}
	return row.Get_role_id() == brow.Get_role_id()
}

//角色表
type Table_role_info struct {
	orm.OrmTable
}
func (this *Table_role_info) Init() {
	this.OrmTableDef = orm.NewOrmTableDef()
	this.TableName = "role_info"
	this.FieldKeys = append(this.FieldKeys, "role_id")
	this.Fields["role_id"] = orm.NewOrmField("role_id", "uint64")
	this.PrimaryKey = append(this.PrimaryKey, this.Fields["role_id"])
	this.FieldKeys = append(this.FieldKeys, "user_id")
	this.Fields["user_id"] = orm.NewOrmField("user_id", "int64")
}

func (this *Table_role_info) NewRow() *Row_role_info {
	r := &Row_role_info{}
	r.TableDef = this.OrmTableDef
	return r
}
