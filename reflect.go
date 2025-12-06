package Fcload

import (
	"fmt"
	"reflect"
	"strconv"

	"gopkg.in/yaml.v3"
)

func MergeYamlNode(node *yaml.Node, cfg *any) {
	v := reflect.ValueOf(*cfg)
	if v.Kind() != reflect.Struct {
		return
	}

	mergeYamlNodeRecursive(node, v)
}

func mergeYamlNodeRecursive(node *yaml.Node, value reflect.Value) {
	if node == nil {
		return
	}

	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return
		}
		value = value.Elem()
	}

	//nolint:exhaustive
	switch node.Kind {
	case yaml.DocumentNode:
		for _, n := range node.Content {
			mergeYamlNodeRecursive(n, value)
		}
	case yaml.MappingNode:
		if value.Kind() != reflect.Struct {
			return
		}

		// Precompute struct fields into a map for faster lookups.
		fields := make(map[string]reflect.Value)
		t := value.Type()
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			tag := field.Tag.Get("yaml")
			if tag != "" {
				fields[tag] = value.Field(i)
			}
		}

		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			valueNode := node.Content[i+1]
			fieldName := keyNode.Value

			if fieldValue, found := fields[fieldName]; found {
				mergeYamlNodeRecursive(valueNode, fieldValue)
			}
		}
	case yaml.ScalarNode:
		if value.CanSet() {
			var newValue string
			//nolint:exhaustive
			switch value.Kind() {
			case reflect.String:
				newValue = value.String()
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				newValue = strconv.FormatInt(value.Int(), 10)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				newValue = strconv.FormatUint(value.Uint(), 10)
			case reflect.Float32, reflect.Float64:
				// Use %g for cleaner float output, avoiding trailing zeros.
				newValue = fmt.Sprintf("%g", value.Float())
			case reflect.Bool:
				newValue = strconv.FormatBool(value.Bool())
			default:
				// Unsupported types are ignored to prevent incorrect mutations.
				return
			}

			if node.Value != newValue && newValue != "" {
				node.Value = newValue
			}
		}
	case yaml.SequenceNode:
		if value.Kind() == reflect.Slice {
			for i := 0; i < value.Len() && i < len(node.Content); i++ {
				mergeYamlNodeRecursive(node.Content[i], value.Index(i))
			}
		}
	default:
	}
}
