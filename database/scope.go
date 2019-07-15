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
		if field := m.GetField(foreignKey); field != nil {
			fields = append(fields, field.DBName)
		}
		
		if len(joinTableDBNames) > i {
			dbName = append(dbName, joinTablesDBNames[i])
		} else {
			dbName = append(
				dbName, DBName(s.StructField.Type.Name()) + "_" + field.DBName,
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

func (s *StructField) DstModelStruct() *ModelStruct {
	return LoadModelStruct(reflect.New(s.StructField.Type).Interface())
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
				var a            *ModelStruct = s.DstModelStruct()
				var relationship *Relationship
				var mForeignKeys []string
				var aForeignKeys []string
				var elemType     reflect.Type = s.StructField.Type
				
				if value, ok := s.GetSetting("FOREIGNKEY"); ok {
					mForeignKeys = strings.Split(value, ",")
				}
				
				if value, ok := s.GetSetting("ASSOCIATION_FOREIGNKEY"); ok {
					aForeignKeys = strings.Split(value, ",")
				} else if value, ok := s.GetSetting("ASSOCIATIONFOREIGNKEY"); ok {
					aForeignKeys = strings.Split(value, ",")
				}
				
				//
				for elemType.Kind() == reflect.Slice || elemType.Kind() == reflect.Ptr {
					elemType = elemType.Elem()
				}

				if elemType.Kind() == reflect.Struct {
					if value, ok := s.GetSetting("MANY2MANY"); ok {
						relationship.Kind = "many_to_many"

						{ // Source
							mForeignKeys, foreignFields, foreignDBNames = relatedFieldDBNames(
								m, s, "JOINTABLE_FOREIGNKEY", mForeignKeys,
							)

							relationship.ForeignFieldNames = foreignFields
							relationship.ForeignDBNames    = foreignDBNames
						}

						{ // Destination
							aForeignKeys, foreignFields, foreignDBNames = relatedFieldDBNames(
								a, s, "ASSOCIATION_JOINTABLE_FOREIGNKEY", aForeignKeys,
							)

							relationship.AssociationForeignFieldNames = foreignFields
							relationship.AssociationForeignDBNames    = foreignDBNames
						}

						s.Relationship = relationship
					} else {
						var associationType = m.Value.Type.Name()
						relationship.Kind = "has_many"
						
						if value, ok := s.GetSetting("POLYMORPHIC"); ok {
							if polyField := d.GetField(value + "Type"); polyType != nil {
								associationType                = value
								relationship.PolymorphicType   = polyField.Name
								relationship.PolymorphicDBName = polyField.DBName
								
								if value, ok := s.GetSetting("POLYMORPHIC_VALUE"); ok {
									relationship.PolymorphicValue = value
								} else {
									relationship.PolymorphicValue = m.TableName()
								}
								
								polyField.IsForeignKey = true
							}
						}
						
						if len(mForeignKeys) < 1 {
							if len(aForeignKeys) < 1 {
								a.PrimaryFields(func(field *StructField) {
									mForeignKeys = append(mForeignKeys, s.Name + field.Name)
									aForeignKeys = append(aForeignKeys, field.Name)
								})
							} else {
								for _, aForeignKey := range aForeignKeys {
									if field := a.GetField(aForeignKey); field != nil {
										mForeignKeys = append(mForeignKeys, s.Name + field.Name)
										aForeignKeys = append(aForeignKeys, field.Name)
									}
								}
							}
						} else {
							if len(aForeignKeys) < 1 {
								for _, mForeignKey := range mForeignKeys {
									if strings.HasPrefix(mForeignKey, s.Name) {
										aForeignKey := strings.TrimPrefix(mForeignKey, s.Name)
										if field := a.GetField(aForeignKey); field != nil {
											aForeignKeys = append(aForeignKeys, aForeignKey)
										}
									}
								}
								
								if len(aForeignKeys) == 0 && len(mForeignKeys) == 1 {
									aForeignKeys = []string{ a.PrimaryKey() }
								}
							} else if len(mForeignKeys) != len(aForeignKeys) {
								panic("Invalid ForeignKeys: a/m should express same length.")
							}
						}
						
						for i, mForeignKey := range mForeignKeys {
							if mField := m.GetField(mForeignKey); mField != nil {
								if aField := a.GetField(aForeignKeys[i]); aField != nil {
									mField.IsForeignKey = true
									
									relationship.ForeignFieldNames = append(
										relationship.ForeignFieldNames, mField.Name,
									)
									
									relationship.ForeignDBNames = append(
										relationship.ForeignDBNames, mField.DBName,
									)
									
									relationship.AssociationForeignFieldNames = append(
										relationship.AssociationForeignFieldNames, aField.Name,
									)
									
									relationship.AssociationForeignDBNames = append(
										relationship.AssociationForeignDBNames, aField.DBName,
									)
								}
							}
						}
						
						if len(relationship.ForeignFieldNames) != 0 {
							s.Relationship = relationship
						}
					}
				} else {
					s.IsNormal = true
				}
			} // (m, s)
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

func (s *ModelStruct) PrimaryFields(
	fns ...func(*StructField),
) (fields []*StructField) {
	for _, field := range s.Fields {
		if field.IsPrimaryKey {
			for _, fn := range fns {
				fn(field)
			}

			fields = append(fields, field)
		}
	}
	
	return
}

func (s *ModelStruct) PrimaryField() *StructField {
	if fields := s.PrimaryFields(); len(fields) > 0 {
		if len(fields) > 1 {
			for _, field := range fields {
				if field.DBName == "id" {
					return field
				}
			}
		}

		return fields[0]
	}

	return nil
}

func (s *ModelStruct) PrimaryKey() string {
	if field := s.PrimaryField(); field != nil {
		return field.DBName
	}

	return ""
}

func (s *ModelStruct) GetField(name string) *StructField {
	for _, field := range s.Fields {
		if field.Name == name || field.DBName == DBName(name) {
			return field
		}
	}
	
	return nil
}
