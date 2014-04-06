package jsonast

import (
	yt "github.com/bmatsuo/yup/yuptype"

	"fmt"
	"testing"
)

func TestParseEmpty(t *testing.T) {
	roots, err := Parse(nil)
	yt.Nil(t, err)
	yt.Equal(t, 0, len(roots))
}

func TestParsePrimativeSimple(t *testing.T) {
	for i, test := range []struct {
		raw string
		typ Type
	}{
		{`123`, TNumber},
		{`null`, TNull},
		{`""`, TString},
		{`true`, TBoolean},
		{`false`, TBoolean},
	} {
		desc := fmt.Sprintf("test %d: %q", i, test.raw)
		roots, err := Parse([]byte(test.raw))
		yt.Nil(t, err, desc)
		yt.Equal(t, 1, len(roots), desc)
		yt.Equal(t, test.typ, roots[0].Type(), desc)
	}
}

func TestParseObjectSimple(t *testing.T) {
	for i, test := range []struct {
		raw string
		typ Type
	}{
		{`{"":123}`, TNumber},
		{`{"":null}`, TNull},
		{`{"":""}`, TString},
		{`{"":true}`, TBoolean},
		{`{"":false}`, TBoolean},
		{`{"":{}}`, TObject},
		{`{"":[]}`, TArray},
	} {
		desc := fmt.Sprintf("test %d: %q", i, test.raw)
		roots, err := Parse([]byte(test.raw))
		yt.Nil(t, err, desc)
		yt.Equal(t, 1, len(roots), desc)
		yt.Equal(t, TObject, roots[0].Type(), desc)
		children := roots[0].Children()
		yt.Equal(t, 2, len(children), desc)
		yt.Equal(t, TString, children[0].Type(), desc)
		yt.Equal(t, test.typ, children[1].Type(), desc)
	}
}

func TestParseArraySimple(t *testing.T) {
	for i, test := range []struct {
		raw string
		typ Type
	}{
		{`[123]`, TNumber},
		{`[null]`, TNull},
		{`[""]`, TString},
		{`[true]`, TBoolean},
		{`[false]`, TBoolean},
		{`[{}]`, TObject},
		{`[[]]`, TArray},
	} {
		desc := fmt.Sprintf("test %d: %q", i, test.raw)
		roots, err := Parse([]byte(test.raw))
		yt.Nil(t, err, desc)
		yt.Equal(t, 1, len(roots), desc)
		yt.Equal(t, TArray, roots[0].Type(), desc)
		children := roots[0].Children()
		yt.Equal(t, 1, len(children), desc)
		yt.Equal(t, test.typ, children[0].Type(), desc)
	}
}
