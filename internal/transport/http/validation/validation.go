package validation

import (
	"github.com/gookit/validate"
)

func ValidateStruct(s interface{}) []byte {
	v := validate.Struct(s)
	if v.Validate() {
		return nil
	}
	return v.Errors.JSON()
}
