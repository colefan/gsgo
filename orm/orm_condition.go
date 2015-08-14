package orm

type QueryCondition struct {
	condtion string
}

func CreateQueryCondition() *QueryCondition {
	return &QueryCondition{}
}

func (this *QueryCondition) And(condstr string) *QueryCondition {
	this.condtion += " and " + condstr
	return this
}

func (this *QueryCondition) Add(filed OrmField) *QueryCondition {
	return this
}

func (this *QueryCondition) Or(connstr string) *QueryCondition {
	this.condtion += " or " + connstr
	return this
}
