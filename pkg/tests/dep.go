package tests

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"go.uber.org/dig"
)

var container = dig.New()

func Provide(constructors ...interface{}) {
	for _, constructor := range constructors {
		if err := container.Provide(constructor); err != nil {
			panic(err)
		}
	}
}

func Dep(msg string, t *testing.T, f interface{}) {
	convey.Convey(fmt.Sprintf("%s after loading dependencies: %v", msg, getArguments(f)), t, func() {
		if err := container.Invoke(f); err != nil {
			t.Error(err)
		}
	})
}

func getArguments(f interface{}) []string {
	x := reflect.TypeOf(f)
	numIn := x.NumIn()
	result := make([]string, numIn)
	for i := 0; i < numIn; i++ {
		result[i] = x.In(i).String()
	}
	return result
}
