package main

import (
	"log"
	"time"

	"github.com/lutfiandri/go-specgen"
)

type Address struct {
	Street string `json:"street" validate:"required"`
	City   string `json:"city" validate:"required"`
	State  string `json:"state" validate:"required"`
	Zip    string `json:"zip" validate:"required"`
}

type CreateUserRequest struct {
	Name      string     `json:"name" validate:"required"`
	Email     string     `json:"email" validate:"required,email,min=3,max=100"`
	Age       int        `json:"age" validate:"required,min=18,max=120"`
	Gender    string     `json:"gender" validate:"required,oneof=male female"`
	Hobbies   []string   `json:"hobbies" validate:"required,min=1,max=10"`
	Birthday  *time.Time `json:"birthday" validate:"datetime"`
	Addresses []Address  `json:"addresses" validate:"required,min=1"`
}

func main() {
	// Configure the OpenAPI spec
	title := "User Management API"
	description := "A comprehensive API for managing users with CRUD operations"
	version := "1.0.0"

	config := specgen.SpecConfig{
		Title:                   &title,
		Description:             &description,
		Version:                 &version,
		WithBearerTokenSecurity: true,
	}

	// Define routes
	routes := []specgen.Route{
		{
			Tags:        []string{"users"},
			Summary:     "Create a new user",
			Description: "Create a new user with the provided information",
			Path:        "/users",
			Method:      "POST",
			Request:     CreateUserRequest{},
		},
	}

	// Generate the OpenAPI spec
	outputFile := "example/with_validator/openapi.yaml"
	if err := specgen.GenerateOpenAPISpec(config, outputFile, routes); err != nil {
		log.Fatalf("Failed to generate OpenAPI spec: %v", err)
	}

	log.Printf("✅ OpenAPI specification generated successfully!\n")
	log.Printf("📄 Output file: %s\n", outputFile)
}
