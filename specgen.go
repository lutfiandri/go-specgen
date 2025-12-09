package specgen

import (
	"fmt"
	"os"

	"github.com/swaggest/openapi-go"
	"github.com/swaggest/openapi-go/openapi3"
)

type SpecConfig struct {
	Title                   *string
	Description             *string
	Version                 *string
	WithBearerTokenSecurity bool
}

func GenerateOpenAPISpec(config SpecConfig, outputFile string, routes []Route) error {
	reflector := openapi3.NewReflector()

	if config.Title != nil {
		reflector.Spec.Info.WithTitle(*config.Title)
	}
	if config.Description != nil {
		reflector.Spec.Info.WithDescription(*config.Description)
	}
	if config.Version != nil {
		reflector.Spec.Info.WithVersion(*config.Version)
	}

	if config.WithBearerTokenSecurity {
		reflector.Spec.SetHTTPBearerTokenSecurity("Bearer Auth", "Bearer token authentication", "")
	}

	for _, route := range routes {
		ParseValidatorV10(reflector, route.Request)

		op, err := reflector.NewOperationContext(route.Method, route.Path)
		if err != nil {
			return fmt.Errorf("failed to create operation context: %w", err)
		}

		// TODO: parse params, query, etc. tags from Request struct
		op.AddReqStructure(route.Request)

		for _, response := range route.Responses {
			op.AddRespStructure(response.Response, func(cu *openapi.ContentUnit) {
				cu.HTTPStatus = response.StatusCode
			})
		}

		if err := reflector.AddOperation(op); err != nil {
			return fmt.Errorf("failed to add operation: %w", err)
		}
	}

	yaml, err := reflector.Spec.MarshalYAML()
	if err != nil {
		return fmt.Errorf("failed to marshal yaml spec: %w", err)
	}

	os.WriteFile(outputFile, yaml, 0644)
	return nil
}
