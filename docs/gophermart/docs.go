// Package gophermart Code generated by swaggo/swag. DO NOT EDIT
package gophermart

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
        "/api/goods": {
            "post": {
                "description": "Хендлер используется менеджерами для добавления механик вознаграждения за покупки\nПолученные системой расчёта начислений составы чеков проверяются на совпадение с зарегистрированными в данном хендлере вознаграждениями",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Админ"
                ],
                "summary": "Регистрация информации о вознаграждении за товар",
                "parameters": [
                    {
                        "description": "Механика вознаграждения",
                        "name": "reward",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/Reward"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Вознаграждение успешно зарегистрировано",
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
                        "description": "Ключ поиска уже зарегистрирован",
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
        "/api/orders": {
            "post": {
                "description": "Для начисления баллов состав заказа должен быть проверен на совпадения с зарегистрированными записями вознаграждений за товары\nНачисляется сумма совпадений\nПринятый заказ не обязан браться в обработку непосредственно в момент получения запроса",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Админ"
                ],
                "summary": "Регистрация нового совершённого заказа",
                "parameters": [
                    {
                        "description": "Регистрация нового совершённого заказа",
                        "name": "reward",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/Order"
                        }
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Заказ успешно принят в обработку",
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
                        "description": "Заказ уже принят в обработку",
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
        "/api/orders/{number}": {
            "get": {
                "description": "Получение информации о расчёте начислений баллов лояльности за совершённый заказ\nНомером заказа является последовательность цифр произвольной длины.\nНомер заказа может быть проверен на корректность ввода с помощью алгоритма Луна.",
                "consumes": [
                    "text/plain"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Сервис"
                ],
                "summary": "Получение информации о расчёте начислений",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Номером заказа",
                        "name": "number",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешная обработка запроса",
                        "schema": {
                            "$ref": "#/definitions/storagedefault.Order"
                        }
                    },
                    "204": {
                        "description": "Заказ не зарегистрирован в системе расчёта",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "429": {
                        "description": "Превышено количество запросов к сервису",
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
                            "$ref": "#/definitions/Balance"
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
        "/api/user/balance/withdraw": {
            "post": {
                "description": "Хендлер доступен только авторизованному пользователю.\nНомер заказа представляет собой гипотетический номер\nнового заказа пользователя, в счёт оплаты которого списываются баллы.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Заказы"
                ],
                "responses": {
                    "200": {
                        "description": "Успешная обработка запроса",
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
                    "402": {
                        "description": "На счету недостаточно средств",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "422": {
                        "description": "Неверный номер заказа",
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
        "/api/user/withdrawals": {
            "get": {
                "description": "Хендлер доступен только авторизованному пользователю.\nФакты выводов в выдаче должны быть отсортированы по времени вывода от самых старых к самым новым.\nФормат даты — RFC3339.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Заказы"
                ],
                "responses": {
                    "200": {
                        "description": "Успешная обработка запроса",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "204": {
                        "description": "Нет ни одного списания",
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
        "Balance": {
            "type": "object",
            "properties": {
                "current": {
                    "type": "number"
                },
                "withdrawn": {
                    "type": "number"
                }
            }
        },
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
        "Good": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                }
            }
        },
        "Order": {
            "type": "object",
            "properties": {
                "goods": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/Good"
                    }
                },
                "order": {
                    "type": "string"
                }
            }
        },
        "Reward": {
            "type": "object",
            "properties": {
                "match": {
                    "type": "string"
                },
                "reward": {
                    "type": "number"
                },
                "reward_type": {
                    "$ref": "#/definitions/storageaccrual.RewardType"
                }
            }
        },
        "storageaccrual.RewardType": {
            "type": "string",
            "enum": [
                "pt",
                "%"
            ],
            "x-enum-varnames": [
                "RewardTypePt",
                "RewardTypePercent"
            ]
        },
        "storagedefault.Order": {
            "type": "object",
            "properties": {
                "accrual": {
                    "type": "number"
                },
                "number": {
                    "type": "string",
                    "example": "12345678903"
                },
                "status": {
                    "$ref": "#/definitions/storagedefault.OrderStatus"
                }
            }
        },
        "storagedefault.OrderStatus": {
            "type": "string",
            "enum": [
                "NEW",
                "REGISTERED",
                "PROCESSING",
                "INVALID",
                "PROCESSED"
            ],
            "x-enum-varnames": [
                "StatusNew",
                "StatusRegistered",
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
