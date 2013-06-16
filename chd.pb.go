// Code generated by protoc-gen-go.
// source: chd.proto
// DO NOT EDIT!

package mph

import proto "code.google.com/p/goprotobuf/proto"
import json "encoding/json"
import math "math"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type CHDProto struct {
	R                []uint64             `protobuf:"varint,1,rep,name=r" json:"r,omitempty"`
	Indicies         []uint64             `protobuf:"varint,2,rep,name=indicies" json:"indicies,omitempty"`
	Table            []*CHDProto_KeyValue `protobuf:"bytes,6,rep,name=table" json:"table,omitempty"`
	XXX_unrecognized []byte               `json:"-"`
}

func (m *CHDProto) Reset()         { *m = CHDProto{} }
func (m *CHDProto) String() string { return proto.CompactTextString(m) }
func (*CHDProto) ProtoMessage()    {}

func (m *CHDProto) GetR() []uint64 {
	if m != nil {
		return m.R
	}
	return nil
}

func (m *CHDProto) GetIndicies() []uint64 {
	if m != nil {
		return m.Indicies
	}
	return nil
}

func (m *CHDProto) GetTable() []*CHDProto_KeyValue {
	if m != nil {
		return m.Table
	}
	return nil
}

type CHDProto_KeyValue struct {
	Key              []byte `protobuf:"bytes,1,req,name=key" json:"key,omitempty"`
	Value            []byte `protobuf:"bytes,2,req,name=value" json:"value,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *CHDProto_KeyValue) Reset()         { *m = CHDProto_KeyValue{} }
func (m *CHDProto_KeyValue) String() string { return proto.CompactTextString(m) }
func (*CHDProto_KeyValue) ProtoMessage()    {}

func (m *CHDProto_KeyValue) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *CHDProto_KeyValue) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func init() {
}
