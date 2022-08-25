package openapi

import (
	"context"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rog-golang-buddies/dto/fileresource"
	mock_logger "github.com/rog-golang-buddies/internal/logger/mocks"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()
	log := mock_logger.NewMockLogger(ctrl)
	content, err := os.ReadFile("./mocks/github_mock.yml")
	assert.Nil(t, err)
	assert.NotNil(t, content)
	oParser := NewOpenApi(log)
	openAPI, err := oParser.Parse(ctx, content)
	assert.NotNil(t, openAPI)
	assert.Nil(t, err)
}

func TestNewOpenApi(t *testing.T) {
	ctrl := gomock.NewController(t)
	log := mock_logger.NewMockLogger(ctrl)
	oParser := NewOpenApi(log)
	assert.NotNil(t, oParser)
	assert.Equal(t, fileresource.OpenApi, oParser.GetType())
}
