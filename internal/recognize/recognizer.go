package recognize

import (
	"strings"

	"github.com/rog-golang-buddies/dto/fileresource"
	"github.com/rog-golang-buddies/internal/logger"
)

// Recognizer provide functionality to recognize file type by content
//
//go:generate mockgen -source=recognizer.go -destination=./mocks/recognizer.go -package=recognize
type Recognizer interface {
	//RecognizeFileType recognizes type of the file by content. Probably we may combine it with validation
	//Also not sure name is needed here. Better to recognize by content (check is file yaml, if yaml - check version;
	//if json - check openApi version) But it is easier to use file extension as starting point to check content.
	RecognizeFileType(resource *fileresource.FileResource) (fileresource.AsdFileType, error)
}

type RecognizerImpl struct {
	log logger.Logger
}

func (r *RecognizerImpl) RecognizeFileType(resource *fileresource.FileResource) (fileresource.AsdFileType, error) {
	r.log.Infof("start file '%s' recognizing", resource.Link)
	// Initially, probably simple extension recognition will be enough.
	if strings.HasSuffix(resource.Link, ".json") ||
		strings.HasSuffix(resource.Link, ".yml") ||
		strings.HasSuffix(resource.Link, ".yaml") {
		return fileresource.OpenApi, nil
	}

	return fileresource.Undefined, nil
}

func NewRecognizer(log logger.Logger) Recognizer {
	return &RecognizerImpl{
		log: log,
	}
}
