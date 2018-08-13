package db

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hailongz/kk-lib/dynamic"
)

type IScheme interface {
	Open(db Database) error
	InstallSQL(db Database, object IObject, prefix string, autoIncrement int64) (string, error)
	Install(db Database, object IObject, prefix string, autoIncrement int64) error
}

type Scheme struct {
}

func NewScheme() IScheme {
	v := Scheme{}
	return &v
}

func (S *Scheme) Open(db Database) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS __kk_go_scheme (id BIGINT NOT NULL AUTO_INCREMENT,name VARCHAR(64) NULL,scheme TEXT NULL,PRIMARY KEY (id),INDEX name (name ASC) ) AUTO_INCREMENT=1;")
	return err
}

func (S *Scheme) InstallSQL(db Database, object IObject, prefix string, autoIncrement int64) (string, error) {

	var tbname = prefix + object.GetName()

	var rs, err = db.Query("SELECT * FROM __kk_go_scheme WHERE name=?", tbname)

	if err != nil {
		return "", err
	}

	defer rs.Close()

	b := bytes.NewBuffer(nil)

	if rs.Next() {

		var id int64
		var name string
		var scheme string
		rs.Scan(&id, &name, &scheme)
		var tb interface{} = nil
		json.Unmarshal([]byte(scheme), &tb)
		var table = map[string]Field{}

		Each(object, func(field Field) bool {

			if field.Name == "id" {
				return true
			}

			fd := dynamic.Get(tb, field.Name)

			if fd == nil {
				b.WriteString(fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN `%s` %s %s;", tbname, field.Name, field.DBType, field.DBValue))
			} else if dynamic.StringValue(dynamic.Get(fd, "dbtype"), "") != field.DBType ||
				dynamic.StringValue(dynamic.Get(fd, "dbvalue"), "") != field.DBValue {
				b.WriteString(fmt.Sprintf("ALTER TABLE `%s` CHANGE `%s` `%s` %s %s;", tbname, field.Name, field.Name, field.DBType, field.DBValue))
			}

			if field.DBIndex != "" && dynamic.StringValue(dynamic.Get(fd, "dbindex"), "") == "" {
				b.WriteString(fmt.Sprintf("CREATE INDEX `%s` ON `%s` (`%s` %s);", field.Name, tbname, field.Name, field.DBIndex))
			}

			table[field.Name] = field

			return true
		})

	} else {

		var i int = 0

		b.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (", tbname))

		b.WriteString("id BIGINT NOT NULL AUTO_INCREMENT")

		i += 1

		indexs := []Field{}
		var table = map[string]Field{}

		Each(object, func(field Field) bool {

			if field.Name == "id" {
				return true
			}

			if field.DBIndex != "" {
				indexs = append(indexs, field)
			}

			if i != 0 {
				b.WriteString(",")
			}

			b.WriteString(fmt.Sprintf("`%s` %s %s", field.Name, field.DBType, field.DBValue))

			i = i + 1

			table[field.Name] = field

			return true
		})

		b.WriteString(", PRIMARY KEY(id) ")

		for _, index := range indexs {
			b.WriteString(fmt.Sprintf(",INDEX `%s` (`%s` %s)", index.Name, index.Name, index.DBIndex))
		}

		if autoIncrement < 1 {
			b.WriteString(" ) ;")
		} else {
			b.WriteString(fmt.Sprintf(" ) AUTO_INCREMENT = %d;", autoIncrement))
		}

	}

	return b.String(), nil

}

func (S *Scheme) Install(db Database, object IObject, prefix string, autoIncrement int64) error {

	sql, err := S.InstallSQL(db, object, prefix, autoIncrement)

	if err != nil {
		return err
	}

	if sql == "" {
		return nil
	}

	_, err = db.Exec(sql)

	if err != nil {

		var table = map[string]Field{}

		Each(object, func(field Field) bool {

			if field.Name == "id" {
				return true
			}

			table[field.Name] = field

			return true
		})

		var tbname = prefix + object.GetName()

		rs, err := db.Query("SELECT id FROM __kk_go_scheme WHERE name=?", tbname)

		if err != nil {
			return err
		}

		defer rs.Close()

		if rs.Next() {

			var id int64

			err = rs.Scan(&id)

			if err != nil {
				return err
			}

			var b, _ = json.Marshal(table)
			_, err = db.Exec("UPDATE __kk_go_scheme SET scheme=? WHERE id=?", string(b), id)
			if err != nil {
				return err
			}
		} else {

			var b, _ = json.Marshal(table)

			_, err = db.Exec("INSERT INTO __kk_go_scheme(name,scheme) VALUES(?,?)", tbname, string(b))

			if err != nil {
				return err
			}
		}

		return err
	}

	return nil
}
