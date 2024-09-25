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
        "/api/user/balance": {
            "get": {
                "description": "Хендлер доступен только авторизованному пользователю.\nВ ответе должны содержаться данные о текущей сумме баллов лояльности,\nа также сумме использованных за весь период регистрации баллов.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Пользователь"
                ],
                "responses": {
                    "200": {
                        "description": "Успешная обработка запроса",
                        "schema": {
                            "$ref": "#/definitions/Wallet"
                        }
                    },
                    "401": {
                        "description": "Пользователь не авторизован",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/user/login": {
            "post": {
                "description": "Для передачи аутентификационных данных используется механизм cookies",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Пользователь"
                ],
                "summary": "Аутентификация пользователя",
                "parameters": [
                    {
                        "description": "Логин и пароль зарегистрированного пользователя",
                        "name": "creds",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/Creds"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Пользователь успешно зарегистрирован и аутентифицирован",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Неверный формат запроса",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Неверная пара логин/пароль",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/user/orders": {
            "get": {
                "description": "Хендлер доступен только авторизованному пользователю\nНомера заказа в выдаче должны быть отсортированы по времени загрузки от самых старых к самым новым\nФормат даты — RFC3339.",
                "consumes": [
                    "text/plain"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Заказы"
                ],
                "summary": "Получение списка загруженных номеров заказов",
                "responses": {
                    "200": {
                        "description": "Успешная обработка запроса",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/Order"
                            }
                        }
                    },
                    "204": {
                        "description": "Нет данных для ответа",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Пользователь не авторизован",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "Хендлер доступен только аутентифицированным пользователям\nНомером заказа является последовательность цифр произвольной длины",
                "consumes": [
                    "text/plain"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Заказы"
                ],
                "summary": "Загрузка номера заказа",
                "parameters": [
                    {
                        "description": "Трек номер заказа",
                        "name": "number",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Номер заказа уже был загружен этим пользователем",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "202": {
                        "description": "Новый номер заказа принят в обработку",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Неверный формат запроса",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Пользователь не аутентифицирован",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "409": {
                        "description": "Номер заказа уже был загружен другим пользователем",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "422": {
                        "description": "Неверный формат номера заказа",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/user/register": {
            "post": {
                "description": "Для передачи аутентификационных данных используется механизм cookies",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Пользователь"
                ],
                "summary": "Регистрация пользователя",
                "parameters": [
                    {
                        "description": "Логин и пароль не зарегистрированного пользователя",
                        "name": "creds",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/Creds"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Пользователь успешно зарегистрирован и аутентифицирован",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Неверный формат запроса",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "409": {
                        "description": "Логин уже занят",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/ping": {
            "get": {
                "description": "Простая проверка состояния сервера",
                "tags": [
                    "Разное"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "Creds": {
            "type": "object",
            "properties": {
                "login": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "Order": {
            "type": "object",
            "properties": {
                "accrual": {
                    "type": "number"
                },
                "is_withdrawn": {
                    "type": "boolean"
                },
                "number": {
                    "type": "string",
                    "example": "12345678903"
                },
                "status": {
                    "$ref": "#/definitions/storage.OrderStatus"
                },
                "sum": {
                    "type": "number"
                },
                "uploaded_at": {
                    "type": "string",
                    "format": "date-time",
                    "example": "2020-12-10T15:15:45+03:00"
                }
            }
        },
        "Wallet": {
            "type": "object",
            "properties": {
                "balance": {
                    "type": "number"
                },
                "withdrawn": {
                    "type": "number"
                }
            }
        },
        "storage.OrderStatus": {
            "type": "string",
            "enum": [
                "NEW",
                "PROCESSING",
                "INVALID",
                "PROCESSED"
            ],
            "x-enum-varnames": [
                "StatusNew",
                "StatusProcessing",
                "StatusInvalid",
                "StatusProcessed"
            ]
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Gophermart",
	Description:      "Swagger for Gopher Market API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
