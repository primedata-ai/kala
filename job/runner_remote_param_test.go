package job

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestParamsStruct struct {
	Foo           string
	Bar           *TestParamsNestedStruct
	Renamed       int     `json:"changed"`
	AlwaysVisible float64 `json:"AlwaysVisible,omitempty"`
	Zero          string  `json:",omitempty"`
	Approved      bool    `json:"approved,omitempty"`
}

type TestParamsNestedStruct struct {
	AAA int
	BBB string
	CCC bool
}

func TestParamsEncode(t *testing.T) {
	var params Params
	buf := &bytes.Buffer{}

	if mime, err := params.Encode(buf); err != nil || mime != mimeFormURLEncoded || buf.Len() != 0 {
		t.Fatalf("empty params must encode to an empty string. actual is [e:%v] [str:%v] [mime:%v]", err, buf.String(), mime)
	}

	buf.Reset()
	params = Params{}
	params["need_escape"] = "&=+"
	expectedEncoding := "need_escape=%26%3D%2B"

	if mime, err := params.Encode(buf); err != nil || mime != mimeFormURLEncoded || buf.String() != expectedEncoding {
		t.Fatalf("wrong params encode result. expected is '%v'. actual is '%v'. [e:%v] [mime:%v]", expectedEncoding, buf.String(), err, mime)
	}

	buf.Reset()
	data := &TestParamsStruct{
		Foo: "hello, world!",
		Bar: &TestParamsNestedStruct{
			AAA: 1234,
			BBB: "bbb",
			CCC: true,
		},
		Renamed:       123,
		AlwaysVisible: 0,
		Zero:          "", // should be a zero value.
	}
	params = MakeParams(data)
	expectedParams := Params{
		"Foo": "hello, world!",
		"Bar": Params{
			"AAA": 1234,
			"BBB": "bbb",
			"CCC": true,
		},
		"changed":       123, // field name should be changed due to struct field tag.
		"AlwaysVisible": 0.0,
		"approved":      false,
	}

	if params == nil {
		t.Fatalf("make params error.")
	}

	if !assert.ObjectsAreEqualValues(params, expectedParams) {
		t.Fatalf("invalid encoded params. [expected:%#v] [actual:%#v]", expectedParams, params)
	}

	// Test against issue #148 in which a field with nil value causes panic.
	params["nil"] = nil

	mime, err := params.Encode(buf)

	if err != nil || buf.Len() == 0 {
		t.Fatalf("complex encode result is '%v'. [e:%v] [mime:%v]", buf.String(), err, mime)
	}
}