package cmd

import (
	"container/list"
	"github.com/alexzorin/onapp/cmd/log"
	"reflect"
	"strconv"
	"strings"
)

type search struct {
	name  string
	query string
}

// This searches via reflect
// It takes q.name, finds the field of that name (case sensitive) and returns any matches on q.value
// String fields use strings.contains
// Ints and Bools are just equality
func (c *cli) Search(q search, items list.List) list.List {
	out := list.New()
	for item := items.Front(); item != nil; item = item.Next() {
		rv := reflect.ValueOf(item.Value)
		f := rv.FieldByName(q.name)
		if !f.IsValid() {
			log.Errorf("Field %s doesn't exist\n", q.name)
			break
		} else {
			switch f.Kind() {
			case reflect.String:
				sValue := f.String()
				if strings.Contains(sValue, q.query) {
					out.PushBack(item.Value)
				}
			case reflect.Int:
				qInt, err := strconv.Atoi(q.query)
				if err != nil {
					log.Errorf("%s is an int field, %s is not an int", q.name, q.query)
					break
				} else {
					if f.Int() == int64(qInt) {
						out.PushBack(item.Value)
					}
				}
			case reflect.Bool:
				qBoo, err := strconv.ParseBool(q.query)
				if err != nil {
					log.Errorf("%s is a bool field, %s is not a bool", q.name, q.query)
					break
				} else {
					if f.Bool() == qBoo {
						out.PushBack(item.Value)
					}
				}
			}
		}
	}
	return *out
}
