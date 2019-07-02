package database

type Query struct {
  TableSelect []string
  FieldSelect []string
  
  QueryConditions []string
  OrderConditions []string
  GroupConditinos []string
  
  JoinTable []string
  
  Join, HasDistinctConditions bool
  
  Limit, Offset int
}

func (q *Query) Query()

func (q *Query) Join(type, table, conditions string) *Query {
  q.JoinTable = append(q.JoinTable, fmt.Sprintf(
    ? JOIN ? ON ?, type, table, conditions,
  ))
  
  return q
}
