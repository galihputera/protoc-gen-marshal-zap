package main

import (
	"fmt"

	pbzap "github.com/kei2100/protoc-gen-marshal-zap"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

const (
	zapcorePkg = protogen.GoImportPath("go.uber.org/zap/zapcore")
	fmtPkg     = protogen.GoImportPath("fmt")
	stringsPkg = protogen.GoImportPath("strings")
)

// wellKnownType represents the category of a well-known protobuf type
type wellKnownType int

const (
	notWellKnown wellKnownType = iota
	wktTimestamp
	wktDuration
	wktStruct
	wktValue
	wktListValue
	wktBoolValue
	wktStringValue
	wktBytesValue
	wktInt32Value
	wktInt64Value
	wktUInt32Value
	wktUInt64Value
	wktFloatValue
	wktDoubleValue
	wktAny
	wktFieldMask
	wktEmpty
)

// getWellKnownType returns the well-known type category for a message field
func getWellKnownType(f *protogen.Field) wellKnownType {
	if f.Desc.Kind() != protoreflect.MessageKind {
		return notWellKnown
	}
	fullName := string(f.Desc.Message().FullName())
	switch fullName {
	case "google.protobuf.Timestamp":
		return wktTimestamp
	case "google.protobuf.Duration":
		return wktDuration
	case "google.protobuf.Struct":
		return wktStruct
	case "google.protobuf.Value":
		return wktValue
	case "google.protobuf.ListValue":
		return wktListValue
	case "google.protobuf.BoolValue":
		return wktBoolValue
	case "google.protobuf.StringValue":
		return wktStringValue
	case "google.protobuf.BytesValue":
		return wktBytesValue
	case "google.protobuf.Int32Value":
		return wktInt32Value
	case "google.protobuf.Int64Value":
		return wktInt64Value
	case "google.protobuf.UInt32Value":
		return wktUInt32Value
	case "google.protobuf.UInt64Value":
		return wktUInt64Value
	case "google.protobuf.FloatValue":
		return wktFloatValue
	case "google.protobuf.DoubleValue":
		return wktDoubleValue
	case "google.protobuf.Any":
		return wktAny
	case "google.protobuf.FieldMask":
		return wktFieldMask
	case "google.protobuf.Empty":
		return wktEmpty
	default:
		return notWellKnown
	}
}

// getWellKnownTypeFromMapValue returns the well-known type category for a map value
func getWellKnownTypeFromMapValue(f *protogen.Field) wellKnownType {
	if f.Desc.MapValue().Kind() != protoreflect.MessageKind {
		return notWellKnown
	}
	fullName := string(f.Desc.MapValue().Message().FullName())
	switch fullName {
	case "google.protobuf.Timestamp":
		return wktTimestamp
	case "google.protobuf.Duration":
		return wktDuration
	case "google.protobuf.Struct":
		return wktStruct
	case "google.protobuf.Value":
		return wktValue
	case "google.protobuf.ListValue":
		return wktListValue
	case "google.protobuf.BoolValue":
		return wktBoolValue
	case "google.protobuf.StringValue":
		return wktStringValue
	case "google.protobuf.BytesValue":
		return wktBytesValue
	case "google.protobuf.Int32Value":
		return wktInt32Value
	case "google.protobuf.Int64Value":
		return wktInt64Value
	case "google.protobuf.UInt32Value":
		return wktUInt32Value
	case "google.protobuf.UInt64Value":
		return wktUInt64Value
	case "google.protobuf.FloatValue":
		return wktFloatValue
	case "google.protobuf.DoubleValue":
		return wktDoubleValue
	case "google.protobuf.Any":
		return wktAny
	case "google.protobuf.FieldMask":
		return wktFieldMask
	case "google.protobuf.Empty":
		return wktEmpty
	default:
		return notWellKnown
	}
}

func generateListField(g *protogen.GeneratedFile, f *protogen.Field) {
	fname := f.Desc.Name()
	g.P(fname, "ArrMarshaller := func(enc ", g.QualifiedGoIdent(zapcorePkg.Ident("ArrayEncoder")), ") error {")
	g.P("for _, v := range x.", f.GoName, " {")
	switch f.Desc.Kind() {
	case protoreflect.BoolKind:
		g.P("enc.AppendBool(v)")
	case protoreflect.BytesKind:
		g.P("enc.AppendByteString(v)")
	case protoreflect.DoubleKind:
		g.P("enc.AppendFloat64(v)")
	case protoreflect.EnumKind:
		g.P("enc.AppendString(v.String())")
	case protoreflect.Fixed32Kind, protoreflect.Uint32Kind:
		g.P("enc.AppendUint32(v)")
	case protoreflect.Fixed64Kind, protoreflect.Uint64Kind:
		g.P("enc.AppendUint64(v)")
	case protoreflect.FloatKind:
		g.P("enc.AppendFloat32(v)")
	case protoreflect.Int32Kind, protoreflect.Sfixed32Kind, protoreflect.Sint32Kind:
		g.P("enc.AppendInt32(v)")
	case protoreflect.Int64Kind, protoreflect.Sfixed64Kind, protoreflect.Sint64Kind:
		g.P("enc.AppendInt64(v)")
	case protoreflect.GroupKind:
		g.P("enc.AppendReflected(v)")
	case protoreflect.MessageKind:
		generateListWellKnownOrMessage(g, f)
	case protoreflect.StringKind:
		g.P("enc.AppendString(v)")
	default:
		g.P("enc.AppendReflected(v)")
	}
	g.P("}")
	g.P("return nil")
	g.P("}")
	g.P("enc.AddArray(\"", fname, "\",", g.QualifiedGoIdent(zapcorePkg.Ident("ArrayMarshalerFunc")), "(", fname, "ArrMarshaller))")
}

// generateListWellKnownOrMessage generates array append code for message kinds,
// handling well-known types with their JSON representation
func generateListWellKnownOrMessage(g *protogen.GeneratedFile, f *protogen.Field) {
	wkt := getWellKnownType(f)
	switch wkt {
	case wktTimestamp:
		// Timestamp -> RFC 3339 string
		g.P("if v != nil {")
		g.P("enc.AppendString(v.AsTime().Format(\"2006-01-02T15:04:05.999999999Z07:00\"))")
		g.P("} else {")
		g.P("enc.AppendReflected(nil)")
		g.P("}")
	case wktDuration:
		// Duration -> string like "1.5s"
		g.P("if v != nil {")
		g.P("enc.AppendString(v.AsDuration().String())")
		g.P("} else {")
		g.P("enc.AppendReflected(nil)")
		g.P("}")
	case wktStruct, wktValue, wktListValue, wktAny:
		// Use AddReflected for complex types
		g.P("enc.AppendReflected(v)")
	case wktFieldMask:
		// FieldMask -> comma-separated paths
		g.P("if v != nil {")
		g.P("enc.AppendString(", g.QualifiedGoIdent(stringsPkg.Ident("Join")), "(v.GetPaths(), \",\"))")
		g.P("} else {")
		g.P("enc.AppendReflected(nil)")
		g.P("}")
	case wktEmpty:
		// Empty -> empty object representation
		g.P("enc.AppendReflected(struct{}{})")
	case wktBoolValue:
		g.P("if v != nil {")
		g.P("enc.AppendBool(v.GetValue())")
		g.P("} else {")
		g.P("enc.AppendReflected(nil)")
		g.P("}")
	case wktStringValue:
		g.P("if v != nil {")
		g.P("enc.AppendString(v.GetValue())")
		g.P("} else {")
		g.P("enc.AppendReflected(nil)")
		g.P("}")
	case wktBytesValue:
		g.P("if v != nil {")
		g.P("enc.AppendByteString(v.GetValue())")
		g.P("} else {")
		g.P("enc.AppendReflected(nil)")
		g.P("}")
	case wktInt32Value:
		g.P("if v != nil {")
		g.P("enc.AppendInt32(v.GetValue())")
		g.P("} else {")
		g.P("enc.AppendReflected(nil)")
		g.P("}")
	case wktInt64Value:
		g.P("if v != nil {")
		g.P("enc.AppendInt64(v.GetValue())")
		g.P("} else {")
		g.P("enc.AppendReflected(nil)")
		g.P("}")
	case wktUInt32Value:
		g.P("if v != nil {")
		g.P("enc.AppendUint32(v.GetValue())")
		g.P("} else {")
		g.P("enc.AppendReflected(nil)")
		g.P("}")
	case wktUInt64Value:
		g.P("if v != nil {")
		g.P("enc.AppendUint64(v.GetValue())")
		g.P("} else {")
		g.P("enc.AppendReflected(nil)")
		g.P("}")
	case wktFloatValue:
		g.P("if v != nil {")
		g.P("enc.AppendFloat32(v.GetValue())")
		g.P("} else {")
		g.P("enc.AppendReflected(nil)")
		g.P("}")
	case wktDoubleValue:
		g.P("if v != nil {")
		g.P("enc.AppendFloat64(v.GetValue())")
		g.P("} else {")
		g.P("enc.AppendReflected(nil)")
		g.P("}")
	default:
		// Regular message - check for ObjectMarshaler interface
		g.P("if obj, ok := interface{}(v).(", g.QualifiedGoIdent(zapcorePkg.Ident("ObjectMarshaler")), "); ok {")
		g.P("enc.AppendObject(obj)")
		g.P("} else {")
		g.P("enc.AppendReflected(v)")
		g.P("}")
	}
}

func generateMapField(g *protogen.GeneratedFile, f *protogen.Field) {
	fname := f.Desc.Name()
	g.P("enc.AddObject(\"", fname, "\", ", g.QualifiedGoIdent(zapcorePkg.Ident("ObjectMarshalerFunc")), "(func(enc ", g.QualifiedGoIdent(zapcorePkg.Ident("ObjectEncoder")), ") error {")
	g.P("for k, v := range x.", f.GoName, " {")
	switch f.Desc.MapValue().Kind() {
	case protoreflect.BoolKind:
		g.P("enc.AddBool(", g.QualifiedGoIdent(fmtPkg.Ident("Sprintf")), "(\"%v\", k), v)")
	case protoreflect.BytesKind:
		g.P("enc.AddBinary(", g.QualifiedGoIdent(fmtPkg.Ident("Sprintf")), "(\"%v\", k), v)")
	case protoreflect.DoubleKind:
		g.P("enc.AddFloat64(", g.QualifiedGoIdent(fmtPkg.Ident("Sprintf")), "(\"%v\", k), v)")
	case protoreflect.EnumKind:
		g.P("enc.AddString(", g.QualifiedGoIdent(fmtPkg.Ident("Sprintf")), "(\"%v\", k), v.String())")
	case protoreflect.Fixed32Kind, protoreflect.Uint32Kind:
		g.P("enc.AddUint32(", g.QualifiedGoIdent(fmtPkg.Ident("Sprintf")), "(\"%v\", k), v)")
	case protoreflect.Fixed64Kind, protoreflect.Uint64Kind:
		g.P("enc.AddUint64(", g.QualifiedGoIdent(fmtPkg.Ident("Sprintf")), "(\"%v\", k), v)")
	case protoreflect.FloatKind:
		g.P("enc.AddFloat32(", g.QualifiedGoIdent(fmtPkg.Ident("Sprintf")), "(\"%v\", k), v)")
	case protoreflect.Int32Kind, protoreflect.Sfixed32Kind, protoreflect.Sint32Kind:
		g.P("enc.AddInt32(", g.QualifiedGoIdent(fmtPkg.Ident("Sprintf")), "(\"%v\", k), v)")
	case protoreflect.Int64Kind, protoreflect.Sfixed64Kind, protoreflect.Sint64Kind:
		g.P("enc.AddInt64(", g.QualifiedGoIdent(fmtPkg.Ident("Sprintf")), "(\"%v\", k), v)")
	case protoreflect.GroupKind:
		g.P("enc.AddReflected(", g.QualifiedGoIdent(fmtPkg.Ident("Sprintf")), "(\"%v\", k), v)")
	case protoreflect.MessageKind:
		generateMapWellKnownOrMessage(g, f)
	case protoreflect.StringKind:
		g.P("enc.AddString(", g.QualifiedGoIdent(fmtPkg.Ident("Sprintf")), "(\"%v\", k), v)")
	default:
		g.P("enc.AddReflected(", g.QualifiedGoIdent(fmtPkg.Ident("Sprintf")), "(\"%v\", k), v)")
	}
	g.P("}")
	g.P("return nil")
	g.P("}))")
}

// generateMapWellKnownOrMessage generates map value encoding for message kinds,
// handling well-known types with their JSON representation
func generateMapWellKnownOrMessage(g *protogen.GeneratedFile, f *protogen.Field) {
	wkt := getWellKnownTypeFromMapValue(f)
	keyFmt := g.QualifiedGoIdent(fmtPkg.Ident("Sprintf")) + "(\"%v\", k)"
	switch wkt {
	case wktTimestamp:
		// Timestamp -> RFC 3339 string
		g.P("if v != nil {")
		g.P("enc.AddString(", keyFmt, ", v.AsTime().Format(\"2006-01-02T15:04:05.999999999Z07:00\"))")
		g.P("} else {")
		g.P("enc.AddReflected(", keyFmt, ", nil)")
		g.P("}")
	case wktDuration:
		// Duration -> string like "1.5s"
		g.P("if v != nil {")
		g.P("enc.AddString(", keyFmt, ", v.AsDuration().String())")
		g.P("} else {")
		g.P("enc.AddReflected(", keyFmt, ", nil)")
		g.P("}")
	case wktStruct, wktValue, wktListValue, wktAny:
		// Use AddReflected for complex types
		g.P("enc.AddReflected(", keyFmt, ", v)")
	case wktFieldMask:
		// FieldMask -> comma-separated paths
		g.P("if v != nil {")
		g.P("enc.AddString(", keyFmt, ", ", g.QualifiedGoIdent(stringsPkg.Ident("Join")), "(v.GetPaths(), \",\"))")
		g.P("} else {")
		g.P("enc.AddReflected(", keyFmt, ", nil)")
		g.P("}")
	case wktEmpty:
		// Empty -> empty object representation
		g.P("enc.AddReflected(", keyFmt, ", struct{}{})")
	case wktBoolValue:
		g.P("if v != nil {")
		g.P("enc.AddBool(", keyFmt, ", v.GetValue())")
		g.P("} else {")
		g.P("enc.AddReflected(", keyFmt, ", nil)")
		g.P("}")
	case wktStringValue:
		g.P("if v != nil {")
		g.P("enc.AddString(", keyFmt, ", v.GetValue())")
		g.P("} else {")
		g.P("enc.AddReflected(", keyFmt, ", nil)")
		g.P("}")
	case wktBytesValue:
		g.P("if v != nil {")
		g.P("enc.AddBinary(", keyFmt, ", v.GetValue())")
		g.P("} else {")
		g.P("enc.AddReflected(", keyFmt, ", nil)")
		g.P("}")
	case wktInt32Value:
		g.P("if v != nil {")
		g.P("enc.AddInt32(", keyFmt, ", v.GetValue())")
		g.P("} else {")
		g.P("enc.AddReflected(", keyFmt, ", nil)")
		g.P("}")
	case wktInt64Value:
		g.P("if v != nil {")
		g.P("enc.AddInt64(", keyFmt, ", v.GetValue())")
		g.P("} else {")
		g.P("enc.AddReflected(", keyFmt, ", nil)")
		g.P("}")
	case wktUInt32Value:
		g.P("if v != nil {")
		g.P("enc.AddUint32(", keyFmt, ", v.GetValue())")
		g.P("} else {")
		g.P("enc.AddReflected(", keyFmt, ", nil)")
		g.P("}")
	case wktUInt64Value:
		g.P("if v != nil {")
		g.P("enc.AddUint64(", keyFmt, ", v.GetValue())")
		g.P("} else {")
		g.P("enc.AddReflected(", keyFmt, ", nil)")
		g.P("}")
	case wktFloatValue:
		g.P("if v != nil {")
		g.P("enc.AddFloat32(", keyFmt, ", v.GetValue())")
		g.P("} else {")
		g.P("enc.AddReflected(", keyFmt, ", nil)")
		g.P("}")
	case wktDoubleValue:
		g.P("if v != nil {")
		g.P("enc.AddFloat64(", keyFmt, ", v.GetValue())")
		g.P("} else {")
		g.P("enc.AddReflected(", keyFmt, ", nil)")
		g.P("}")
	default:
		// Regular message - check for ObjectMarshaler interface
		g.P("if obj, ok := interface{}(v).(", g.QualifiedGoIdent(zapcorePkg.Ident("ObjectMarshaler")), "); ok {")
		g.P("enc.AddObject(", keyFmt, ", obj)")
		g.P("} else {")
		g.P("enc.AddReflected(", keyFmt, ", v)")
		g.P("}")
	}
}

func generatePrimitiveField(g *protogen.GeneratedFile, f *protogen.Field) {
	fname := f.Desc.Name()
	var gname string
	if f.Oneof != nil {
		gname = fmt.Sprintf("Get%s()", f.GoName)
	} else {
		gname = f.GoName
	}
	switch f.Desc.Kind() {
	case protoreflect.BoolKind:
		g.P("enc.AddBool(\"", fname, "\", x.", gname, ")")
	case protoreflect.BytesKind:
		g.P("enc.AddBinary(\"", fname, "\", x.", gname, ")")
	case protoreflect.DoubleKind:
		g.P("enc.AddFloat64(\"", fname, "\", x.", gname, ")")
	case protoreflect.EnumKind:
		g.P("enc.AddString(\"", fname, "\", x.", gname, ".String())")
	case protoreflect.Fixed32Kind, protoreflect.Uint32Kind:
		g.P("enc.AddUint32(\"", fname, "\", x.", gname, ")")
	case protoreflect.Fixed64Kind, protoreflect.Uint64Kind:
		g.P("enc.AddUint64(\"", fname, "\", x.", gname, ")")
	case protoreflect.FloatKind:
		g.P("enc.AddFloat32(\"", fname, "\", x.", gname, ")")
	case protoreflect.Int32Kind, protoreflect.Sfixed32Kind, protoreflect.Sint32Kind:
		g.P("enc.AddInt32(\"", fname, "\", x.", gname, ")")
	case protoreflect.Int64Kind, protoreflect.Sfixed64Kind, protoreflect.Sint64Kind:
		g.P("enc.AddInt64(\"", fname, "\", x.", gname, ")")
	case protoreflect.GroupKind:
		g.P("enc.AddReflected(\"", fname, "\", x.", gname, ")")
	case protoreflect.MessageKind:
		generatePrimitiveWellKnownOrMessage(g, f, fname, gname)
	case protoreflect.StringKind:
		g.P("enc.AddString(\"", fname, "\", x.", gname, ")")
	default:
		g.P("enc.AddReflected(\"", fname, "\", x.", gname, ")")
	}
}

// generatePrimitiveWellKnownOrMessage generates field encoding for message kinds,
// handling well-known types with their JSON representation
func generatePrimitiveWellKnownOrMessage(g *protogen.GeneratedFile, f *protogen.Field, fname protoreflect.Name, gname string) {
	wkt := getWellKnownType(f)
	switch wkt {
	case wktTimestamp:
		// Timestamp -> RFC 3339 string
		g.P("if x.", gname, " != nil {")
		g.P("enc.AddString(\"", fname, "\", x.", gname, ".AsTime().Format(\"2006-01-02T15:04:05.999999999Z07:00\"))")
		g.P("}")
	case wktDuration:
		// Duration -> string like "1.5s"
		g.P("if x.", gname, " != nil {")
		g.P("enc.AddString(\"", fname, "\", x.", gname, ".AsDuration().String())")
		g.P("}")
	case wktStruct, wktValue, wktListValue, wktAny:
		// Use AddReflected for complex types
		g.P("enc.AddReflected(\"", fname, "\", x.", gname, ")")
	case wktFieldMask:
		// FieldMask -> comma-separated paths
		g.P("if x.", gname, " != nil {")
		g.P("enc.AddString(\"", fname, "\", ", g.QualifiedGoIdent(stringsPkg.Ident("Join")), "(x.", gname, ".GetPaths(), \",\"))")
		g.P("}")
	case wktEmpty:
		// Empty -> empty object representation
		g.P("if x.", gname, " != nil {")
		g.P("enc.AddReflected(\"", fname, "\", struct{}{})")
		g.P("}")
	case wktBoolValue:
		g.P("if x.", gname, " != nil {")
		g.P("enc.AddBool(\"", fname, "\", x.", gname, ".GetValue())")
		g.P("}")
	case wktStringValue:
		g.P("if x.", gname, " != nil {")
		g.P("enc.AddString(\"", fname, "\", x.", gname, ".GetValue())")
		g.P("}")
	case wktBytesValue:
		g.P("if x.", gname, " != nil {")
		g.P("enc.AddBinary(\"", fname, "\", x.", gname, ".GetValue())")
		g.P("}")
	case wktInt32Value:
		g.P("if x.", gname, " != nil {")
		g.P("enc.AddInt32(\"", fname, "\", x.", gname, ".GetValue())")
		g.P("}")
	case wktInt64Value:
		g.P("if x.", gname, " != nil {")
		g.P("enc.AddInt64(\"", fname, "\", x.", gname, ".GetValue())")
		g.P("}")
	case wktUInt32Value:
		g.P("if x.", gname, " != nil {")
		g.P("enc.AddUint32(\"", fname, "\", x.", gname, ".GetValue())")
		g.P("}")
	case wktUInt64Value:
		g.P("if x.", gname, " != nil {")
		g.P("enc.AddUint64(\"", fname, "\", x.", gname, ".GetValue())")
		g.P("}")
	case wktFloatValue:
		g.P("if x.", gname, " != nil {")
		g.P("enc.AddFloat32(\"", fname, "\", x.", gname, ".GetValue())")
		g.P("}")
	case wktDoubleValue:
		g.P("if x.", gname, " != nil {")
		g.P("enc.AddFloat64(\"", fname, "\", x.", gname, ".GetValue())")
		g.P("}")
	default:
		// Regular message - check for ObjectMarshaler interface
		g.P("if obj, ok := interface{}(x.", gname, ").(", g.QualifiedGoIdent(zapcorePkg.Ident("ObjectMarshaler")), "); ok {")
		g.P("enc.AddObject(\"", fname, "\", obj)")
		g.P("} else {")
		g.P("enc.AddReflected(\"", fname, "\", x.", gname, ")")
		g.P("}")
	}
}

func isMasked(opts *descriptorpb.FieldOptions) bool {
	return proto.GetExtension(opts, pbzap.E_Mask).(bool) || opts.GetDebugRedact()
}

func handleExplicitPresence(g *protogen.GeneratedFile, f *protogen.Field, generateFunc func(*protogen.GeneratedFile, *protogen.Field)) {
	// Omit the fields that are defined as `Explicit Presence` and the value is not present.
	// https://protobuf.dev/programming-guides/field_presence/#presence-in-proto3-apis
	switch {
	case f.Oneof != nil && f.Desc.HasOptionalKeyword():
		// handle optional fields
		g.P("if x.", f.GoName, " != nil {")
		defer g.P("}")
	case f.Oneof != nil && !f.Desc.HasOptionalKeyword():
		// handle oneof fields
		g.P("if _, ok := x.Get", f.Oneof.GoName, "().(*", f.GoIdent, "); ok {")
		defer g.P("}")
	case f.Desc.Kind() == protoreflect.MessageKind || f.Desc.Kind() == protoreflect.GroupKind:
		// handle message fields
		g.P("if x.", f.GoName, " != nil {")
		defer g.P("}")
	}
	generateFunc(g, f)
}

func generateMessage(g *protogen.GeneratedFile, m *protogen.Message) {
	ident := g.QualifiedGoIdent(m.GoIdent)
	g.P("func (x *", ident, ") MarshalLogObject(enc ", g.QualifiedGoIdent(zapcorePkg.Ident("ObjectEncoder")), ") error {")
	g.P("if x == nil {")
	g.P("return nil")
	g.P("}")
	g.P()
	for _, f := range m.Fields {
		if isMasked(f.Desc.Options().(*descriptorpb.FieldOptions)) {
			g.P("enc.AddString(\"", f.Desc.Name(), "\", \"[MASKED]\")")
		} else if f.Desc.IsList() {
			generateListField(g, f)
		} else if f.Desc.IsMap() {
			generateMapField(g, f)
		} else {
			handleExplicitPresence(g, f, generatePrimitiveField)
		}
		g.P()
	}
	g.P("return nil")
	g.P("}")
	g.P()
	for _, submsg := range m.Messages {
		if submsg.Desc.IsMapEntry() {
			continue
		}
		generateMessage(g, submsg)
	}
}

func generateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	if len(file.Messages) == 0 {
		return nil
	}

	filename := fmt.Sprintf("%s.pb.marshal-zap.go", file.GeneratedFilenamePrefix)
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("// Code generated by protoc-gen-marshal-zap. DO NOT EDIT.")
	g.P("//")
	g.P("// source: ", file.Desc.Path())
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()

	for _, m := range file.Messages {
		generateMessage(g, m)
	}

	return g
}

func main() {
	protogen.Options{}.Run(func(plugin *protogen.Plugin) error {
		plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, file := range plugin.FilesByPath {
			if !file.Generate {
				continue
			}
			generateFile(plugin, file)
		}
		return nil
	})
}
