package config

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func InitValidator(){
	validate = validator.New()
}

func ValidateStruct(s interface{}) (error){
	if validate == nil{
		return nil
	}
	return validate.Struct(s)
}