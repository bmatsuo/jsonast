package main

import (
	"github.com/bitly/go-simplejson"
	"github.com/bmatsuo/jsonast"

	"encoding/json"
	"testing"
	"unicode/utf8"
)

func BenchmarkDecodeRunes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := gists
		pos := 0
		n := len(gists)
		for pos < n {
			c, width := utf8.DecodeRune(p[pos:])
			if c == utf8.RuneError && width == 1 {
				b.Fatal("decode error at %d", pos)
			}
			pos += width
		}
	}
}

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
func BenchmarkEncodingJson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var v interface{}
		err := json.Unmarshal(gists, &v)
		if err != nil {
			b.Fatal("parse error: ", err)
		}
	}
}
