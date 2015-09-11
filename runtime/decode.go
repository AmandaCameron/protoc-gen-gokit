package runtime

import (
	"fmt"

	"strconv"
)

// Decode decodes the specified val into the specified target.
func Decode(target interface{}, val string) error {
	switch target.(type) {
	case *string:
		*(target.(*string)) = val
		return nil

	case *int32:
		tmp, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			return err
		}

		*(target.(*int32)) = int32(tmp)
		return nil

	case *uint32:
		tmp, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			return err
		}

		*(target.(*uint32)) = uint32(tmp)
		return nil

	case *int64:
		tmp, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}

		*(target.(*int64)) = int64(tmp)
		return nil

	case *uint64:
		tmp, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}

		*(target.(*uint64)) = uint64(tmp)
		return nil

	default:
		return fmt.Errorf("Unacceptable type: %T", target)
	}
}
