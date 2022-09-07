package apispecdoc

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/dto"
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

func NewService(log logger.Logger, repo AsdRepository) Service {
	return &ServiceImpl{log: log, asdRepo: repo}
}

type ServiceImpl struct {
	log     logger.Logger
	asdRepo AsdRepository
}

func (s *ServiceImpl) Search(ctx context.Context, req *apispecproto.SearchRequest) (*apispecproto.SearchResponse, error) {
	if req == nil {
		s.log.Error("nil request body received")
		return nil, errors.New("request body must not be nil")
	}
	pageReq := dto.PageRequest{}
	if req.Page != nil {
		pageReq.Page = int(*req.Page)
	}
	if req.PerPage != nil {
		pageReq.PerPage = int(*req.PerPage)
	} else {
		pageReq.PerPage = 10
	}
	asdPage, err := s.asdRepo.SearchShort(ctx, req.Search, pageReq)
	if err != nil {
		return nil, err
	}
	res := new(apispecproto.SearchResponse)
	resDocs := make([]*apispecproto.ShortASD, 0)
	for _, asd := range asdPage.Data {
		resDocs = append(resDocs, &apispecproto.ShortASD{
			Id:          uint32(asd.ID),
			Name:        asd.Title,
			Description: asd.Description,
		})
	}
	res.ShortSpecDocs = resDocs
	res.Page = &apispecproto.Page{
		Total:   int32(asdPage.Total),
		Current: int32(asdPage.Page),
		PerPage: int32(asdPage.PerPage),
	}
	return res, nil
}

func (s *ServiceImpl) Get(ctx context.Context, req *apispecproto.GetRequest) (*apispecproto.GetResponse, error) {
	return nil, errors.New("not implemented")
}

func (s *ServiceImpl) Save(ctx context.Context, asd *apispecdoc.ApiSpecDoc) (uint, error) {
	if asd == nil {
		return 0, errors.New("nil asd model received")
	}
	asdEntity, err := asdToEntity(asd)
	if err != nil {
		return 0, err
	}
	//Check records by md5 hash sum - if exists than all methods the same and update not required
	asdByHash, err := s.asdRepo.FindByHash(ctx, asd.Md5Sum)
	if err != nil {
		return 0, err
	}
	if asdByHash != nil {
		s.log.Infof("record '%s' hash '%s' no changes", asd.Title, asd.Md5Sum)
		return asdByHash.ID, nil
	}
	//Check records by file url - if exists than need to update ASD in db (prev step didn't find matched hash - so hash changed)
	asdByUrl, err := s.asdRepo.FindByUrl(ctx, asd.Url)
	if err != nil {
		return 0, err
	}
	if asdByUrl != nil {
		asdByUrl.Title = asdEntity.Title
		asdByUrl.Description = asdEntity.Description
		asdByUrl.Type = asdEntity.Type
		asdByUrl.Groups = asdEntity.Groups
		asdByUrl.ApiMethods = asdEntity.ApiMethods
		asdByUrl.Md5sum = asdEntity.Md5sum
		asdByUrl.Url = asdEntity.Url
		//clear and reattach all dependencies
		err = s.asdRepo.Update(ctx, asdByUrl)
		if err != nil {
			return 0, err
		}
		s.log.Infof("record '%s' with hash '%s' updated", asd.Title, asd.Md5Sum)
		return asdByUrl.ID, nil
	}
	s.log.Infof("create new record for '%s' hash '%s'", asd.Title, asd.Md5Sum)
	return s.asdRepo.Save(ctx, asdEntity)
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
		Url:         dto.Url,
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
