# Go-Playground Validator V10 Parsing

| Go Validator Tag | OpenAPI Schema Keyword                 | Applies To          |
| :--------------- | :------------------------------------- | :------------------ |
| `required`       | `required` (in parent object)          | All Types           |
| `min=X`          | `minLength: X`                         | `string`            |
| `max=X`          | `maxLength: X`                         | `string`            |
| `len=X`          | `minLength: X`, `maxLength: X`         | `string`            |
| `min=X`          | `minimum: X`                           | `number`, `integer` |
| `max=X`          | `maximum: X`                           | `number`, `integer` |
| `gte=X`          | `minimum: X`                           | `number`, `integer` |
| `gt=X`           | `exclusiveMinimum: true`, `minimum: X` | `number`, `integer` |
| `lte=X`          | `maximum: X`                           | `number`, `integer` |
| `lt=X`           | `exclusiveMaximum: true`, `maximum: X` | `number`, `integer` |
| `oneof=A B C`    | `enum: [A, B, C]`                      | `string`, `number`  |
| `email`          | `format: email`                        | `string`            |
| `url`            | `format: uri`                          | `string`            |
| `uuid`           | `format: uuid`                         | `string`            |
| `datetime`       | `format: date-time`                    | `string`            |
| `min=X`          | `minItems: X`                          | `array`             |
| `max=X`          | `maxItems: X`                          | `array`             |
