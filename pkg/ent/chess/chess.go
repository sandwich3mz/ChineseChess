// Code generated by ent, DO NOT EDIT.

package chess

import (
	"entgo.io/ent/dialect/sql"
)

const (
	// Label holds the string label denoting the chess type in the database.
	Label = "chess"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldBefore holds the string denoting the before field in the database.
	FieldBefore = "before"
	// FieldAfter holds the string denoting the after field in the database.
	FieldAfter = "after"
	// FieldCount holds the string denoting the count field in the database.
	FieldCount = "count"
	// Table holds the table name of the chess in the database.
	Table = "chess"
)

// Columns holds all SQL columns for chess fields.
var Columns = []string{
	FieldID,
	FieldBefore,
	FieldAfter,
	FieldCount,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

// OrderOption defines the ordering options for the Chess queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByBefore orders the results by the before field.
func ByBefore(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldBefore, opts...).ToFunc()
}

// ByAfter orders the results by the after field.
func ByAfter(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAfter, opts...).ToFunc()
}

// ByCount orders the results by the count field.
func ByCount(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCount, opts...).ToFunc()
}
