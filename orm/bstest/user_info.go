package bstest
import (
	"github.com/colefan/gsgo/orm"
)
type Row_user_info struct {
	orm.OrmModule
	user_id	int64	//用户ID
	account	string	//账户名称
	pwd	string	//密码
	addtime	int32	//账号添加时间
	ip	string	//登录IP
	pt_code	string	//平台编码
	role_id	int64	//角色ID
	lasttime	int32	//最后登录时间
}
func (row *Row_user_info) Set_user_id(v int64) {
	row.user_id = v
	row.SaveValue("user_id", v)
}

func (row *Row_user_info) Get_user_id() int64 {
	return	row.user_id
}

func (row *Row_user_info) Set_account(v string) {
	row.account = v
	row.SaveValue("account", v)
}

func (row *Row_user_info) Get_account() string {
	return	row.account
}

func (row *Row_user_info) Set_pwd(v string) {
	row.pwd = v
	row.SaveValue("pwd", v)
}

func (row *Row_user_info) Get_pwd() string {
	return	row.pwd
}

func (row *Row_user_info) Set_addtime(v int32) {
	row.addtime = v
	row.SaveValue("addtime", v)
}

func (row *Row_user_info) Get_addtime() int32 {
	return	row.addtime
}

func (row *Row_user_info) Set_ip(v string) {
	row.ip = v
	row.SaveValue("ip", v)
}

func (row *Row_user_info) Get_ip() string {
	return	row.ip
}

func (row *Row_user_info) Set_pt_code(v string) {
	row.pt_code = v
	row.SaveValue("pt_code", v)
}

func (row *Row_user_info) Get_pt_code() string {
	return	row.pt_code
}

func (row *Row_user_info) Set_role_id(v int64) {
	row.role_id = v
	row.SaveValue("role_id", v)
}

func (row *Row_user_info) Get_role_id() int64 {
	return	row.role_id
}

func (row *Row_user_info) Set_lasttime(v int32) {
	row.lasttime = v
	row.SaveValue("lasttime", v)
}

func (row *Row_user_info) Get_lasttime() int32 {
	return	row.lasttime
}

func (row *Row_user_info) Equals(b orm.OrmModuleInf) bool {
	brow, ok := b.(*Row_user_info)
	if !ok {
		return false
	}
	return row.Get_user_id() == brow.Get_user_id()
}

//
type Table_user_info struct {
	orm.OrmTable
}
func (this *Table_user_info) Init() {
	this.OrmTableDef = orm.NewOrmTableDef()
	this.TableName = "user_info"
	this.FieldKeys = append(this.FieldKeys, "user_id")
	this.Fields["user_id"] = orm.NewOrmField("user_id", "int64")
	this.PrimaryKey = append(this.PrimaryKey, this.Fields["user_id"])
	this.FieldKeys = append(this.FieldKeys, "account")
	this.Fields["account"] = orm.NewOrmField("account", "string")
	this.FieldKeys = append(this.FieldKeys, "pwd")
	this.Fields["pwd"] = orm.NewOrmField("pwd", "string")
	this.FieldKeys = append(this.FieldKeys, "addtime")
	this.Fields["addtime"] = orm.NewOrmField("addtime", "int32")
	this.FieldKeys = append(this.FieldKeys, "ip")
	this.Fields["ip"] = orm.NewOrmField("ip", "string")
	this.FieldKeys = append(this.FieldKeys, "pt_code")
	this.Fields["pt_code"] = orm.NewOrmField("pt_code", "string")
	this.FieldKeys = append(this.FieldKeys, "role_id")
	this.Fields["role_id"] = orm.NewOrmField("role_id", "int64")
	this.FieldKeys = append(this.FieldKeys, "lasttime")
	this.Fields["lasttime"] = orm.NewOrmField("lasttime", "int32")
}

func (this *Table_user_info) NewRow() *Row_user_info {
	r := &Row_user_info{}
	r.TableDef = this.OrmTableDef
	return r
}
