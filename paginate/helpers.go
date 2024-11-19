package paginate

import (
	"context"
	"fmt"
	"math"
	"net/url"

	eu "github.com/eidng8/go-url"
	"github.com/gin-gonic/gin"
)

type Paginator[V any, Q any, T PQ[V, Q]] struct {
	params   PaginatedParams
	BaseUrl  string
	Query    T
	GinCtx   *gin.Context
	QueryCtx context.Context
}

// GetPaginationParams returns the PaginatedParams from the gin.Context,
// with default values of page `1` and `10` items per page.
func GetPaginationParams(gc *gin.Context) PaginatedParams {
	return GetPaginationParamsWithDefault(gc, 1, 10)
}

// GetPaginationParamsWithDefault returns the PaginatedParams from the
// gin.Context with given default values.
func GetPaginationParamsWithDefault(
	gc *gin.Context, defaultPage, defaultPerPage int,
) PaginatedParams {
	var page PaginatedParams
	if gc.ShouldBind(&page) != nil {
		page.Page = defaultPage
		page.PerPage = defaultPerPage
	}
	if page.Page < 1 {
		page.Page = defaultPage
	}
	if page.PerPage < 1 {
		page.PerPage = defaultPerPage
	}
	return page
}

// GetPage returns a paginated list of items. `V` is the type of items in the
// paginated list. `Q` is the query type to be used to retrieve items, which in
// most cases can be inferred. So in most cases, only the `V` needs to be
// provided.
//
// The `gc` parameter is the gin.Context to be used to generate various links in
// the paginated list; `qc` is the context to be used in query execution;
// `query` is the ent query instance to be executed; and `params` is the
// PaginatedParams to be used in pagination.
//
// Please remember to explicitly add the `ORDER` clause to the query before
// calling this function.
func (p Paginator[V, Q, T]) GetPage() (*PaginatedList[V], error) {
	var nextUrl, prevUrl string
	p.params = GetPaginationParams(p.GinCtx)
	firstIdx := 1
	perPage := p.params.GetPerPage()
	pageIdx := p.params.GetPage()
	nextIdx := pageIdx + 1
	prevIdx := pageIdx - 1
	count, err := p.Query.Count(p.QueryCtx)
	if err != nil {
		return nil, err
	}
	if 0 == count {
		return &PaginatedList[V]{
			Total:        0,
			PerPage:      perPage,
			CurrentPage:  1,
			LastPage:     1,
			FirstPageUrl: p.UrlWithPage(1, perPage).String(),
			LastPageUrl:  "",
			NextPageUrl:  "",
			PrevPageUrl:  "",
			Path:         p.UrlWithoutQuery().String(),
			From:         0,
			To:           0,
			Data:         []*V{},
		}, nil
	}
	from := prevIdx*perPage + 1
	to := int(math.Min(float64(pageIdx*perPage), float64(count)))
	p.Query.Offset(prevIdx * perPage)
	p.Query.Limit(perPage)
	rows, err := p.Query.All(p.QueryCtx)
	if err != nil {
		return nil, err
	}
	lastIdx := int(math.Ceil(float64(count) / float64(perPage)))
	firstUrl := p.UrlWithPage(firstIdx, perPage).String()
	var lastUrl string
	if lastIdx <= 1 {
		lastIdx = 1
		lastUrl = ""
	} else {
		lastUrl = p.UrlWithPage(lastIdx, perPage).String()
	}
	if nextIdx > lastIdx {
		nextIdx = lastIdx
		nextUrl = ""
	} else {
		nextUrl = p.UrlWithPage(nextIdx, perPage).String()
	}
	if prevIdx < 1 {
		prevIdx = 1
		prevUrl = ""
	} else {
		prevUrl = p.UrlWithPage(prevIdx, perPage).String()
	}
	return &PaginatedList[V]{
		Total:        count,
		PerPage:      perPage,
		CurrentPage:  pageIdx,
		LastPage:     lastIdx,
		FirstPageUrl: firstUrl,
		LastPageUrl:  lastUrl,
		NextPageUrl:  nextUrl,
		PrevPageUrl:  prevUrl,
		Path:         p.UrlWithoutQuery().String(),
		From:         from,
		To:           to,
		Data:         rows,
	}, nil
}

func (p Paginator[V, Q, T]) setSchemeHost(url *url.URL) *url.URL {
	u, err := url.Parse(p.BaseUrl)
	if err != nil {
		return url
	}
	url.Host = u.Host
	url.Scheme = u.Scheme
	return url
}

// UrlWithPage returns a URL with the page and per_page query parameters set.
func (p Paginator[V, Q, T]) UrlWithPage(page int, perPage int) *url.URL {
	u := eu.RequestUrlWithQueryParams(
		p.GinCtx.Request, PageQueryParams(page, perPage),
	)
	return p.setSchemeHost(u)
}

// UrlWithoutPageParams returns a URL without the page and per_page query
// parameters.
func (p Paginator[V, Q, T]) UrlWithoutPageParams() *url.URL {
	u := eu.RequestUrlWithoutQueryParams(
		p.GinCtx.Request, ParamPage, ParamPerPage,
	)
	return p.setSchemeHost(u)
}

// UrlWithoutQuery returns a URL without query string.
func (p Paginator[V, Q, T]) UrlWithoutQuery() *url.URL {
	return p.setSchemeHost(eu.RequestBaseUrl(p.GinCtx.Request))
}

// PageQueryParams sets the page and per_page query parameters.
func PageQueryParams(page int, perPage int) map[string]string {
	params := make(map[string]string, 2)
	params[ParamPage] = fmt.Sprintf("%d", page)
	params[ParamPerPage] = fmt.Sprintf("%d", perPage)
	return params
}
