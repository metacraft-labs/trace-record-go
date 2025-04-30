package trace_record

// import "encoding/json"

type TypeId uint64

type TypeKind uint8

// TODO
// type TypeKind enum {
// None
// }

type TypeSpecificInfo interface {
	IsTypeSpecificInfo() bool
}

type NoneTypeSpecificInfo struct {
	Kind string `json:"kind"`
}

func (i NoneTypeSpecificInfo) IsTypeSpecificInfo() bool {
	return true
}

func NewNonTypeSpecificInfo() NoneTypeSpecificInfo {
	return NoneTypeSpecificInfo{"None"}
}

const INT_TYPE_KIND = TypeKind(7)

type TypeRecord struct {
	Kind         TypeKind         `json:"kind"`
	LangType     string           `json:"lang_type"`
	SpecificInfo TypeSpecificInfo `json:"specific_info"`
}

func NewSimpleTypeRecord(kind TypeKind, langType string) TypeRecord {
	return TypeRecord{kind, langType, NewNonTypeSpecificInfo()}
}

type ValueRecord interface {
	IsValueRecord() bool
	// MarshalJson() ([]byte, error)
}

type NilValueRecord struct {
	Kind   string `json:"kind"`
	TypeId TypeId `json:"type_id"`
}

func (n NilValueRecord) IsValueRecord() bool {
	return true
}

func NilValue() NilValueRecord {
	return NilValueRecord{"None", TypeId(0)}
}

type IntValueRecord struct {
	Kind   string `json:"kind"`
	I      int64  `json:"i"`
	TypeId TypeId `json:"type_id"`
}

func (i IntValueRecord) IsValueRecord() bool {
	return true
}

func IntValue(i int64, typeId TypeId) IntValueRecord {
	return IntValueRecord{"Int", i, typeId}
}

// func (receiver NilValueRecord) MarshalJson() ([]byte, error) {
// 	return json.Marshal(
// }
