package parse

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	mock_logger "github.com/internal/logger/mocks"
	parse "github.com/internal/parse"
	"github.com/rog-golang-buddies/api_hub_common/apispecdoc"
	"github.com/rog-golang-buddies/dto/fileresource"
	"github.com/rog-golang-buddies/internal/parse/openapi"
	"github.com/stretchr/testify/assert"
)

func TestNewConverter(t *testing.T) {
	ctrl := gomock.NewController(t)
	log := mock_logger.NewMockLogger(ctrl)
	parsers := []Parser{openapi.NewOpenApi(log)}
	converter := NewConverter(log, parsers)
	assert.NotNil(t, converter)
}

func TestConverterImpl_Convert_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	log := mock_logger.NewMockLogger(ctrl)
	parser := parse.NewMockParser(ctrl)
	parsers := []Parser{parser}
	ctx := context.Background()
	testLink := "test link"
	testContent := []byte("test content")
	fType := fileresource.OpenApi
	fResource := fileresource.FileResource{
		Link:    testLink,
		Content: testContent,
		Type:    fType,
	}
	expectedApiSpec := apispecdoc.ApiSpecDoc{}
	getTypeMethod := parser.EXPECT().GetType().Times(1).Return(fileresource.OpenApi)
	parser.EXPECT().Parse(ctx, testContent).Times(1).Return(&expectedApiSpec, nil).After(getTypeMethod)
	log.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

	converter := newConverterImpl(log, parsers)
	assert.NotNil(t, converter)
	apiSpec, err := converter.Convert(ctx, &fResource)
	assert.Nil(t, err)
	assert.Equal(t, expectedApiSpec, *apiSpec)
}

func TestConverterImpl_Convert_parseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	log := mock_logger.NewMockLogger(ctrl)
	parser := parse.NewMockParser(ctrl)
	parsers := []Parser{parser}
	ctx := context.Background()
	testLink := "test link"
	testContent := []byte("test content")
	fType := fileresource.OpenApi
	fResource := fileresource.FileResource{
		Link:    testLink,
		Content: testContent,
		Type:    fType,
	}
	expErr := errors.New("test error")
	getTypeMethod := parser.EXPECT().GetType().Times(1).Return(fileresource.OpenApi)
	parser.EXPECT().Parse(ctx, testContent).Times(1).Return(nil, expErr).After(getTypeMethod)
	log.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

	converter := newConverterImpl(log, parsers)
	assert.NotNil(t, converter)
	apiSpec, err := converter.Convert(ctx, &fResource)
	assert.Nil(t, apiSpec)
	assert.NotNil(t, err)
	assert.Equal(t, expErr, err)
}

func TestConverterImpl_Convert_noTypeErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	log := mock_logger.NewMockLogger(ctrl)
	parser := parse.NewMockParser(ctrl)
	parsers := []Parser{parser}
	ctx := context.Background()
	testLink := "test link"
	testContent := []byte("test content")
	fType := fileresource.OpenApi
	fResource := fileresource.FileResource{
		Link:    testLink,
		Content: testContent,
		Type:    fType,
	}
	parser.EXPECT().GetType().Times(1).Return(fileresource.Undefined)
	log.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

	converter := newConverterImpl(log, parsers)
	assert.NotNil(t, converter)
	apiSpec, err := converter.Convert(ctx, &fResource)
	assert.Nil(t, apiSpec)
	assert.NotNil(t, err)
}
