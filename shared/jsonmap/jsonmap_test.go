package jsonmap

import (
	"reflect"
	"testing"
)

func TestNewJsonMap(t *testing.T) {
	jm := NewJsonMap()
	if jm == nil {
		t.Fatal("NewJsonMap() returned nil")
	}
	if len(jm) != 0 {
		t.Fatalf("expected empty map, got len=%d", len(jm))
	}
}

func TestJsonMapString(t *testing.T) {
	jm := JsonMap{"foo": "bar", "num": 42}
	str := jm.String()
	if str == "" {
		t.Fatal("String() returned empty string")
	}
	if !(contains(str, "foo") && contains(str, "bar") && contains(str, "num")) {
		t.Fatalf("String() output missing expected keys/values: %s", str)
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && (reflect.ValueOf(s).String() == substr || len(s) >= len(substr) && (s[0:len(substr)] == substr || contains(s[1:], substr)))
}

func TestFromString(t *testing.T) {
	jsonStr := `{"a":1,"b":"two"}`
	jm, err := FromString(jsonStr)
	if err != nil {
		t.Fatalf("FromString error: %v", err)
	}
	if jm["a"] != float64(1) || jm["b"] != "two" {
		t.Fatalf("unexpected map contents: %v", jm)
	}
}

func TestAssign(t *testing.T) {
	jm := JsonMap{"a": 1, "b": 2}
	src1 := JsonMap{"b": 3, "c": 4}
	src2 := JsonMap{"d": 5}
	jm.Assign(src1, src2)
	if jm["a"] != 1 || jm["b"] != 3 || jm["c"] != 4 || jm["d"] != 5 {
		t.Fatalf("Assign failed: %v", jm)
	}
}

func TestGet(t *testing.T) {
	jm := JsonMap{"foo": "bar", "num": float64(42)}
	val, ok := Get[string](jm, "foo")
	if !ok || val != "bar" {
		t.Fatalf("Get failed for 'foo': %v, %v", val, ok)
	}
	fval, ok := Get[float64](jm, "num")
	if !ok || fval != 42 {
		t.Fatalf("Get failed for 'num': %v, %v", fval, ok)
	}
	ival, ok := Get[int](jm, "num")
	if ok {
		t.Fatalf("Get should fail for int type, got: %v", ival)
	}
}

func TestGetOrDefaultMethod(t *testing.T) {
	jm := JsonMap{"foo": "bar"}
	if jm.GetOrDefault("foo", "baz") != "bar" {
		t.Fatal("GetOrDefault did not return value")
	}
	if jm.GetOrDefault("missing", "baz") != "baz" {
		t.Fatal("GetOrDefault did not return default")
	}
}

func TestGetOrDefaultGeneric(t *testing.T) {
	jm := JsonMap{"foo": "bar", "num": float64(42)}
	v := GetOrDefault[string](jm, "foo", "baz")
	if v != "bar" {
		t.Fatalf("GetOrDefault generic failed: %v", v)
	}
	v2 := GetOrDefault[float64](jm, "num", 99)
	if v2 != 42 {
		t.Fatalf("GetOrDefault generic failed for float64 type, got: %v", v2)
	}
	v3 := GetOrDefault[int](jm, "num", 99)
	if v3 != 99 {
		t.Fatalf("GetOrDefault generic should fail for int type, got: %v", v3)
	}
	v4 := GetOrDefault[string](jm, "missing", "default")
	if v4 != "default" {
		t.Fatalf("GetOrDefault generic did not return default: %v", v4)
	}
}
