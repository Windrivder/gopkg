package structs

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Applies the default values to the struct object, the struct type must have
// the StructTag with name "default" and the directed value.
//
// Usage
//     type ExampleBasic struct {
//         Foo bool   `default:"true"`
//         Bar string `default:"33"`
//         Qux int8
//         Dur time.Duration `default:"2m3s"`
//     }
//
//      foo := &ExampleBasic{}
//      SetDefaults(foo)
func SetDefaults(variable interface{}) {
	getDefaultFiller().Fill(variable)
}

var (
	defaultTagName      string   = "default"
	defaultTimeTypeHash TypeHash = "time.Duration"
	defaultFiller       *Filler  = nil
)

func getDefaultFiller() *Filler {
	if defaultFiller == nil {
		defaultFiller = newDefaultFiller()
	}

	return defaultFiller
}

func newDefaultFiller() *Filler {
	funcs := make(map[reflect.Kind]FillerFunc, 0)

	funcs[reflect.Bool] = func(field *FieldData) {
		value, _ := strconv.ParseBool(field.TagValue)
		field.Value.SetBool(value)
	}

	funcs[reflect.Int] = func(field *FieldData) {
		value, _ := strconv.ParseInt(field.TagValue, 10, 64)
		field.Value.SetInt(value)
	}
	funcs[reflect.Int8] = funcs[reflect.Int]
	funcs[reflect.Int16] = funcs[reflect.Int]
	funcs[reflect.Int32] = funcs[reflect.Int]

	funcs[reflect.Int64] = func(field *FieldData) {
		if field.Field.Type == reflect.TypeOf(time.Second) {
			value, _ := time.ParseDuration(field.TagValue)
			field.Value.Set(reflect.ValueOf(value))
		} else {
			value, _ := strconv.ParseInt(field.TagValue, 10, 64)
			field.Value.SetInt(value)
		}
	}

	funcs[reflect.Float32] = func(field *FieldData) {
		value, _ := strconv.ParseFloat(field.TagValue, 64)
		field.Value.SetFloat(value)
	}
	funcs[reflect.Float64] = funcs[reflect.Float32]

	funcs[reflect.Uint] = func(field *FieldData) {
		value, _ := strconv.ParseUint(field.TagValue, 10, 64)
		field.Value.SetUint(value)
	}
	funcs[reflect.Uint8] = funcs[reflect.Uint]
	funcs[reflect.Uint16] = funcs[reflect.Uint]
	funcs[reflect.Uint32] = funcs[reflect.Uint]
	funcs[reflect.Uint64] = funcs[reflect.Uint]

	funcs[reflect.String] = func(field *FieldData) {
		field.Value.SetString(field.TagValue)
	}

	funcs[reflect.Struct] = func(field *FieldData) {
		fields := getDefaultFiller().GetFieldsFromValue(field.Value, nil)
		getDefaultFiller().SetDefaultValues(fields)
	}

	types := make(map[TypeHash]FillerFunc, 1)

	types[defaultTimeTypeHash] = func(field *FieldData) {
		d, _ := time.ParseDuration(field.TagValue)
		field.Value.Set(reflect.ValueOf(d))
	}

	funcs[reflect.Slice] = func(field *FieldData) {
		k := field.Value.Type().Elem().Kind()
		switch k {
		case reflect.Uint8:
			if field.Value.Bytes() != nil {
				return
			}
			field.Value.SetBytes([]byte(field.TagValue))
		case reflect.Struct:
			count := field.Value.Len()
			for i := 0; i < count; i++ {
				fields := getDefaultFiller().GetFieldsFromValue(field.Value.Index(i), nil)
				getDefaultFiller().SetDefaultValues(fields)
			}
		default:
			// 处理形如 [1,2,3,4]
			reg := regexp.MustCompile(`^\[(.*)\]$`)
			matchs := reg.FindStringSubmatch(field.TagValue)
			if len(matchs) != 2 {
				return
			}
			if matchs[1] == "" {
				field.Value.Set(reflect.MakeSlice(field.Value.Type(), 0, 0))
			} else {
				defaultValue := strings.Split(matchs[1], ",")
				result := reflect.MakeSlice(field.Value.Type(), len(defaultValue), len(defaultValue))
				for i := 0; i < len(defaultValue); i++ {
					itemValue := result.Index(i)
					item := &FieldData{
						Value:    itemValue,
						Field:    reflect.StructField{},
						TagValue: defaultValue[i],
						Parent:   nil,
					}
					funcs[k](item)
				}
				field.Value.Set(result)
			}
		}
	}

	return &Filler{FuncByKind: funcs, FuncByType: types, Tag: defaultTagName}
}
