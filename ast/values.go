/*
Copyright 2024 The Apex Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ast

import "github.com/apexlang/apex-go/kinds"

type Value interface {
	Node
	GetValue() interface{}
}

// Ensure that all value types implements Value interface
var _ Value = (*IntValue)(nil)
var _ Value = (*FloatValue)(nil)
var _ Value = (*StringValue)(nil)
var _ Value = (*BooleanValue)(nil)
var _ Value = (*EnumValue)(nil)
var _ Value = (*ListValue)(nil)
var _ Value = (*ObjectValue)(nil)

// IntValue implements Node, Value
type IntValue struct {
	BaseNode
	Value int `json:"value"`
}

func NewIntValue(loc *Location, value int) *IntValue {
	return &IntValue{
		BaseNode: BaseNode{kinds.IntValue, loc},
		Value:    value,
	}
}

func (v *IntValue) GetValue() interface{} {
	return v.Value
}

// FloatValue implements Node, Value
type FloatValue struct {
	BaseNode
	Value float64 `json:"value"`
}

func NewFloatValue(loc *Location, value float64) *FloatValue {
	return &FloatValue{
		BaseNode: BaseNode{kinds.FloatValue, loc},
		Value:    value,
	}
}

func (v *FloatValue) GetValue() interface{} {
	return v.Value
}

// StringValue implements Node, Value
type StringValue struct {
	BaseNode
	Value string `json:"value"`
}

func NewStringValue(loc *Location, value string) *StringValue {
	return &StringValue{
		BaseNode: BaseNode{kinds.StringValue, loc},
		Value:    value,
	}
}

func (v *StringValue) GetValue() interface{} {
	return v.Value
}

// BooleanValue implements Node, Value
type BooleanValue struct {
	BaseNode
	Value bool `json:"value"`
}

func NewBooleanValue(loc *Location, value bool) *BooleanValue {
	return &BooleanValue{
		BaseNode: BaseNode{kinds.BooleanValue, loc},
		Value:    value,
	}
}

func (v *BooleanValue) GetValue() interface{} {
	return v.Value
}

// EnumValue implements Node, Value
type EnumValue struct {
	BaseNode
	Value string `json:"value"`
}

func NewEnumValue(loc *Location, value string) *EnumValue {
	return &EnumValue{
		BaseNode: BaseNode{kinds.EnumValue, loc},
		Value:    value,
	}
}

func (v *EnumValue) GetValue() interface{} {
	return v.Value
}

// ListValue implements Node, Value
type ListValue struct {
	BaseNode
	Values []Value `json:"values"`
}

func NewListValue(loc *Location, values []Value) *ListValue {
	return &ListValue{
		BaseNode: BaseNode{kinds.ListValue, loc},
		Values:   values,
	}
}

// GetValue alias to ListValue.GetValues()
func (v *ListValue) GetValue() interface{} {
	values := make([]interface{}, len(v.Values))
	for i, v := range v.Values {
		values[i] = v.GetValue()
	}
	return values
}

// ObjectValue implements Node, Value
type ObjectValue struct {
	BaseNode
	Fields []*ObjectField `json:"fields"`
}

func NewObjectValue(loc *Location, fields []*ObjectField) *ObjectValue {
	return &ObjectValue{
		BaseNode: BaseNode{kinds.ObjectValue, loc},
		Fields:   fields,
	}
}

func (v *ObjectValue) GetValue() interface{} {
	obj := make(map[string]interface{}, len(v.Fields))
	for _, field := range v.Fields {
		obj[field.Name.Value] = field.Value.GetValue()
	}
	return obj
}

// ObjectField implements Node, Value
type ObjectField struct {
	BaseNode
	Name  *Name `json:"name"`
	Value Value `json:"value"`
}

func NewObjectField(loc *Location, name *Name, value Value) *ObjectField {
	return &ObjectField{
		BaseNode: BaseNode{kinds.ObjectField, loc},
		Name:     name,
		Value:    value,
	}
}

func (f *ObjectField) GetValue() interface{} {
	return f.Value.GetValue()
}
