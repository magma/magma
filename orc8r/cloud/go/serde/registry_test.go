package serde_test

import (
	"testing"

	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/serde/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSerialize(t *testing.T) {
	serde.UnregisterAllSerdes(t)
	defer func() {
		serde.UnregisterAllSerdes(t)
	}()

	mockSerde := &mocks.Serde{}
	mockSerde.On("GetDomain").Return("foo")
	mockSerde.On("GetType").Return("bar")
	mockSerde.On("Serialize", mock.Anything).Return([]byte("hello world"), nil)

	err := serde.RegisterSerdes(mockSerde)
	assert.NoError(t, err)
	actual, err := serde.Serialize("foo", "bar", "baz")
	assert.NoError(t, err)
	assert.Equal(t, []byte("hello world"), actual)

	_, err = serde.Serialize("bar", "foo", "baz")
	assert.EqualError(t, err, "No serdes registered for domain bar")

	_, err = serde.Serialize("foo", "baz", "bar")
	assert.EqualError(t, err, "No Serde found for type baz")
}

func TestDeserialize(t *testing.T) {
	serde.UnregisterAllSerdes(t)
	defer func() {
		serde.UnregisterAllSerdes(t)
	}()

	mockSerde := &mocks.Serde{}
	mockSerde.On("GetDomain").Return("foo")
	mockSerde.On("GetType").Return("bar")
	mockSerde.On("Deserialize", mock.Anything).Return("hello world", nil)

	err := serde.RegisterSerdes(mockSerde)
	assert.NoError(t, err)
	actual, err := serde.Deserialize("foo", "bar", []byte("baz"))
	assert.NoError(t, err)
	assert.Equal(t, "hello world", actual)

	_, err = serde.Serialize("bar", "foo", "baz")
	assert.EqualError(t, err, "No serdes registered for domain bar")

	_, err = serde.Serialize("foo", "baz", "bar")
	assert.EqualError(t, err, "No Serde found for type baz")

}
