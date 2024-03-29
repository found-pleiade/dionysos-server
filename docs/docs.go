// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "API Support",
            "url": "https://github.com/Brawdunoir/dionysos-server/issues"
        },
        "license": {
            "name": "GNU General Public License v3.0",
            "url": "https://www.gnu.org/licenses/gpl-3.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/rooms": {
            "post": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Rooms"
                ],
                "summary": "Creates a room.",
                "parameters": [
                    {
                        "description": "Room object",
                        "name": "room",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.RoomUpdate"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Room created",
                        "schema": {
                            "$ref": "#/definitions/utils.CreateResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "User not authorized",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Invalid user in auth method",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/rooms/{id}": {
            "get": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Rooms"
                ],
                "summary": "Gets a room.",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Room ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Room"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "User not authorized",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Room not found or invalid user in auth method",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            },
            "patch": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Rooms"
                ],
                "summary": "Updates a room.",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Room ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Room object",
                        "name": "room",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.RoomUpdate"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "User not authorized",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Room not found or invalid user in auth method",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/rooms/{id}/connect": {
            "patch": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Rooms"
                ],
                "summary": "Connects a user to a room.",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Room ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "User not authorized",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Room not found or invalid user in auth method",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "409": {
                        "description": "User already in room",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/rooms/{id}/disconnect": {
            "patch": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Rooms"
                ],
                "summary": "Disconnects a user from a room.",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Room ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "User not authorized",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Room not found or invalid user in auth method",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/rooms/{id}/kick/{userid}": {
            "patch": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Rooms"
                ],
                "summary": "Kicks a user from a room.",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Room ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "userid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "User not authorized",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Room not found or invalid user in auth method",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/rooms/{id}/stream": {
            "get": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "This endpoint is used to subscribe to a SSE stream for a given room.\nThe stream will send an event when a room is updated.\nA room is updated when a user connects or disconnects from it, or when we have a owner change, and so on.",
                "produces": [
                    "text/event-stream"
                ],
                "tags": [
                    "Rooms",
                    "SSE"
                ],
                "summary": "SSE stream of a room for any updates.",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Room ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Send \\\"RoomUpdate\\\" event each time room is updated. Send 200 when stream is closed"
                    },
                    "401": {
                        "description": "User not authorized",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Room not found or invalid user in auth method",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/users": {
            "post": {
                "description": "Creates a user. You will need to use BasicAuth to authenticate with the created user, using its ID and password produced by this endpoint.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Creates a user. Needed for further request during authentication.",
                "parameters": [
                    {
                        "description": "User object",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.UserUpdate"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "User created",
                        "schema": {
                            "$ref": "#/definitions/utils.CreateResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/users/{id}": {
            "get": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Gets a user.",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "User not authorized",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "User not found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Deletes a user. Should be used when disconnecting a user.",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "User not authorized",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "User not found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            },
            "patch": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Updates a user.",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "User object",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.UserUpdate"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "User not authorized",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "User not found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/version": {
            "get": {
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Misc"
                ],
                "summary": "Get the api version.",
                "responses": {
                    "200": {
                        "description": "Version as x.x.x",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Version has not been set correctly on the server side",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Room": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "name": {
                    "type": "string",
                    "maxLength": 20,
                    "minLength": 2,
                    "example": "BirthdayParty"
                },
                "ownerID": {
                    "type": "integer"
                },
                "users": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.User"
                    }
                }
            }
        },
        "models.RoomUpdate": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string",
                    "maxLength": 20,
                    "minLength": 2,
                    "example": "BirthdayParty"
                }
            }
        },
        "models.User": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string",
                    "maxLength": 20,
                    "minLength": 2,
                    "example": "Diablox9"
                }
            }
        },
        "models.UserUpdate": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string",
                    "maxLength": 20,
                    "minLength": 2,
                    "example": "Diablox9"
                }
            }
        },
        "utils.CreateResponse": {
            "type": "object",
            "properties": {
                "password": {
                    "description": "Only set for user creation",
                    "type": "string"
                },
                "uri": {
                    "type": "string"
                }
            }
        },
        "utils.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BasicAuth": {
            "type": "basic"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Dionysos",
	Description:      "API instance for the Dionysos client application.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
