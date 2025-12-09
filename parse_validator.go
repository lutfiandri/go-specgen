package specgen

import (
	"strconv"
	"strings"

	"github.com/swaggest/jsonschema-go"
	"github.com/swaggest/openapi-go/openapi3"
)

func ParseValidatorV10(reflector *openapi3.Reflector, structure any) {
	reflector.DefaultOptions = append(reflector.DefaultOptions,
		jsonschema.InterceptProp(
			// Property-level (field-level) validation
			func(params jsonschema.InterceptPropParams) error {
				if params.PropertySchema == nil {
					return nil
				}

				field := params.Field

				validationTagKey := "validate"
				validationInfo := ParseValidatorV10Tag(field.Tag.Get(validationTagKey))

				// required
				if validationInfo.Required {
					requiredFields := params.ParentSchema.Required
					requiredFields = append(requiredFields, params.Name)
					params.ParentSchema.WithRequired(requiredFields...)
				}

				// enum
				if len(validationInfo.OneOf) > 0 {
					enumValues := make([]any, len(validationInfo.OneOf))
					for i, v := range validationInfo.OneOf {
						enumValues[i] = v
					}
					params.PropertySchema.Enum = enumValues
				}

				// format
				if validationInfo.Format != "" {
					format := validationInfo.Format
					params.PropertySchema.Format = &format
				}

				// String validators: min, max, len → minLength, maxLength
				if params.PropertySchema.HasType(jsonschema.String) {
					if validationInfo.Len != nil {
						// len=X sets both minLength and maxLength to the same value
						length := *validationInfo.Len
						params.PropertySchema.MinLength = length
						params.PropertySchema.MaxLength = &length
					} else {
						if validationInfo.Min != nil {
							minLen := int64(*validationInfo.Min)
							params.PropertySchema.MinLength = minLen
						}
						if validationInfo.Max != nil {
							maxLen := int64(*validationInfo.Max)
							params.PropertySchema.MaxLength = &maxLen
						}
					}
				}

				// Number validators: min/gte → minimum, max/lte → maximum, gt → minimum+exclusiveMinimum, lt → maximum+exclusiveMaximum
				if params.PropertySchema.HasType(jsonschema.Number) || params.PropertySchema.HasType(jsonschema.Integer) {
					// Handle minimum (from min or gte)
					if validationInfo.Min != nil {
						minVal := *validationInfo.Min
						params.PropertySchema.Minimum = &minVal
					} else if validationInfo.Gte != nil {
						gteVal := *validationInfo.Gte
						params.PropertySchema.Minimum = &gteVal
					}
					// Handle gt (exclusive minimum)
					if validationInfo.Gt != nil {
						gtVal := *validationInfo.Gt
						params.PropertySchema.Minimum = &gtVal
						params.PropertySchema.ExclusiveMinimum = &gtVal
					}

					// Handle maximum (from max or lte)
					if validationInfo.Max != nil {
						maxVal := *validationInfo.Max
						params.PropertySchema.Maximum = &maxVal
					} else if validationInfo.Lte != nil {
						lteVal := *validationInfo.Lte
						params.PropertySchema.Maximum = &lteVal
					}
					// Handle lt (exclusive maximum)
					if validationInfo.Lt != nil {
						ltVal := *validationInfo.Lt
						params.PropertySchema.Maximum = &ltVal
						params.PropertySchema.ExclusiveMaximum = &ltVal
					}
				}

				// Array validators: min → minItems, max → maxItems
				if params.PropertySchema.HasType(jsonschema.Array) {
					if validationInfo.Min != nil {
						minItems := int64(*validationInfo.Min)
						params.PropertySchema.MinItems = minItems
					}
					if validationInfo.Max != nil {
						maxItems := int64(*validationInfo.Max)
						params.PropertySchema.MaxItems = &maxItems
					}
				}

				return nil
			},
		),
	)
}

type ValidationInfo struct {
	Required bool
	Format   string // email, uri, uuid, date-time
	OneOf    []string
	Min      *float64
	Max      *float64
	Len      *int64
	Gt       *float64
	Lt       *float64
	Gte      *float64
	Lte      *float64
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
			continue
		}

		// Check for oneof
		if enumStr, ok := strings.CutPrefix(part, "oneof="); ok {
			info.OneOf = strings.Fields(enumStr)
			continue
		}

		// Check for format validators
		switch part {
		case "email":
			info.Format = "email"
			continue
		case "url":
			info.Format = "uri"
			continue
		case "uuid":
			info.Format = "uuid"
			continue
		case "datetime":
			info.Format = "date-time"
			continue
		}

		// Check for min validator
		if minStr, ok := strings.CutPrefix(part, "min="); ok {
			if val, err := strconv.ParseFloat(minStr, 64); err == nil {
				info.Min = &val
			}
			continue
		}
		// Check for max validator
		if maxStr, ok := strings.CutPrefix(part, "max="); ok {
			if val, err := strconv.ParseFloat(maxStr, 64); err == nil {
				info.Max = &val
			}
			continue
		}
		// Check for len validator
		if lenStr, ok := strings.CutPrefix(part, "len="); ok {
			if val, err := strconv.ParseInt(lenStr, 10, 64); err == nil {
				info.Len = &val
			}
			continue
		}
		// Check for gte validator
		if gteStr, ok := strings.CutPrefix(part, "gte="); ok {
			if val, err := strconv.ParseFloat(gteStr, 64); err == nil {
				info.Gte = &val
			}
			continue
		}
		// Check for lte validator
		if lteStr, ok := strings.CutPrefix(part, "lte="); ok {
			if val, err := strconv.ParseFloat(lteStr, 64); err == nil {
				info.Lte = &val
			}
			continue
		}
		// Check for gt validator
		if gtStr, ok := strings.CutPrefix(part, "gt="); ok {
			if val, err := strconv.ParseFloat(gtStr, 64); err == nil {
				info.Gt = &val
			}
			continue
		}
		// Check for lt validator
		if ltStr, ok := strings.CutPrefix(part, "lt="); ok {
			if val, err := strconv.ParseFloat(ltStr, 64); err == nil {
				info.Lt = &val
			}
			continue
		}
	}

	return info
}
