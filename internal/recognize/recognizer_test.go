package recognize

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rog-golang-buddies/dto/fileresource"
	mock_logger "github.com/rog-golang-buddies/internal/logger/mocks"
	"github.com/stretchr/testify/assert"
)

func TestRecognizeJson_OpenApiType(t *testing.T) {
	ctrl := gomock.NewController(t)
	log := mock_logger.NewMockLogger(ctrl)
	log.EXPECT().Infof(gomock.Any(), gomock.Any())
	rec := NewRecognizer(log)
	resource := fileresource.FileResource{Link: "https://github.com/github/rest-api-description/blob/main/descriptions/ghes-3.6/ghes-3.6.json"}
	fileType, err := rec.RecognizeFileType(&resource)
	assert.Nil(t, err)
	assert.Equal(t, fileresource.OpenApi, fileType)
}

func TestRecognizeYaml_OpenApiType(t *testing.T) {
	ctrl := gomock.NewController(t)
	log := mock_logger.NewMockLogger(ctrl)
	log.EXPECT().Infof(gomock.Any(), gomock.Any())
	rec := NewRecognizer(log)
	resource := fileresource.FileResource{Link: "https://github.com/github/rest-api-description/blob/main/descriptions/ghes-3.6/ghes-3.6.yaml"}
	fileType, err := rec.RecognizeFileType(&resource)
	assert.Nil(t, err)
	assert.Equal(t, fileresource.OpenApi, fileType)
}

func TestRecognizeWrongExtension_UnknownType(t *testing.T) {
	ctrl := gomock.NewController(t)
	log := mock_logger.NewMockLogger(ctrl)
	log.EXPECT().Infof(gomock.Any(), gomock.Any())
	rec := NewRecognizer(log)
	resource := fileresource.FileResource{Link: "https://github.com/github/rest-api-description/blob/main/descriptions/ghes-3.6/ghes-3.6.txt"}
	fileType, err := rec.RecognizeFileType(&resource)
	assert.Nil(t, err)
	assert.Equal(t, fileresource.Undefined, fileType)
}
