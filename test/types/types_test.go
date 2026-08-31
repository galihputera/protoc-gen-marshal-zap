package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestTypes_MarshalLogObject(t *testing.T) {
	testTime := time.Date(2024, 1, 15, 10, 30, 0, 123456789, time.UTC)
	testDuration := 90*time.Second + 500*time.Millisecond

	anyVal, _ := anypb.New(&OtherType3{Val: "any_value"})

	m := &Types{
		SecretVal:   "secret",
		DoubleVal:   0.1,
		FloatVal:    0.1,
		Int32Val:    1,
		Int64Val:    1,
		Uint32Val:   1,
		Uint64Val:   1,
		Sint32Val:   2,
		Sint64Val:   2,
		Fixed32Val:  2,
		Fixed64Val:  2,
		Sfixed32Val: 3,
		Sfixed64Val: 3,
		BoolVal:     true,
		StringVal:   "string",
		BytesVal:    []byte{1, 2, 3},
		EnumVal:     Types_ENUM_1,
		OtherTypeVal: &OtherType1{
			OtherStringVal: "other_string",
			OtherSecretVal: "other_secret",
		},
		NestedTypeVal: &Types_NestedType{
			NestedStringVal: "nested_string",
			NestedSecretVal: "nested_secret",
		},
		OtherTypeNestedTypeVal: &OtherType2_NestedType{
			NestedStringVal: "other_nested_string",
			NestedSecretVal: "other_nested_secret",
		},
		OneofVal: &Types_OneofStringVal{
			OneofStringVal: "", // set zero value explicitly
		},
		MapVal1: map[string]string{
			"foo": "bar",
		},
		MapVal2: map[string]*OtherType1{
			"foo": {
				OtherStringVal: "other_string",
				OtherSecretVal: "other_secret",
			},
		},
		RepeatedVal1: []string{
			"foo", "bar",
		},
		RepeatedVal2: []Types_Enum{
			Types_ENUM_1, Types_ENUM_2,
		},
		RepeatedVal3: []*OtherType1{
			{
				OtherStringVal: "other_string",
				OtherSecretVal: "other_secret",
			},
		},
		RepeatedEmptyVal: []string{},
		StructVal: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"struct_string": {
					Kind: &structpb.Value_StringValue{StringValue: "foo"},
				},
				"struct_number": {
					Kind: &structpb.Value_NumberValue{NumberValue: 0.1},
				},
				"struct_bool": {
					Kind: &structpb.Value_BoolValue{BoolValue: true},
				},
			},
		},
		// Well-known types
		TimestampVal: timestamppb.New(testTime),
		DurationVal:  durationpb.New(testDuration),
		AnyVal:       anyVal,
		FieldMaskVal: &fieldmaskpb.FieldMask{Paths: []string{"field1", "field2", "nested.field"}},
		EmptyVal:     &emptypb.Empty{},
		ValueVal: &structpb.Value{
			Kind: &structpb.Value_StringValue{StringValue: "value_string"},
		},
		ListValueVal: &structpb.ListValue{
			Values: []*structpb.Value{
				{Kind: &structpb.Value_StringValue{StringValue: "item1"}},
				{Kind: &structpb.Value_NumberValue{NumberValue: 42}},
			},
		},
		// Wrapper types
		BoolWrapperVal:   wrapperspb.Bool(true),
		StringWrapperVal: wrapperspb.String("wrapped_string"),
		BytesWrapperVal:  wrapperspb.Bytes([]byte{4, 5, 6}),
		Int32WrapperVal:  wrapperspb.Int32(100),
		Int64WrapperVal:  wrapperspb.Int64(200),
		Uint32WrapperVal: wrapperspb.UInt32(300),
		Uint64WrapperVal: wrapperspb.UInt64(400),
		FloatWrapperVal:  wrapperspb.Float(1.5),
		DoubleWrapperVal: wrapperspb.Double(2.5),
		// Repeated well-known types
		RepeatedTimestampVal: []*timestamppb.Timestamp{
			timestamppb.New(testTime),
			timestamppb.New(testTime.Add(time.Hour)),
		},
		RepeatedDurationVal: []*durationpb.Duration{
			durationpb.New(time.Minute),
			durationpb.New(time.Hour),
		},
		RepeatedStringWrapperVal: []*wrapperspb.StringValue{
			wrapperspb.String("wrapped1"),
			wrapperspb.String("wrapped2"),
		},
		// Map with well-known types
		MapTimestampVal: map[string]*timestamppb.Timestamp{
			"time1": timestamppb.New(testTime),
		},
		MapStringWrapperVal: map[string]*wrapperspb.StringValue{
			"key1": wrapperspb.String("map_wrapped_string"),
		},
		XString:                   "foo",
		OptionalVal:               refs("foo"),
		OptionalNotPresentVal:     nil,
		OptionalEnum:              refs(Types_ENUM_1),
		OptionalNotPresentEnum:    nil,
		OptionalMessage:           &OtherType3{Val: "foo"},
		OptionalNotPresentMessage: nil,
		PresentMessage:            &OtherType3{Val: "foo"},
		NotPresentMessage:         nil,
	}
	enc := zapcore.NewMapObjectEncoder()
	err := m.MarshalLogObject(enc)
	if !assert.NoError(t, err) {
		return
	}

	// Check non-JSON fields with exact matching
	assert.Equal(t, "[MASKED]", enc.Fields["secret_val"])
	assert.Equal(t, "[MASKED]", enc.Fields["secret_val2"])
	assert.Equal(t, float64(0.1), enc.Fields["double_val"])
	assert.Equal(t, float32(0.1), enc.Fields["float_val"])
	assert.Equal(t, int32(1), enc.Fields["int32_val"])
	assert.Equal(t, int64(1), enc.Fields["int64_val"])
	assert.Equal(t, uint32(1), enc.Fields["uint32_val"])
	assert.Equal(t, uint64(1), enc.Fields["uint64_val"])
	assert.Equal(t, int32(2), enc.Fields["sint32_val"])
	assert.Equal(t, int64(2), enc.Fields["sint64_val"])
	assert.Equal(t, uint32(2), enc.Fields["fixed32_val"])
	assert.Equal(t, uint64(2), enc.Fields["fixed64_val"])
	assert.Equal(t, int32(3), enc.Fields["sfixed32_val"])
	assert.Equal(t, int64(3), enc.Fields["sfixed64_val"])
	assert.Equal(t, true, enc.Fields["bool_val"])
	assert.Equal(t, "string", enc.Fields["string_val"])
	assert.Equal(t, []byte{1, 2, 3}, enc.Fields["bytes_val"])
	assert.Equal(t, Types_ENUM_1.String(), enc.Fields["enum_val"])

	// Nested message types
	assert.EqualValues(t, map[string]interface{}{
		"other_string_val": "other_string",
		"other_secret_val": "[MASKED]",
	}, enc.Fields["other_type_val"])
	assert.EqualValues(t, map[string]interface{}{
		"nested_string_val": "nested_string",
		"nested_secret_val": "[MASKED]",
	}, enc.Fields["nested_type_val"])
	assert.EqualValues(t, map[string]interface{}{
		"nested_string_val": "other_nested_string",
		"nested_secret_val": "[MASKED]",
	}, enc.Fields["other_type_nested_type_val"])

	// Oneof
	assert.Equal(t, "", enc.Fields["oneof_string_val"])

	// Maps
	assert.EqualValues(t, map[string]interface{}{"foo": "bar"}, enc.Fields["map_val1"])
	assert.EqualValues(t, map[string]interface{}{
		"foo": map[string]interface{}{
			"other_string_val": "other_string",
			"other_secret_val": "[MASKED]",
		},
	}, enc.Fields["map_val2"])
	assert.EqualValues(t, map[string]interface{}{}, enc.Fields["map_empty_val"])

	// Repeated
	assert.EqualValues(t, []interface{}{"foo", "bar"}, enc.Fields["repeated_val1"])
	assert.EqualValues(t, []interface{}{Types_ENUM_1.String(), Types_ENUM_2.String()}, enc.Fields["repeated_val2"])
	assert.EqualValues(t, []interface{}{
		map[string]interface{}{
			"other_string_val": "other_string",
			"other_secret_val": "[MASKED]",
		},
	}, enc.Fields["repeated_val3"])
	assert.EqualValues(t, []interface{}{}, enc.Fields["repeated_empty_val"])

	// Well-known types serialized with AddReflected - compare actual protobuf objects
	assert.Equal(t, m.StructVal, enc.Fields["struct_val"])
	assert.Equal(t, m.AnyVal, enc.Fields["any_val"])
	assert.Equal(t, m.ValueVal, enc.Fields["value_val"])
	assert.Equal(t, m.ListValueVal, enc.Fields["list_value_val"])

	// Well-known types with native representations
	assert.Equal(t, testTime.Format("2006-01-02T15:04:05.999999999Z07:00"), enc.Fields["timestamp_val"])
	assert.Equal(t, testDuration.String(), enc.Fields["duration_val"])
	assert.Equal(t, "field1,field2,nested.field", enc.Fields["field_mask_val"])
	assert.Equal(t, struct{}{}, enc.Fields["empty_val"])

	// Wrapper types - unwrapped to native types
	assert.Equal(t, true, enc.Fields["bool_wrapper_val"])
	assert.Equal(t, "wrapped_string", enc.Fields["string_wrapper_val"])
	assert.Equal(t, []byte{4, 5, 6}, enc.Fields["bytes_wrapper_val"])
	assert.Equal(t, int32(100), enc.Fields["int32_wrapper_val"])
	assert.Equal(t, int64(200), enc.Fields["int64_wrapper_val"])
	assert.Equal(t, uint32(300), enc.Fields["uint32_wrapper_val"])
	assert.Equal(t, uint64(400), enc.Fields["uint64_wrapper_val"])
	assert.Equal(t, float32(1.5), enc.Fields["float_wrapper_val"])
	assert.Equal(t, float64(2.5), enc.Fields["double_wrapper_val"])

	// Repeated well-known types
	assert.EqualValues(t, []interface{}{
		testTime.Format("2006-01-02T15:04:05.999999999Z07:00"),
		testTime.Add(time.Hour).Format("2006-01-02T15:04:05.999999999Z07:00"),
	}, enc.Fields["repeated_timestamp_val"])
	assert.EqualValues(t, []interface{}{
		time.Minute.String(),
		time.Hour.String(),
	}, enc.Fields["repeated_duration_val"])
	assert.EqualValues(t, []interface{}{
		"wrapped1",
		"wrapped2",
	}, enc.Fields["repeated_string_wrapper_val"])

	// Map with well-known types
	assert.EqualValues(t, map[string]interface{}{
		"time1": testTime.Format("2006-01-02T15:04:05.999999999Z07:00"),
	}, enc.Fields["map_timestamp_val"])
	assert.EqualValues(t, map[string]interface{}{
		"key1": "map_wrapped_string",
	}, enc.Fields["map_string_wrapper_val"])

	// Other fields
	assert.Equal(t, "foo", enc.Fields["_String"])
	assert.Equal(t, "foo", enc.Fields["optional_val"])
	assert.Equal(t, Types_ENUM_1.String(), enc.Fields["optional_enum"])
	assert.EqualValues(t, map[string]interface{}{"val": "foo"}, enc.Fields["optional_message"])
	assert.EqualValues(t, map[string]interface{}{"val": "foo"}, enc.Fields["present_message"])

	assert.NotContains(t, enc.Fields, "oneof_int64_val")
	assert.NotContains(t, enc.Fields, "oneof_bool_val")
	assert.NotContains(t, enc.Fields, "optional_not_present_val")
	assert.NotContains(t, enc.Fields, "optional_not_present_enum")
	assert.NotContains(t, enc.Fields, "optional_not_present_message")
	assert.NotContains(t, enc.Fields, "not_present_message")
}

func refs[T any](v T) *T {
	return &v
}
