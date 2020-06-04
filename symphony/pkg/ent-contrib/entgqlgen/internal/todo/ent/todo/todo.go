// Code generated (@generated) by entc, DO NOT EDIT.

package todo

const (
	// Label holds the string label denoting the todo type in the database.
	Label = "todo"
	// FieldID holds the string denoting the id field in the database.
	FieldID   = "id" // FieldText holds the string denoting the text vertex property in the database.
	FieldText = "text"

	// EdgeParent holds the string denoting the parent edge name in mutations.
	EdgeParent = "parent"
	// EdgeChildren holds the string denoting the children edge name in mutations.
	EdgeChildren = "children"

	// Table holds the table name of the todo in the database.
	Table = "todos"
	// ParentTable is the table the holds the parent relation/edge.
	ParentTable = "todos"
	// ParentColumn is the table column denoting the parent relation/edge.
	ParentColumn = "todo_children"
	// ChildrenTable is the table the holds the children relation/edge.
	ChildrenTable = "todos"
	// ChildrenColumn is the table column denoting the children relation/edge.
	ChildrenColumn = "todo_children"
)

// Columns holds all SQL columns for todo fields.
var Columns = []string{
	FieldID,
	FieldText,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the Todo type.
var ForeignKeys = []string{
	"todo_children",
}

var (
	// TextValidator is a validator for the "text" field. It is called by the builders before save.
	TextValidator func(string) error
)
