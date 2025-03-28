// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/users/register": {
            "post": {
                "description": "Register a new user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Register",
                "operationId": "do-translate",
                "parameters": [
                    {
                        "description": "Registration payload",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.registerRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1.registerResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.ErrResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/response.ErrResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "response.ErrResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "v1.registerRequest": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "email@example.com"
                },
                "password": {
                    "type": "string",
                    "example": "y0urP@ssw0rd"
                }
            }
        },
        "v1.registerResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string",
                    "example": "d13a75e2-3d21-4e57-9dc0-3a7f5bee4c25"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Universal authorization service API",
	Description:      "Authorization, registration, etc...",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
