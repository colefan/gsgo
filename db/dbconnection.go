package db

import (
	"database/sql"
	"sync"
)

type Connection interface {
	GetID() int
	SetID(id int)
	Close() error
	//执行语句，是否返回主键
	Create(sql string) error
	Read(sql string) (*sql.Rows, error)
	Update(sql string) error
	Delete(sql string) error
	IsUsed() bool
	Used(b bool)
}

type DbConnection struct {
	*sql.DB
	id int
	sync.RWMutex
	used bool
}

func NewDbConnection(db *sql.DB) *DbConnection {
	conn := &DbConnection{}
	conn.DB = db
	return conn
}

func (conn *DbConnection) GetID() int {
	return conn.id
}

func (conn *DbConnection) SetID(id int) {
	conn.id = id
}

func (conn *DbConnection) Close() error {
	return conn.DB.Close()
}

func (conn *DbConnection) Create(sql string) error {
	result, err := conn.DB.Query(sql)
	if err != nil {
		return err
	} else {
		defer result.Close()
		return err
	}

}

func (conn *DbConnection) Read(sql string) (*sql.Rows, error) {
	result, err := conn.DB.Query(sql)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (conn *DbConnection) Update(sql string) error {
	result, err := conn.DB.Query(sql)
	if err != nil {
		return err
	} else {
		defer result.Close()
		return err
	}

}

func (conn *DbConnection) Delete(sql string) error {
	result, err := conn.DB.Query(sql)
	if err != nil {
		return err
	} else {
		defer result.Close()
		return err
	}

}

func (conn *DbConnection) IsUsed() bool {
	conn.RWMutex.RLock()
	defer conn.RWMutex.RUnlock()
	return conn.used
}

func (conn *DbConnection) Used(b bool) {
	conn.RWMutex.Lock()
	defer conn.RWMutex.Unlock()
	conn.used = b
}
