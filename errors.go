package gogination

type constError string

func (err constError) Error() string {
	return string(err)
}

// errors
const (
	ErrExpectedStruct = constError("struct type object expected")
	ErrNoIdField      = constError("given object does not have an ID field")
	ErrInvalidOrderBy = constError("given order by document is invalid")
)
