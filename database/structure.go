package database

import (
	"binary/bytes"
	"strings"
	"fmt"
)

type Query struct {
	TableSelect []string
	FieldSelect []string

	QueryConditions []string
	OrderConditions []string
	GroupConditions []string

	JoinTables []string

	Join bool

	Onset, Limit int
}

func (q *Query) Join(join, table, condition string) *Query {
	q.JoinTables = append(q.JoinTables, fmt.Sprintf(
		"%s JOIN %s ON (%s)", join, table, condition,
	))

	return q
}

func (q *Query) Condition(condition string) *Query {
	q.QueryConditions = append(q.QueryConditions, condition)
	return q
}

func (q *Query) Order(field, direction ...string) *Query {
	q.OrderConditions = append(q.OrderConditions, fmt.Sprintf(
		"? ?", field, len(direction) == 1 ? direction[0] : 'ASC',
	))
	
	return q
}

func (q *Query) Group(field string) *Query {
	q.GroupConditions = append(q.GroupConditions, field)
	return q
}

func (q *Query) Limit(onset, limit int) *Query {
	q.Onset = onset
	q.Limit = limit
	return q
}

func (q *Query) ToSQL() string {
	var buf bytes.Buffer
	
	if(len(q.FieldSelect) == 0) {
		fmt.Fprintf(buf, "SELECT *\n"))
	} else {
		fmt.Fprintf(buf, "SELECT (%s)\n", strings.Join(q.FieldSelect, ", "))
	}
	
	if(len(q.TableSelect) == 0) {
		panic("Table name is required in query builder")
	} else {
		fmt.Fprintf(buf, "FROM (%s)\n", strings.Join(q.TableSelect, ", "))
	}
	
	if(len(q.JoinConditions) > 0) {
		fmt.Fprintf(buf, "\t%s\n", strings.Join(q.JoinTables, "\n\t")) 
	}
	
	if(len(q.QueryConditions) > 0) {
		fmt.Fprintf(buf, "WHERE (%s)\n", strings.Join(q.QueryConditions, " && "))
	}
	
	if(len(q.GroupConditions) > 0) {
		fmt.Fprintf(buf, "GROUP BY (%s)\n", strings.Join(q.GroupConditions, ", "))
	}
	
	if(len(q.OrderConditions) > 0) {
		fmt.Fprintf(buf, "ORDER BY (%s)\n", strings.Join(q.OrderConditions, ", "))
	}
	
	if(q.Limit > 0) {
		fmt.Fprintf(buf, "LIMIT %d, %d\n", q.Onset, q.Limit)
	}
	
	return buf.String()
}
