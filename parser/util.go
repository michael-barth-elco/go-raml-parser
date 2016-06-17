package parser

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// LoadRAMLFromDir load RAML data from directory, concat *.raml
func LoadRAMLFromDir(dirPath string) (ramlData []byte, err error) {
	var filenames []string
	if filenames, err = filepath.Glob(filepath.Join(dirPath, "*.raml")); err != nil {
		return
	}
	sort.Strings(filenames)

	buffer := &bytes.Buffer{}
	for _, filename := range filenames {
		var filedata []byte
		if filedata, err = ioutil.ReadFile(filename); err != nil {
			return
		}
		if _, err = buffer.Write(filedata); err != nil {
			return
		}
		if _, err = buffer.WriteRune('\n'); err != nil {
			return
		}
	}

	return buffer.Bytes(), nil
}

// GetAPITypeName return type name from APIType, and isArray
func GetAPITypeName(apiType APIType) (typeName string, isArray bool) {
	typeName = apiType.Type
	isArray = strings.HasSuffix(apiType.Type, "[]")
	if isArray {
		typeName = apiType.Type[:len(apiType.Type)-2]
	}
	return
}

// CheckValueOption for changing CheckValueAPIType behavior
type CheckValueOption interface{}

// CheckValueOptionAllowIntegerToBeNumber allow type integer to be type number,
// e.g. APIType need a integer, but value is a number
// default: false
type CheckValueOptionAllowIntegerToBeNumber bool

// CheckValueAPIType check value is valid for apiType
func CheckValueAPIType(apiType APIType, value Value, options ...CheckValueOption) (err error) {
	if value.IsEmpty() {
		// no need to check if value is empty
		return
	}

	allowIntegerToBeNumber := CheckValueOptionAllowIntegerToBeNumber(false)

	for _, option := range options {
		switch option.(type) {
		case CheckValueOptionAllowIntegerToBeNumber:
			allowIntegerToBeNumber = option.(CheckValueOptionAllowIntegerToBeNumber)
		}
	}

	switch apiType.Type {
	case TypeBoolean, TypeInteger, TypeNumber, TypeString:
		if apiType.Type != value.Type {
			if allowIntegerToBeNumber &&
				apiType.Type == TypeNumber &&
				value.Type == TypeInteger {
				break
			}
			return ErrorPropertyTypeMismatch2.New(nil, apiType.Type, value.Type)
		}
	case TypeFile:
		// no type check for file type
		return
	default:
		if isInlineAPIType(apiType) {
			// no type check if declared by JSON
			return
		}

		switch value.Type {
		case TypeArray, TypeObject:
			break
		default:
			return ErrorPropertyTypeMismatch2.New(nil, apiType.Type, value.Type)
		}

		for name, property := range apiType.Properties {
			if property.Required {
				if !isValueContainKey(value, name) {
					return ErrorRequiredProperty2.New(nil, name, apiType.Type)
				}
			}

			if v, exist := value.Map[name]; exist && v != nil {
				if err = CheckValueAPIType(property.APIType, *v, options...); err != nil {
					if ErrorPropertyTypeMismatch2.Match(err) {
						return ErrorPropertyTypeMismatch3.New(nil, name, property.Type, v.Type)
					}
					return
				}
			}
		}
	}

	return nil
}

func isInlineAPIType(apiType APIType) bool {
	regValidType := regexp.MustCompile(`^[\w]+(\[\])?$`)
	return !regValidType.MatchString(apiType.Type)
}

func isValueContainKey(value Value, key string) bool {
	switch value.Type {
	case TypeArray:
		for _, v := range value.Array {
			if !isValueContainKey(*v, key) {
				return false
			}
		}
		return true
	case TypeObject:
		_, exist := value.Map[key]
		return exist
	}
	return false
}