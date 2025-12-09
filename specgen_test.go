package specgen_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lutfiandri/go-specgen"
)

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
