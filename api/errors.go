package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func ValidationError(vErrs validator.ValidationErrors) map[string]string {
	errs := make(map[string]string)

	for _, f := range vErrs {
		err := f.ActualTag()
		if f.Param() != "" {
			err = fmt.Sprintf("%s=%s", err, f.Param())
		}
		errs[f.Field()] = err
	}

	return errs
}
