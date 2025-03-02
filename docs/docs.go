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
        "/v1/auth/login": {
            "post": {
                "description": "Authenticate a user account",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Login a user",
                "parameters": [
                    {
                        "description": "Login request",
                        "name": "loginRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_adapters_http_handlers.loginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Login response",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.Response-go-starter_internal_adapters_http_responses_LoginResponse"
                        }
                    },
                    "400": {
                        "description": "Bad request error",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized / credentials error",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    },
                    "422": {
                        "description": "Validation error",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/v1/auth/logout": {
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Logout an authenticated user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Logout an authenticated user",
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.EmptyResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized error",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/v1/auth/password-reset": {
            "post": {
                "description": "Send a password reset email to the user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Send a password reset email",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User's email address",
                        "name": "email",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.EmptyResponse"
                        }
                    },
                    "400": {
                        "description": "Bad request error",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/v1/auth/password-reset/{token}": {
            "get": {
                "description": "Verify a password reset token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Verify a password reset token",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Password reset token",
                        "name": "token",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.EmptyResponse"
                        }
                    },
                    "400": {
                        "description": "Bad request error",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized error / invalid token",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    }
                }
            },
            "patch": {
                "description": "Reset a user's password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Reset a user's password",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Password reset token",
                        "name": "token",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Reset password request",
                        "name": "resetPasswordRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_adapters_http_handlers.resetPasswordRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.EmptyResponse"
                        }
                    },
                    "400": {
                        "description": "Bad request error",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized error / invalid token",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    },
                    "422": {
                        "description": "Validation error",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/v1/auth/register": {
            "post": {
                "description": "Create a new user account",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Register a new user",
                "parameters": [
                    {
                        "description": "Register request",
                        "name": "registerRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_adapters_http_handlers.registerRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Created user",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.Response-go-starter_internal_adapters_http_responses_LoginResponse"
                        }
                    },
                    "400": {
                        "description": "Bad request error",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    },
                    "409": {
                        "description": "Duplication error",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    },
                    "422": {
                        "description": "Validation error",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/v1/health/postgres": {
            "get": {
                "description": "Get database health information",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Health"
                ],
                "summary": "Get database health information",
                "responses": {
                    "200": {
                        "description": "Postgres health information",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.HealthResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/v1/mailer": {
            "get": {
                "description": "Send an example email",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Mail"
                ],
                "summary": "Send an example email",
                "responses": {
                    "200": {
                        "description": "Success"
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/v1/users/me": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Get information of logged-in user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Get authenticated user information",
                "responses": {
                    "200": {
                        "description": "User displayed",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.Response-go-starter_internal_adapters_http_responses_UserResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized error",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/v1/users/me/password": {
            "patch": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Update user password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Update user password",
                "parameters": [
                    {
                        "description": "Update user password request",
                        "name": "updatePasswordRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_adapters_http_handlers.updatePasswordRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.EmptyResponse"
                        }
                    },
                    "400": {
                        "description": "Bad request error",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized error",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    },
                    "422": {
                        "description": "Validation error",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/v1/users/me/verify-email/resend": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Resend user email verification",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Resend user email verification",
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.EmptyResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized error",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    },
                    "409": {
                        "description": "Conflict error / already verified",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/v1/users/me/verify-email/{token}": {
            "get": {
                "description": "Verify user email",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Verify user email",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Verification token",
                        "name": "token",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User displayed",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.EmptyResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized error / invalid token",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    },
                    "409": {
                        "description": "Conflict error / already verified by another user",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/v1/users/{uuid}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Get a user by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Get a user",
                "parameters": [
                    {
                        "type": "string",
                        "format": "uuid",
                        "description": "User ID",
                        "name": "uuid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User displayed",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.Response-go-starter_internal_adapters_http_responses_GetUserByIDResponse"
                        }
                    },
                    "400": {
                        "description": "Incorrect User ID",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Data not found error",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/go-starter_internal_adapters_http_responses.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "go-starter_internal_adapters_http_responses.EmptyResponse": {
            "type": "object",
            "properties": {
                "success": {
                    "type": "boolean",
                    "example": true
                }
            }
        },
        "go-starter_internal_adapters_http_responses.ErrorResponse": {
            "type": "object",
            "properties": {
                "messages": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "Error message 1",
                        "Error message 2"
                    ]
                },
                "success": {
                    "type": "boolean",
                    "example": false
                }
            }
        },
        "go-starter_internal_adapters_http_responses.GetUserByIDResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string",
                    "example": "6b947a32-8919-4974-9ef3-048a556b0b75"
                },
                "name": {
                    "type": "string",
                    "example": "John Doe"
                },
                "username": {
                    "type": "string",
                    "example": "john"
                }
            }
        },
        "go-starter_internal_adapters_http_responses.HealthResponse": {
            "type": "object",
            "properties": {
                "idle": {
                    "type": "string",
                    "example": "1"
                },
                "in_use": {
                    "type": "string",
                    "example": "0"
                },
                "max_idle_closed": {
                    "type": "string",
                    "example": "0"
                },
                "max_lifetime_closed": {
                    "type": "string",
                    "example": "0"
                },
                "message": {
                    "type": "string",
                    "example": "It's healthy'"
                },
                "open_connections": {
                    "type": "string",
                    "example": "1"
                },
                "status": {
                    "type": "string",
                    "example": "up"
                },
                "wait_count": {
                    "type": "string",
                    "example": "0"
                },
                "wait_duration": {
                    "type": "string",
                    "example": "0s"
                }
            }
        },
        "go-starter_internal_adapters_http_responses.LoginResponse": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/go-starter_internal_adapters_http_responses.UserResponse"
                }
            }
        },
        "go-starter_internal_adapters_http_responses.Response-go-starter_internal_adapters_http_responses_GetUserByIDResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/go-starter_internal_adapters_http_responses.GetUserByIDResponse"
                },
                "success": {
                    "type": "boolean"
                }
            }
        },
        "go-starter_internal_adapters_http_responses.Response-go-starter_internal_adapters_http_responses_LoginResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/go-starter_internal_adapters_http_responses.LoginResponse"
                },
                "success": {
                    "type": "boolean"
                }
            }
        },
        "go-starter_internal_adapters_http_responses.Response-go-starter_internal_adapters_http_responses_UserResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/go-starter_internal_adapters_http_responses.UserResponse"
                },
                "success": {
                    "type": "boolean"
                }
            }
        },
        "go-starter_internal_adapters_http_responses.UserResponse": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string",
                    "example": "2024-08-15T16:23:33.455225Z"
                },
                "email": {
                    "type": "string",
                    "example": "john@example.com"
                },
                "id": {
                    "type": "string",
                    "example": "6b947a32-8919-4974-9ef3-048a556b0b75"
                },
                "is_email_verified": {
                    "type": "boolean",
                    "example": true
                },
                "name": {
                    "type": "string",
                    "example": "John Doe"
                },
                "updated_at": {
                    "type": "string",
                    "example": "2025-01-15T14:29:33.455225Z"
                },
                "username": {
                    "type": "string",
                    "example": "john"
                }
            }
        },
        "internal_adapters_http_handlers.loginRequest": {
            "type": "object",
            "required": [
                "password"
            ],
            "properties": {
                "password": {
                    "type": "string",
                    "example": "secret123"
                },
                "username": {
                    "type": "string",
                    "example": "john"
                }
            }
        },
        "internal_adapters_http_handlers.registerRequest": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "john@example.com"
                },
                "name": {
                    "type": "string",
                    "maxLength": 50,
                    "example": "John Doe"
                },
                "password": {
                    "type": "string",
                    "minLength": 8,
                    "example": "secret123"
                },
                "username": {
                    "type": "string",
                    "maxLength": 15,
                    "minLength": 4,
                    "example": "john"
                }
            }
        },
        "internal_adapters_http_handlers.resetPasswordRequest": {
            "type": "object",
            "required": [
                "password",
                "password_confirmation"
            ],
            "properties": {
                "password": {
                    "type": "string",
                    "minLength": 8,
                    "example": "secret123"
                },
                "password_confirmation": {
                    "type": "string",
                    "example": "secret123"
                }
            }
        },
        "internal_adapters_http_handlers.updatePasswordRequest": {
            "type": "object",
            "required": [
                "password",
                "password_confirmation"
            ],
            "properties": {
                "password": {
                    "type": "string",
                    "minLength": 8,
                    "example": "secret123"
                },
                "password_confirmation": {
                    "type": "string",
                    "example": "secret123"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "description": "Type \"Bearer\" followed by a space and the access token.",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Go Starter API",
	Description:      "This is a simple starter API written in Go using net/http, Postgres database, and Redis cache.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
