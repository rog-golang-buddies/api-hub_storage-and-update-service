package apispecdoc

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/logger"
	"github.com/rog-golang-buddies/api_hub_common/apispecdoc"
	"github.com/rog-golang-buddies/api_hub_common/apispecproto"
	"time"
)

//go:generate mockgen -source=service.go -destination=./mocks/service.go
type Service interface {
	Search(context.Context, *apispecproto.SearchRequest) (*apispecproto.SearchResponse, error)
	Get(context.Context, *apispecproto.GetRequest) (*apispecproto.GetResponse, error)
	Save(context.Context, *apispecdoc.ApiSpecDoc) (uint, error)
}

func NewService(log logger.Logger, repo Repository) Service {
	return &ServiceImpl{log: log, repo: repo}
}

type ServiceImpl struct {
	log  logger.Logger
	repo Repository
}

func (s *ServiceImpl) Search(ctx context.Context, req *apispecproto.SearchRequest) (*apispecproto.SearchResponse, error) {
	return nil, errors.New("not implemented")
}

func (s *ServiceImpl) Get(ctx context.Context, req *apispecproto.GetRequest) (*apispecproto.GetResponse, error) {
	return nil, errors.New("not implemented")
}

func (s *ServiceImpl) Save(ctx context.Context, asd *apispecdoc.ApiSpecDoc) (uint, error) {
	//TODO validate md5 sum
	if asd == nil {
		return 0, errors.New("nil asd model received")
	}
	asdEntity, err := asdToEntity(asd)
	if err != nil {
		return 0, err
	}
	return s.repo.Save(ctx, asdEntity)
}

func asdToEntity(dto *apispecdoc.ApiSpecDoc) (*ApiSpecDoc, error) {
	groups := make([]*Group, 0, len(dto.Groups))
	for _, group := range dto.Groups {
		methods, err := methodsToEntities(group.Methods)
		if err != nil {
			return nil, err
		}
		groups = append(groups, &Group{
			Name:        group.Name,
			Description: group.Description,
			ApiMethods:  methods,
		})
	}
	methods, err := methodsToEntities(dto.Methods)
	if err != nil {
		return nil, err
	}
	return &ApiSpecDoc{
		Title:       dto.Title,
		Description: dto.Description,
		Type:        string(dto.Type),
		Groups:      groups,
		ApiMethods:  methods,
		Md5sum:      dto.Md5Sum,
		FetchedAt:   time.Now(),
	}, nil
}

func methodToEntity(method *apispecdoc.ApiMethod) (*ApiMethod, error) {
	if method == nil {
		return &ApiMethod{}, nil
	}
	var params []byte
	var err error
	if method.Parameters != nil {
		params, err = json.Marshal(method.Parameters)
		if err != nil {
			return nil, err
		}
	}
	var body []byte
	if method.RequestBody != nil {
		body, err = json.Marshal(method.RequestBody)
		if err != nil {
			return nil, err
		}
	}
	var servers []*Server
	if method.Servers != nil {
		servers = make([]*Server, 0, len(method.Servers))
		for _, server := range method.Servers {
			servers = append(servers, &Server{
				URL:         server.Url,
				Description: server.Description,
			})
		}
	}
	var extDoc ExternalDoc
	if method.ExternalDoc != nil {
		extDoc = ExternalDoc{
			Description: method.ExternalDoc.Description,
			URL:         method.ExternalDoc.Url,
		}
	}
	return &ApiMethod{
		Path:        method.Path,
		Name:        method.Name,
		Description: method.Description,
		Type:        string(method.Type),
		Parameters:  string(params),
		Servers:     servers,
		RequestBody: string(body),
		ExternalDoc: &extDoc,
	}, nil
}

func methodsToEntities(methods []*apispecdoc.ApiMethod) ([]*ApiMethod, error) {
	if methods == nil {
		return make([]*ApiMethod, 0), nil
	}
	resMeth := make([]*ApiMethod, 0, len(methods))
	for _, method := range methods {
		methEntity, err := methodToEntity(method)
		if err != nil {
			return nil, err
		}
		resMeth = append(resMeth, methEntity)
	}
	return resMeth, nil
}
