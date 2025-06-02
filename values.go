package trace_record

// import "encoding/json"

type TypeId uint64

type TypeKind uint8

// TODO
// type TypeKind enum {
// None
// }

type TypeSpecificInfo interface {
	IsTypeSpecificInfo()
}

type NoneTypeSpecificInfo struct {
	Kind string `json:"kind"`
}

func (i NoneTypeSpecificInfo) IsTypeSpecificInfo() {}

func NewNonTypeSpecificInfo() NoneTypeSpecificInfo {
	return NoneTypeSpecificInfo{"None"}
}

const INT_TYPE_KIND = TypeKind(7)
const FLOAT_TYPE_KIND = TypeKind(8)
const POINTER_TYPE_KIND = TypeKind(23)
const TUPLE_TYPE_KIND = TypeKind(27)
const ARRAY_TYPE_KIND = TypeKind(4)
const SLICE_TYPE_KIND = TypeKind(33)
const BOOL_TYPE_KIND = TypeKind(12)
const STRING_TYPE_KIND = TypeKind(9)
const STRUCT_TYPE_KIND = TypeKind(6)

type TypeRecord struct {
	Kind         TypeKind         `json:"kind"`
	LangType     string           `json:"lang_type"`
	SpecificInfo TypeSpecificInfo `json:"specific_info"`
}

func NewSimpleTypeRecord(kind TypeKind, langType string) TypeRecord {
	return TypeRecord{kind, langType, NewNonTypeSpecificInfo()}
}

func NewTypeRecord(kind TypeKind, langType string, specificInfo TypeSpecificInfo) TypeRecord {
	return TypeRecord{kind, langType, specificInfo}
}

type FieldTypeRecord struct {
	Name   string `json:"name"`
	TypeId TypeId `json:"type_id"`
}

func NewFieldTypeRecord(name string, typeId TypeId) FieldTypeRecord {
	return FieldTypeRecord{name, typeId}
}

type StructTypeInfo struct {
	Kind   string            `json:"kind"`
	Fields []FieldTypeRecord `json:"fields"`
}

func (i StructTypeInfo) IsTypeSpecificInfo() {}

func NewStructTypeInfo(fields []FieldTypeRecord) StructTypeInfo {
	return StructTypeInfo{Kind: "Struct", Fields: fields}
}

type PointerTypeInfo struct {
	Kind              string `json:"kind"`
	DereferenceTypeId TypeId `json:"dereference_type_id"`
}

func (i PointerTypeInfo) IsTypeSpecificInfo() {}

func NewPointerTypeInfo(typeId TypeId) PointerTypeInfo {
	return PointerTypeInfo{"Pointer", typeId}
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

type FloatValueRecord struct {
	Kind   string  `json:"kind"`
	F      float64 `json:"f"`
	TypeId TypeId  `json:"type_id"`
}

func (i FloatValueRecord) IsValueRecord() {}

func FloatValue(f float64, typeId TypeId) FloatValueRecord {
	return FloatValueRecord{"Float", f, typeId}
}

type BoolValueRecord struct {
	Kind   string `json:"kind"`
	B      bool   `json:"b"`
	TypeId TypeId `json:"type_id"`
}

func (b BoolValueRecord) IsValueRecord() {}

func BoolValue(b bool, typeId TypeId) BoolValueRecord {
	return BoolValueRecord{"Bool", b, typeId}
}

type StringValueRecord struct {
	Kind   string `json:"kind"`
	Text   string `json:"text"`
	TypeId TypeId `json:"type_id"`
}

func (s StringValueRecord) IsValueRecord() {}

func StringValue(text string, typeId TypeId) StringValueRecord {
	return StringValueRecord{"String", text, typeId}
}

type StructValueRecord struct {
	Kind   string        `json:"kind"`
	Fields []ValueRecord `json:"field_values"`
	TypeId TypeId        `json:"type_id"`
}

func (s StructValueRecord) IsValueRecord() {}

func StructValue(fields []ValueRecord, typeId TypeId) StructValueRecord {
	return StructValueRecord{"Struct", fields, typeId}
}

type SequenceValueRecord struct {
	Kind     string        `json:"kind"`
	Elements []ValueRecord `json:"elements"`
	IsSlice  bool          `json:"is_slice"`
	TypeId   TypeId        `json:"type_id"`
}

func (s SequenceValueRecord) IsValueRecord() {}

func SequenceValue(elements []ValueRecord, isSlice bool, typeId TypeId) SequenceValueRecord {
	return SequenceValueRecord{"Sequence", elements, isSlice, typeId}
}

type ReferenceValueRecord struct {
	Kind         string      `json:"kind"`
	Dereferenced ValueRecord `json:"dereferenced"`
	Address      uint32      `json:"address"`
	Mutable      bool        `json:"mutable"`
	TypeId       TypeId      `json:"type_id"`
}

func (s ReferenceValueRecord) IsValueRecord() {}

func ReferenceValue(dereferenced ValueRecord, address uint32, mutable bool, typeId TypeId) ReferenceValueRecord {
	return ReferenceValueRecord{"Reference", dereferenced, address, mutable, typeId}
}

type TupleValueRecord struct {
	Kind     string        `json:"kind"`
	Elements []ValueRecord `json:"elements"`
	TypeId   TypeId        `json:"type_id"`
}

func (s TupleValueRecord) IsValueRecord() {}

func TupleValue(elements []ValueRecord, typeId TypeId) TupleValueRecord {
	return TupleValueRecord{"Tuple", elements, typeId}
}
