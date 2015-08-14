package orm

import (
	"fmt"

	"github.com/colefan/gsgo/config"
)

type OrmConfig struct {
	dbtype string
	dbname string
	dblink string
}

func (this *OrmConfig) LoadConfig(filepath string) error {
	cfg, err := config.NewConfig("ini", filepath)
	if err != nil {
		return fmt.Errorf("load file.error->", err.Error())
	}
	this.dbtype = cfg.String("dbtype")
	this.dbname = cfg.String("dbname")
	this.dblink = cfg.String("dblink")
	return nil
}
