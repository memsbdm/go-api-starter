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
                            "$ref": "#/definitions/http.loginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Access and refresh tokens",
                        "schema": {
                            "$ref": "#/definitions/http.response-http_loginResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized / credentials error",
                        "schema": {
                            "$ref": "#/definitions/http.errorResponse"
                        }
                    },
                    "403": {
                        "description": "Forbidden error",
                        "schema": {
                            "$ref": "#/definitions/http.errorResponse"
                        }
                    },
                    "422": {
                        "description": "Validation error",
                        "schema": {
                            "$ref": "#/definitions/http.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/http.errorResponse"
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
                "parameters": [
                    {
                        "description": "Refresh token request",
                        "name": "refreshTokenRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/http.refreshTokenRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/http.emptyResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized error",
                        "schema": {
                            "$ref": "#/definitions/http.errorResponse"
                        }
                    },
                    "422": {
                        "description": "Validation error",
                        "schema": {
                            "$ref": "#/definitions/http.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/http.errorResponse"
                        }
                    }
                }
            }
        },
        "/v1/auth/refresh": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Generate a new access token and refresh token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Generate a new access token and refresh token",
                "parameters": [
                    {
                        "description": "Refresh token request",
                        "name": "refreshTokenRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/http.refreshTokenRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Access and refresh tokens",
                        "schema": {
                            "$ref": "#/definitions/http.response-http_loginResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized error",
                        "schema": {
                            "$ref": "#/definitions/http.errorResponse"
                        }
                    },
                    "422": {
                        "description": "Validation error",
                        "schema": {
                            "$ref": "#/definitions/http.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/http.errorResponse"
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
                        "name": "registerUserRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/http.registerUserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Created user",
                        "schema": {
                            "$ref": "#/definitions/http.response-http_userResponse"
                        }
                    },
                    "403": {
                        "description": "Forbidden error",
                        "schema": {
                            "$ref": "#/definitions/http.errorResponse"
                        }
                    },
                    "409": {
                        "description": "Duplication error",
                        "schema": {
                            "$ref": "#/definitions/http.errorResponse"
                        }
                    },
                    "422": {
                        "description": "Validation error",
                        "schema": {
                            "$ref": "#/definitions/http.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/http.errorResponse"
                        }
                    }
                }
            }
        },
        "/v1/health": {
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
                        "description": "DB information",
                        "schema": {
                            "$ref": "#/definitions/http.healthResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/http.errorResponse"
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
                            "$ref": "#/definitions/http.response-http_userResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized error",
                        "schema": {
                            "$ref": "#/definitions/http.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/http.errorResponse"
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
                            "$ref": "#/definitions/http.response-http_getUserByIDResponse"
                        }
                    },
                    "400": {
                        "description": "Incorrect User ID",
                        "schema": {
                            "$ref": "#/definitions/http.errorResponse"
                        }
                    },
                    "404": {
                        "description": "Data not found error",
                        "schema": {
                            "$ref": "#/definitions/http.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/http.errorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "http.emptyResponse": {
            "type": "object",
            "properties": {
                "success": {
                    "type": "boolean",
                    "example": true
                }
            }
        },
        "http.errorResponse": {
            "type": "object",
            "properties": {
                "messages": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "Error message 1",
                        " Error message 2"
                    ]
                },
                "success": {
                    "type": "boolean",
                    "example": false
                }
            }
        },
        "http.getUserByIDResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string",
                    "example": "6b947a32-8919-4974-9ef3-048a556b0b75"
                },
                "username": {
                    "type": "string",
                    "example": "john"
                }
            }
        },
        "http.healthResponse": {
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
        "http.loginRequest": {
            "type": "object",
            "required": [
                "password",
                "username"
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
        "http.loginResponse": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "refresh_token": {
                    "type": "string"
                }
            }
        },
        "http.refreshTokenRequest": {
            "type": "object",
            "properties": {
                "refreshToken": {
                    "type": "string"
                }
            }
        },
        "http.registerUserRequest": {
            "type": "object",
            "required": [
                "password",
                "username"
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
        "http.response-http_getUserByIDResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/http.getUserByIDResponse"
                },
                "success": {
                    "type": "boolean"
                }
            }
        },
        "http.response-http_loginResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/http.loginResponse"
                },
                "success": {
                    "type": "boolean"
                }
            }
        },
        "http.response-http_userResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/http.userResponse"
                },
                "success": {
                    "type": "boolean"
                }
            }
        },
        "http.userResponse": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string",
                    "example": "2024-08-15T16:23:33.455225Z"
                },
                "id": {
                    "type": "string",
                    "example": "6b947a32-8919-4974-9ef3-048a556b0b75"
                },
                "is_email_verified": {
                    "type": "boolean",
                    "example": true
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
	Description:      "This is a simple starter API written in Go using net/http, PostgresSQL database, and Redis cache.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
