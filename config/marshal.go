package config

import (
	"bytes"
	"fmt"
	"reflect"
)

// write out a config in YAML
func (c *Type) ToYaml() *bytes.Buffer {
	cv := reflect.ValueOf(*c)
	buf := new(bytes.Buffer)
	printStruct("", &cv, buf)
	return buf
}

// write out a struct in YAML
func printStruct(leader string, val *reflect.Value, buf *bytes.Buffer) {
	tval := val.Type()
	nval := val.NumField()
	_, _ = fmt.Fprintln(buf)
	for i := 0; i < nval; i++ {
		sf := tval.Field(i)
		sv := val.Field(i)
		_, _ = fmt.Fprintf(buf, "%s%s: ", leader, sf.Name)
		switch sf.Type.Kind() {
		case reflect.Struct:
			printStruct(leader+"  ", &sv, buf)
		case reflect.Bool:
			_, _ = fmt.Fprintf(buf, "%v\n", sv.Bool())
		case reflect.String:
			_, _ = fmt.Fprintln(buf, sv.String())
		case reflect.Uint:
			_, _ = fmt.Fprintf(buf, "%d\n", sv.Uint())
		case reflect.Float32:
			_, _ = fmt.Fprintf(buf, "%f\n", sv.Float())
		default:
			panic("Unknown type in config struct")
		}
	}
}
