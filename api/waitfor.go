package api

import (
	"fmt"
	"strconv"

	"github.com/VictorLowther/jsonpatch2/utils"
)

// TestItem creates a test function to see if a value in the
// passed interface is true.
func TestItem(field, value string) func(interface{}) (bool, error) {
	return func(ref interface{}) (bool, error) {
		var err error
		fields := map[string]interface{}{}
		if err := utils.Remarshal(ref, fields); err != nil {
			return false, err
		}
		matched := false
		if d, ok := fields[field]; ok {
			switch v := d.(type) {
			case bool:
				var bval bool
				bval, err = strconv.ParseBool(value)
				if err == nil {
					if v == bval {
						matched = true
					}
				}
			case string:
				if v == value {
					matched = true
				}
			case int:
				var ival int64
				ival, err = strconv.ParseInt(value, 10, 64)
				if err == nil {
					if int(ival) == v {
						matched = true
					}
				}
			default:
				err = fmt.Errorf("Unsupported field type: %T\n", d)
			}
		}
		return matched, err
	}
}
