package specgen_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/lutfiandri/go-specgen"
	"gopkg.in/yaml.v3"
)

// OpenAPISpec represents the structure of an OpenAPI 3.0 specification
type OpenAPISpec struct {
	OpenAPI    string                 `yaml:"openapi"`
	Info       Info                   `yaml:"info"`
	Paths      map[string]interface{} `yaml:"paths"`
	Components Components             `yaml:"components"`
}

type Info struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
}

type Components struct {
	Schemas map[string]Schema `yaml:"schemas"`
}

type Schema struct {
	Type       string              `yaml:"type"`
	Properties map[string]Property `yaml:"properties"`
	Required   []string            `yaml:"required"`
}

type Property struct {
	Type      *string       `yaml:"type"`
	Format    *string       `yaml:"format"`
	MinLength *int64        `yaml:"minLength"`
	MaxLength *int64        `yaml:"maxLength"`
	Minimum   *float64      `yaml:"minimum"`
	Maximum   *float64      `yaml:"maximum"`
	MinItems  *int64        `yaml:"minItems"`
	MaxItems  *int64        `yaml:"maxItems"`
	Enum      []interface{} `yaml:"enum"`
	Items     *Property     `yaml:"items"`
	Nullable  *bool         `yaml:"nullable"`
}

// Test request and response structs
type CreateUserRequest struct {
	Name  string `json:"name" required:"true"`
	Email string `json:"email" required:"true"`
	Age   int    `json:"age"`
}

type UserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func TestGenerateOpenAPISpec_Basic(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "spec.yaml")

	title := "Test API"
	description := "Test API Description"
	version := "1.0.0"

	config := specgen.SpecConfig{
		Title:       &title,
		Description: &description,
		Version:     &version,
	}

	routes := []specgen.Route{
		{
			Path:   "/users",
			Method: "POST",
			Request: CreateUserRequest{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   30,
			},
			Responses: []specgen.RouteResponse{
				{
					StatusCode: 201,
					Response: UserResponse{
						ID:    1,
						Name:  "John Doe",
						Email: "john@example.com",
					},
				},
			},
		},
	}

	err := specgen.GenerateOpenAPISpec(config, outputFile, routes)
	if err != nil {
		t.Fatalf("GenerateOpenAPISpec failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Output file was not created: %s", outputFile)
	}

	// Read and verify content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	yamlContent := string(content)
	if !strings.Contains(yamlContent, "Test API") {
		t.Error("YAML should contain title 'Test API'")
	}
	if !strings.Contains(yamlContent, "Test API Description") {
		t.Error("YAML should contain description")
	}
	if !strings.Contains(yamlContent, "1.0.0") {
		t.Error("YAML should contain version")
	}
	if !strings.Contains(yamlContent, "/users") {
		t.Error("YAML should contain route path '/users'")
	}
	if !strings.Contains(yamlContent, "post") {
		t.Error("YAML should contain HTTP method 'post'")
	}
}

func TestGenerateOpenAPISpec_WithBearerToken(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "spec.yaml")

	title := "Secure API"
	config := specgen.SpecConfig{
		Title:                   &title,
		WithBearerTokenSecurity: true,
	}

	routes := []specgen.Route{
		{
			Path:    "/protected",
			Method:  "GET",
			Request: struct{}{},
			Responses: []specgen.RouteResponse{
				{
					StatusCode: 200,
					Response:   UserResponse{ID: 1, Name: "Test"},
				},
			},
		},
	}

	err := specgen.GenerateOpenAPISpec(config, outputFile, routes)
	if err != nil {
		t.Fatalf("GenerateOpenAPISpec failed: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	yamlContent := string(content)
	if !strings.Contains(yamlContent, "Bearer") {
		t.Error("YAML should contain Bearer token security")
	}
}

func TestGenerateOpenAPISpec_MultipleRoutes(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "spec.yaml")

	title := "Multi-Route API"
	config := specgen.SpecConfig{
		Title: &title,
	}

	routes := []specgen.Route{
		{
			Path:    "/users",
			Method:  "GET",
			Request: struct{}{},
			Responses: []specgen.RouteResponse{
				{
					StatusCode: 200,
					Response:   []UserResponse{},
				},
			},
		},
		{
			Path:   "/users",
			Method: "POST",
			Request: CreateUserRequest{
				Name:  "Test",
				Email: "test@example.com",
			},
			Responses: []specgen.RouteResponse{
				{
					StatusCode: 201,
					Response:   UserResponse{ID: 1, Name: "Test"},
				},
				{
					StatusCode: 400,
					Response:   ErrorResponse{Message: "Invalid input", Code: "INVALID"},
				},
			},
		},
		{
			Path:   "/users/{id}",
			Method: "GET",
			Request: struct {
				ID int `path:"id"`
			}{ID: 1},
			Responses: []specgen.RouteResponse{
				{
					StatusCode: 200,
					Response:   UserResponse{ID: 1, Name: "Test"},
				},
				{
					StatusCode: 404,
					Response:   ErrorResponse{Message: "Not found", Code: "NOT_FOUND"},
				},
			},
		},
	}

	err := specgen.GenerateOpenAPISpec(config, outputFile, routes)
	if err != nil {
		t.Fatalf("GenerateOpenAPISpec failed: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	yamlContent := string(content)
	// Verify all routes are present
	if !strings.Contains(yamlContent, "/users") {
		t.Error("YAML should contain '/users' path")
	}
	if !strings.Contains(yamlContent, "get") || !strings.Contains(yamlContent, "post") {
		t.Error("YAML should contain both GET and POST methods")
	}
	if !strings.Contains(yamlContent, "201") || !strings.Contains(yamlContent, "400") {
		t.Error("YAML should contain multiple status codes")
	}

}

func TestGenerateOpenAPISpec_MinimalConfig(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "spec.yaml")

	config := specgen.SpecConfig{}

	routes := []specgen.Route{
		{
			Path:    "/health",
			Method:  "GET",
			Request: struct{}{},
			Responses: []specgen.RouteResponse{
				{
					StatusCode: 200,
					Response: struct {
						Status string `json:"status"`
					}{Status: "ok"},
				},
			},
		},
	}

	err := specgen.GenerateOpenAPISpec(config, outputFile, routes)
	if err != nil {
		t.Fatalf("GenerateOpenAPISpec failed: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if len(content) == 0 {
		t.Error("Generated YAML should not be empty")
	}
}

func TestGenerateOpenAPISpec_InvalidMethod(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "spec.yaml")

	config := specgen.SpecConfig{}

	routes := []specgen.Route{
		{
			Path:    "/test",
			Method:  "INVALID_METHOD",
			Request: struct{}{},
			Responses: []specgen.RouteResponse{
				{
					StatusCode: 200,
					Response:   struct{}{},
				},
			},
		},
	}

	err := specgen.GenerateOpenAPISpec(config, outputFile, routes)
	if err == nil {
		t.Error("Expected error for invalid HTTP method, but got nil")
	}
	if !strings.Contains(err.Error(), "failed to create operation context") {
		t.Errorf("Expected error about operation context, got: %v", err)
	}
}

func TestGenerateOpenAPISpec_EmptyRoutes(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "spec.yaml")

	title := "Empty Routes API"
	config := specgen.SpecConfig{
		Title: &title,
	}

	routes := []specgen.Route{}

	err := specgen.GenerateOpenAPISpec(config, outputFile, routes)
	if err != nil {
		t.Fatalf("GenerateOpenAPISpec should succeed with empty routes: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	yamlContent := string(content)
	if !strings.Contains(yamlContent, "Empty Routes API") {
		t.Error("YAML should contain title even with empty routes")
	}
}

type CreateUserWithValidateRequest struct {
	Name     string     `json:"name" validate:"required"`
	Email    string     `json:"email" validate:"required,email,min=3,max=100"`
	Age      int        `json:"age" validate:"required,min=18,max=120"`
	Gender   string     `json:"gender" validate:"required,oneof=male female"`
	Hobbies  []string   `json:"hobbies" validate:"required,min=1,max=10"`
	Birthday *time.Time `json:"birthday" validate:"datetime"`
}

func TestGenerateOpenAPISpec_Validator(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "spec.yaml")

	title := "Test API"
	description := "Test API with validators"
	version := "1.0.0"

	config := specgen.SpecConfig{
		Title:       &title,
		Description: &description,
		Version:     &version,
	}

	birthday := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	routes := []specgen.Route{
		{
			Path:   "/users",
			Method: "POST",
			Request: CreateUserWithValidateRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Age:      30,
				Gender:   "male",
				Hobbies:  []string{"reading", "coding"},
				Birthday: &birthday,
			},
			Responses: []specgen.RouteResponse{
				{
					StatusCode: 201,
					Response:   struct{}{},
				},
			},
		},
	}

	err := specgen.GenerateOpenAPISpec(config, outputFile, routes)
	if err != nil {
		t.Fatalf("GenerateOpenAPISpec failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Output file was not created: %s", outputFile)
	}

	// Read and unmarshal YAML content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	var spec OpenAPISpec
	if err := yaml.Unmarshal(content, &spec); err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	if len(spec.Components.Schemas) == 0 {
		t.Fatalf("No schemas found in spec")
	}

	schema := Schema{}
	for name, s := range spec.Components.Schemas {
		if strings.Contains(name, "CreateUserWithValidateRequest") {
			schema = s
			break
		}
	}

	if schema.Type != "object" {
		t.Fatalf("Schema type should be object, got: %s", schema.Type)
	}

	// Verify required fields
	// Note: This will fail if validator parsing is not integrated into GenerateOpenAPISpec
	requiredFields := []string{"name", "email", "age", "gender", "hobbies"}
	if len(schema.Required) == 0 {
		t.Log("Warning: No required fields found. Validator parsing may not be integrated yet.")
	}
	for _, field := range requiredFields {
		found := false
		for _, req := range schema.Required {
			if req == field {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Required field '%s' not found in schema required array. Found: %v", field, schema.Required)
		}
	}

	// Verify string validators for email field
	emailProp, ok := schema.Properties["email"]
	if !ok {
		t.Fatal("Email property not found in schema")
	}
	// Note: These checks will fail if validator parsing is not integrated
	if emailProp.MinLength == nil {
		t.Error("Email field should have minLength: 3 (validator not applied)")
	} else if *emailProp.MinLength != 3 {
		t.Errorf("Email field should have minLength: 3, got: %d", *emailProp.MinLength)
	}
	if emailProp.MaxLength == nil {
		t.Error("Email field should have maxLength: 100 (validator not applied)")
	} else if *emailProp.MaxLength != 100 {
		t.Errorf("Email field should have maxLength: 100, got: %d", *emailProp.MaxLength)
	}
	if emailProp.Format == nil {
		t.Error("Email field should have format: email (validator not applied)")
	} else if *emailProp.Format != "email" {
		t.Errorf("Email field should have format: email, got: %s", *emailProp.Format)
	}

	// Verify number validators for age field
	ageProp, ok := schema.Properties["age"]
	if !ok {
		t.Fatal("Age property not found in schema")
	}
	// Note: These checks will fail if validator parsing is not integrated
	if ageProp.Minimum == nil {
		t.Error("Age field should have minimum: 18 (validator not applied)")
	} else if *ageProp.Minimum != 18 {
		t.Errorf("Age field should have minimum: 18, got: %f", *ageProp.Minimum)
	}
	if ageProp.Maximum == nil {
		t.Error("Age field should have maximum: 120 (validator not applied)")
	} else if *ageProp.Maximum != 120 {
		t.Errorf("Age field should have maximum: 120, got: %f", *ageProp.Maximum)
	}

	// Verify enum for gender field
	genderProp, ok := schema.Properties["gender"]
	if !ok {
		t.Fatal("Gender property not found in schema")
	}
	// Note: This check will fail if validator parsing is not integrated
	if len(genderProp.Enum) == 0 {
		t.Error("Gender field should have enum (validator not applied)")
	} else {
		enumValues := make(map[string]bool)
		for _, val := range genderProp.Enum {
			if str, ok := val.(string); ok {
				enumValues[str] = true
			}
		}
		if !enumValues["male"] || !enumValues["female"] {
			t.Errorf("Gender enum should contain 'male' and 'female', got: %v", genderProp.Enum)
		}
	}

	// Verify array validators for hobbies field
	hobbiesProp, ok := schema.Properties["hobbies"]
	if !ok {
		t.Fatal("Hobbies property not found in schema")
	}
	// Note: These checks will fail if validator parsing is not integrated
	if hobbiesProp.MinItems == nil {
		t.Error("Hobbies field should have minItems: 1 (validator not applied)")
	} else if *hobbiesProp.MinItems != 1 {
		t.Errorf("Hobbies field should have minItems: 1, got: %d", *hobbiesProp.MinItems)
	}
	if hobbiesProp.MaxItems == nil {
		t.Error("Hobbies field should have maxItems: 10 (validator not applied)")
	} else if *hobbiesProp.MaxItems != 10 {
		t.Errorf("Hobbies field should have maxItems: 10, got: %d", *hobbiesProp.MaxItems)
	}

	// Verify datetime format for birthday field
	birthdayProp, ok := schema.Properties["birthday"]
	if !ok {
		t.Fatal("Birthday property not found in schema")
	}
	// Note: This check will fail if validator parsing is not integrated
	if birthdayProp.Format == nil {
		t.Error("Birthday field should have format: date-time (validator not applied)")
	} else if *birthdayProp.Format != "date-time" {
		t.Errorf("Birthday field should have format: date-time, got: %s", *birthdayProp.Format)
	}
}
