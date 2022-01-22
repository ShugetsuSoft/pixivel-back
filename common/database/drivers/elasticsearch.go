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

func (es *ElasticSearch) CreateIndex(name string, mapping string) error {
	exists, err := es.cli.IndexExists(name).Do(es.ctx)
	if err != nil {
		return err
	}
	if !exists {
		_, err := es.cli.CreateIndex(name).BodyString(mapping).Do(es.ctx)
		if err != nil {
			return err
		}
	} else {
		return models.ErrorIndexExist
	}
	return nil
}

func (es *ElasticSearch) DeleteIndex(name string) error {
	_, err := es.cli.DeleteIndex(name).Do(es.ctx)
	if err != nil {
		return err
	}
	return nil
}

func (es *ElasticSearch) InsertDocument(index string, id string, object interface{}) error {
	_, err := es.cli.Index().
		Index(index).
		Id(id).
		BodyJson(object).
		Do(es.ctx)
	if err != nil {
		return err
	}
	return nil
}

func (es *ElasticSearch) Search(index string) *elastic.SearchService  {
	return es.cli.Search().Index(index)
}

func (es *ElasticSearch) MultiSearch() *elastic.MultiSearchService  {
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

func (es *ElasticSearch) DoSearch(s *elastic.SearchService) (*elastic.SearchResult, error) {
	return s.Do(es.ctx)
}

func (es *ElasticSearch) DoMultiSearch(s *elastic.MultiSearchService) (*elastic.MultiSearchResult, error) {
	return s.Do(es.ctx)
}
