package main

import (
	"log"

	"github.com/lutfiandri/go-specgen"
)

// Example request and response types
type CreateUserRequest struct {
	Name  string `json:"name" required:"true"`
	Email string `json:"email" required:"true" format:"email"`
	Age   int    `json:"age" minimum:"18" maximum:"120"`
}

type UpdateUserRequest struct {
	Name  *string `json:"name"`
	Email *string `json:"email" format:"email"`
	Age   *int    `json:"age" minimum:"18" maximum:"120"`
}

type UserResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Age       int    `json:"age"`
	CreatedAt string `json:"created_at"`
}

type UsersListResponse struct {
	Users []UserResponse `json:"users"`
	Total int            `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Details string `json:"details,omitempty"`
}

type GetUserParams struct {
	ID int `path:"id" required:"true"`
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
			Summary:     "List all users",
			Description: "Retrieve a paginated list of all users",
			Path:        "/users",
			Method:      "GET",
			Request: struct {
				Page  int `query:"page" default:"1" minimum:"1"`
				Limit int `query:"limit" default:"10" minimum:"1" maximum:"100"`
			}{Page: 1, Limit: 10},
			Responses: []specgen.RouteResponse{
				{
					StatusCode: 200,
					Response: UsersListResponse{
						Users: []UserResponse{},
						Total: 0,
						Page:  1,
						Limit: 10,
					},
				},
				{
					StatusCode: 401,
					Response: ErrorResponse{
						Message: "Unauthorized",
						Code:    "UNAUTHORIZED",
					},
				},
			},
		},
		{
			Tags:        []string{"users"},
			Summary:     "Create a new user",
			Description: "Create a new user with the provided information",
			Path:        "/users",
			Method:      "POST",
			Request: CreateUserRequest{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   30,
			},
			Responses: []specgen.RouteResponse{
				{
					StatusCode: 201,
					Response: UserResponse{
						ID:        1,
						Name:      "John Doe",
						Email:     "john@example.com",
						Age:       30,
						CreatedAt: "2024-01-01T00:00:00Z",
					},
				},
				{
					StatusCode: 400,
					Response: ErrorResponse{
						Message: "Invalid input",
						Code:    "VALIDATION_ERROR",
						Details: "Email format is invalid",
					},
				},
				{
					StatusCode: 401,
					Response: ErrorResponse{
						Message: "Unauthorized",
						Code:    "UNAUTHORIZED",
					},
				},
			},
		},
		{
			Tags:        []string{"users"},
			Summary:     "Get user by ID",
			Description: "Retrieve a specific user by their ID",
			Path:        "/users/{id}",
			Method:      "GET",
			Request: GetUserParams{
				ID: 1,
			},
			Responses: []specgen.RouteResponse{
				{
					StatusCode: 200,
					Response: UserResponse{
						ID:        1,
						Name:      "John Doe",
						Email:     "john@example.com",
						Age:       30,
						CreatedAt: "2024-01-01T00:00:00Z",
					},
				},
				{
					StatusCode: 401,
					Response: ErrorResponse{
						Message: "Unauthorized",
						Code:    "UNAUTHORIZED",
					},
				},
				{
					StatusCode: 404,
					Response: ErrorResponse{
						Message: "User not found",
						Code:    "NOT_FOUND",
					},
				},
			},
		},
		{
			Tags:        []string{"users"},
			Summary:     "Update user",
			Description: "Update an existing user's information",
			Path:        "/users/{id}",
			Method:      "PUT",
			Request: struct {
				ID int `path:"id" required:"true"`
				UpdateUserRequest
			}{
				ID: 1,
				UpdateUserRequest: UpdateUserRequest{
					Name:  stringPtr("Jane Doe"),
					Email: stringPtr("jane@example.com"),
				},
			},
			Responses: []specgen.RouteResponse{
				{
					StatusCode: 200,
					Response: UserResponse{
						ID:        1,
						Name:      "Jane Doe",
						Email:     "jane@example.com",
						Age:       30,
						CreatedAt: "2024-01-01T00:00:00Z",
					},
				},
				{
					StatusCode: 400,
					Response: ErrorResponse{
						Message: "Invalid input",
						Code:    "VALIDATION_ERROR",
					},
				},
				{
					StatusCode: 401,
					Response: ErrorResponse{
						Message: "Unauthorized",
						Code:    "UNAUTHORIZED",
					},
				},
				{
					StatusCode: 404,
					Response: ErrorResponse{
						Message: "User not found",
						Code:    "NOT_FOUND",
					},
				},
			},
		},
		{
			Tags:        []string{"users"},
			Summary:     "Delete user",
			Description: "Delete a user by their ID",
			Path:        "/users/{id}",
			Method:      "DELETE",
			Request: GetUserParams{
				ID: 1,
			},
			Responses: []specgen.RouteResponse{
				{
					StatusCode: 204,
					Response:   struct{}{},
				},
				{
					StatusCode: 401,
					Response: ErrorResponse{
						Message: "Unauthorized",
						Code:    "UNAUTHORIZED",
					},
				},
				{
					StatusCode: 404,
					Response: ErrorResponse{
						Message: "User not found",
						Code:    "NOT_FOUND",
					},
				},
			},
		},
	}

	// Generate the OpenAPI spec
	outputFile := "example/basic/openapi.yaml"
	if err := specgen.GenerateOpenAPISpec(config, outputFile, routes); err != nil {
		log.Fatalf("Failed to generate OpenAPI spec: %v", err)
	}

	log.Printf("✅ OpenAPI specification generated successfully!\n")
	log.Printf("📄 Output file: %s\n", outputFile)
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
