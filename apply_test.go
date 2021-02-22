package migrations_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/reecerussell/migrations"
	"github.com/reecerussell/migrations/mock"
)

func TestApply_GivenUnappliedMigrations_ReturnsNoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testContent := "My Migration Content"
	testCtx := context.Background()
	testMigration := &migrations.Migration{
		Name:   "MyMigration",
		UpFile: "MyFile",
	}

	mockProvider := mock.NewMockProvider(ctrl)
	mockProvider.EXPECT().GetAppliedMigrations(testCtx).Return([]*migrations.Migration{}, nil)
	mockProvider.EXPECT().Apply(testCtx, "MyMigration", testContent).Return(nil)

	mockFileReader := mock.NewMockFileReader(ctrl)
	mockFileReader.EXPECT().Read("MyFile").Return(testContent, nil)

	err := migrations.Apply(testCtx, []*migrations.Migration{testMigration}, mockProvider, mockFileReader, "")
	assert.NoError(t, err)
}

func TestApply_GivenAppliedMigration_SkipsMigration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCtx := context.Background()
	testMigration := &migrations.Migration{}

	mockProvider := mock.NewMockProvider(ctrl)
	mockProvider.EXPECT().GetAppliedMigrations(testCtx).Return([]*migrations.Migration{testMigration}, nil)

	err := migrations.Apply(testCtx, []*migrations.Migration{testMigration}, mockProvider, nil, "")
	assert.NoError(t, err)
}

func TestApply_FailsToGetAppliedMigrations_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCtx := context.Background()
	testError := errors.New("an error occured")

	mockProvider := mock.NewMockProvider(ctrl)
	mockProvider.EXPECT().GetAppliedMigrations(testCtx).Return(nil, testError)

	err := migrations.Apply(testCtx, nil, mockProvider, nil, "")
	assert.Equal(t, testError, err)
}

func TestApply_GivenMigrationWithMissingFile_ReturnsIsNotExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCtx := context.Background()
	testMigration := &migrations.Migration{
		Name:   "MyMigration",
		UpFile: "MyFile",
	}

	mockProvider := mock.NewMockProvider(ctrl)
	mockProvider.EXPECT().GetAppliedMigrations(testCtx).Return([]*migrations.Migration{}, nil)

	mockFileReader := mock.NewMockFileReader(ctrl)
	mockFileReader.EXPECT().Read("MyFile").Return("", os.ErrNotExist)

	err := migrations.Apply(testCtx, []*migrations.Migration{testMigration}, mockProvider, mockFileReader, "")
	assert.True(t, os.IsNotExist(err))
}

func TestApply_FailsToApplyMigration_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testContent := "My Migration Content"
	testCtx := context.Background()
	testMigration := &migrations.Migration{
		Name:   "MyMigration",
		UpFile: "MyFile",
	}
	testError := errors.New("an error occured")

	mockProvider := mock.NewMockProvider(ctrl)
	mockProvider.EXPECT().GetAppliedMigrations(testCtx).Return([]*migrations.Migration{}, nil)
	mockProvider.EXPECT().Apply(testCtx, "MyMigration", testContent).Return(testError)

	mockFileReader := mock.NewMockFileReader(ctrl)
	mockFileReader.EXPECT().Read("MyFile").Return(testContent, nil)

	err := migrations.Apply(testCtx, []*migrations.Migration{testMigration}, mockProvider, mockFileReader, "")
	assert.Equal(t, testError, err)
}

func TestApply_ApplyTargetMigration_SkipsFurtherMigrations(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testContent := "My Migration Content"
	testCtx := context.Background()
	testMigrations := []*migrations.Migration{
		&migrations.Migration{Name: "One", UpFile: "MyFile"},
		&migrations.Migration{Name: "Two"},
	}

	mockProvider := mock.NewMockProvider(ctrl)
	mockProvider.EXPECT().GetAppliedMigrations(testCtx).Return([]*migrations.Migration{}, nil)
	mockProvider.EXPECT().Apply(testCtx, "One", testContent).Return(nil)

	mockFileReader := mock.NewMockFileReader(ctrl)
	mockFileReader.EXPECT().Read("MyFile").Return(testContent, nil)

	err := migrations.Apply(testCtx, testMigrations, mockProvider, mockFileReader, "One")
	assert.NoError(t, err)
}
