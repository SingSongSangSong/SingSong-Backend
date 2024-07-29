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
        "/recommend/home": {
            "post": {
                "description": "태그에 해당하는 노래를 추천합니다.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Recommendation"
                ],
                "summary": "노래 추천 by 태그",
                "parameters": [
                    {
                        "description": "태그 목록",
                        "name": "songs",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.homeRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "성공",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/pkg.BaseResponseStruct"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/handler.homeResponse"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/recommend/refresh": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "태그에 해당하는 노래를 새로고침합니다.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Recommendation"
                ],
                "summary": "새로고침 노래 추천",
                "parameters": [
                    {
                        "description": "태그",
                        "name": "songs",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.refreshRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "성공",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/pkg.BaseResponseStruct"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/handler.refreshResponse"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/recommend/songs": {
            "post": {
                "description": "노래 번호 목록을 보내면 유사한 노래들을 추천합니다.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Recommendation"
                ],
                "summary": "노래 추천 by 노래 번호 목록",
                "parameters": [
                    {
                        "description": "노래 번호 목록",
                        "name": "songs",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.songRecommendRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "성공",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/pkg.BaseResponseStruct"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/handler.songRecommendResponse"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/tags": {
            "get": {
                "description": "ssss 태그 목록을 조회합니다.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Tags"
                ],
                "summary": "ssss 태그 목록 가져오기",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/pkg.BaseResponseStruct"
                        }
                    }
                }
            }
        },
        "/user/login": {
            "post": {
                "description": "IdToken을 이용한 회원가입 및 로그인",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Signup and Login"
                ],
                "summary": "회원가입 및 로그인",
                "parameters": [
                    {
                        "description": "idToken 및 Provider",
                        "name": "songs",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "성공",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/pkg.BaseResponseStruct"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/handler.LoginResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/user/reissue": {
            "post": {
                "description": "AccessToken 재발급 및 RefreshToken 재발급 (RTR Refresh Token Rotation)",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Reissue"
                ],
                "summary": "AccessToken RefreshToken 재발급",
                "parameters": [
                    {
                        "description": "accessToken 및 refreshToken",
                        "name": "songs",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.ReissueRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "성공",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/pkg.BaseResponseStruct"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/handler.LoginResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handler.LoginRequest": {
            "type": "object",
            "properties": {
                "IdToken": {
                    "type": "string"
                },
                "Provider": {
                    "type": "string"
                }
            }
        },
        "handler.LoginResponse": {
            "type": "object",
            "properties": {
                "accessToken": {
                    "type": "string"
                },
                "refreshToken": {
                    "type": "string"
                }
            }
        },
        "handler.ReissueRequest": {
            "type": "object",
            "properties": {
                "accessToken": {
                    "type": "string"
                },
                "refreshToken": {
                    "type": "string"
                }
            }
        },
        "handler.homeRequest": {
            "type": "object",
            "properties": {
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "handler.homeResponse": {
            "type": "object",
            "properties": {
                "songs": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/handler.songHomeResponse"
                    }
                },
                "tag": {
                    "type": "string"
                }
            }
        },
        "handler.refreshRequest": {
            "type": "object",
            "properties": {
                "tag": {
                    "type": "string"
                }
            }
        },
        "handler.refreshResponse": {
            "type": "object",
            "properties": {
                "singerName": {
                    "type": "string"
                },
                "songName": {
                    "type": "string"
                },
                "songNumber": {
                    "type": "integer"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "handler.songHomeResponse": {
            "type": "object",
            "properties": {
                "singerName": {
                    "type": "string"
                },
                "songName": {
                    "type": "string"
                },
                "songNumber": {
                    "type": "integer"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "handler.songRecommendRequest": {
            "type": "object",
            "properties": {
                "songs": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                }
            }
        },
        "handler.songRecommendResponse": {
            "type": "object",
            "properties": {
                "singerName": {
                    "type": "string"
                },
                "songName": {
                    "type": "string"
                },
                "songNumber": {
                    "type": "integer"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "pkg.BaseResponseStruct": {
            "type": "object",
            "properties": {
                "data": {},
                "message": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
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
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "싱송생송 API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
