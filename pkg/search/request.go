package search

import (
	"github.com/blevesearch/bleve"
)

type Page struct {
	Count  int `json:"count"`
	Offset int `json:"offset"`
}

type Request struct {
	Filter *Filter `json:"filter"`
	Sort   *Sort   `json:"sort"`
	Page   *Page   `json:"page"`
}

func (req *Request) ToBleveSearchRequest(index *Index) (*bleve.SearchRequest, error) {
	q, err := ParseFilter(req.Filter, index)
	if err != nil {
		return nil, err
	}
	if q == nil {
		q = bleve.NewMatchAllQuery()
	}
	out := bleve.NewSearchRequest(q)

	s, err := ParseSort(req.Sort, index)
	if err != nil {
		return nil, err
	}
	if s != nil {
		out.SortByCustom(s)
	}

	if req.Page != nil && req.Page.Offset != 0 {
		out.From = req.Page.Offset
	}
	if req.Page != nil && req.Page.Count != 0 {
		out.Size = req.Page.Count
	}
	return out, nil
}
