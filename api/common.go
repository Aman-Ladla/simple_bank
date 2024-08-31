package api

import (
	"example.com/simple_bank/db/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func errorResponse(err error) gin.H {
	return gin.H{"err": err.Error()}
}

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsCurrencyValid(currency)
	}
	return false
}
