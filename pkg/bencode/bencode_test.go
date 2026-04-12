package bencode

import (
	"reflect"
	"testing"
)

func TestUnmarshalInt(t *testing.T) {
	data := []byte("i42e")
	val, err := Unmarshal(data)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if i, ok := val.(Int); !ok {
		t.Fatalf("Expected Int, got %T", val)
	} else if int64(i) != 42 {
		t.Fatalf("Expected 42, got %d", i)
	}
}

func TestUnmarshalString(t *testing.T) {
	data := []byte("4:test")
	val, err := Unmarshal(data)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if s, ok := val.(String); !ok {
		t.Fatalf("Expected String, got %T", val)
	} else if string(s) != "test" {
		t.Fatalf("Expected 'test', got '%s'", s)
	}
}

func TestUnmarshalList(t *testing.T) {
	data := []byte("li1ei2ei3ee")
	val, err := Unmarshal(data)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	list, ok := val.(List)
	if !ok {
		t.Fatalf("Expected List, got %T", val)
	}

	if len(list) != 3 {
		t.Fatalf("Expected 3 elements, got %d", len(list))
	}

	for i, expected := range []int64{1, 2, 3} {
		if l, ok := list[i].(Int); !ok {
			t.Fatalf("Element %d: expected Int, got %T", i, list[i])
		} else if int64(l) != expected {
			t.Fatalf("Element %d: expected %d, got %d", i, expected, l)
		}
	}
}

func TestUnmarshalDict(t *testing.T) {
	data := []byte("d3:foo3:bar4:teste")
	val, err := Unmarshal(data)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	dict, ok := val.(Dict)
	if !ok {
		t.Fatalf("Expected Dict, got %T", val)
	}

	if str, exists := dict["foo"]; !exists {
		t.Fatal("Key 'foo' not found")
	} else if s, ok := str.(String); !ok {
		t.Fatalf("Expected String, got %T", str)
	} else if string(s) != "bar" {
		t.Fatalf("Expected 'bar', got '%s'", s)
	}
}

func TestMarshalInt(t *testing.T) {
	val := Int(42)
	data, err := Marshal(val)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	expected := []byte("i42e")
	if !reflect.DeepEqual(data, expected) {
		t.Fatalf("Expected %s, got %s", expected, data)
	}
}

func TestMarshalString(t *testing.T) {
	val := String("test")
	data, err := Marshal(val)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	expected := []byte("4:test")
	if !reflect.DeepEqual(data, expected) {
		t.Fatalf("Expected %s, got %s", expected, data)
	}
}

func TestMarshalList(t *testing.T) {
	val := List{Int(1), Int(2), Int(3)}
	data, err := Marshal(val)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	expected := []byte("li1ei2ei3ee")
	if !reflect.DeepEqual(data, expected) {
		t.Fatalf("Expected %s, got %s", expected, data)
	}
}

func TestMarshalDict(t *testing.T) {
	val := Dict{
		"foo": String("bar"),
		"baz": Int(42),
	}
	data, err := Marshal(val)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Ключи должны быть отсортированы (baz < foo)
	expected := []byte("d3:bazi42e3:foo3:bare")
	if !reflect.DeepEqual(data, expected) {
		t.Fatalf("Expected %s, got %s", expected, data)
	}
}

func TestRoundTrip(t *testing.T) {
	original := Dict{
		"name":   String("test torrent"),
		"length": Int(1024),
		"list":   List{Int(1), Int(2)},
	}

	data, err := Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	decoded, err := Unmarshal(data)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if !reflect.DeepEqual(original, decoded) {
		t.Fatalf("Round trip failed:\nOriginal: %v\nDecoded:  %v", original, decoded)
	}
}
