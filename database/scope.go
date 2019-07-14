package database

import (
	"reflect"
	"strings"
	"go/ast"
	"database/sql"
	//"database/sql/driver"
	"github.com/gobuffalo/flect"
)

type Relationship struct {
	Kind                         string
	
	//PolymorphicType              string
	//PolymorphicDBName            string
	//PolymorphicValue             string
	
	ForeignFieldNames            []string
	ForeignDBNames               []string
	AssociationForeignFieldNames []string
	AssociationForeignDBNames    []string
}

func relationshipFieldDBNames(
	m *ModelStruct, s *StructField, 
	joinTableSetting string,
	foreignKeys []string
) (
	keys   []string,
	fields []string, 
	dbName []string,
) {
	var joinTableDBNames []string
	
	if value, ok := s.GetSetting(joinTableSetting); ok {
		joinTableDBNames = strings.Split(value, ",")
	}

	if len(foreignKeys) < 1 {
		for _, field := range m.Fields {
			if field.IsPrimaryKey {
				foreignKeys = append(foreignKeys, field.DBName)
			}
		}
	}

	for i, foreignKey := range foreignKeys {
		if field := m.getField(foreignKey); field != nil {
			fields = append(fields, field.DBName)
		}
		
		if len(joinTableDBNames) > i {
			dbName = append(dbName, joinTablesDBNames[i])
		} else {
			dbName = append(
				dbName, DBName(s.StructField.Type.Name())+"_"+field.DBName,
			)
		}
	}
	
	return foreignKeys, fields, dbName
}

type StructField struct {
	reflect.StructField
	Relationship
	
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

func DBName(name string) string {
	return flect.Underscore(name)
}

func (s *StructField) Setup() func(*StructField) {
	if value, ok := s.GetSetting("COLUMN"); ok {
		s.DBName = value
	} else {
		s.DBName = DBName(s.Name)
	}

	if _, ok := s.GetSetting("-"); ok {
		s.IsIgnored = true
		return
	}

	//
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
	indirectType := s.StructField.Type
	for indirectType.Kind() == reflect.Ptr {
		indirectType = indirectType.Elem()
	}
	
	fieldValue := reflect.New(indirectType).Interface()
	
	//
	if _, isScanner := fieldValue.(sql.Scanner); isScanner {
		s.IsScanner, s.IsNormal = true, true
	} else if _, isTime := fieldValue.(*time.Time); isTime {
		s.IsNormal = true
	} else if _, ok := s.GetSetting("EMBEDDED"); ok || s.StructField.Anonymous {
		for _, subField := range LoadModelStruct(fieldValue).Fields() {
			// Preloaded SubFields
		}
		
		return
	} else {
		switch indirectType.Kind() {
		case reflect.Slice:
			return func(m *ModelStruct, s *StructField) {
				var relationship *Relationship
				var foreignKeys []string
				var associationForeignKeys []string
				var elemType reflect.Type = s.StructField.Type
				
				if value, ok := s.GetSetting("FOREIGNKEY"); ok {
					foreignKeys = strings.Split(value, ",")
				}
				
				if value, ok := s.GetSetting("ASSOCIATION_FOREIGNKEY"); ok {
					associationForeignKeys = strings.Split(value, ",")
				} else if value, ok := s.GetSetting("ASSOCIATIONFOREIGNKEY"); ok {
					associationForeignKeys = strings.Split(value, ",")
				}
				
				//
				for elemType.Kind() == reflect.Slice || elemType.Kind() == reflect.Ptr {
					elemType = elemType.Elem()
				}

				if elemType.Kind() == reflect.Struct {
					if value, ok := s.GetSetting("MANY2MANY"); ok {
						relationship.Kind = "many_to_many"

						{ // Source
							foreignKeys, foreignFields, foreignDBNames = relatedFieldDBNames(
								m, s, "JOINTABLE_FOREIGNKEY", foreignKeys,
							)

							relationship.ForeignFieldNames = foreignFields
							relationship.ForeignDBNames    = foreignDBNames
						}

						{ // Destination
							associationForeignKeys, foreignFields, foreignDBNames = relatedFieldDBNames(
								LoadModelStruct(reflect.New(s.StructField.Type).Interface()), 
								s, "ASSOCIATION_JOINTABLE_FOREIGNKEY", associationForeignKeys,
							)

							relationship.AssociationForeignFieldNames = foreignFields
							relationship.AssociationForeignDBNames    = foreignDBNames
						}

						s.Relationship = relationship
					} else {
						var associationType = m.Value.Type.Name()
						var fields = LoadModelStruct(
							reflect.New(s.StructField.Type).Interface(),
						).Fields()
						
						relationship.Kind = "has_many"
						
						// ---https://github.com/jinzhu/gorm/blob/836fb2c19d84dac7b0272958dfb9af7cf0d0ade4/model_struct.go#L365
					}
				}

			}
		case reflect.Struct:
		default:
			s.IsNormal = true
		}
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

func (s *ModelStruct) Fields() []*StructField {
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

func (s *ModelStruct) getField(name string) *StructField {
	for _, field := range s.Fields {
		if field.Name == name || field.DBName == DBName(name) {
			return field
		}
	}
	
	return nil
}
