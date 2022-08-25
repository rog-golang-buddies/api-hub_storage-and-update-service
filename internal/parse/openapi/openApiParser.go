package openapi

import (
	"context"

	"github.com/rog-golang-buddies/api_hub_common/apispecdoc"
	"github.com/rog-golang-buddies/dto/fileresource"
	"github.com/rog-golang-buddies/internal/logger"
)

type Parser struct {
	log logger.Logger
}

func (p *Parser) Parse(ctx context.Context, content []byte) (*apispecdoc.ApiSpecDoc, error) {
	openapi, err := parseOpenAPI(ctx, content)
	if err != nil {
		return nil, err
	}
	return openapiToApiSpec(p.log, openapi), nil
}

func (p *Parser) GetType() fileresource.AsdFileType {
	return fileresource.OpenApi
}

func NewOpenApi(log logger.Logger) *Parser {
	return &Parser{log: log}
}
