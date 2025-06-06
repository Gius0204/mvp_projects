// Code generated by ent, DO NOT EDIT.

package task

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the task type in the database.
	Label = "task"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldParentID holds the string denoting the parent_id field in the database.
	FieldParentID = "parent_id"
	// FieldProgress holds the string denoting the progress field in the database.
	FieldProgress = "progress"
	// FieldEstimatedTime holds the string denoting the estimated_time field in the database.
	FieldEstimatedTime = "estimated_time"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// EdgeParent holds the string denoting the parent edge name in mutations.
	EdgeParent = "parent"
	// EdgeSubtasks holds the string denoting the subtasks edge name in mutations.
	EdgeSubtasks = "subtasks"
	// Table holds the table name of the task in the database.
	Table = "tasks"
	// ParentTable is the table that holds the parent relation/edge.
	ParentTable = "tasks"
	// ParentColumn is the table column denoting the parent relation/edge.
	ParentColumn = "parent_id"
	// SubtasksTable is the table that holds the subtasks relation/edge.
	SubtasksTable = "tasks"
	// SubtasksColumn is the table column denoting the subtasks relation/edge.
	SubtasksColumn = "parent_id"
)

// Columns holds all SQL columns for task fields.
var Columns = []string{
	FieldID,
	FieldParentID,
	FieldProgress,
	FieldEstimatedTime,
	FieldCreatedAt,
	FieldUpdatedAt,
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

var (
	// DefaultProgress holds the default value on creation for the "progress" field.
	DefaultProgress float64
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() time.Time
	// UpdateDefaultUpdatedAt holds the default value on update for the "updated_at" field.
	UpdateDefaultUpdatedAt func() time.Time
)

// OrderOption defines the ordering options for the Task queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByParentID orders the results by the parent_id field.
func ByParentID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldParentID, opts...).ToFunc()
}

// ByProgress orders the results by the progress field.
func ByProgress(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldProgress, opts...).ToFunc()
}

// ByEstimatedTime orders the results by the estimated_time field.
func ByEstimatedTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldEstimatedTime, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByUpdatedAt orders the results by the updated_at field.
func ByUpdatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdatedAt, opts...).ToFunc()
}

// ByParentField orders the results by parent field.
func ByParentField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newParentStep(), sql.OrderByField(field, opts...))
	}
}

// BySubtasksCount orders the results by subtasks count.
func BySubtasksCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newSubtasksStep(), opts...)
	}
}

// BySubtasks orders the results by subtasks terms.
func BySubtasks(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newSubtasksStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newParentStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(Table, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, ParentTable, ParentColumn),
	)
}
func newSubtasksStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(Table, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, SubtasksTable, SubtasksColumn),
	)
}
