package jsonast

import (
	"fmt"
)

type Selector struct {
	Index int
	Key   string
}

type JSON struct {
	ast *ASTNode
	err error
	a   []JSON
	m   map[string]*JSON
}

func (js *JSON) Type() Type {
	return js.ast.Type()
}

func (js *JSON) String() (string, error) {
	if js.Type() != TString {
		return "", fmt.Errorf("not a string")
	}
	return js.ast.String(), nil
}

func (js *JSON) Float64() (float64, error) {
	if js.Type() != TNumber {
		return 0, fmt.Errorf("not a number")
	}
	return js.ast.Float64()
}

func (js *JSON) Int64() (float64, error) {
	if js.Type() != TNumber {
		return 0, fmt.Errorf("not a number")
	}
	return js.ast.Float64()
}

func (js *JSON) Bool() (bool, error) {
	if js.Type() != TNumber {
		return false, fmt.Errorf("not a bool")
	}
	return js.ast.Bool(), nil
}

// panics if js.Type() is not TArray or TObject
func (js *JSON) Len() (int, error) {
	if js.err != nil {
		return 0, js.err
	}
	typ := js.Type()
	if typ == TObject {
		if js.m == nil {
			return len(js.ast.Children()) / 2, nil
		}
		return len(js.m), nil
	} else if typ == TArray {
		if js.a == nil {
			return len(js.ast.Children()), nil
		}
		return len(js.a), nil
	} else {
		panic("invalid type to call Len()")
	}
}

func (js *JSON) Get(path []Selector) (*JSON, error) {
	_js := js
	var _path []Selector
	for i, sel := range path {
		if sel.Index != 0 && sel.Key != "" {
			_js = &JSON{
				ast: _js.ast,
				err: fmt.Errorf("selector %d: both Index and Key given", i),
			}
			break
		}
		typ := _js.Type()
		if typ == TObject {
			if sel.Index != 0 {
				_js = &JSON{ast: _js.ast, err: fmt.Errorf("not an array")}
				break
			}
			var err error
			_js, err = _js.getk(sel.Key, _path)
			if err != nil {
				_js = &JSON{ast: _js.ast, err: err}
				break
			}
		} else if typ == TArray {
			if sel.Key != "" {
				_js = &JSON{ast: _js.ast, err: fmt.Errorf("not an object")}
				break
			}
			var err error
			_js, err = _js.geti(sel.Index, _path)
			if err != nil {
				_js = &JSON{ast: _js.ast, err: err}
				break
			}
		} else {
			_js = &JSON{ast: _js.ast, err: fmt.Errorf("invalid type")}
			break
		}
		_path = append(_path, sel)
	}
	return _js, _js.err
}

func (js *JSON) getk(k string, path []Selector) (*JSON, error) {
	if js.err != nil {
		return js, js.err
	}
	if js.ast.Type() != TObject {
		err := fmt.Errorf("not an object: %v", path)
		_js := &JSON{
			ast: js.ast,
			err: err,
		}
		return _js, err
	}
	if js.m == nil {
		err := js.build()
		if err != nil {
			_js := &JSON{
				ast: js.ast,
				err: err,
			}
			return _js, err
		}
	}
	_js, ok := js.m[k]
	if !ok {
		err := fmt.Errorf("not found: %q", k)
		_js = &JSON{
			ast: js.ast,
			err: err,
		}
		return _js, err
	}
	return _js, nil
}

func (js *JSON) geti(i int, path []Selector) (*JSON, error) {
	if js.err != nil {
		return js, js.err
	}
	if js.ast.Type() != TArray {
		err := fmt.Errorf("not an array: %v", path)
		_js := &JSON{
			ast: js.ast,
			err: err,
		}
		return _js, err
	}
	if js.a == nil {
		err := js.build()
		if err != nil {
			_js := &JSON{
				ast: js.ast,
				err: err,
			}
			return _js, err
		}
	}
	if i < 0 || i >= len(js.a) {
		err := fmt.Errorf("index out of range: %d", i)
		_js := &JSON{
			ast: js.ast,
			err: err,
		}
		return _js, err
	}
	return &js.a[i], nil
}

func (js *JSON) build() error {
	if js.ast.Type() == TObject {
		cs := js.ast.Children()
		n := len(cs)
		js.m = make(map[string]*JSON, n/2)
		if n%2 == 1 {
			return fmt.Errorf("invalid object: key without value")
		}
		for i := 0; i < n; i += 2 {
			kast, vast := cs[i], cs[i+1]
			if kast.Type() != TString {
				return fmt.Errorf("invalid object: non-string key")
			}
			js.m[kast.String()] = &JSON{ast: vast}
		}
	}
	if js.ast.Type() == TArray {
		cs := js.ast.Children()
		js.a = make([]JSON, len(cs))
		for i := range cs {
			js.a[i].ast = cs[i]
		}
	}
	return nil
}

func ParseMultiJSON(p []byte) ([]JSON, error) {
	roots, err := Parse(p)
	if len(roots) == 0 {
		return nil, err
	}
	js := make([]JSON, len(roots))
	for i := range roots {
		js[i].ast = roots[i]
	}
	return js, err
}

func ParseJSON(p []byte) (JSON, error) {
	roots, err := Parse(p)
	if err != nil {
		return JSON{}, err
	}
	if len(roots) == 0 {
		return JSON{}, fmt.Errorf("no json value found")
	}
	if len(roots) > 1 {
		return JSON{}, fmt.Errorf("multiple json values found")
	}
	return JSON{ast: roots[0]}, nil
}

func (js *JSON) UnmarshalJSON(p []byte) error {
	if js == nil {
		return fmt.Errorf("nil JSON")
	}
	var err error
	*js, err = ParseJSON(p)
	return err
}

func (js JSON) MarshalJSON() ([]byte, error) {
	if js.ast == nil {
		return nil, fmt.Errorf("unitialized value")
	}
	return js.ast.JSON(nil), nil
}
