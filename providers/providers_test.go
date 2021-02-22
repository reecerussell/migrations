package providers

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/reecerussell/migrations"
	"github.com/reecerussell/migrations/mock"
)

func TestAdd_GivenNewProvider_AddsSuccessfully(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProvider := mock.NewMockProvider(ctrl)

	Add("TestAdd_GivenNewProvider_AddsSuccessfully", func(conf migrations.ConfigMap) migrations.Provider {
		return mockProvider
	})

	p, ok := prvs["TestAdd_GivenNewProvider_AddsSuccessfully"]
	assert.Equal(t, mockProvider, p(nil))
	assert.True(t, ok)
}

func TestGet_GivenValidName_ReturnsProvider(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProvider := mock.NewMockProvider(ctrl)
	prvs["TestGet_GivenValidName_ReturnsProvider"] = func(conf migrations.ConfigMap) migrations.Provider {
		assert.Equal(t, "test", conf["value"])

		return mockProvider
	}

	testConfig := migrations.ConfigMap{
		"value": "test",
	}

	p := Get("TestGet_GivenValidName_ReturnsProvider", testConfig)
	assert.Equal(t, mockProvider, p)
}

func TestGet_GivenInvalidName_Panics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	defer func() {
		r := recover()
		assert.Equal(t, "no provider 'TestGet_GivenInvalidName_Panics' is registered", r)
	}()

	_ = Get("TestGet_GivenInvalidName_Panics", nil)
}
