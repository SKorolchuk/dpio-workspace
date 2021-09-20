package validation

import (
	"context"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translation "github.com/go-playground/validator/v10/translations/en"
)

// v contains rules and settings for validation values.
var v *validator.Validate

// tr contains cache for locale translation.
var tr ut.Translator

func init() {
	v = validator.New()

	tr, _ = ut.New(en.New(), en.New()).GetTranslator("en")

	// Register default English locale
	en_translation.RegisterDefaultTranslations(v, tr)

	// Use JSON tag names of fields instead of Go structure field names.
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}

		return name
	})
}

// Check query the model and validate it against declared tags.
func Check(ctx context.Context, value interface{}) error {
	if err := v.StructCtx(ctx, value); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)

		if !ok {
			return err
		}

		var fields FieldErrors
		for _, validationError := range validationErrors {
			field := FieldError{
				FieldName:    validationError.Field(),
				ErrorMessage: validationError.Translate(tr),
			}
			fields = append(fields, field)
		}

		return fields
	}

	return nil
}
