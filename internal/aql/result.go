package aql

type Result struct {
	lastInsertId int64
	rowsAffected int64
	result       interface{}
	err          error
}

func NewResult(lastInsertId int64, rowsAffected int64, result interface{}, err error) Result {
	return Result{
		lastInsertId: lastInsertId,
		rowsAffected: rowsAffected,
		result:       result,
		err:          err,
	}
}

func (r Result) LastInsertId() (int64, error) {
	return r.lastInsertId, r.err
}

func (r Result) RowsAffected() (int64, error) {
	return r.rowsAffected, r.err
}

func (r Result) Result() interface{} {
	return r.result
}
