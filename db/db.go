package db

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/hailongz/kk-lib/dynamic"
)

type Database interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type IObject interface {
	GetTitle() string
	GetName() string
	GetId() int64
	SetId(id int64)
}

type Object struct {
	Id int64 `json:"id" title:"ID"`
}

type Field struct {
	F        reflect.StructField `json:"-"`
	V        reflect.Value       `json:"-"`
	Name     string              `json:"name"`
	IsObject bool                `json:"-"`
}

func (O *Object) GetTitle() string {
	return "数据对象"
}

func (O *Object) GetName() string {
	return "object"
}

func (O *Object) GetId() int64 {
	return O.Id
}

func (O *Object) SetId(id int64) {
	O.Id = id
}

func TableName(prefix string, object IObject) string {
	return prefix + object.GetName()
}

func Transaction(db *sql.DB, fn func(conn Database) error) error {

	tx, err := db.Begin()

	err = fn(tx)

	if err == nil {
		err = tx.Commit()
	}

	if err != nil {
		tx.Rollback()
	}

	return err
}

func Query(db Database, object IObject, prefix string, sql string, args ...interface{}) (*sql.Rows, error) {
	var tbname = prefix + object.GetName()
	return db.Query(fmt.Sprintf("SELECT * FROM `%s` %s", tbname, sql), args...)
}

func QueryWithKeys(db Database, object IObject, prefix string, keys map[string]bool, sql string, args ...interface{}) (*sql.Rows, error) {

	s := bytes.NewBuffer(nil)

	if keys == nil {
		s.WriteString("SELECT *")
	} else {

		s.WriteString("SELECT id")

		Each(object, func(field Field) bool {

			if keys[field.Name] {
				s.WriteString(fmt.Sprintf(",`%s`", field.Name))
			}

			return true
		})

	}

	s.WriteString(fmt.Sprintf(" FROM `%s%s` %s", prefix, object.GetName(), sql))

	return db.Query(s.String(), args...)
}

func Delete(db Database, object IObject, prefix string) (sql.Result, error) {
	var tbname = prefix + object.GetName()
	return db.Exec(fmt.Sprintf("DELETE FROM `%s` WHERE id=?", tbname), object.GetId())
}

func DeleteWithSQL(db Database, object IObject, prefix string, sql string, args ...interface{}) (sql.Result, error) {
	var tbname = prefix + object.GetName()
	return db.Exec(fmt.Sprintf("DELETE FROM `%s` %s", tbname, sql), args...)
}

func Count(db Database, object IObject, prefix string, sql string, args ...interface{}) (int64, error) {

	var tbname = prefix + object.GetName()

	var rows, err = db.Query(fmt.Sprintf("SELECT COUNT(*) as c FROM `%s` %s", tbname, sql), args...)

	if err != nil {
		return 0, err
	}

	defer rows.Close()

	if rows.Next() {
		var v int64 = 0
		err = rows.Scan(&v)
		if err != nil {
			return 0, err
		}
		return v, nil
	}

	return 0, nil
}

func Update(db Database, object IObject, prefix string) (sql.Result, error) {
	return UpdateWithKeys(db, object, prefix, nil)
}

func UpdateWithKeys(db Database, object IObject, prefix string, keys map[string]bool) (sql.Result, error) {

	var tbname = prefix + object.GetName()
	var s bytes.Buffer
	var fs = []interface{}{}
	var n = 0

	s.WriteString(fmt.Sprintf("UPDATE `%s` SET ", tbname))

	Each(object, func(field Field) bool {

		if field.Name == "id" {
			return true
		}

		if keys == nil || keys[field.Name] {
			if n != 0 {
				s.WriteString(",")
			}
			s.WriteString(fmt.Sprintf(" `%s`=?", field.Name))
			if field.IsObject {
				b, _ := json.Marshal(field.V.Interface())
				fs = append(fs, string(b))
			} else {
				fs = append(fs, field.V.Interface())
			}
			n += 1
		}

		return true
	})

	s.WriteString(" WHERE id=?")

	fs = append(fs, object.GetId())

	n += 1

	// log.Printf("%s %s\n", s.String(), fs)

	return db.Exec(s.String(), fs...)
}

func Insert(db Database, object IObject, prefix string) (sql.Result, error) {
	var tbname = prefix + object.GetName()
	var s bytes.Buffer
	var w bytes.Buffer
	var fs = []interface{}{}
	var n = 0

	s.WriteString(fmt.Sprintf("INSERT INTO `%s`(", tbname))
	w.WriteString(" VALUES (")

	Each(object, func(field Field) bool {

		if field.Name == "id" && object.GetId() == 0 {
			return true
		}

		if n != 0 {
			s.WriteString(",")
			w.WriteString(",")
		}
		s.WriteString("`" + field.Name + "`")
		w.WriteString("?")
		if field.IsObject {
			b, _ := json.Marshal(field.V.Interface())
			fs = append(fs, string(b))
		} else {
			fs = append(fs, field.V.Interface())
		}

		n += 1

		return true
	})

	s.WriteString(")")

	w.WriteString(")")

	s.Write(w.Bytes())

	// log.Printf("%s %s\n", s.String(), fs)

	var rs, err = db.Exec(s.String(), fs...)

	if err == nil && object.GetId() == 0 {
		id, err := rs.LastInsertId()
		if err == nil {
			object.SetId(id)
		}
	}

	return rs, err
}

type booleanValue struct {
	v        reflect.Value
	intValue int
}

type jsonValue struct {
	v         reflect.Value
	byteValue interface{}
}

type Scaner struct {
	object        IObject
	fields        []interface{}
	jsonObjects   []*jsonValue
	nilValue      interface{}
	booleanValues []*booleanValue
}

func NewScaner(object IObject) *Scaner {
	return &Scaner{object, nil, nil, nil, nil}
}

func (o *Scaner) Scan(rows *sql.Rows) error {

	if o.fields == nil {

		var columns, err = rows.Columns()

		if err != nil {
			return err
		}

		var fdc = len(columns)
		var mi = map[string]int{}

		for i := 0; i < fdc; i += 1 {
			mi[columns[i]] = i
		}

		o.booleanValues = []*booleanValue{}
		o.jsonObjects = []*jsonValue{}
		o.fields = make([]interface{}, fdc)

		for i := 0; i < fdc; i += 1 {
			o.fields[i] = &o.nilValue
		}

		Each(o.object, func(field Field) bool {

			idx, ok := mi[field.Name]

			if ok {
				if field.F.Type.Kind() == reflect.Bool {
					b := booleanValue{}
					b.v = field.V
					b.intValue = 0
					o.fields[idx] = &b.intValue
					o.booleanValues = append(o.booleanValues, &b)
				} else if field.IsObject {
					b := jsonValue{}
					b.v = field.V
					b.byteValue = nil
					o.fields[idx] = &b.byteValue
					o.jsonObjects = append(o.jsonObjects, &b)
				} else {
					o.fields[idx] = field.V.Addr().Interface()
				}
			}

			return true
		})

	}

	err := rows.Scan(o.fields...)

	if err != nil {
		return err
	}

	for _, fd := range o.jsonObjects {

		if fd.byteValue == nil {
			dynamic.SetReflectValue(fd.v, nil)
			continue
		}

		{
			b, ok := (fd.byteValue).([]byte)
			if ok {
				_ = json.Unmarshal(b, fd.v.Addr().Interface())
				continue
			}
		}

		{
			s, ok := (fd.byteValue).(string)
			if ok {
				_ = json.Unmarshal([]byte(s), fd.v.Addr().Interface())
				continue
			}
		}

		dynamic.SetReflectValue(fd.v, nil)
	}

	for _, fd := range o.booleanValues {
		dynamic.SetReflectValue(fd.v, fd.intValue != 0)
	}

	return nil
}
