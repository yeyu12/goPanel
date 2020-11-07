package validations

import (
	"errors"
	zh_cn "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
)

var Validate *validator.Validate
var trans ut.Translator

func init() {
	zh := zh_cn.New()
	uni := ut.New(zh, zh)
	trans, _ = uni.GetTranslator("zh")

	Validate = validator.New()
	Validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		if label == "" {
			return field.Name
		}

		return label
	})

	_ = zh_translations.RegisterDefaultTranslations(Validate, trans)
}

func Translate(errs validator.ValidationErrors) error {
	for _, e := range errs {
		return errors.New(e.Translate(trans))
	}

	return nil
}
