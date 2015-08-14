package orm

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/colefan/gsgo/db"
)

//orm 命令行工具
//Useage:set db = dbtype/dbname/dbuser:dbpwd@dbaddress:dbport
//Useage:orm -p packagename -table tablesneedtoorm [-path filestorepath]
//Useage:loadconf config.ini
//Useage:help
type OrmCmd struct {
	dbtype      string
	dbname      string
	dblink      string
	packagename string
	storepath   string
	ormtables   string
	dbpool      db.DbPool
}

const (
	SET  = "set "
	ORM  = "orm "
	CONF = "loadconf "
	HELP = "help"
)

type cmdline struct {
	word   string
	params string
	kv     map[string]string
}

type tableO struct {
	tablename string
	tabledesc string
	fields    []fieldO
}

func (t *tableO) docTable(packname string, storepath string) {
	filename := ""
	seppath := string(os.PathSeparator)
	if storepath == "" {
		filename = t.tablename + ".go"
	} else {
		if strings.HasSuffix(storepath, seppath) {
			filename = storepath + t.tablename + ".go"
		} else {
			filename = storepath + seppath + t.tablename + ".go"
		}
	}

	var strpackline string = "package " + packname + "\n"
	var strimportline string = "import (" + "\n"
	strimportline += "\t" + "\"github.com/colefan/gsgo/orm\"" + "\n"
	strimportline += ")\n"

	var strrow string = "//" + t.tabledesc + "\n"
	strrow = "type Row_" + t.tablename + " struct {" + "\n"
	strrow += "\t" + "orm.OrmModule" + "\n"
	for _, f := range t.fields {
		strrow += "\t" + f.col_name + "\t" + f.getFieldType() + "\t//" + f.col_comment + "\n"
	}
	strrow += "}\n"

	for _, f := range t.fields {
		strrow += "func (row *Row_" + t.tablename + ") Set_" + f.col_name + "(v " + f.getFieldType() + ") {\n"
		strrow += "\trow." + f.col_name + " = v" + "\n"
		strrow += "\t" + "row.SaveValue(\"" + f.col_name + "\", " + "v)\n"
		strrow += "}\n\n"
		strrow += "func (row *Row_" + t.tablename + ") Get_" + f.col_name + "() " + f.getFieldType() + " {\n"
		strrow += "\treturn\trow." + f.col_name + "\n"
		strrow += "}\n\n"
	}

	strrow += "func " + "(row *Row_" + t.tablename + ") " + "Equals(b orm.OrmModuleInf) bool {\n"
	strrow += "\t" + "brow, ok := b.(*Row_" + t.tablename + ")\n"
	strrow += "\tif !ok {\n"
	strrow += "\t\treturn false\n"
	strrow += "\t}\n"

	strequal := ""
	for _, kf := range t.fields {
		if kf.col_key == "PRI" {
			if strequal != "" {
				strequal += " && "
			}
			strequal += "row.Get_" + kf.col_name + "() == brow.Get_" + kf.col_name + "()"
		}
	}
	strrow += "\treturn " + strequal + "\n"
	strrow += "}\n\n"

	strtable := "//" + t.tabledesc + "\n"

	strtable += "type Table_" + t.tablename + " struct {\n"
	strtable += "\t" + "orm.OrmTable" + "\n"
	strtable += "}\n"
	strtable += "func (this *Table_" + t.tablename + ") Init() {\n"
	strtable += "\t" + "this.OrmTableDef = orm.NewOrmTableDef()" + "\n"
	strtable += "\t" + "this.TableName = \"" + t.tablename + "\"" + "\n"

	for _, f := range t.fields {
		strtable += "\t" + "this.FieldKeys = append(this.FieldKeys, \"" + f.col_name + "\")\n"
		strtable += "\t" + "this.Fields[\"" + f.col_name + "\"] = orm.NewOrmField(\"" + f.col_name + "\", \"" + f.getFieldType() + "\")\n"
		if f.col_key == "PRI" {
			strtable += "\t" + "this.PrimaryKey = append(this.PrimaryKey, this.Fields[\"" + f.col_name + "\"])\n"
		}
	}
	strtable += "}\n\n"

	strtable += "func (this *Table_" + t.tablename + ") NewRow() *Row_" + t.tablename + " {\n"
	strtable += "\t" + "r := &Row_" + t.tablename + "{}\n"
	strtable += "\t" + "r.TableDef = this.OrmTableDef\n"
	strtable += "\t" + "return r\n"
	strtable += "}\n"
	//	fmt.Println(filename)
	//	fmt.Print(strpackline)
	//	fmt.Print(strimportline)
	//	fmt.Print(strrow)
	//	fmt.Print(strtable)

	file, err := os.Create(filename)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		defer file.Close()
		file.WriteString(strpackline + strimportline + strrow + strtable)
	}

}

type fieldO struct {
	col_name     string
	col_datatype string
	col_type     string
	col_key      string
	col_comment  string
}

func (this *fieldO) getFieldType() string {
	ftype := ""
	if strings.HasSuffix(this.col_type, "unsigned") {
		ftype = "u"
	}

	switch this.col_datatype {
	case "bigint":
		ftype += "int64"
	case "varchar":
		ftype = "string"
	case "int":
		ftype += "int32"
	case "tinyint":
		ftype += "int8"
	case "smallint":
		ftype += "int16"
	}
	return ftype
}

type rowO struct {
}

func newCmdLine() *cmdline {
	c := &cmdline{}
	c.kv = make(map[string]string)
	return c
}

func (this *OrmCmd) Run(storepath string) {
	this.storepath = storepath
	fmt.Println("请输入命令执行ORM操作")
	reader := bufio.NewReader(os.Stdin)
	for {

		data, _, _ := reader.ReadLine()
		command := string(data)
		if command == "exit" {
			break
		}
		fmt.Println("cmd = ", command)
		result, cmd := this.parsecmd(command)
		if cmd == nil {
			fmt.Println(result)
		} else {
			//执行命令
			this.executecmd(cmd)
		}

	}

	if this.dbpool != nil {
		this.dbpool.DestroyPool()
	}

}

func (this *OrmCmd) executecmd(cmd *cmdline) {
	switch cmd.word {
	case SET:
		if this.dbpool == nil {
			this.dbpool = db.NewDefaultDbPool()
		}
		strconn := this.dblink
		if this.dbtype == "mysql" {
			strconn += "/information_schema"
		}
		err := this.dbpool.InitPool(this.dbname, this.dbtype, strconn, 1)
		if err != nil {
			fmt.Println("数据库连接失败，请重新设置数据库信息")
		} else {
			fmt.Println("亲，数据库已经连上了")
		}

	case HELP:
		fmt.Println("Useage:set db = dbtype/dbname/dbuser:dbpwd@dbaddress:dbport")
		fmt.Println("Useage:orm -p packagename -table tablesneedtoorm [-path filestorepath]")
		fmt.Println("Useage:loadconf config.ini")
		fmt.Println("Useage:help")
	case ORM:
		if this.dblink == "" {
			fmt.Println("数据库尚未设置，无法做ORM操作，请先设置数据库")
		} else {
			this.packagename = cmd.kv["-p"]
			if cmd.kv["-path"] != "" {
				this.storepath = cmd.kv["-path"]
			}
			this.ormtables = cmd.kv["-table"]
			this.orm(this.ormtables)

		}
	case CONF:
		configfile := cmd.kv["conf"]
		conf := OrmConfig{}
		err := conf.LoadConfig(configfile)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			this.dbtype = conf.dbtype
			this.dbname = conf.dbname
			this.dblink = conf.dblink

			if this.dbpool == nil {
				this.dbpool = db.NewDefaultDbPool()
			}
			strconn := this.dblink
			if this.dbtype == "mysql" {
				strconn += "/information_schema"
			}
			err := this.dbpool.InitPool(this.dbname, this.dbtype, strconn, 1)
			if err != nil {
				fmt.Println("数据库连接失败，请重新设置数据库信息")
			} else {
				fmt.Println("亲，数据库已经连上了")
			}

		}

	}

}

func (this *OrmCmd) orm(tables string) {
	//fmt.Println(tables)
	tableList := make([]*tableO, 0)
	if strings.ToLower(tables) == "all" {
		tableList = this.getAllTables()
	} else {
		ts := strings.Split(tables, ";")
		for _, tmp := range ts {
			t := &tableO{}
			t.tablename = tmp
			tableList = append(tableList, t)
		}
	}

	for _, t := range tableList {
		fmt.Println("开始生成:", t.tablename)
		if this.getFields(t) {
			t.docTable(this.packagename, this.storepath)
			fmt.Println("生成成功:", t.tablename)
		} else {
			fmt.Println("生成失败:（", t.tablename)
		}
	}

}

func (this *OrmCmd) getFields(table *tableO) bool {
	sql := "select column_name,data_type,column_type,column_comment,column_key from columns where table_schema='" + this.dbname + "' and table_name='" + table.tablename + "'"
	conn := this.dbpool.GetConnection()
	defer this.dbpool.ReleaseConnection(conn)
	rows, err := conn.Read(sql)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	for rows.Next() {
		f := fieldO{}
		err := rows.Scan(&f.col_name, &f.col_datatype, &f.col_type, &f.col_comment, &f.col_key)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
		//fmt.Println(f)
		if table.fields == nil {
			table.fields = make([]fieldO, 0)
		}
		table.fields = append(table.fields, f)
	}

	return true
}

func (this *OrmCmd) getAllTables() []*tableO {
	sql := "select table_name,table_comment from tables where table_schema='" + this.dbname + "' and table_type='BASE TABLE' "
	conn := this.dbpool.GetConnection()
	//fmt.Println(sql)
	defer this.dbpool.ReleaseConnection(conn)
	rows, err := conn.Read(sql)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	tableList := make([]*tableO, 0)
	for rows.Next() {
		tablename := ""
		tabledesc := ""
		err := rows.Scan(&tablename, &tabledesc)
		if err != nil {
			fmt.Println(err.Error())
			return tableList
		}
		t := &tableO{}
		t.tablename = tablename
		t.tabledesc = tabledesc
		//fmt.Println("name = ", tablename, ",desc = ", tabledesc)
		tableList = append(tableList, t)

	}
	return tableList

}

func (this *OrmCmd) parsecmd(command string) (string, *cmdline) {
	if strings.HasPrefix(command, SET) {
		cmd := newCmdLine()
		cmd.word = SET
		cmd.params = command[len(SET):]
		cmd.params = strings.TrimSpace(cmd.params)
		kvs := strings.Split(cmd.params, "=")
		if len(kvs) != 2 {
			return "set命令格式错，请用help查看命令帮助", nil
		}
		if strings.TrimSpace(kvs[0]) == "db" {
		} else {
			return "set命令格式错，请用help查看命令帮助", nil
		}

		kvs = strings.Split(kvs[1], "/")
		if len(kvs) != 3 {
			return "set命令格式错，请用help查看命令帮助", nil
		}
		cmd.kv["dbtype"] = strings.TrimSpace(kvs[0])
		cmd.kv["dbname"] = strings.TrimSpace(kvs[1])
		cmd.kv["dblink"] = strings.TrimSpace(kvs[2])
		this.dbtype = cmd.kv["dbtype"]
		this.dbname = cmd.kv["dbname"]
		this.dblink = cmd.kv["dblink"]

		return "", cmd
	} else if strings.HasPrefix(command, CONF) {
		cmd := newCmdLine()
		cmd.word = CONF
		cmd.params = command[len(CONF):]
		cmd.params = strings.TrimSpace(cmd.params)
		cmd.kv["conf"] = cmd.params
		return "", cmd
	} else if strings.HasPrefix(command, ORM) {
		cmd := newCmdLine()
		cmd.word = ORM
		cmd.params = command[len(ORM):]
		cmd.params = strings.TrimSpace(cmd.params)
		//-p packagename -table tablesneedtoorm [-path filestorepath]
		kvs := strings.Split(cmd.params, " ")
		if len(kvs)%2 != 0 {
			return "orm命令格式错，请用help查看命令帮助", nil
		}
		key := ""
		val := ""
		for _, s := range kvs {

			if strings.HasPrefix(s, "-") {
				key = strings.TrimSpace(s)
				if key != "-p" && key != "-table" && key != "-path" {
					return "orm命令格式错，请用help查看命令帮助", nil
				}
			} else {
				val = strings.TrimSpace(s)
				cmd.kv[key] = val
			}
		}

		return "", cmd
	} else if strings.HasPrefix(command, HELP) {
		cmd := newCmdLine()
		cmd.word = HELP
		return "", cmd
	} else {
		return "命令不存在，请用help查看命令帮助", nil
	}

}
