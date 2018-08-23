package db

import (
	"bytes"
	"fmt"

	"github.com/hailongz/kk-lib/dynamic"
)

func fieldSQLType(stype string, length int64, defaultValue string) string {

	b := bytes.NewBuffer(nil)

	switch stype {
	case "int":
		if length > 0 {
			b.WriteString(fmt.Sprintf("INT(%d)", length))
		} else {
			b.WriteString("INT")
		}
		b.WriteString(" DEFAULT ")
		b.WriteString(dynamic.StringValue(dynamic.IntValue(defaultValue, 0), "0"))
		break
	case "long":
		if length > 0 {
			b.WriteString(fmt.Sprintf("BIGINT(%d)", length))
		} else {
			b.WriteString("BIGINT")
		}
		b.WriteString(" DEFAULT ")
		b.WriteString(dynamic.StringValue(dynamic.IntValue(defaultValue, 0), "0"))
		break
	case "double":
		if length > 0 {
			b.WriteString(fmt.Sprintf("DOUBLE(%d)", length))
		} else {
			b.WriteString("DOUBLE")
		}
		b.WriteString(" DEFAULT ")
		b.WriteString(dynamic.StringValue(dynamic.IntValue(defaultValue, 0), "0"))
		break
	case "boolean":
		if length > 0 {
			b.WriteString(fmt.Sprintf("DOUBLE(%d)", length))
		} else {
			b.WriteString("DOUBLE")
		}
		b.WriteString(" DEFAULT ")
		if dynamic.BooleanValue(defaultValue, false) {
			b.WriteString("0")
		} else {
			b.WriteString("1")
		}
		break
	case "string":
		if length == 0 {
			b.WriteString("VARCHAR(64)")
			b.WriteString(" DEFAULT ")
			b.WriteString(fmt.Sprintf("'%s'", dynamic.StringValue(defaultValue, "")))
		} else if length > 65535 {
			b.WriteString(fmt.Sprintf("LONGTEXT(%d)", length))
		} else if length > 4096 {
			b.WriteString(fmt.Sprintf("TEXT(%d)", length))
		} else if length == -1 {
			b.WriteString("TEXT")
		} else if length == -2 {
			b.WriteString("LONGTEXT")
		} else {
			b.WriteString(fmt.Sprintf("VARCHAR(%d)", length))
			b.WriteString(" DEFAULT ")
			b.WriteString(fmt.Sprintf("'%s'", dynamic.StringValue(defaultValue, "")))
		}
	default:
		if length == 0 {
			b.WriteString("TEXT")
		} else {
			b.WriteString(fmt.Sprintf("TEXT(%d)", length))
		}
	}

	return b.String()
}

func InstallSQL(table interface{}, prefix string, autoIncrement int64, ver interface{}) (string, interface{}) {

	tbname := prefix + dynamic.StringValue(dynamic.Get(table, "name"), "")
	tbtitle := dynamic.StringValue(dynamic.Get(table, "title"), "")

	b := bytes.NewBuffer(nil)

	if ver == nil {

		var i int = 0

		b.WriteString(fmt.Sprintf("#%s\r\n", tbtitle))
		b.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (", tbname))

		b.WriteString("\r\n\tid BIGINT NOT NULL")

		if autoIncrement != 0 {
			b.WriteString(" AUTO_INCREMENT")
		}

		b.WriteString("\t#ID\r\n")

		i += 1

		indexs := []interface{}{}

		var tb = map[string]interface{}{}

		dynamic.Each(dynamic.Get(table, "fields"), func(key interface{}, field interface{}) bool {

			name := dynamic.StringValue(dynamic.Get(field, "name"), "")
			stype := dynamic.StringValue(dynamic.Get(field, "type"), "")
			length := dynamic.IntValue(dynamic.Get(field, "length"), 0)
			index := dynamic.StringValue(dynamic.Get(field, "index"), "")
			title := dynamic.StringValue(dynamic.Get(field, "title"), "")
			defaultValue := dynamic.StringValue(dynamic.Get(field, "default"), "")

			if index != "" {
				indexs = append(indexs, field)
			}

			if i != 0 {
				b.WriteString("\t,")
			}

			b.WriteString(fmt.Sprintf("`%s` %s", name, fieldSQLType(stype, length, defaultValue)))
			b.WriteString(fmt.Sprintf("\t#[字段] %s\r\n", title))
			i = i + 1

			tb[name] = field

			return true
		})

		b.WriteString("\t, PRIMARY KEY(id) \r\n")

		for _, field := range indexs {
			name := dynamic.StringValue(dynamic.Get(field, "name"), "")
			index := dynamic.StringValue(dynamic.Get(field, "index"), "")
			title := dynamic.StringValue(dynamic.Get(field, "title"), "")
			b.WriteString(fmt.Sprintf("\t,INDEX `%s` (`%s` %s)", name, name, index))
			b.WriteString(fmt.Sprintf("\t#[索引] %s\r\n", title))
		}

		if autoIncrement == 0 {
			b.WriteString(" ) ;\r\n")
		} else {
			b.WriteString(fmt.Sprintf(" ) AUTO_INCREMENT = %d;\r\n", autoIncrement))
		}

		return b.String(), tb

	} else {

		var tb = map[string]interface{}{}

		dynamic.Each(dynamic.Get(table, "fields"), func(key interface{}, field interface{}) bool {

			name := dynamic.StringValue(dynamic.Get(field, "name"), "")
			stype := dynamic.StringValue(dynamic.Get(field, "type"), "")
			length := dynamic.IntValue(dynamic.Get(field, "length"), 0)
			index := dynamic.StringValue(dynamic.Get(field, "index"), "")
			title := dynamic.StringValue(dynamic.Get(field, "title"), "")
			defaultValue := dynamic.StringValue(dynamic.Get(field, "default"), "")

			fd := dynamic.Get(ver, name)

			if fd == nil {
				b.WriteString(fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN `%s` %s;", tbname, name, fieldSQLType(stype, length, defaultValue)))
				b.WriteString(fmt.Sprintf("\t#[增加字段] %s\r\n", title))
			} else if dynamic.StringValue(dynamic.Get(fd, "type"), "") != stype ||
				dynamic.IntValue(dynamic.Get(fd, "length"), 0) != length {
				b.WriteString(fmt.Sprintf("ALTER TABLE `%s` CHANGE `%s` `%s` %s;", tbname, name, name, fieldSQLType(stype, length, defaultValue)))
				b.WriteString(fmt.Sprintf("\t#[修改字段] %s\r\n", title))
			}

			if index != "" && dynamic.StringValue(dynamic.Get(fd, "index"), "") == "" {
				b.WriteString(fmt.Sprintf("CREATE INDEX `%s` ON `%s` (`%s` %s);", name, tbname, name, index))
				b.WriteString(fmt.Sprintf("\t#[创建索引] %s\r\n", title))
			}

			tb[name] = field

			return true
		})

		return b.String(), tb
	}
}

func Install(db Database, object IObject, prefix string, autoIncrement int64, table interface{}) (error, interface{}) {

	sql, tb := InstallSQL(object, prefix, autoIncrement, table)

	if sql == "" {
		return nil, tb
	}

	_, err := db.Exec(sql)

	return err, tb
}
