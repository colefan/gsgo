package bstest
import (
	"github.com/colefan/gsgo/orm"
)
type Row_user_pt_login struct {
	orm.OrmModule
	pt_code	string	//平台编码
	account	string	//账户
	createtime	int32	//创建时间
	visittime	int32	//访问时间
	param1	string	//访问参数
	param2	string	//访问参数
	param3	string	//访问参数
	status	int8	//账户状态
}
func (row *Row_user_pt_login) Set_pt_code(v string) {
	row.pt_code = v
	row.SaveValue("pt_code", v)
}

func (row *Row_user_pt_login) Get_pt_code() string {
	return	row.pt_code
}

func (row *Row_user_pt_login) Set_account(v string) {
	row.account = v
	row.SaveValue("account", v)
}

func (row *Row_user_pt_login) Get_account() string {
	return	row.account
}

func (row *Row_user_pt_login) Set_createtime(v int32) {
	row.createtime = v
	row.SaveValue("createtime", v)
}

func (row *Row_user_pt_login) Get_createtime() int32 {
	return	row.createtime
}

func (row *Row_user_pt_login) Set_visittime(v int32) {
	row.visittime = v
	row.SaveValue("visittime", v)
}

func (row *Row_user_pt_login) Get_visittime() int32 {
	return	row.visittime
}

func (row *Row_user_pt_login) Set_param1(v string) {
	row.param1 = v
	row.SaveValue("param1", v)
}

func (row *Row_user_pt_login) Get_param1() string {
	return	row.param1
}

func (row *Row_user_pt_login) Set_param2(v string) {
	row.param2 = v
	row.SaveValue("param2", v)
}

func (row *Row_user_pt_login) Get_param2() string {
	return	row.param2
}

func (row *Row_user_pt_login) Set_param3(v string) {
	row.param3 = v
	row.SaveValue("param3", v)
}

func (row *Row_user_pt_login) Get_param3() string {
	return	row.param3
}

func (row *Row_user_pt_login) Set_status(v int8) {
	row.status = v
	row.SaveValue("status", v)
}

func (row *Row_user_pt_login) Get_status() int8 {
	return	row.status
}

func (row *Row_user_pt_login) Equals(b orm.OrmModuleInf) bool {
	brow, ok := b.(*Row_user_pt_login)
	if !ok {
		return false
	}
	return row.Get_pt_code() == brow.Get_pt_code() && row.Get_account() == brow.Get_account()
}

//账户平台登录表
type Table_user_pt_login struct {
	orm.OrmTable
}
func (this *Table_user_pt_login) Init() {
	this.OrmTableDef = orm.NewOrmTableDef()
	this.TableName = "user_pt_login"
	this.FieldKeys = append(this.FieldKeys, "pt_code")
	this.Fields["pt_code"] = orm.NewOrmField("pt_code", "string")
	this.PrimaryKey = append(this.PrimaryKey, this.Fields["pt_code"])
	this.FieldKeys = append(this.FieldKeys, "account")
	this.Fields["account"] = orm.NewOrmField("account", "string")
	this.PrimaryKey = append(this.PrimaryKey, this.Fields["account"])
	this.FieldKeys = append(this.FieldKeys, "createtime")
	this.Fields["createtime"] = orm.NewOrmField("createtime", "int32")
	this.FieldKeys = append(this.FieldKeys, "visittime")
	this.Fields["visittime"] = orm.NewOrmField("visittime", "int32")
	this.FieldKeys = append(this.FieldKeys, "param1")
	this.Fields["param1"] = orm.NewOrmField("param1", "string")
	this.FieldKeys = append(this.FieldKeys, "param2")
	this.Fields["param2"] = orm.NewOrmField("param2", "string")
	this.FieldKeys = append(this.FieldKeys, "param3")
	this.Fields["param3"] = orm.NewOrmField("param3", "string")
	this.FieldKeys = append(this.FieldKeys, "status")
	this.Fields["status"] = orm.NewOrmField("status", "int8")
}

func (this *Table_user_pt_login) NewRow() *Row_user_pt_login {
	r := &Row_user_pt_login{}
	r.TableDef = this.OrmTableDef
	return r
}
