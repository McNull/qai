package jsonmap

import (
	"encoding/json"
	"io"
	"maps"
	"os"
	"path/filepath"
)

type JsonMap map[string]any

func NewJsonMap() JsonMap {
	return make(JsonMap)
}

func (jm JsonMap) String() string {
	b, err := json.MarshalIndent(jm, "", "  ")
	if err != nil {
		return ""
	}
	return string(b)
}

func FromString(txt string) (JsonMap, error) {
	var m JsonMap
	err := json.Unmarshal([]byte(txt), &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func FromFile(filePath string) (JsonMap, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var jm JsonMap
	err = json.Unmarshal(data, &jm)
	if err != nil {
		return nil, err
	}

	return jm, nil
}

func ToFile(jm JsonMap, filePath string) error {
	dirPath := filepath.Dir(filePath)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			return err
		}
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(jm, "", "  ")
	if err != nil {
		return err
	}

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}

func Get[T any](jm JsonMap, key string) (T, bool) {
	if value, ok := jm[key]; ok {
		if typedValue, ok := value.(T); ok {
			return typedValue, true
		}
	}
	var zero T
	return zero, false
}

func (jm JsonMap) GetOrDefault(key string, defaultValue any) any {
	if value, ok := jm[key]; ok {
		return value
	}
	return defaultValue
}

func GetOrDefault[T any](jm JsonMap, key string, defaultValue T) T {
	if value, ok := jm[key]; ok {
		if typedValue, ok := value.(T); ok {
			return typedValue
		}
	}
	return defaultValue
}

// Assign copies all properties from one or more source JsonMaps to the target JsonMap (jm), similar to Object.assign in JavaScript.
func (jm JsonMap) Assign(sources ...JsonMap) (JsonMap, error) {
	return Assign(jm, sources...)
}

// Assign copies all properties from one or more source JsonMaps to the target JsonMap (target), similar to Object.assign in JavaScript.
// It returns the target JsonMap after the assignment.
func Assign(target JsonMap, sources ...JsonMap) (JsonMap, error) {
	if target == nil {
		return nil, os.ErrInvalid
	}
	if len(sources) == 0 {
		return target, nil
	}
	for _, src := range sources {
		maps.Copy(target, src)
	}
	return target, nil
}

func (jm JsonMap) ToStruct(s any) error {
	if s == nil {
		return nil
	}
	data, err := json.Marshal(jm)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, s)
}

// func FromStruct(s any, defaults JsonMap) (JsonMap, error) {
// 	if s == nil {
// 		return nil, nil
// 	}
// 	data, err := json.Marshal(s)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var jm JsonMap
// 	err = json.Unmarshal(data, &jm)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if defaults != nil {
// 		// remove all keys from the result where values match the default values
// 		for k, v := range defaults {
// 			val, ok := jm[k]

// 			if ok && val == v {
// 				delete(jm, k)
// 			}

// 		}
// 	}

// 	return jm, nil
// }
