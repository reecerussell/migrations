package migrations_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/reecerussell/migrations"
	"github.com/reecerussell/migrations/mock"
)

func TestApply_GivenUnappliedMigrations_ReturnsNoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCtx := context.Background()
	testMigration := &migrations.Migration{}

	mockProvider := mock.NewMockProvider(ctrl)
	mockProvider.EXPECT().GetAppliedMigrations(testCtx).Return([]*migrations.Migration{}, nil)
	mockProvider.EXPECT().Apply(testCtx, testMigration).Return(nil)

	err := migrations.Apply(testCtx, []*migrations.Migration{testMigration}, mockProvider, "")
	assert.NoError(t, err)
}

func TestApply_GivenAppliedMigration_SkipsMigration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCtx := context.Background()
	testMigration := &migrations.Migration{}

	mockProvider := mock.NewMockProvider(ctrl)
	mockProvider.EXPECT().GetAppliedMigrations(testCtx).Return([]*migrations.Migration{testMigration}, nil)

	err := migrations.Apply(testCtx, []*migrations.Migration{testMigration}, mockProvider, "")
	assert.NoError(t, err)
}

func TestApply_FailsToGetAppliedMigrations_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCtx := context.Background()
	testError := errors.New("an error occured")

	mockProvider := mock.NewMockProvider(ctrl)
	mockProvider.EXPECT().GetAppliedMigrations(testCtx).Return(nil, testError)

	err := migrations.Apply(testCtx, nil, mockProvider, "")
	assert.Equal(t, testError, err)
}

func TestApply_FailsToApplyMigration_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCtx := context.Background()
	testMigration := &migrations.Migration{}
	testError := errors.New("an error occured")

	mockProvider := mock.NewMockProvider(ctrl)
	mockProvider.EXPECT().GetAppliedMigrations(testCtx).Return([]*migrations.Migration{}, nil)
	mockProvider.EXPECT().Apply(testCtx, testMigration).Return(testError)

	err := migrations.Apply(testCtx, []*migrations.Migration{testMigration}, mockProvider, "")
	assert.Equal(t, testError, err)
}

func TestApply_ApplyTargetMigration_SkipsFurtherMigrations(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCtx := context.Background()
	testMigrations := []*migrations.Migration{
		&migrations.Migration{Name: "One"},
		&migrations.Migration{Name: "Two"},
	}

	mockProvider := mock.NewMockProvider(ctrl)
	mockProvider.EXPECT().GetAppliedMigrations(testCtx).Return([]*migrations.Migration{}, nil)
	mockProvider.EXPECT().Apply(testCtx, testMigrations[0]).Return(nil)

	err := migrations.Apply(testCtx, testMigrations, mockProvider, "One")
	assert.NoError(t, err)
}
