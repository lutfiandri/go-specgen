package specgen

import "reflect"

type Tag struct {
	Key   string
	Value string
}

func ExtractStructFieldTags(field reflect.StructField, tagKeys []string) []Tag {
	tags := make([]Tag, 0, len(tagKeys))

	for _, tagKey := range tagKeys {
		if value, ok := field.Tag.Lookup(tagKey); ok {
			tags = append(tags, Tag{Key: tagKey, Value: value})
		}
	}

	return tags
}

type StructTags struct {
	Name string
	Tags []Tag
}

func ExtractStructTags(structure any, tagKeys []string) []StructTags {
	val := reflect.ValueOf(structure)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil
	}

	typ := val.Type()

	structTags := make([]StructTags, 0)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		// Extract embedded structs
		if field.Anonymous {
			embeddedStructTags := ExtractStructTags(val.Field(i).Interface(), tagKeys)
			structTags = append(structTags, embeddedStructTags...)
			continue
		}

		fieldTags := ExtractStructFieldTags(field, tagKeys)
		structTags = append(structTags, StructTags{Name: field.Name, Tags: fieldTags})
	}

	return structTags
}
