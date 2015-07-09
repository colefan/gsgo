package db

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

const (
	MAX_DB_CONNECTIONS = 100
)

type DbPool interface {
	//初始化链接池
	//数据库名称，类型，连接串，初始化连接数
	InitPool(database string, dbtype string, connstr string, initConnsNum int) error
	//从连接池中获取一个连接
	GetConnection() Connection
	//将连接还给连接池
	ReleaseConnection(conn Connection) bool
	//销毁数据连接池
	DestroyPool()
}

type DefaultDbPool struct {
	dbtype    string
	dbname    string
	dbconnstr string
	connNums  int
	freeConns []Connection
	curConnId int
	lock      sync.Mutex
}

func NewDefaultDbPool() DbPool {
	pool := &DefaultDbPool{}
	pool.freeConns = make([]Connection, 0)
	return pool
}

//database:testdb,dbtype:mysql,connstr:username:pwd@tcp(192.168.13.21:3306)/testdb
func (this *DefaultDbPool) InitPool(database string, dbtype string, connstr string, initConnsNum int) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.dbname = database
	this.dbtype = dbtype
	this.dbconnstr = connstr
	if initConnsNum <= 0 {
		initConnsNum = 5
	}
	this.connNums = 0

	if "mysql" == strings.ToLower(dbtype) {
		for i := 0; i < initConnsNum; i++ {
			_, err := this.createMysqlConnection()

			if err != nil {
				return err
			}
		}
	} else {
		return fmt.Errorf("unknown db type: ", dbtype)
	}

	return nil
}

func (this *DefaultDbPool) createMysqlConnection() (Connection, error) {
	db, err := sql.Open("mysql", this.dbconnstr)
	if err != nil {
		return nil, err
	}
	conn := NewDbConnection(db)
	conn.SetID(this.getNextConnId())
	conn.Used(false)
	this.connNums++
	this.freeConns = append(this.freeConns, conn)
	return conn, nil
}

func (this *DefaultDbPool) GetConnection() Connection {
	this.lock.Lock()
	defer this.lock.Unlock()
	fmt.Println("free len ,", len(this.freeConns))
	for _, tmp := range this.freeConns {

		if false == tmp.IsUsed() {
			tmp.Used(true)
			return tmp
		}
	}

	if this.connNums < MAX_DB_CONNECTIONS {
		conn, err := this.createMysqlConnection()
		if err != nil {
			return nil
		}
		conn.Used(true)
		return conn
	} else {
		return nil
	}

	return nil
}

func (this *DefaultDbPool) ReleaseConnection(conn Connection) bool {
	if conn == nil {
		return false
	}
	this.lock.Lock()
	defer this.lock.Unlock()

	for _, tmp := range this.freeConns {
		if tmp.GetID() == conn.GetID() {
			tmp.Used(false)
			return true
		}
	}
	return false
}
func (this *DefaultDbPool) DestroyPool() {
	this.lock.Lock()
	defer this.lock.Unlock()
	for _, tmp := range this.freeConns {
		tmp.Close()
	}

}

func (this *DefaultDbPool) getNextConnId() int {
	this.curConnId++
	return this.curConnId
}
