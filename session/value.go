package session

import (
	"strconv"
)

type SessionValue interface {
	String() string
	Raw() interface{}
}

type Int64SessionValue int64

func (i Int64SessionValue) String() string {
	return strconv.FormatInt(int64(i), 10)
}

func (i Int64SessionValue) Raw() interface{} {
	return int64(i)
}

//TODO: ValueFromInt64... rename
func FromInt64(i int64) SessionValue {
	return Int64SessionValue(i)
}

type StringSessionValue string

func (s StringSessionValue) String() string {
	return string(s)
}

func (s StringSessionValue) Raw() interface{} {
	return string(s)
}

func FromString(s string) SessionValue {
	return StringSessionValue(s)
}
