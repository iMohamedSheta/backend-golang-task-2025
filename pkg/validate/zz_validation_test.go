package validate

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name string `validate:"required"`
	Age  int    `validate:"gte=0,lte=130"`
}

func (t TestStruct) Messages() map[string]string {
	return map[string]string{"Name.required": "Name is required"}
}
func (t TestStruct) GetRequestSentFields() map[string]any { return nil }
func (t TestStruct) SetRequestSentFields(map[string]any)  {}

func TestValidator_ValidateRequest(t *testing.T) {
	v := New(validator.New())
	ts := TestStruct{Age: 25}
	err := v.ValidateRequest(ts)

	assert.Error(t, err)
}

func TestValidator_ValidStruct(t *testing.T) {
	v := New(validator.New())
	ts := TestStruct{Name: "John", Age: 30}
	err := v.ValidateRequest(ts)

	assert.NoError(t, err)
}
