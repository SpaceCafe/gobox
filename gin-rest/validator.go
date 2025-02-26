package rest

import (
	"errors"
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

const (
	// sortByAttributesRegexString is a Regex pattern to validate sorting attributes,
	// allowing optional '+' or '-' prefix and comma-separated values.
	sortByAttributesRegexString = "^(?:[+-]?[a-z][a-z0-9_]*[a-z0-9](?:,|$))*"
)

var (
	ErrNoValidatorEngine = errors.New("no validator engine found")
)

// InitializeValidators create custom validators and registers them with the validation engine.
// It compiles the regex for sorting attributes and registers the validators with the engine.
func InitializeValidators() error {
	sortByAttributesRegex := regexp.MustCompile(sortByAttributesRegexString)

	validators := map[string]validator.Func{
		"sort_by_attributes": func(fl validator.FieldLevel) bool {
			return sortByAttributesRegex.MatchString(fl.Field().String())
		},
		"readonly": func(_ validator.FieldLevel) bool { return false },
	}

	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		for tag, fn := range validators {
			if err := validate.RegisterValidation(tag, fn); err != nil {
				return err
			}
		}
		return nil
	}
	return ErrNoValidatorEngine
}
