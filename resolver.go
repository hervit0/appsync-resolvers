package resolvers

import (
	awsContext "context"
	"encoding/json"
	"reflect"
)

type resolver struct {
	// TODO: Enforce by type instead of reflection with indexes
	function interface{}
}

func (r *resolver) hasArguments() bool {
	return reflect.TypeOf(r.function).NumIn() == 2
}

func (r *resolver) call(ctx awsContext.Context, p json.RawMessage) (interface{}, error) {
	var args []reflect.Value
	var err error

	if r.hasArguments() {
		pld := payload{p}

		args, err = pld.parse(reflect.TypeOf(r.function).In(1))
		if err != nil {
			return nil, err
		}

		args = append([]reflect.Value{reflect.ValueOf(ctx)}, args...)
	}

	returnValues := reflect.ValueOf(r.function).Call(args)
	var returnData interface{}
	var returnError error

	if len(returnValues) == 2 {
		returnData = returnValues[0].Interface()
	}

	if err := returnValues[len(returnValues)-1].Interface(); err != nil {
		returnError = returnValues[len(returnValues)-1].Interface().(error)
	}

	return returnData, returnError
}
