package runtime

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

// Decode decodes the specified val into the specified target.
func Decode(target interface{}, val string) error {
	return decode(reflect.ValueOf(target).Elem(), val)
}

func decode(target reflect.Value, inputValue string) error {
	targetType := target.Type()

	if target.Kind() == reflect.Ptr {
		target.Set(reflect.New(targetType.Elem()))

		return decode(target.Elem(), inputValue)
	}

	if targetType.Kind() == reflect.String {
		target.Set(reflect.ValueOf(inputValue))
		return nil
	}

	if targetType.Kind() == reflect.Struct {
		if targetProto, ok := target.Addr().Interface().(proto.Message); ok {
			return jsonpb.UnmarshalString(inputValue, targetProto)
		}

		return fmt.Errorf("Unacceptable type %s", targetType)
	}

	return json.Unmarshal([]byte(inputValue), target.Addr().Interface())
}
