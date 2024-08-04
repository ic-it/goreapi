package goreapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
)

type vartype struct {
	name       string
	kind       any
	isInternal bool
}

type param struct {
	name string
	typ  vartype
}

type body struct {
	name string
	typ  vartype
}

type query struct {
	name string
	typ  vartype
}

type response struct {
	typ vartype
}

// H is a dev-time finction that returns a http.HandlerFunc
func H(v ...any) http.HandlerFunc {
	params := []param{}
	bodies := []body{}
	queries := []query{}
	var resp response
	var fn any

	for _, vv := range v {
		switch vvv := vv.(type) {
		case param:
			params = append(params, vvv)
		case body:
			bodies = append(bodies, vvv)
		case query:
			queries = append(queries, vvv)
		case response:
			resp = vvv
		default:
			v := reflect.TypeOf(vvv)
			if v.Kind() != reflect.Func {
				panic("expected func")
			}
			fn = vvv
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		type argval struct {
			name string
			val  reflect.Value
		}

		args := []argval{}

		for _, q := range queries {
			rawdata := r.URL.Query().Get(q.name)
			if q.typ.isInternal {
				switch (q.typ.kind).(type) {
				case *Int:
					data, err := strconv.Atoi(rawdata)
					if err != nil {
						panic(err)
					}
					val := reflect.New(reflect.TypeOf(1))
					val.Elem().Set(reflect.ValueOf(data))
					args = append(args, argval{q.name, val})
				case *String:
					data := rawdata
					val := reflect.New(reflect.TypeOf("string"))
					val.Elem().Set(reflect.ValueOf(data))
					args = append(args, argval{q.name, val})
				default:
					panic("unexpected type")
				}
			} else {
				data := reflect.New(reflect.TypeOf(q.typ.kind).Elem()).Interface()
				err := json.Unmarshal([]byte(rawdata), &data)
				if err != nil {
					panic(err)
				}
				val := reflect.New(reflect.TypeOf(q.typ.kind))
				val.Elem().Set(reflect.ValueOf(data))
				args = append(args, argval{q.name, val})
			}
		}

		for _, p := range params {
			rawdata := r.PathValue(p.name)
			if p.typ.isInternal {
				switch (p.typ.kind).(type) {
				case *Int:
					data, err := strconv.Atoi(rawdata)
					if err != nil {
						panic(err)
					}
					val := reflect.New(reflect.TypeOf(1))
					val.Elem().Set(reflect.ValueOf(data))
					args = append(args, argval{p.name, val})
				case *String:
					data := rawdata
					val := reflect.New(reflect.TypeOf("string"))
					val.Elem().Set(reflect.ValueOf(data))
					args = append(args, argval{p.name, val})
				default:
					panic("unexpected type")
				}
			} else {
				data := reflect.New(reflect.TypeOf(p.typ.kind).Elem()).Interface()
				err := json.Unmarshal([]byte(rawdata), &data)
				if err != nil {
					panic(err)
				}
				val := reflect.New(reflect.TypeOf(p.typ.kind))
				val.Elem().Set(reflect.ValueOf(data))
				args = append(args, argval{p.name, val})
			}
		}

		for _, b := range bodies {
			rawdata := r.Body
			if b.typ.isInternal {
				switch (b.typ.kind).(type) {
				case *Int:
					rawdata := bytes.NewBuffer([]byte{})
					rawdata.ReadFrom(r.Body)
					data, err := strconv.Atoi(string(rawdata.Bytes()))
					if err != nil {
						panic(err)
					}
					val := reflect.New(reflect.TypeOf(1))
					val.Elem().Set(reflect.ValueOf(data))
					args = append(args, argval{b.name, val})
				case *String:
					data := rawdata
					val := reflect.New(reflect.TypeOf("string"))
					val.Elem().Set(reflect.ValueOf(data))
					args = append(args, argval{b.name, val})
				default:
					panic("unexpected type")
				}
			} else {
				data := reflect.New(reflect.TypeOf(b.typ.kind).Elem()).Interface()
				err := json.NewDecoder(r.Body).Decode(&data)
				if err != nil {
					panic(err)
				}
				val := reflect.New(reflect.TypeOf(b.typ.kind))
				val.Elem().Set(reflect.ValueOf(data))
				args = append(args, argval{b.name, val})
			}
		}

		fnArgsF := make([]reflect.Value, len(args))
		// match the order of the args
		for i, arg := range args {
			fnArgsF[i] = arg.val.Elem()
		}

		ret := reflect.ValueOf(fn).Call(fnArgsF)
		if resp.typ.isInternal {
			switch (resp.typ.kind).(type) {
			case *Int:
				w.Write([]byte(strconv.Itoa(ret[0].Interface().(int))))
			case *String:
				w.Write([]byte(ret[0].Interface().(string)))
			default:
				panic("unexpected type")
			}
		} else {
			data, err := json.Marshal(ret[0].Interface())
			if err != nil {
				panic(err)
			}
			w.Write(data)
		}
	}
}

func Param(
	name string,
	kind any,
) param {
	switch (kind).(type) {
	case *Int, *String:
		return param{name, vartype{name, kind, true}}
	default:
		return param{name, vartype{name, kind, false}}
	}
}

func Body(
	name string,
	kind any,
) body {
	switch (kind).(type) {
	case *Int, *String:
		return body{name, vartype{name, kind, true}}
	default:
		return body{name, vartype{name, kind, false}}
	}
}

func Query(
	name string,
	kind any,
) query {
	switch (kind).(type) {
	case *Int, *String:
		return query{name, vartype{name, kind, true}}
	default:
		return query{name, vartype{name, kind, false}}
	}
}

func Response(
	kind any,
) response {
	switch (kind).(type) {
	case *Int, *String:
		return response{vartype{"", kind, true}}
	default:
		return response{vartype{"", kind, false}}
	}
}

type Int struct {
	Value int
}

type String struct {
	Value string
}
