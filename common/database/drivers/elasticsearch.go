package drivers

import (
	"context"

	"github.com/ShugetsuSoft/pixivel-back/common/models"
	elastic "github.com/olivere/elastic/v7"
)

type ElasticSearch struct {
	cli *elastic.Client
	ctx context.Context
}

func NewElasticSearchClient(ctx context.Context, uri string, username string, password string) (*ElasticSearch, error) {
	client, err := elastic.NewClient(elastic.SetURL(uri), elastic.SetBasicAuth(username, password), elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}
	return &ElasticSearch{
		cli: client,
		ctx: ctx,
	}, nil
}

func (es *ElasticSearch) CreateIndex(ctx context.Context, name string, mapping string) error {
	exists, err := es.cli.IndexExists(name).Do(ctx)
	if err != nil {
		return err
	}
	if !exists {
		_, err := es.cli.CreateIndex(name).BodyString(mapping).Do(ctx)
		if err != nil {
			return err
		}
	} else {
		return models.ErrorIndexExist
	}
	return nil
}

func (es *ElasticSearch) DeleteIndex(ctx context.Context, name string) error {
	_, err := es.cli.DeleteIndex(name).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (es *ElasticSearch) InsertDocument(ctx context.Context, index string, id string, object interface{}) error {
	_, err := es.cli.Index().
		Index(index).
		Id(id).
		BodyJson(object).
		Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (es *ElasticSearch) Search(index string) *elastic.SearchService {
	return es.cli.Search().Index(index)
}

func (es *ElasticSearch) MultiSearch() *elastic.MultiSearchService {
	return es.cli.MultiSearch()
}

func (es *ElasticSearch) Query(key string, val interface{}) *elastic.FuzzyQuery {
	return elastic.NewFuzzyQuery(key, val)
}

func (es *ElasticSearch) TermsQuery(key string, val []string) *elastic.TermsQuery {
	return elastic.NewTermsQueryFromStrings(key, val...)
}

func (es *ElasticSearch) Suggest(key string) *elastic.CompletionSuggester {
	return elastic.NewCompletionSuggester(key)
}

func (es *ElasticSearch) TermSuggest(key string) *elastic.TermSuggester {
	return elastic.NewTermSuggester(key)
}

func (es *ElasticSearch) BoolQuery() *elastic.BoolQuery {
	return elastic.NewBoolQuery()
}

func (es *ElasticSearch) DoSearch(ctx context.Context, s *elastic.SearchService) (*elastic.SearchResult, error) {
	return s.Do(ctx)
}

func (es *ElasticSearch) DoMultiSearch(ctx context.Context, s *elastic.MultiSearchService) (*elastic.MultiSearchResult, error) {
	return s.Do(ctx)
}
