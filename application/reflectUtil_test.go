package application

import (
	"testing"
)

type Lain struct {
	Name string
	Age  int64
}

func TestGetObjectParamPath(t *testing.T) {
	s := Lain{
		Name: "13",
		Age:  14,
	}
	GetObjectParamPath(s)
}
