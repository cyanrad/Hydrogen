package object

const (
	INT_OBJ      = "INT_OBJ"     // integers: 1,2,3,...
	BOOLEAN_OBJ  = "BOOLEAN_OBJ" // true or false
	NULL_OBJ     = "NULL_OBJ"
	FUNCTION_OBJ = "FUNCTION_OBJ" // function object
	BUILTIN_OBJ  = "BUILTIN_OBJ"  // built-in function object
	STRING_OBJ   = "STRING_OBJ"   // "hello"
	ARRAY_OBJ    = "ARRAY_OBJ"    // [1,2,3]
	HASH_OBJ     = "HASH_OBJ"     // {"key": "value"}
)

type ObjectType string
