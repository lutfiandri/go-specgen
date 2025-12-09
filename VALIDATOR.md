# Validator Tags Reference

This document describes the validator tags supported by go-specgen and how they map to OpenAPI schema keywords. These tags follow the [go-playground/validator](https://github.com/go-playground/validator) v10 format and are parsed from the `validate` struct tag.

## General Validators

| Go Validator Tag | OpenAPI Schema Keyword        | Description                                                          |
| :--------------- | :---------------------------- | :------------------------------------------------------------------- |
| `required`       | `required` (in parent object) | The field must be present in the request body. Applies to all types. |

## String Validators

| Go Validator Tag | OpenAPI Schema Keyword         | Description                                                   |
| :--------------- | :----------------------------- | :------------------------------------------------------------ |
| `min=X`          | `minLength: X`                 | Minimum length of the string.                                 |
| `max=X`          | `maxLength: X`                 | Maximum length of the string.                                 |
| `len=X`          | `minLength: X`, `maxLength: X` | Length must be exactly X (sets both minLength and maxLength). |
| `oneof=A B C`    | `enum: [A, B, C]`              | Value must be one of the specified options.                   |
| `email`          | `format: email`                | Must be a valid email format.                                 |
| `url`            | `format: uri`                  | Must be a valid URI (URL).                                    |
| `uuid`           | `format: uuid`                 | Must be a valid UUID.                                         |
| `datetime`       | `format: date-time`            | Must be a date and time string adhering to RFC 3339.          |

## Number/Integer Validators

| Go Validator Tag | OpenAPI Schema Keyword                 | Description                                       |
| :--------------- | :------------------------------------- | :------------------------------------------------ |
| `min=X`          | `minimum: X`                           | Minimum numeric value (inclusive).                |
| `max=X`          | `maximum: X`                           | Maximum numeric value (inclusive).                |
| `gte=X`          | `minimum: X`                           | Greater Than or Equal to X (equivalent to `min`). |
| `gt=X`           | `exclusiveMinimum: true`, `minimum: X` | Strictly Greater Than X.                          |
| `lte=X`          | `maximum: X`                           | Less Than or Equal to X (equivalent to `max`).    |
| `lt=X`           | `exclusiveMaximum: true`, `maximum: X` | Strictly Less Than X.                             |
| `oneof=A B C`    | `enum: [A, B, C]`                      | Value must be one of the specified options.       |

## Array Validators

| Go Validator Tag | OpenAPI Schema Keyword | Description                           |
| :--------------- | :--------------------- | :------------------------------------ |
| `min=X`          | `minItems: X`          | Minimum number of items in the array. |
| `max=X`          | `maxItems: X`          | Maximum number of items in the array. |
