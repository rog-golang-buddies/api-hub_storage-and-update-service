package process

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/rog-golang-buddies/api_hub_common/apispecdoc"
	"github.com/rog-golang-buddies/dto/fileresource"

	"github.com/golang/mock/gomock"
	load "github.com/rog-golang-buddies/internal/load/mocks"
	parse "github.com/rog-golang-buddies/internal/parse/mocks"
	recognize "github.com/rog-golang-buddies/internal/recognize/mocks"
	"github.com/stretchr/testify/assert"
)

func TestProcess_RecognizeFail_processReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)

	contentLoader := load.NewMockContentLoader(ctrl)
	recognizer := recognize.NewMockRecognizer(ctrl)
	converter := parse.NewMockConverter(ctrl)

	ctx := context.Background()
	url := "test_url"
	expectedErr := errors.New("recognize error")
	fileResource := new(fileresource.FileResource)

	loadCall := contentLoader.EXPECT().Load(ctx, url).Times(1).Return(fileResource, nil)
	recognizer.EXPECT().RecognizeFileType(fileResource).After(loadCall).Times(1).Return(fileresource.Undefined, expectedErr)

	processor, err := NewProcessor(recognizer, converter, contentLoader)
	assert.Nil(t, err)
	assert.NotNil(t, processor, "Processor must not be nil")

	asd, err := processor.Process(ctx, url)
	assert.Nil(t, asd)
	assert.Equal(t, expectedErr, err, "Should return error from recognizer")
}

func TestProcess_ContentLoaderFail_processReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)

	contentLoader := load.NewMockContentLoader(ctrl)
	recognizer := recognize.NewMockRecognizer(ctrl)
	converter := parse.NewMockConverter(ctrl)

	ctx := context.Background()
	url := "test_url"
	expectedErr := errors.New("contentload error")

	contentLoader.EXPECT().Load(ctx, url).Times(1).Return(nil, expectedErr)

	processor, err := NewProcessor(recognizer, converter, contentLoader)
	assert.Nil(t, err)

	assert.NotNil(t, processor, "Processor must not be nil")

	asd, err := processor.Process(ctx, url)
	assert.Nil(t, asd)
	assert.Equal(t, expectedErr, err, "Should return error from contentLoader")
}

func TestProcess_ConverterFail_processReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)

	contentLoader := load.NewMockContentLoader(ctrl)
	recognizer := recognize.NewMockRecognizer(ctrl)
	converter := parse.NewMockConverter(ctrl)

	ctx := context.Background()
	url := "test_url_from_yaml_openapi_file"
	expectedErr := errors.New("convert error")
	fileResource := new(fileresource.FileResource)

	loadCall := contentLoader.EXPECT().Load(ctx, url).Times(1).Return(fileResource, nil)
	recognizeCall := recognizer.EXPECT().RecognizeFileType(fileResource).After(loadCall).Times(1).Return(fileresource.OpenApi, nil)
	converter.EXPECT().Convert(ctx, gomock.Any()).Times(1).After(recognizeCall).Return(nil, expectedErr)

	processor, err := NewProcessor(recognizer, converter, contentLoader)
	assert.Nil(t, err)
	assert.NotNil(t, processor, "Processor must not be nil")

	asd, err := processor.Process(ctx, url)
	assert.Nil(t, asd)
	assert.Equal(t, expectedErr, err, "Should return error from converter")
}

func TestProcess_completed_hashPopulated(t *testing.T) {
	ctrl := gomock.NewController(t)

	contentLoader := load.NewMockContentLoader(ctrl)
	recognizer := recognize.NewMockRecognizer(ctrl)
	converter := parse.NewMockConverter(ctrl)

	ctx := context.Background()
	url := "test_url_from_yaml_openapi_file"

	fileContent := []byte("Test content")
	fileResource := new(fileresource.FileResource)
	fileResource.Content = fileContent

	result := &apispecdoc.ApiSpecDoc{}
	loadCall := contentLoader.EXPECT().Load(ctx, url).Times(1).Return(fileResource, nil)
	recognizeCall := recognizer.EXPECT().RecognizeFileType(fileResource).After(loadCall).Times(1).Return(fileresource.OpenApi, nil)
	converter.EXPECT().Convert(ctx, gomock.Any()).Times(1).After(recognizeCall).Return(result, nil)

	processor, err := NewProcessor(recognizer, converter, contentLoader)
	assert.Nil(t, err)
	assert.NotNil(t, processor, "Processor must not be nil")

	asd, err := processor.Process(ctx, url)
	assert.Nil(t, err)
	assert.NotNil(t, asd)
	assert.Equal(t, result, asd)
	hashSum := md5.Sum(fileContent)
	expectedHash := hex.EncodeToString(hashSum[:])
	assert.Equal(t, expectedHash, asd.Md5Sum)
}
