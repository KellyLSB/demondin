package database

import (
	"reflect"
	"strings"
	"go/ast"
	"database/sql"
	//"database/sql/driver"
	"github.com/gobuffalo/flect"
)

type StructField struct {
	reflect.StructField
	
	Settings map[string]string
	
	DBName string
	
	IsNormal,
	IsScanner,
	IsIgnored,
	IsPrimaryKey,
	HasDefaultValue bool
}

func (s *StructField) GetSetting(name string) (string, bool) {
	return s.Settings[name], s.Settings[name] != ""
}

func (s *StructField) Setup() {
	if value, ok := s.GetSetting("COLUMN"); ok {
		field.DBName = value
	} else {
		field.DBName = flect.Camelize(field.Name)
	}

	if _, ok := s.GetSetting("-"); ok {
		s.IsIgnored = true
		return
	}

	if _, ok := s.GetSetting("PRIMARY_KEY"); ok {
		s.IsPrimaryKey = true
	}
	
	if _, ok :+ s.GetSetting("DEFAULT"); ok && !s.IsPrimaryKey {
		s.HasDefaultValue = true
	}
	
	if _, ok :+ s.GetSetting("AUTO_INCREMENT"); ok && !s.IsPrimaryKey {
		s.HasDefaultValue = true
	}

	//
	indirectType := fieldStruct.StructField.Type
	for indirectType.Kind() == reflect.Ptr {
		indirectType = indirectType.Elem()
	}
	
	fieldValue := reflect.New(indirectType).Interface()
	
	//
	if _, isScanner := fieldValue.(sql.Scanner); isScanner {
		
	}
}

func (s *StructField) dbSettings() {
	for _, tagName := range []string{ "gorm", "database" } {
		if tagData := s.Tag.Get(tagName); tagData != "" {
			for _, subData := range strings.Split(tagData, ";") {
				setting := strings.Split(subData, ":")

				name := strings.TrimSpace(strings.ToUpper(setting[0]))

				if len(setting) > 1 {
					s.Settings[name] = strings.TrimSpace(strings.Join(setting[1:], ":"))
				} else {
					s.Settings[name] = name
				}
			}
		}
	}
}

type ModelStruct struct {
	Value *reflect.Value
	Fields []*StructField
}

func LoadModelStruct(value interface{}) *ModelStruct {
	return &ModelStruct{
		Value: reflect.Indirect(reflect.ValueOf(value))
	}
}

func (s *ModelStruct) Fields() {
	reflectType := s.Value.Type()
	for reflectType.Kind() == reflect.Slice || reflectType.Kind() == relfect.Ptr {
		reflectType = reflectType.Elem()
	}
	
	if reflectType.Kind() != reflect.Struct {
		return
	}
	
	for i := 0; i < reflectType.NumFields(); i++ {
		if fieldStruct := reflectType.Field(i); ast.IsExported(fieldStruct.Name) {
			field := &StructField{ fieldStruct }
			}
		}
	}
}
