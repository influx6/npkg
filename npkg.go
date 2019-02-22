package npkg

//************************************************************
// Decoder
//************************************************************

// DecodableObject defines what we expect the a decodable type should provide.
type DecodableObject interface {
	DecodeOject(decoder Decoder)
}

// DecodableList defines what we expect the a decodable list type should provide.
type DecodableList interface {
	DecodeList(decoder ListDecoder)
}

// ListDecoder defines an interface for what we expect a array decoder
// to provide.
type ListDecoder interface {
	Int(v *int)
	Bool(v *bool)
	UInt(v *uint)
	Int8(v *int8)
	Int16(v *int16)
	Int32(v *int32)
	Int64(v *int64)
	UInt8(v *uint8)
	UInt16(v *uint16)
	UInt32(v *uint32)
	UInt64(v *uint64)
	String(v *string)
	Float64(v *float64)
	Float32(v *float32)
	Base64(v *int64, b int)

	List(DecodableList)
	Object(DecodableObject)
	ObjectWith(fn func(decoder Decoder))
	ListWith(fn func(decoder ListDecoder))
}

// Decoder defines an interface for what we expect a object decoder
// to provide for key value pairs.
type Decoder interface {
	Int(key string, v *int)
	UInt(key string, v *uint)
	Bool(key string, v *bool)
	Int8(key string, v *int8)
	Hex(key string, v *string)
	UInt8(key string, v *uint8)
	Int16(key string, v *int16)
	UInt16(key string, v *uint16)
	Int32(key string, v *int32)
	UInt32(key string, v *uint32)
	Int64(key string, v *int64)
	UInt64(key string, v *uint64)
	String(key string, v *string)
	Float64(key string, v *float64)
	Float32(key string, v *float32)
	Base64(key string, v *int64, b int)

	ListKey(key string, object DecodableList)
	ObjectKey(key string, object DecodableObject)
	ObjectFor(key string, fx func(decoder Decoder))
	ListFor(key string, fx func(decoder ListDecoder))
}

//************************************************************
// Encoder
//************************************************************

// EncodableObject defines what we expect the a Encodable type should provide.
type EncodableObject interface {
	EncodeObject(encoder Encoder)
}

// EncodableList defines what we expect the a Encodable list type should provide.
type EncodableList interface {
	EncodeList(encoder ListEncoder)
}

// ListEncoder defines an interface which defines methods for items into
// a underline list encoding.
type ListEncoder interface {
	AddInt(v int)
	AddBool(v bool)
	AddUInt(v uint)
	AddInt8(v int8)
	AddInt16(v int16)
	AddInt32(v int32)
	AddInt64(v int64)
	AddUInt8(v uint8)
	AddUInt16(v uint16)
	AddUInt32(v uint32)
	AddUInt64(v uint64)
	AddString(v string)
	AddFloat64(v float64)
	AddFloat32(v float32)
	AddBase64(v int64, b int)

	AddList(list EncodableList)
	AddObject(object EncodableObject)
	AddObjectWith(fn func(encoder Encoder))
	AddListWith(fn func(encoder ListEncoder))
}

// Encoder encodes giving key - value pairs as it's implementer
// decides providing a encapsulated formatting of data as it sees fit.
type Encoder interface {
	Int(k string, v int)
	UInt(k string, v uint)
	Bool(k string, v bool)
	Int8(k string, v int8)
	Hex(k string, v string)
	UInt8(k string, v uint8)
	Int16(k string, v int16)
	UInt16(k string, v uint16)
	Int32(k string, v int32)
	UInt32(k string, v uint32)
	Int64(k string, v int64)
	UInt64(k string, v uint64)
	String(k string, v string)
	Float64(k string, v float64)
	Float32(k string, v float32)
	Base64(k string, v int64, b int)

	List(k string, list EncodableList)
	ObjectFor(k string, fx func(Encoder))
	Object(k string, object EncodableObject)
	ListFor(k string, fx func(ListEncoder))
}

// Encodable defines a type which exposes a method to encode it's internals
// with provided encoder.
type Encodable interface {
	Encode(encoder Encoder)
}
