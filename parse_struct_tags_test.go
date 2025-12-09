package specgen_test

import (
	"reflect"
	"testing"

	"github.com/lutfiandri/go-specgen"
)

type TestStruct struct {
	HasTags   string `json:"has_tags" validate:"required" params:"id"`
	NoTags    string
	OtherTags string `params:"id" query:"search"`
}

func TestExtractStructFieldTags(t *testing.T) {
	typ := reflect.TypeOf(TestStruct{})
	field, _ := typ.FieldByName("HasTags")
	result := specgen.ExtractStructFieldTags(field, []string{"json", "validate", "params"})

	expectedTags := map[string]string{
		"json":     "has_tags",
		"validate": "required",
		"params":   "id",
	}

	if len(result) != len(expectedTags) {
		t.Errorf("Expected %d tags, got %d", len(expectedTags), len(result))
	}

	for _, tag := range result {
		expectedValue, ok := expectedTags[tag.Key]
		if !ok {
			t.Errorf("Unexpected tag key: %s", tag.Key)
			continue
		}
		if tag.Value != expectedValue {
			t.Errorf("Expected value '%s' for key '%s', got '%s'", expectedValue, tag.Key, tag.Value)
		}
	}
}

func TestExtractStructTags(t *testing.T) {
	result := specgen.ExtractStructTags(TestStruct{}, []string{"json", "validate", "params", "query"})

	expectedTags := map[string]map[string]string{
		"HasTags": {
			"json":     "has_tags",
			"validate": "required",
			"params":   "id",
		},
		"NoTags": {},
		"OtherTags": {
			"params": "id",
			"query":  "search",
		},
	}

	if len(result) != len(expectedTags) {
		t.Fatalf("Expected %d fields, got %d", len(expectedTags), len(result))
	}

	for _, structTag := range result {
		expectedFieldTags, ok := expectedTags[structTag.Name]
		if !ok {
			t.Errorf("Unexpected field name: %s", structTag.Name)
			continue
		}

		if len(structTag.Tags) != len(expectedFieldTags) {
			t.Errorf("Field '%s': Expected %d tags, got %d", structTag.Name, len(expectedFieldTags), len(structTag.Tags))
			continue
		}

		for _, tag := range structTag.Tags {
			expectedValue, ok := expectedFieldTags[tag.Key]
			if !ok {
				t.Errorf("Field '%s': Unexpected tag key: %s", structTag.Name, tag.Key)
				continue
			}
			if tag.Value != expectedValue {
				t.Errorf("Field '%s': Expected value '%s' for key '%s', got '%s'", structTag.Name, expectedValue, tag.Key, tag.Value)
			}
		}
	}
}
