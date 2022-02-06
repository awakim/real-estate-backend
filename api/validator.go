package api

import (
	db "github.com/awakim/immoblock-backend/db/sqlc"
	"github.com/go-playground/validator/v10"
)

// IsSupportedGender returns true if the gender is supported
func IsSupportedGender(gender db.Gender) bool {
	switch gender {
	case db.GenderF, db.GenderM:
		return true
	}
	return false
}

var validGender validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if gender, ok := fieldLevel.Field().Interface().(db.Gender); ok {
		return IsSupportedGender(gender)
	}
	return false
}
