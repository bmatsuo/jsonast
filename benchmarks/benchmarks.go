package main

import (
	"github.com/bmatsuo/jsonast"

	"flag"
	"log"
	"os"
	"runtime/pprof"
)

func main() {
	_n := flag.Int("n", 10000, "number of parses")
	cpuprofile := flag.String("cpuprofile", "cpu.prof", "cpu profile output path")
	flag.Parse()
	pprofFile, err := os.Create(*cpuprofile)
	if err != nil {
		log.Fatal("unable to open cpuprofile path: %v", err)
	}
	pprof.StartCPUProfile(pprofFile)
	defer pprof.StopCPUProfile()

	n := *_n
	for i := 0; i < n; i++ {
		_, err := jsonast.ParseJSON(gists)
		if err != nil {
			log.Fatal("parse error: ", err)
		}
		_, err = jsonast.ParseJSON(swagger)
		if err != nil {
			log.Fatal("parse error: ", err)
		}
	}
}

// https://github.com/wordnik/swagger-core/blob/master/schemas/api-declaration-schema.json
var swagger = []byte(`
{
	"type": "object",
	"$schema": "http://json-schema.org/draft-04/schema",
	"required": [
	"swaggerVersion",
	"resourcePath",
	"apis",
	"basePath"
	],
	"properties": {
		"apiVersion": {
			"type": "string"
		},
		"basePath": {
			"type": "string"
		},
		"swaggerVersion": {
			"enum": [
			"1.2"
			]
		},
		"consumes": {
			"type": "array",
			"items": {
				"type": "string"
			}
		},
		"produces": {
			"type": "array",
			"items": {
				"type": "string"
			}
		},
		"resourcePath": {
			"type": "string"
		},
		"apis": {
			"type": "array",
			"items": [
			{
				"type": "object",
				"required": [
				"path",
				"operations"
				],
				"properties": {
					"path": {
						"type": "string"
					},
					"operations": {
						"type": "array",
						"items": [
						{
							"type": "object",
							"required": [
							"method",
							"nickname",
							"summary",
							"type"
							],
							"properties": {
								"authorizations": {
									"type": "array",
									"items": {
										"type": "string"
									}
								},
								"method": {
									"type": "string",
									"enum": [
									"GET",
									"PUT",
									"POST",
									"DELETE",
									"OPTIONS",
									"PATCH",
									"LINK"
									]
								},
								"nickname": {
									"type": "string"
								},
								"summary": {
									"type": "string"
								},
								"notes": {
									"type": "string"
								},
								"type": {
									"type": "string"
								},
								"parameters": {
									"type": "array",
									"items": {
										"type": "object",
										"required": [
										"name",
										"paramType",
										"required",
										"type"
										],
										"properties": {
											"allowMultiple": {
												"type": "boolean",
												"enum": [
												true,
												false
												]
											},
											"description": {
												"type": "string"
											},
											"name": {
												"type": "string"
											},
											"paramType": {
												"type": "string",
												"enum": [
												"query",
												"path",
												"body",
												"header"
												]
											},
											"required": {
												"type": "boolean",
												"enum": [
												true,
												false
												]
											},
											"type": {
												"type": "string"
											},
											"items": {
												"anyOf": [
												{
													"$ref": "#"
												},
												{
													"$ref": "#/definitions/schemaArray"
												}
												],
												"default": {}
											}
										}
									}
								},
								"produces": {
									"type": "array",
									"items": {
										"type": "string"
									}
								},
								"responseMessages": {
									"type": "array",
									"items": {
										"type": "object",
										"properties": {
											"code": {
												"type": "number"
											},
											"message": {
												"type": "string"
											}
										}
									}
								}
							}
						}
						]
					}
				}
			}
			]
		},
		"models": {

		}
	}
}
`)

var gists = []byte(`[
	{
		"url": "https://api.github.com/gists/3abd1c18fbe3ff2b97c9",
		"forks_url": "https://api.github.com/gists/90f509bfd7b2a15e6702/forks",
		"commits_url": "https://api.github.com/gists/ae5ba1b108de4e762837/commits",
		"id": "1",
		"description": "description of gist",
		"public": true,
		"user": {
			"login": "octocat",
			"id": 1,
			"avatar_url": "https://github.com/images/error/octocat_happy.gif",
			"gravatar_id": "somehexcode",
			"url": "https://api.github.com/users/octocat",
			"html_url": "https://github.com/octocat",
			"followers_url": "https://api.github.com/users/octocat/followers",
			"following_url": "https://api.github.com/users/octocat/following{/other_user}",
			"gists_url": "https://api.github.com/users/octocat/gists{/gist_id}",
			"starred_url": "https://api.github.com/users/octocat/starred{/owner}{/repo}",
			"subscriptions_url": "https://api.github.com/users/octocat/subscriptions",
			"organizations_url": "https://api.github.com/users/octocat/orgs",
			"repos_url": "https://api.github.com/users/octocat/repos",
			"events_url": "https://api.github.com/users/octocat/events{/privacy}",
			"received_events_url": "https://api.github.com/users/octocat/received_events",
			"type": "User",
			"site_admin": false
		},
		"files": {
			"ring.erl": {
				"size": 932,
				"raw_url": "https://gist.githubusercontent.com/raw/365370/8c4d2d43d178df44f4c03a7f2ac0ff512853564e/ring.erl",
				"type": "text/plain",
				"language": "Erlang"
			}
		},
		"comments": 0,
		"comments_url": "https://api.github.com/gists/7a96cc4d6c6a6d66a662/comments/",
		"html_url": "https://gist.github.com/1",
		"git_pull_url": "git://gist.github.com/1.git",
		"git_push_url": "git@gist.github.com:1.git",
		"created_at": "2010-04-14T02:15:15Z",
		"updated_at": "2011-06-20T11:34:15Z"
	}
	]`)
