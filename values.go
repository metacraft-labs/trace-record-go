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
	IsValueRecord()
	// MarshalJson() ([]byte, error)
}

type NilValueRecord struct {
	Kind   string `json:"kind"`
	TypeId TypeId `json:"type_id"`
}

func (n NilValueRecord) IsValueRecord() {}

func NilValue() NilValueRecord {
	return NilValueRecord{"None", TypeId(0)}
}

type IntValueRecord struct {
	Kind   string `json:"kind"`
	I      int64  `json:"i"`
	TypeId TypeId `json:"type_id"`
}

func (i IntValueRecord) IsValueRecord() {}

func IntValue(i int64, typeId TypeId) IntValueRecord {
	return IntValueRecord{"Int", i, typeId}
}

type BoolValueRecord struct {
	Kind   string `json:"kind"`
	I      bool   `json:"i"`
	TypeId TypeId `json:"type_id"`
}

func (b BoolValueRecord) IsValueRecord() {}

func BoolValue(i bool, typeId TypeId) BoolValueRecord {
	return BoolValueRecord{"Bool", i, typeId}
}

type StringValueRecord struct {
	Kind   string `json:"kind"`
	I      string `json:"i"`
	TypeId TypeId `json:"type_id"`
}

func (b StringValueRecord) IsValueRecord() {}

func StringValue(i string, typeId TypeId) StringValueRecord {
	return StringValueRecord{"String", i, typeId}
}

type StructValueRecord struct {
	Kind   string `json:"kind"`
	Fields []ValueRecord
}

func (s StructValueRecord) IsValueRecord() {}

func StructValue(fields []ValueRecord) StructValueRecord {
	return StructValueRecord{"Struct", fields}
}

// func (receiver NilValueRecord) MarshalJson() ([]byte, error) {
// 	return json.Marshal(
// }
