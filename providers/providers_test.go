package providers

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/reecerussell/migrations/mock"
)

func TestAdd_GivenNewProvider_AddsSuccessfully(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProvider := mock.NewMockProvider(ctrl)

	Add("TestAdd_GivenNewProvider_AddsSuccessfully", mockProvider)

	p, ok := prvs["TestAdd_GivenNewProvider_AddsSuccessfully"]
	assert.Equal(t, mockProvider, p)
	assert.True(t, ok)
}

func TestGet_GivenValidName_ReturnsProvider(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProvider := mock.NewMockProvider(ctrl)
	prvs["TestGet_GivenValidName_ReturnsProvider"] = mockProvider

	p := Get("TestGet_GivenValidName_ReturnsProvider")
	assert.Equal(t, mockProvider, p)
}

func TestGet_GivenInvalidName_Panics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	defer func() {
		r := recover()
		assert.Equal(t, "no provider 'TestGet_GivenInvalidName_Panics' is registered", r)
	}()

	_ = Get("TestGet_GivenInvalidName_Panics")
}
