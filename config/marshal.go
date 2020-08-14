package config

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
)

// write out a config in GoB (Go Binary)
func (c *Type) ToGoB() (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(c)
	if err != nil {
		return nil, fmt.Errorf("error encoding config: %v", err)
	}
	return buf, nil
}

// Read a GoB config in
func FromGoB(in []byte) (*Type, error) {
	buf := bytes.NewBuffer(in)
	dec := gob.NewDecoder(buf)
	cfg := new(Type)
	err := dec.Decode(cfg)
	if err != nil {
		return nil, fmt.Errorf("error decoding config: %v", err)
	}
	return cfg, nil
}

// write out a config in YAML
func (c *Type) ToYaml(isWindows bool) *bytes.Buffer {
	cv := reflect.ValueOf(*c)
	buf := new(bytes.Buffer)
	printStruct("", &cv, buf, isWindows)
	return buf
}

// write out a struct in YAML
func printStruct(leader string, val *reflect.Value, buf *bytes.Buffer, isWindows bool) {
	tval := val.Type()
	nval := val.NumField()
	for i := 0; i < nval; i++ {
		sv := val.Field(i)
		sf := tval.Field(i)
		sn := sf.Name
		// If field is tagged omitwindows and we are on Windows, don't write it at all
		_, ok := sf.Tag.Lookup("omitwindows")
		if ok && isWindows {
			continue
		}
		snl := string(sn[0]+('a'-'A')) + sn[1:] // lowercase first letter
		switch sf.Type.Kind() {
		case reflect.Struct:
			if checkStruct(&sv, isWindows) {
				_, _ = fmt.Fprintf(buf, "%s%s:\n", leader, snl)
				printStruct(leader+"  ", &sv, buf, isWindows)
			}
		case reflect.Bool:
			if !isDefault(&sv, &sf) {
				_, _ = fmt.Fprintf(buf, "%s%s: %t\n", leader, snl, sv.Bool())
			}
		case reflect.String:
			if !isDefault(&sv, &sf) {
				_, _ = fmt.Fprintf(buf, "%s%s: %s\n", leader, snl, sv.String())
			}
		case reflect.Uint:
			if !isDefault(&sv, &sf) {
				_, _ = fmt.Fprintf(buf, "%s%s: %d\n", leader, snl, sv.Uint())
			}
		case reflect.Float32:
			if !isDefault(&sv, &sf) {
				_, _ = fmt.Fprintf(buf, "%s%s: %.2f\n", leader, snl, sv.Float())
			}
		default:
			panic("Unknown type in config struct")
		}
	}
}

// Check if struct needs to be put into the YAML output. If it has all sub-elements with values that should not be output, return false, otherwise return true
// this a lookahead when we encounter a struct field
func checkStruct(val *reflect.Value, isWindows bool) bool {
	ret := false
	tval := val.Type()
	nval := val.NumField()
	for i := 0; i < nval; i++ {
		sf := tval.Field(i)
		sv := val.Field(i)
		_, ok := sf.Tag.Lookup("omitwindows")
		if ok && isWindows {
			return false
		}
		switch sf.Type.Kind() {
		case reflect.Struct:
			ret = checkStruct(&sv, isWindows)
		case reflect.Bool, reflect.String, reflect.Uint, reflect.Float32:
			ret = !isDefault(&sv, &sf)
		default:
			panic("Unknown type in config struct")
		}
		if ret {
			break
		}
	}
	return ret
}

// Check if a struct field is the default value or not, return true if it is
// the default value is the zero value unless there is a tag on the struct field specifying a different default
func isDefault(val *reflect.Value, field *reflect.StructField) (ret bool) {
	def, ok := field.Tag.Lookup("default")
	if ok {
		switch field.Type.Kind() {
		case reflect.String:
			ret = def == val.String()
		case reflect.Bool:
			var v bool
			_, _ = fmt.Sscanf(def, "%t", &v)
			ret = v == val.Bool()
		case reflect.Uint:
			var v uint64
			_, _ = fmt.Sscanf(def, "%d", &v)
			ret = v == val.Uint()
		case reflect.Float32:
			var v float64
			_, _ = fmt.Sscanf(def, "%f", &v)
			ret = v == val.Float()
		}
	} else { // no default specified, compare check against zero value
		switch field.Type.Kind() {
		case reflect.Bool:
			ret = val.Bool() == false
		case reflect.String:
			ret = val.String() == ""
		case reflect.Uint:
			ret = val.Uint() == 0
		case reflect.Float32:
			ret = val.Float() == 0.0
		default:
			panic("isDefault: Unknown type in config struct")
		}
	}
	return ret
}
