package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

type envDecoder struct {
	cfg      *Config
	mappings []fieldMapping
}

// fieldMapping represents a mapping between YAML path and struct field.
type fieldMapping struct {
	EnvName   string
	EnvAlias  string
	FieldPath []int
	FieldType reflect.Type
}

// LoadFromEnv populates a configuration struct using environment variables based on field mappings and provided context.
func LoadFromEnv(cfg *Config, config Configure) error {
	decoder := envDecoder{
		cfg:      cfg,
		mappings: []fieldMapping{},
	}
	rootValue := reflect.ValueOf(config).Elem()
	decoder.extractFieldMappings(rootValue.Type(), []string{}, []int{})

	for _, mapping := range decoder.mappings {
		envValue, ok := decoder.getEnvValue(mapping.EnvName)
		if !ok {
			if mapping.EnvAlias != "" {
				if envValue, ok = decoder.getEnvValue(mapping.EnvAlias); !ok {
					continue
				}
			} else {
				continue
			}
		}

		if err := decoder.setFieldValue(rootValue, mapping, envValue); err != nil {
			return fmt.Errorf("error setting %s: %w", mapping.EnvName, err)
		}
	}
	return nil
}

// extractFieldMappings recursively processes struct fields.
func (r *envDecoder) extractFieldMappings(t reflect.Type, yamlPath []string, fieldPath []int) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if !field.IsExported() {
			continue
		}

		yamlTag := field.Tag.Get("yaml")
		if yamlTag == "" || yamlTag == "-" {
			continue
		}

		yamlName := strings.Split(yamlTag, ",")[0]
		currentYAMLPath := append(yamlPath, yamlName) //nolint:gocritic
		currentFieldPath := append(fieldPath, i)      //nolint:gocritic

		fieldType := field.Type

		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		if fieldType.Kind() == reflect.Struct {
			r.extractFieldMappings(fieldType, currentYAMLPath, currentFieldPath)
		} else {
			r.mappings = append(r.mappings, fieldMapping{
				EnvName:   r.yamlPathToEnvName(currentYAMLPath),
				EnvAlias:  strings.ToUpper(r.cfg.EnvAliases[strings.Join(currentYAMLPath, ".")]),
				FieldPath: currentFieldPath,
				FieldType: field.Type,
			})
		}
	}
}

// yamlPathToEnvName converts a YAML path to an environment variable name
// e.g., "log.level" -> "LOG_LEVEL" or "APP_LOG_LEVEL" (with prefix).
func (r *envDecoder) yamlPathToEnvName(yamlPath []string) string {
	envName := strings.ToUpper(strings.Join(yamlPath, "_"))
	if r.cfg.EnvPrefix != "" {
		return r.cfg.EnvPrefix + "_" + envName
	}
	return envName
}

// getEnvValue retrieves the environment variable value for a given name.
func (r *envDecoder) getEnvValue(envName string) (string, bool) {
	if value, ok := os.LookupEnv(envName); ok {
		return value, true
	}

	if r.cfg.EnvFileLoading {
		if filePath, ok := os.LookupEnv(envName + "_FILE"); ok {
			content, err := os.ReadFile(filepath.Clean(filePath))
			if err != nil {
				r.cfg.Logger.Warn(err)
				return "", false
			}
			return strings.TrimSpace(string(content)), true
		}
	}

	return "", false
}

// setFieldValue sets a value on a nested field using reflection.
func (r *envDecoder) setFieldValue(rootValue reflect.Value, mapping fieldMapping, value string) error {
	fieldValue := rootValue
	for _, idx := range mapping.FieldPath {
		fieldValue = fieldValue.Field(idx)

		if fieldValue.Kind() == reflect.Ptr {
			if fieldValue.IsNil() {
				fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
			}
			fieldValue = fieldValue.Elem()
		}
	}

	if !fieldValue.CanSet() {
		return ErrFieldNotSettable
	}

	return r.setConvertedValue(fieldValue, mapping.FieldType, value)
}

// setConvertedValue converts the string value to the appropriate type and sets it.
func (r *envDecoder) setConvertedValue(fieldValue reflect.Value, fieldType reflect.Type, value string) error {
	if fieldType.Kind() == reflect.Ptr {
		if fieldValue.IsNil() {
			fieldValue.Set(reflect.New(fieldType.Elem()))
		}
		fieldValue = fieldValue.Elem()
		fieldType = fieldType.Elem()
	}

	switch fieldType.Kind() {
	case reflect.String:
		fieldValue.SetString(value)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("cannot convert '%s' to int: %w", value, err)
		}
		fieldValue.SetInt(intVal)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("cannot convert '%s' to uint: %w", value, err)
		}
		fieldValue.SetUint(uintVal)

	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("cannot convert '%s' to float: %w", value, err)
		}
		fieldValue.SetFloat(floatVal)

	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("cannot convert '%s' to bool: %w", value, err)
		}
		fieldValue.SetBool(boolVal)

	case reflect.Slice:
		if fieldType.Elem().Kind() == reflect.String {
			parts := strings.Split(value, ",")
			slice := reflect.MakeSlice(fieldType, len(parts), len(parts))
			for i, part := range parts {
				slice.Index(i).SetString(strings.TrimSpace(part))
			}
			fieldValue.Set(slice)
		} else {
			return fmt.Errorf("unsupported slice type: %s", fieldType.Elem().Kind())
		}

	default:
		return fmt.Errorf("unsupported field type: %s", fieldType.Kind())
	}

	return nil
}
