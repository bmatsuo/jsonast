package main

import (
	"github.com/bitly/go-simplejson"
	"github.com/bmatsuo/jsonast"

	"encoding/json"
	"testing"
)

func BenchmarkSimplejson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var js simplejson.Json
		err := json.Unmarshal(gists, &js)
		if err != nil {
			b.Fatal("parse error: ", err)
		}
	}
}
func BenchmarkJsonast(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := jsonast.ParseJSON(gists)
		if err != nil {
			b.Fatal("parse error: ", err)
		}
	}
}
func BenchmarkEncodingJson(b *testing.B) {}
