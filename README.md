# go-specgen

> 🚀 Generate OpenAPI (Swagger) specifications from Go structs with ease

[![Go Reference](https://pkg.go.dev/badge/github.com/lutfiandri/go-specgen.svg)](https://pkg.go.dev/github.com/lutfiandri/go-specgen)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-%3E%3D1.25-blue.svg)](https://golang.org/)

**go-specgen** is a Go library for generating OpenAPI 3.0 specifications from your Go structs. Built on top of [swaggest/openapi-go](https://github.com/swaggest/openapi-go), it provides a clean API to define your routes and automatically generate comprehensive API documentation.

## 📦 Installation

```bash
go get github.com/lutfiandri/go-specgen
```

## 🎯 Quick Start

Here's a simple example with two endpoints:

```go
package main

import (
	"log"
	"github.com/lutfiandri/go-specgen"
)

// Request types
type CreateUserRequest struct {
	Name  string `json:"name" required:"true"`
	Email string `json:"email" required:"true" format:"email"`
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

func main() {
	// Configure the OpenAPI spec
	title := "User API"
	description := "A simple user management API"
	version := "1.0.0"

	config := specgen.SpecConfig{
		Title:                   &title,
		Description:             &description,
		Version:                 &version,
		WithBearerTokenSecurity: true,
	}

	// Define your routes
	routes := []specgen.Route{
		{
			Path:   "/users",
			Method: "GET",
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
			Request: CreateUserRequest{},
			Responses: []specgen.RouteResponse{
				{
					StatusCode: 201,
					Response: UserResponse{},
				},
				{
					StatusCode: 400,
					Response: ErrorResponse{}
				},
			},
		},
	}

	// Generate the OpenAPI spec
	if err := specgen.GenerateOpenAPISpec(config, "openapi.yaml", routes); err != nil {
		log.Fatalf("Failed to generate OpenAPI spec: %v", err)
	}

	log.Println("✅ OpenAPI specification generated successfully!")
}
```

Run it:

```bash
go run main.go
```

This will generate an `openapi.yaml` file with your API specification!

## 📚 More Examples

Check out the [example](example) directory for a complete CRUD API example with multiple endpoints, error handling, and path parameters.

## 🛠️ Usage

### SpecConfig

Configure your OpenAPI specification:

```go
config := specgen.SpecConfig{
	Title:                   stringPtr("My API"),
	Description:             stringPtr("API description"),
	Version:                 stringPtr("1.0.0"),
	WithBearerTokenSecurity: true, // Optional: enable Bearer token auth
}
```

### Defining Routes

Each route requires:

- `Path` - The endpoint path (e.g., `/users/{id}`)
- `Method` - HTTP method (GET, POST, PUT, DELETE, etc.)
- `Request` - Request body/params struct
- `Responses` - Array of possible responses with status codes

```go
route := specgen.Route{
	Tags:        []string{"users"},
	Summary:     "Create user",
	Description: "Create a new user",
	Path:        "/users",
	Method:      "POST",
	Request:     CreateUserRequest{},
	Responses: []specgen.RouteResponse{
		{StatusCode: 201, Response: UserResponse{}},
		{StatusCode: 400, Response: ErrorResponse{}},
	},
}
```

## ✅ Validation

go-specgen supports parsing validation tags from the `validate` struct tag, following the [go-playground/validator](https://github.com/go-playground/validator) v10 format. These validators are automatically converted to OpenAPI schema constraints.

### Example

```go
type CreateUserRequest struct {
	Name     string   `json:"name" validate:"required,min=3,max=50"`
	Email    string   `json:"email" validate:"required,email"`
	Age      int      `json:"age" validate:"required,min=18,max=120"`
	Score    float64  `json:"score" validate:"gte=0,lte=100"`
	Gender   string   `json:"gender" validate:"required,oneof=male female"`
	Tags     []string `json:"tags" validate:"min=1,max=10"`
	Password string   `json:"password" validate:"required,len=8"`
}
```

For a complete reference of all supported validator tags, see [VALIDATOR.md](VALIDATOR.md).

## 🗺️ Roadmap

- [x] Generate OpenAPI in YAML
- [ ] Generate OpenAPI in JSON
- [x] Parse request struct that using `github.com/go-playground/validator`
- [ ] Support for query parameters and path parameters parsing
- [ ] Trim and prefix request/response schema names
- [ ] Support for request headers

## 🙏 Acknowledgments

Built with [swaggest/openapi-go](https://github.com/swaggest/openapi-go) - a powerful OpenAPI 3.0 reflection library for Go.
