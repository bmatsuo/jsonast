package jsonast

import (
	"fmt"
)

// https://developer.github.com/v3/gists/
func ExampleParseJSON_gist() {
	js, err := ParseJSON([]byte(`[
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
	]`))
	if err != nil {
		panic(err)
	}
	jsLogin, err := js.Get([]Selector{
		{Index: 0},
		{Key: "user"},
		{Key: "login"},
	})
	if err != nil {
		fmt.Println("error: ", err)
	}
	fmt.Println(jsLogin.String())

	// Output:
	// octocat <nil>
}
