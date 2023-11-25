// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "https://www.duoc.cl/politica-privacidad/",
        "contact": {
            "name": "Matias \"Hinotori\" Canovas",
            "url": "https://github.com/imhinotori/",
            "email": "hello@hinotori.moe"
        },
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/license/mit/"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/attendance": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Get user attendance",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "responses": {
                    "200": {
                        "description": "Successfully retrieved attendance",
                        "schema": {
                            "$ref": "#/definitions/common.Attendance"
                        }
                    },
                    "400": {
                        "description": "Error getting attendance.",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/auth": {
            "post": {
                "description": "Authenticate",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "Username of the user",
                        "name": "username",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Password of the user",
                        "name": "password",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/common.AuthenticationResponse"
                        }
                    },
                    "400": {
                        "description": "Error reading body.",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/grades": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Get user grades",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "responses": {
                    "200": {
                        "description": "Successfully retrieved grades",
                        "schema": {
                            "$ref": "#/definitions/common.GradesCourses"
                        }
                    },
                    "400": {
                        "description": "Error getting grades.",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/schedule": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Get user schedule",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "responses": {
                    "200": {
                        "description": "Successfully retrieved schedule",
                        "schema": {
                            "$ref": "#/definitions/common.Schedule"
                        }
                    },
                    "400": {
                        "description": "Error getting schedule.",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "common.Attendance": {
            "type": "object",
            "properties": {
                "asistenciaAsignaturas": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/common.SubjectAttendance"
                    }
                },
                "codCarrera": {
                    "type": "string"
                },
                "nomCarrera": {
                    "type": "string"
                }
            }
        },
        "common.AuthenticationResponse": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "access_token_expires_in": {
                    "type": "integer"
                },
                "refresh_expires_in": {
                    "type": "integer"
                },
                "refresh_token": {
                    "type": "string"
                }
            }
        },
        "common.Course": {
            "type": "object",
            "properties": {
                "codPlan": {
                    "type": "string"
                },
                "dia": {
                    "type": "string"
                },
                "horaFin": {
                    "type": "string"
                },
                "horaInicio": {
                    "type": "string"
                },
                "nombre": {
                    "type": "string"
                },
                "nombrePlan": {
                    "type": "string"
                },
                "profesor": {
                    "type": "string"
                },
                "sala": {
                    "type": "string"
                },
                "seccion": {
                    "type": "string"
                },
                "sede": {
                    "type": "string"
                },
                "sigla": {
                    "type": "string"
                }
            }
        },
        "common.Day": {
            "type": "object",
            "properties": {
                "dia": {
                    "type": "string"
                },
                "ramos": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/common.Course"
                    }
                }
            }
        },
        "common.Grade": {
            "type": "object",
            "properties": {
                "nota": {
                    "type": "string"
                },
                "texto": {
                    "type": "string"
                }
            }
        },
        "common.GradesCourses": {
            "type": "object",
            "properties": {
                "asignaturas": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/common.Subject"
                    }
                },
                "codCarrera": {
                    "type": "string"
                },
                "nomCarrera": {
                    "type": "string"
                }
            }
        },
        "common.Schedule": {
            "type": "object",
            "properties": {
                "codCarrera": {
                    "type": "string"
                },
                "dias": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/common.Day"
                    }
                },
                "nomCarrera": {
                    "type": "string"
                }
            }
        },
        "common.Subject": {
            "type": "object",
            "properties": {
                "codAsignatura": {
                    "type": "string"
                },
                "examenes": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/common.Grade"
                    }
                },
                "nombre": {
                    "type": "string"
                },
                "parciales": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/common.Grade"
                    }
                },
                "promedio": {
                    "type": "string"
                }
            }
        },
        "common.SubjectAttendance": {
            "type": "object",
            "properties": {
                "cabecera": {
                    "$ref": "#/definitions/common.SubjectAttendanceHeader"
                },
                "detalle": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/common.SubjectAttendanceDetail"
                    }
                }
            }
        },
        "common.SubjectAttendanceDetail": {
            "type": "object",
            "properties": {
                "asistencia": {
                    "type": "string"
                },
                "fechaLarga": {
                    "type": "string"
                }
            }
        },
        "common.SubjectAttendanceHeader": {
            "type": "object",
            "properties": {
                "clasesAsistente": {
                    "type": "string"
                },
                "clasesRealizadas": {
                    "type": "string"
                },
                "codAsignatura": {
                    "type": "string"
                },
                "nomAsignatura": {
                    "type": "string"
                },
                "porcentaje": {
                    "description": "Why...?",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "api-duoc.hinotori.moe",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Duoc Plus API",
	Description:      "Duoc Plus, is a REST API that allows you to access your grades, schedule and attendance from DuocUC.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
