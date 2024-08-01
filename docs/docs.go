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
        "/keep": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "플레이리스트에 있는 노래들을 가져온다",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Playlist"
                ],
                "summary": "플레이리스트에 노래를 가져온다",
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
                                                "$ref": "#/definitions/handler.PlaylistAddResponse"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "노래들을 하나씩 플레이리스트에 추가한 후 적용된 플레이리스트의 노래들을 리턴한다",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Playlist"
                ],
                "summary": "플레이리스트에 노래를 추가한다",
                "parameters": [
                    {
                        "description": "노래 리스트",
                        "name": "PlaylistAddRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.PlaylistAddRequest"
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
                                                "$ref": "#/definitions/handler.PlaylistAddResponse"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "노래들을 하나씩 플레이리스트에서 삭제한다",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Playlist"
                ],
                "summary": "플레이리스트에 노래를 제거한다",
                "parameters": [
                    {
                        "description": "노래 리스트",
                        "name": "SongDeleteFromPlaylistRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.SongDeleteFromPlaylistRequest"
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
                                            "$ref": "#/definitions/handler.PlaylistAddResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/member": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "사용자 정보 조회",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Member"
                ],
                "summary": "Member의 정보를 가져온다",
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
                                            "$ref": "#/definitions/handler.MemberResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/member/login": {
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
        "/member/logout": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "멤버 회원 로그아웃",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Member"
                ],
                "summary": "멤버 회원 로그아웃",
                "responses": {
                    "200": {
                        "description": "성공",
                        "schema": {
                            "$ref": "#/definitions/pkg.BaseResponseStruct"
                        }
                    }
                }
            }
        },
        "/member/reissue": {
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
        },
        "/member/withdraw": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "멤버 회원 탈퇴",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Member"
                ],
                "summary": "멤버 회원 탈퇴",
                "responses": {
                    "200": {
                        "description": "성공",
                        "schema": {
                            "$ref": "#/definitions/pkg.BaseResponseStruct"
                        }
                    }
                }
            }
        },
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
        "/recommend/home/songs": {
            "get": {
                "description": "앨범 이미지와 함께 노래를 추천",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Recommendation"
                ],
                "summary": "노래 추천 5곡",
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
                                                "$ref": "#/definitions/handler.homeSongResponse"
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
        "/version": {
            "get": {
                "description": "등록되어 있는 모든 버전 확인 가능",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "App Version"
                ],
                "summary": "모든 버전 확인",
                "responses": {
                    "200": {
                        "description": "성공\" {object} pkg.BaseResponseStruct{data=[]versionResponse} \"성공"
                    }
                }
            }
        },
        "/version/check": {
            "post": {
                "description": "헤더에 플랫폼 정보를 포함하고, request body 앱의 버전을 보내면, 최신 버전인지 여부와 강제 업데이트 필요 여부를 응답",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "App Version"
                ],
                "summary": "버전 확인",
                "parameters": [
                    {
                        "description": "현재 앱 버전 정보",
                        "name": "version",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.versionCheckRequest"
                        }
                    }
                ],
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
        "/version/update": {
            "post": {
                "description": "새로운 버전이 나왔을때 버전을 추가할 수 있음 (플랫폼(ios, android), 버전, 이전 버전들을 강제 업데이트 할지 여부)",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "App Version"
                ],
                "summary": "버전 추가",
                "parameters": [
                    {
                        "description": "등록 버전 정보",
                        "name": "version",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.latestVersionUpdateRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "성공"
                    }
                }
            }
        }
    },
    "definitions": {
        "handler.LoginRequest": {
            "type": "object",
            "properties": {
                "birthYear": {
                    "type": "string"
                },
                "gender": {
                    "type": "string"
                },
                "idToken": {
                    "type": "string"
                },
                "provider": {
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
        "handler.MemberResponse": {
            "type": "object",
            "properties": {
                "birthYear": {
                    "type": "integer"
                },
                "email": {
                    "type": "string"
                },
                "gender": {
                    "type": "string"
                },
                "nickname": {
                    "type": "string"
                }
            }
        },
        "handler.PlaylistAddRequest": {
            "type": "object",
            "properties": {
                "songNumbers": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                }
            }
        },
        "handler.PlaylistAddResponse": {
            "type": "object",
            "properties": {
                "singerName": {
                    "type": "string"
                },
                "songId": {
                    "type": "integer"
                },
                "songName": {
                    "type": "string"
                },
                "songNumber": {
                    "type": "integer"
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
        "handler.SongDeleteFromPlaylistRequest": {
            "type": "object",
            "properties": {
                "songNumbers": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
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
        "handler.homeSongResponse": {
            "type": "object",
            "properties": {
                "album": {
                    "type": "string"
                },
                "singerName": {
                    "type": "string"
                },
                "songId": {
                    "type": "integer"
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
        "handler.latestVersionUpdateRequest": {
            "type": "object",
            "properties": {
                "forceUpdate": {
                    "type": "boolean"
                },
                "platform": {
                    "type": "string"
                },
                "version": {
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
                "isKeep": {
                    "type": "boolean"
                },
                "singerName": {
                    "type": "string"
                },
                "songId": {
                    "type": "integer"
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
                "songId": {
                    "type": "integer"
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
                "songNumbers": {
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
        "handler.versionCheckRequest": {
            "type": "object",
            "properties": {
                "version": {
                    "type": "string"
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
