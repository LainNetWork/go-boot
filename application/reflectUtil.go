package application

import (
	"reflect"
)

func GetObjectParamPath(obj interface{}) []string {
	objType := reflect.TypeOf(obj)
	objType.Elem()
	return make([]string, 0)
}
