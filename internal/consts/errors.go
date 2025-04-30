package consts

import "errors"

var (
	ErrTableDoesNotExist  = errors.New("table does not exist")
	ErrColumnDoesNotExist = errors.New("column does not exist")
	ErrInvalidSelectItem  = errors.New("select item is not valid")
	ErrInvalidDatatype    = errors.New("invalid datatype")
	ErrMissingValues      = errors.New("missing values")
	ErrTableAlreadyExists = errors.New("table already exists")
)
