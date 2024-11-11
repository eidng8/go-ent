package paginate

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"net/url"

	eu "github.com/eidng8/go-url"
	"github.com/gin-gonic/gin"
	"github.com/ogen-go/ogen"
)

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
func GetPage[V any, Q any, T PQ[V, Q]](
	gc *gin.Context, qc context.Context, query T, params PaginatedParams,
) (*PaginatedList[V], error) {
	var next, prev string
	fi := 1
	ni := params.Page + 1
	pi := params.Page - 1
	req := gc.Request
	count, err := query.Count(qc)
	if err != nil {
		return nil, err
	}
	if 0 == count {
		return &PaginatedList[V]{
			Total:        0,
			PerPage:      params.PerPage,
			CurrentPage:  1,
			LastPage:     1,
			FirstPageUrl: UrlWithPage(req, 1, params.PerPage).String(),
			LastPageUrl:  "",
			NextPageUrl:  "",
			PrevPageUrl:  "",
			Path:         eu.RequestBaseUrl(req).String(),
			From:         0,
			To:           0,
			Data:         []*V{},
		}, nil
	}
	from := pi*params.PerPage + 1
	to := int(math.Min(float64(params.Page*params.PerPage), float64(count)))
	query.Offset(pi * params.PerPage)
	query.Limit(params.PerPage)
	rows, err := query.All(qc)
	if err != nil {
		return nil, err
	}
	li := int(math.Ceil(float64(count) / float64(params.PerPage)))
	first := UrlWithPage(req, fi, params.PerPage).String()
	var last string
	if li <= 1 {
		li = 1
		last = ""
	} else {
		last = UrlWithPage(req, li, params.PerPage).String()
	}
	if ni > li {
		ni = li
		next = ""
	} else {
		next = UrlWithPage(req, ni, params.PerPage).String()
	}
	if pi < 1 {
		pi = 1
		prev = ""
	} else {
		prev = UrlWithPage(req, pi, params.PerPage).String()
	}
	return &PaginatedList[V]{
		Total:        count,
		PerPage:      params.PerPage,
		CurrentPage:  params.Page,
		LastPage:     li,
		FirstPageUrl: first,
		LastPageUrl:  last,
		NextPageUrl:  next,
		PrevPageUrl:  prev,
		Path:         eu.RequestBaseUrl(req).String(),
		From:         from,
		To:           to,
		Data:         rows,
	}, nil
}

// GetPageMapped returns a paginated list of items. `I` is the type of items
// returned by the query, `V` is the type of items in the paginated list. `Q` is
// the query type to be used to retrieve items, which in most cases can be
// inferred. So in most cases, only the `I` and `V` types need to be provided.
//
// The `gc` parameter is the gin.Context to be used to generate various links in
// the paginated list; `qc` is the context to be used in query execution;
// `query` is the ent query instance to be executed; and `params` is the
// PaginatedParams to be used in pagination; the `mapper` is a function that
// maps the one query result row to an item in the paginated list, the 2nd
// parameter is the index of the item in the result set.
//
// Please remember to explicitly add the `ORDER` clause to the query before
// calling this function.
func GetPageMapped[I any, V any, Q any, T PQ[I, Q]](
	gc *gin.Context, qc context.Context, query T, page PaginatedParams,
	mapper func(*I, int) *V,
) (*PaginatedList[V], error) {
	list, err := GetPage[I, Q, T](gc, qc, query, page)
	if err != nil {
		return nil, err
	}
	data := make([]*V, len(list.Data))
	for i, row := range list.Data {
		data[i] = mapper(row, i)
	}
	return &PaginatedList[V]{
		Total:        list.Total,
		PerPage:      list.PerPage,
		CurrentPage:  list.CurrentPage,
		LastPage:     list.LastPage,
		FirstPageUrl: list.FirstPageUrl,
		LastPageUrl:  list.LastPageUrl,
		NextPageUrl:  list.NextPageUrl,
		PrevPageUrl:  list.PrevPageUrl,
		Path:         list.Path,
		From:         list.From,
		To:           list.To,
		Data:         data,
	}, nil
}

// UrlWithPage returns a URL with the page and per_page query parameters set.
func UrlWithPage(request *http.Request, page int, perPage int) *url.URL {
	return eu.RequestUrlWithQueryParams(request, PageQueryParams(page, perPage))
}

// UrlWithoutPageParams returns a URL without the page and per_page query
// parameters.
func UrlWithoutPageParams(req *http.Request) *url.URL {
	return eu.RequestUrlWithoutQueryParams(req, ParamPage, ParamPerPage)
}

// PageQueryParams sets the page and per_page query parameters.
func PageQueryParams(page int, perPage int) map[string]string {
	params := make(map[string]string, 2)
	params[ParamPage] = fmt.Sprintf("%d", page)
	params[ParamPerPage] = fmt.Sprintf("%d", perPage)
	return params
}

// FixParamNames fixes the parameter names to be `per_page` and `page`, to be
// used with oapi-codegen.
func FixParamNames(params []*ogen.Parameter) {
	FixParamNamesWith(params, "itemsPerPage", "page")
}

// FixParamNamesWith fixes the parameter names to be `per_page` and `page`, to
// be used with oapi-codegen. `pageParam` and `perPageParam` are names generated
// by oapi-codegen to be replaced.
func FixParamNamesWith(
	params []*ogen.Parameter, pageParam string, perPageParam string,
) {
	for _, param := range params {
		if perPageParam == param.Name {
			param.Name = "per_page"
		} else if pageParam == param.Name {
			param.Name = "page"
		}
	}
}

// SetResponse changes the response of given OpenAPI operation to meet the
// paginate response.
func SetResponse(op *ogen.Operation, description string, itemRef string) {
	op.Responses["200"] = &ogen.Response{
		Description: description,
		Content: map[string]ogen.Media{
			"application/json": {
				Schema: &ogen.Schema{
					Type: "object",
					Properties: []ogen.Property{
						{
							Name: "current_page",
							Schema: &ogen.Schema{
								Type:        "integer",
								Description: "Page number (1-based)",
								Minimum:     ogen.Num("1"),
							},
						},
						{
							Name: "total",
							Schema: &ogen.Schema{
								Type:        "integer",
								Description: "Total number of items",
								Minimum:     ogen.Num("0"),
							},
						},
						{
							Name: "per_page",
							Schema: &ogen.Schema{
								Type:        "integer",
								Description: "Number of items per page",
								Minimum:     ogen.Num("1"),
							},
						},
						{
							Name: "last_page",
							Schema: &ogen.Schema{
								Type:        "integer",
								Description: "Last page number",
								Minimum:     ogen.Num("1"),
							},
						},
						{
							Name: "from",
							Schema: &ogen.Schema{
								Type:        "integer",
								Description: "Index (1-based) of the first item in the current page",
								Minimum:     ogen.Num("0"),
							},
						},
						{
							Name: "to",
							Schema: &ogen.Schema{
								Type:        "integer",
								Description: "Index (1-based) of the last item in the current page",
								Minimum:     ogen.Num("0"),
							},
						},
						{
							Name: "first_page_url",
							Schema: &ogen.Schema{
								Type:        "string",
								Description: "URL to the first page",
							},
						},
						{
							Name: "last_page_url",
							Schema: &ogen.Schema{
								Type:        "string",
								Description: "URL to the last page",
							},
						},
						{
							Name: "next_page_url",
							Schema: &ogen.Schema{
								Type:        "string",
								Description: "URL to the next page",
							},
						},
						{
							Name: "prev_page_url",
							Schema: &ogen.Schema{
								Type:        "string",
								Description: "URL to the previous page",
							},
						},
						{
							Name: "path",
							Schema: &ogen.Schema{
								Type:        "string",
								Description: "Base path of the request",
							},
						},
						{
							Name: "data",
							Schema: &ogen.Schema{
								Type:        "array",
								Description: "List of items",
								Items: &ogen.Items{
									Item: &ogen.Schema{Ref: itemRef},
								},
							},
						},
					},
					Required: []string{
						"current_page",
						"total",
						"per_page",
						"last_page",
						"from",
						"to",
						"first_page_url",
						"last_page_url",
						"next_page_url",
						"prev_page_url",
						"path",
						"data",
					},
				},
			},
		},
	}
}
