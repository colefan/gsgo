package orm

import (
	"fmt"
	"strconv"
	"strings"
)

type OrmModuleInf interface {
	Equals(row OrmModuleInf) bool
	SaveValue(filedname string, v interface{})
	DeleteSql(cond *QueryCondition) string
	InsertSql(excludeFields []string) string
	UpdateSql(updateFields []string, cond *QueryCondition) string
	GetTableName() string
	GetField(fieldname string) *OrmField
}

type OrmModule struct {
	TableDef    *OrmTableDef
	FieldValMap map[string]interface{}
}

func (this *OrmModule) SaveValue(fieldname string, v interface{}) {
	if this.FieldValMap == nil {
		this.FieldValMap = make(map[string]interface{})
	}
	this.FieldValMap[strings.ToLower(fieldname)] = v
}

func (this *OrmModule) DeleteSql(cond *QueryCondition) string {
	strsql := "delete from " + this.TableDef.TableName
	if cond == nil {
		strsql += " where " + this.getPrimaryKeyCond()
	} else {
		strsql += " where " + cond.condtion
	}
	return strsql
}

func (this *OrmModule) InsertSql(excludeFields []string) string {
	sql := "insert into " + this.TableDef.TableName + " ("
	value := " values ("

	for _, fieldname := range this.TableDef.FieldKeys {
		jumpbreak := false
		if excludeFields != nil {
			for _, excludef := range excludeFields {
				if strings.ToLower(excludef) == fieldname {
					jumpbreak = true
					break
				}
			}

		}

		if jumpbreak {
			break
		}

		field := this.TableDef.Fields[fieldname]
		sql += fieldname + ","
		switch field.FieldType {
		case "int8":
			v, _ := this.FieldValMap[fieldname].(int8)
			value += strconv.FormatInt(int64(v), 10) + ","
		case "uint8":
			v, _ := this.FieldValMap[fieldname].(uint8)
			value += strconv.FormatUint(uint64(v), 10) + ","
		case "int16":
			v, _ := this.FieldValMap[fieldname].(int16)
			value += strconv.FormatInt(int64(v), 10) + ","
		case "uint16":
			v, _ := this.FieldValMap[fieldname].(uint16)
			value += strconv.FormatUint(uint64(v), 10) + ","
		case "int32":
			v, _ := this.FieldValMap[fieldname].(int32)
			value += strconv.FormatInt(int64(v), 10) + ","
		case "uint32":
			v, _ := this.FieldValMap[fieldname].(uint32)
			value += strconv.FormatUint(uint64(v), 10) + ","
		case "int64":
			v, _ := this.FieldValMap[fieldname].(int64)
			value += strconv.FormatInt(v, 10) + ","
		case "uint64":
			v, _ := this.FieldValMap[fieldname].(uint64)
			value += strconv.FormatUint(v, 10) + ","

		case "string":
			v, _ := this.FieldValMap[fieldname].(string)
			value += "'" + v + "',"
		default:
			//value += "'" + string(this.FieldValMap[fieldname]) + "',"
		}

	}

	value = value[0 : len(value)-1]
	value += ") "
	sql = sql[0 : len(sql)-1]

	sql += ") " + value
	return sql
}

func (this *OrmModule) getPrimaryKeyCond() string {
	cond := " "
	for _, pk := range this.TableDef.PrimaryKey {
		cond += pk.Name + " = "
		switch pk.FieldType {
		case "int64":
			v, _ := this.FieldValMap[pk.Name].(int64)
			cond += strconv.FormatInt(v, 10)
		case "uint64":
			v, _ := this.FieldValMap[pk.Name].(uint64)
			cond += strconv.FormatUint(v, 10)
		case "string":
			v, _ := this.FieldValMap[pk.Name].(string)
			cond += "'" + v + "'"
		case "int32":
			v, _ := this.FieldValMap[pk.Name].(int32)
			cond += strconv.FormatInt(int64(v), 10)
		case "uint32":
			v, _ := this.FieldValMap[pk.Name].(uint32)
			cond += strconv.FormatUint(uint64(v), 10)
		case "int16":
			v, _ := this.FieldValMap[pk.Name].(int16)
			cond += strconv.FormatInt(int64(v), 10)
		case "uint16":
			v, _ := this.FieldValMap[pk.Name].(uint16)
			cond += strconv.FormatUint(uint64(v), 10)
		case "int8":
			v, _ := this.FieldValMap[pk.Name].(int8)
			cond += strconv.FormatInt(int64(v), 10)
		case "uint8":
			v, _ := this.FieldValMap[pk.Name].(uint8)
			cond += strconv.FormatUint(uint64(v), 10)

		}
		cond += " and "
	}
	if len(cond) > 5 {
		cond = cond[0 : len(cond)-5]
	}

	return cond

}

func (this *OrmModule) UpdateSql(updateFields []string, cond *QueryCondition) string {
	strsql := "update " + this.TableDef.TableName + " set "
	if updateFields != nil {
		for _, fn := range updateFields {
			fn = strings.ToLower(fn)
			strsql += fn + " = "
			f := this.TableDef.Fields[fn]
			switch f.FieldType {
			case "int8":
				v, _ := this.FieldValMap[f.Name].(int8)
				strsql += strconv.FormatInt(int64(v), 10) + ","
			case "uint8":
				v, _ := this.FieldValMap[f.Name].(uint8)
				strsql += strconv.FormatUint(uint64(v), 10) + ","
			case "int16":
				v, _ := this.FieldValMap[f.Name].(int16)
				strsql += strconv.FormatInt(int64(v), 10) + ","
			case "uint16":
				v, _ := this.FieldValMap[f.Name].(uint16)
				strsql += strconv.FormatUint(uint64(v), 10) + ","
			case "int32":
				v, _ := this.FieldValMap[f.Name].(int32)
				strsql += strconv.FormatInt(int64(v), 10) + ","
			case "uint32":
				v, _ := this.FieldValMap[f.Name].(uint32)
				strsql += strconv.FormatUint(uint64(v), 10) + ","
			case "int64":
				v, _ := this.FieldValMap[f.Name].(int64)
				strsql += strconv.FormatInt(v, 10) + ","
			case "uint64":
				v, _ := this.FieldValMap[f.Name].(uint64)
				strsql += strconv.FormatUint(v, 10) + ","
			case "string":
				v, _ := this.FieldValMap[f.Name].(string)
				strsql += "'" + v + "',"
			default:
				fmt.Println("unkown fieldtype:", this.FieldValMap[f.Name])

			}
		}

	} else {
		for _, f := range this.TableDef.Fields {
			strsql += f.Name + " = "
			switch f.FieldType {
			case "int8":
				v, _ := this.FieldValMap[f.Name].(int8)
				strsql += strconv.FormatInt(int64(v), 10) + ","
			case "uint8":
				v, _ := this.FieldValMap[f.Name].(uint8)
				strsql += strconv.FormatUint(uint64(v), 10) + ","
			case "int16":
				v, _ := this.FieldValMap[f.Name].(int16)
				strsql += strconv.FormatInt(int64(v), 10) + ","
			case "uint16":
				v, _ := this.FieldValMap[f.Name].(uint16)
				strsql += strconv.FormatUint(uint64(v), 10) + ","
			case "int32":
				v, _ := this.FieldValMap[f.Name].(int32)
				strsql += strconv.FormatInt(int64(v), 10) + ","
			case "uint32":
				v, _ := this.FieldValMap[f.Name].(uint32)
				strsql += strconv.FormatUint(uint64(v), 10) + ","
			case "int64":
				v, _ := this.FieldValMap[f.Name].(int64)
				strsql += strconv.FormatInt(v, 10) + ","
			case "uint64":
				v, _ := this.FieldValMap[f.Name].(uint64)
				strsql += strconv.FormatUint(v, 10) + ","
			case "string":
				v, _ := this.FieldValMap[f.Name].(string)
				strsql += "'" + v + "',"
			default:
				fmt.Println("unkown fieldtype:", this.FieldValMap[f.Name])

			}
		}
	}

	strsql = strsql[0 : len(strsql)-1]

	if cond != nil {
		strsql += " where " + cond.condtion
	} else {
		strsql += " where " + this.getPrimaryKeyCond()
	}
	return strsql
}

func (this *OrmModule) GetTableName() string {
	return this.TableDef.TableName
}

func (this *OrmModule) GetField(fieldname string) *OrmField {
	return this.TableDef.Fields[strings.ToLower(fieldname)]

}

//次方法需要被复写
func (this *OrmModule) Equals(row OrmModuleInf) bool {
	//检查是否主键相同
	return false
}
