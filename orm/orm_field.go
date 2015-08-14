package orm

type OrmField struct {
	Name      string
	FieldType string
}

func NewOrmField(fname, ftype string) *OrmField {
	return &OrmField{Name: fname, FieldType: ftype}
}
