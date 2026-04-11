package base

import (
	"fmt"
	"reflect"
	"regexp"

	"go-framework-guap/core/errors"
)

type Validator struct {
	rules map[string][]ValidationRule
}

type ValidationRule func(field string, value any) error

func NewValidator() *Validator {
	return &Validator{
		rules: make(map[string][]ValidationRule),
	}
}

func (v *Validator) AddRule(field string, rule ValidationRule) {
	v.rules[field] = append(v.rules[field], rule)
}

func (v *Validator) Validate(data interface{}) []error {
	var errs []error
	val := reflect.ValueOf(data)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return []error{fmt.Errorf("expected struct, got %s", val.Kind())}
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldName := field.Name

		rules, ok := v.rules[fieldName]
		if !ok {
			continue
		}

		fieldValue := val.Field(i).Interface()

		for _, rule := range rules {
			if err := rule(fieldName, fieldValue); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return errs
}

func Required(field string, value any) error {
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr && rv.IsNil() {
		return errors.NewValidationError(field, "is required")
	}
	if isEmpty(rv) {
		return errors.NewValidationError(field, "is required")
	}
	return nil
}

func MinLength(min int) ValidationRule {
	return func(field string, value any) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		if len(str) < min {
			return fmt.Errorf("field %s must be at least %d characters", field, min)
		}
		return nil
	}
}

func MaxLength(max int) ValidationRule {
	return func(field string, value any) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		if len(str) > max {
			return fmt.Errorf("field %s must be at most %d characters", field, max)
		}
		return nil
	}
}

func Email(field string, value any) error {
	email, ok := value.(string)
	if !ok {
		return nil
	}
	if email == "" {
		return nil
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.NewValidationError(field, "must be a valid email")
	}
	return nil
}

func Pattern(pattern string) ValidationRule {
	return func(field string, value any) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		re := regexp.MustCompile(pattern)
		if !re.MatchString(str) {
			return fmt.Errorf("field %s does not match pattern %s", field, pattern)
		}
		return nil
	}
}

func Range(min, max int) ValidationRule {
	return func(field string, value any) error {
		num, ok := value.(int)
		if !ok {
			return nil
		}
		if num < min || num > max {
			return fmt.Errorf("field %s must be between %d and %d", field, min, max)
		}
		return nil
	}
}

func Positive(field string, value any) error {
	num, ok := value.(int)
	if !ok {
		fnum, ok := value.(float64)
		if !ok {
			return nil
		}
		if fnum <= 0 {
			return fmt.Errorf("field %s must be positive", field)
		}
		return nil
	}
	if num <= 0 {
		return fmt.Errorf("field %s must be positive", field)
	}
	return nil
}

func Bool(field string, value any) error {
	_, ok := value.(bool)
	if !ok {
		return errors.NewValidationError(field, "must be boolean")
	}
	return nil
}

func isEmpty(rv reflect.Value) bool {
	switch rv.Kind() {
	case reflect.String:
		return rv.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return rv.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return rv.Float() == 0
	case reflect.Bool:
		return !rv.Bool()
	case reflect.Slice, reflect.Map, reflect.Array:
		return rv.Len() == 0
	case reflect.Ptr:
		return rv.IsNil()
	}
	return false
}

type Schema struct {
	fields map[string]fieldSchema
}

type fieldSchema struct {
	required   bool
	minLength  *int
	maxLength  *int
	pattern    *string
	minValue   *int
	maxValue   *int
	email      bool
	boolean    bool
	custom     []ValidationRule
}

func NewSchema() *Schema {
	return &Schema{
		fields: make(map[string]fieldSchema),
	}
}

func (s *Schema) Field(name string) *fieldSchema {
	fs := &fieldSchema{}
	s.fields[name] = *fs
	return fs
}

func (fs *fieldSchema) Required() *fieldSchema {
	fs.required = true
	return fs
}

func (fs *fieldSchema) MinLength(n int) *fieldSchema {
	fs.minLength = &n
	return fs
}

func (fs *fieldSchema) MaxLength(n int) *fieldSchema {
	fs.maxLength = &n
	return fs
}

func (fs *fieldSchema) Pattern(p string) *fieldSchema {
	fs.pattern = &p
	return fs
}

func (fs *fieldSchema) Range(min, max int) *fieldSchema {
	fs.minValue = &min
	fs.maxValue = &max
	return fs
}

func (fs *fieldSchema) Email() *fieldSchema {
	fs.email = true
	return fs
}

func (fs *fieldSchema) Boolean() *fieldSchema {
	fs.boolean = true
	return fs
}

func (fs *fieldSchema) Custom(rule ValidationRule) *fieldSchema {
	fs.custom = append(fs.custom, rule)
	return fs
}

func (s *Schema) Validate(data interface{}) []error {
	var errs []error
	val := reflect.ValueOf(data)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return []error{fmt.Errorf("expected struct, got %s", val.Kind())}
	}

	for name, fs := range s.fields {
		_, ok := val.Type().FieldByName(name)
		if !ok {
			if fs.required {
				errs = append(errs, errors.NewValidationError(name, "is required"))
			}
			continue
		}

		fieldValue := val.FieldByName(name)

		if fs.required && isEmpty(fieldValue) {
			errs = append(errs, errors.NewValidationError(name, "is required"))
		}

		if fs.minLength != nil && fieldValue.Kind() == reflect.String {
			if len(fieldValue.String()) < *fs.minLength {
				errs = append(errs, fmt.Errorf("field %s must be at least %d characters", name, *fs.minLength))
			}
		}

		if fs.maxLength != nil && fieldValue.Kind() == reflect.String {
			if len(fieldValue.String()) > *fs.maxLength {
				errs = append(errs, fmt.Errorf("field %s must be at most %d characters", name, *fs.maxLength))
			}
		}

		if fs.pattern != nil && fieldValue.Kind() == reflect.String {
			re := regexp.MustCompile(*fs.pattern)
			if !re.MatchString(fieldValue.String()) {
				errs = append(errs, fmt.Errorf("field %s does not match pattern", name))
			}
		}

		if fs.email && fieldValue.Kind() == reflect.String {
			email := fieldValue.String()
			if email != "" {
				emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
				if !emailRegex.MatchString(email) {
					errs = append(errs, errors.NewValidationError(name, "must be a valid email"))
				}
			}
		}

		if fs.boolean && fieldValue.Kind() != reflect.Bool {
			errs = append(errs, errors.NewValidationError(name, "must be boolean"))
		}

		for _, rule := range fs.custom {
			if err := rule(name, fieldValue.Interface()); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return errs
}

type JSONSchema map[string]interface{}

func ValidateJSONSchema(data map[string]interface{}, schema JSONSchema) []error {
	var errs []error

	required, _ := schema["required"].([]string)
	properties, _ := schema["properties"].(map[string]interface{})

	for _, field := range required {
		if _, ok := data[field]; !ok {
			errs = append(errs, errors.NewValidationError(field, "is required"))
		}
	}

	for field, value := range data {
		prop, ok := properties[field].(map[string]interface{})
		if !ok {
			continue
		}

		if prop["type"] == "string" {
			if str, ok := value.(string); ok {
				if minLength, ok := prop["minLength"].(float64); ok && len(str) < int(minLength) {
					errs = append(errs, fmt.Errorf("field %s is too short", field))
				}
				if maxLength, ok := prop["maxLength"].(float64); ok && len(str) > int(maxLength) {
					errs = append(errs, fmt.Errorf("field %s is too long", field))
				}
			}
		}

		if prop["type"] == "integer" {
			if num, ok := value.(float64); ok {
				if minimum, ok := prop["minimum"].(float64); ok && num < minimum {
					errs = append(errs, fmt.Errorf("field %s is below minimum", field))
				}
				if maximum, ok := prop["maximum"].(float64); ok && num > maximum {
					errs = append(errs, fmt.Errorf("field %s is above maximum", field))
				}
			}
		}
	}

	return errs
}
