package elasticsearch

import (
	"context"
	"reflect"

	"github.com/olivere/elastic"

	"bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/dto"
)

type SearchVisitorArguments struct {
	Name         string
	RealName     string
	Telephone    string
	Email        string
	LandingTitle string
	LandingURL   string
	Browser      string
}

// CreateOrUpdateVisitorDoc ...
func CreateOrUpdateVisitorDoc(esCli *esClient, doc *dto.SearchVisitor) error {
	return esCli.CreateDoc("visitor", "_doc", doc.ID, doc)
}

// SearchVisitors ...
func SearchVisitors(esCli *esClient, query elastic.Query, from, size int) (res []*dto.SearchVisitor, totalCount int64, err error) {
	searchResult, err := esCli.client.Search().
		Index("visitor").
		Type("_doc").
		Query(query).
		Sort("created_on", false).
		From(from).Size(size).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		var reason, rootReason string
		reason = searchResult.Error.Reason
		if len(searchResult.Error.RootCause) > 0 {
			rootReason = searchResult.Error.RootCause[0].Reason
		}

		log.Logger.Errorf("SearchVisitors reason: %s, root reason: %s, error: %v\n", reason, rootReason, err)
		return nil, 0, err
	}

	for _, hit := range searchResult.Each(reflect.TypeOf(&dto.SearchVisitor{})) {
		v, ok := hit.(*dto.SearchVisitor)
		if ok {
			res = append(res, v)
		}
	}

	return res, searchResult.Hits.TotalHits, err
}
