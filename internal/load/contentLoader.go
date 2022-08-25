package load

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/rog-golang-buddies/dto/fileresource"
	"github.com/rog-golang-buddies/internal/config"
	"github.com/rog-golang-buddies/internal/logger"
)

// ContentLoader loads content by url
//
//go:generate mockgen -source=contentLoader.go -destination=./mocks/contentLoader.go -package=load
type ContentLoader interface {
	Load(ctx context.Context, url string) (*fileresource.FileResource, error)
}

type ContentLoaderImpl struct {
	log    logger.Logger
	config *config.Web
}

// Load gets context and an url of a OpenApi file (Swagger file) string as parameter and returns a FileResource containing the link, optionally name and main content of the file.
func (cl *ContentLoaderImpl) Load(ctx context.Context, url string) (*fileresource.FileResource, error) {
	cl.log.Infof("start loading file from the url: %s", url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			cl.log.Error("error while closing request body: ", err)
		}
	}(resp.Body)
	reader := &io.LimitedReader{R: resp.Body, N: cl.config.RespLimBytes}
	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	if reader.N <= 0 {
		return nil, errors.New("file exceed the limit: " + strconv.FormatInt(cl.config.RespLimBytes, 10))
	}

	return &fileresource.FileResource{
		Link:    url,
		Content: body,
	}, nil
}

func NewContentLoader(log logger.Logger, config *config.Web) ContentLoader {
	return &ContentLoaderImpl{
		log:    log,
		config: config,
	}
}
