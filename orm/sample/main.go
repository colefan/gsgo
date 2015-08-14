package main

import "github.com/colefan/gsgo/orm"

func main() {
	cmd := &orm.OrmCmd{}
	cmd.Run("../bstest")
}
