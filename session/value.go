package session

import "strconv"

type SessionValue interface {
	String() string
}

type Int64SessionValue int64

func (i Int64SessionValue) String() string {
	return strconv.FormatInt(int64(i), 10)
}

func FromInt64(i int64) SessionValue {
	return Int64SessionValue(i)
}
