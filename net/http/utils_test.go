package http

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeJSONObject(t *testing.T) {
	var s struct {
		Key string `json:"key"`
	}

	req, _ := http.NewRequest(http.MethodGet, "/", bytes.NewBuffer([]byte(`{"key": "foo"}`)))
	err := DecodeJSON(req, &s)
	require.NoError(t, err)
	assert.Equal(t, "foo", s.Key)
}

func TestDecodeJSONStringSlice(t *testing.T) {
	var s []string
	req, _ := http.NewRequest(http.MethodGet, "/", bytes.NewBuffer([]byte(`["foo","bar"]`)))
	err := DecodeJSON(req, &s)
	require.NoError(t, err)
	assert.Equal(t, []string{"foo", "bar"}, s)
}

func TestDecodeJSONObjectSlice(t *testing.T) {
	var s []struct {
		Key string `json:"key"`
	}
	req, _ := http.NewRequest(http.MethodGet, "/", bytes.NewBuffer([]byte(`[{"key": "foo"},{"key": "bar"}]`)))
	err := DecodeJSON(req, &s)
	require.NoError(t, err)
	require.Len(t, s, 2)
	assert.Equal(t, s[0].Key, "foo")
	assert.Equal(t, s[1].Key, "bar")
}

func TestDecodeMalformedJSONObject(t *testing.T) {
	var s struct {
		Key string `json:"key"`
	}

	req, _ := http.NewRequest(http.MethodGet, "/", bytes.NewBuffer([]byte(`{"key: "foo"}`)))
	err := DecodeJSON(req, &s)
	require.Error(t, err)

	req, _ = http.NewRequest(http.MethodGet, "/", bytes.NewBuffer([]byte(`{"key":foo"}`)))
	err = DecodeJSON(req, &s)
	require.Error(t, err)

	req, _ = http.NewRequest(http.MethodGet, "/", bytes.NewBuffer([]byte(`{"key": "foo,"key2":"bar"}`)))
	err = DecodeJSON(req, &s)
	require.Error(t, err)

	b := []byte(`[
		{
			"key1" : "foo,
			"key2": true,
			"key3": [
				{
					"subkey1" : "foo",
					"subkey2" : "bar"
				}
			]
		}
	]`)
	req, _ = http.NewRequest(http.MethodGet, "/", bytes.NewBuffer(b))

	var arr []struct {
		Key1 string `json:"key1"`
		Key2 bool   `json:"key2"`
	}
	err = DecodeJSON(req, &arr)
	require.Error(t, err)
}

func TestParseQueryObject(t *testing.T) {
	var s struct {
		Bool bool   `json:"bool,string"`
		Num  int    `json:"num,string"`
		Str  string `json:"str"`
	}

	req, _ := http.NewRequest(http.MethodGet, "/?bool=true&num=42&str=foobar", http.NoBody)
	err := ParseQuery(req, &s)

	require.NoError(t, err)

	assert.Equal(t, true, s.Bool)
	assert.Equal(t, 42, s.Num)
	assert.Equal(t, "foobar", s.Str)
}
