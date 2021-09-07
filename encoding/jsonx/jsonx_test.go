package jsonx

import (
	"testing"
)

func TestEncode(t *testing.T) {
	type user struct {
		Name string
		Age  int
	}

	bytes, err := Encode(user{Name: "alice", Age: 18})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(bytes)
}

func TestDecode(t *testing.T) {
	type user struct {
		Name string
		Age  int
	}

	u := new(user)
	bytes := []byte(`{"Name": 8, "Age": 18}`)
	if err := Decode(bytes, u); err != nil {
		t.Fatal(err)
	}

	t.Log(u)
}
