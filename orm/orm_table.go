package orm

import (
	"strings"
)

//描述了一个表应该具有的基本属性和基本方法

type OrmTableDef struct {
	TableName  string
	PrimaryKey []*OrmField
	Fields     map[string]*OrmField
	FieldKeys  []string
}

func NewOrmTableDef() *OrmTableDef {
	def := &OrmTableDef{}
	def.PrimaryKey = make([]*OrmField, 0)
	def.Fields = make(map[string]*OrmField)
	def.FieldKeys = make([]string, 0)
	return def
}

type OrmTable struct {
	*OrmTableDef
	Rows []OrmModuleInf
}

func (this *OrmTable) Init() {

}

func (this *OrmTable) LoadData(condstr string) {

}

func (this *OrmTable) GetPrimaryKey() []*OrmField {
	return this.PrimaryKey
}

func (this *OrmTable) GetFieldDef(fieldname string) *OrmField {
	fieldname = strings.ToLower(fieldname)

	f := this.Fields[fieldname]
	return f
}

func (this *OrmTable) GetTableName() string {
	return this.TableName
}

func (this *OrmTable) GetRowCount() int {
	return len(this.Rows)
}

func (this *OrmTable) AddRow(row OrmModuleInf) bool {
	if row == nil {
		return false
	}
	for _, myrow := range this.Rows {
		if myrow.Equals(row) {
			return false
		}
	}
	this.Rows = append(this.Rows, row)
	return true
}

func (this *OrmTable) FindRow(row OrmModuleInf) (bool, OrmModuleInf) {
	for _, myrow := range this.Rows {
		if myrow.Equals(row) {
			return true, myrow
		}
	}
	return false, nil
}

func (this *OrmTable) UpdateRow(row OrmModuleInf) bool {
	for i, myrow := range this.Rows {
		if myrow.Equals(row) {
			this.Rows[i] = row
			return true
		}
	}
	return false
}

func (this *OrmTable) DeleteRow(row OrmModuleInf) bool {
	for i, myrow := range this.Rows {
		if myrow.Equals(row) {
			if i == len(this.Rows)-1 {
				this.Rows = this.Rows[:i]
			} else {
				this.Rows = append(this.Rows[:i], this.Rows[i+1:]...)
			}
			return true
		}
	}
	return false
}
