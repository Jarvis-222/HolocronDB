package core

import (
	"log"
	"strings"
	"testing"
)

func TestReadSimpleString(t *testing.T) {

	cases := map[string]string{
		"+OK\r\n":   "OK",
		"+PONG\r\n": "PONG",
	}

	for k, v := range cases {
		val, _ := Decode([]byte(k))

		if val != v {
			log.Printf("The returned string is " + strings.TrimSuffix(val.(string), "\r\n"))
			t.Fail()
		}
	}
}

func TestReadErrorMessage(t *testing.T) {

	cases := map[string]string{
		"-ERR unknown command 'foobar'\r\n": "ERR unknown command 'foobar'",
	}

	for k, v := range cases {
		val, _ := Decode([]byte(k))

		if val != v {
			t.Fail()
		}
	}
}

func TestReadInteger64(t *testing.T) {

	cases := map[string]int{
		":1000\r\n":  1000,
		":-1000\r\n": -1000,
	}

	for k, v := range cases {
		val, _ := Decode([]byte(k))

		if val != v {
			t.Fail()
		}
	}
}

func TestReadBulkString(t *testing.T) {

	cases := map[string]string{
		"$6\r\nfoobar\r\n": "foobar",
		"$0\r\n\r\n":       "",
		"$-1\r\n":          "",
	}

	for k, v := range cases {
		val, _ := Decode([]byte(k))

		if val != v {
			t.Fail()
		}
	}
}

func TestReadArray(t *testing.T) {

	cases := map[string][]interface{}{
		"*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n": {"foo", "bar"},
		"*3\r\n:1\r\n:2\r\n:3\r\n":         {1, 2, 3},
	}

	for k, v := range cases {
		val, _ := Decode([]byte(k))

		if len(val.([]interface{})) != len(v) {
			t.Fail()
		}

		for i := 0; i < len(v); i++ {
			if val.([]interface{})[i] != v[i] {
				t.Fail()
			}
		}
	}
}
