package valid

import "github.com/asaskevich/govalidator"

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
	govalidator.SetNilPtrAllowedByRequired(true)
}
