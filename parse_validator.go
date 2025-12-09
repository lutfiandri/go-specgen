package specgen

import (
	"strings"

	"github.com/swaggest/jsonschema-go"
	"github.com/swaggest/openapi-go/openapi3"
)

func ParseValidatorV10(reflector *openapi3.Reflector, structure any) {
	validationTagKey := "validate"

	reflector.DefaultOptions = append(reflector.DefaultOptions,
		jsonschema.InterceptProp(
			// Property-level (field-level) validation
			func(params jsonschema.InterceptPropParams) error {
				if params.PropertySchema == nil {
					return nil
				}

				field := params.Field
				validationInfo := ParseValidatorV10Tag(field.Tag.Get(validationTagKey))

				// enum
				if len(validationInfo.OneOf) > 0 {
					enumValues := make([]any, len(validationInfo.OneOf))
					for i, v := range validationInfo.OneOf {
						enumValues[i] = v
					}
					params.PropertySchema.Enum = enumValues
				}

				return nil
			},
		),
	)
}

type ValidationInfo struct {
	Required bool
	OneOf    []string
}

func ParseValidatorV10Tag(validateTag string) ValidationInfo {
	info := ValidationInfo{}
	if validateTag == "" {
		return info
	}

	parts := strings.SplitSeq(validateTag, ",")
	for part := range parts {
		part = strings.TrimSpace(part)

		// Check for required
		if part == "required" {
			info.Required = true
		}

		// Check for oneof
		if enumStr, ok := strings.CutPrefix(part, "oneof="); ok {
			info.OneOf = strings.Fields(enumStr)
		}
	}

	return info
}
