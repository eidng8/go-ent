package paginate

import (
	"context"

	"github.com/ogen-go/ogen"
)

const (
	// ParamPage is the query parameter name for the page number.
	ParamPage = "page"

	// ParamPerPage is the query parameter name for the number of items per page.
	ParamPerPage = "per_page"
)

// PaginatedParams is a struct that contains the page and per_page parameters.
type PaginatedParams struct {
	// Page is the current page number.
	Page int `form:"page"`
	// PerPage is the number of items per page.
	PerPage int `form:"per_page"`
}

// GetPage returns the current page number.
func (pp *PaginatedParams) GetPage() int {
	if pp.Page < 1 {
		return 1
	}
	return pp.Page
}

// GetPerPage returns the number of items per page.
func (pp *PaginatedParams) GetPerPage() int {
	if pp.PerPage < 1 {
		return 10
	}
	return pp.PerPage
}

// PaginatedList is a struct that contains the paginated list of items.
type PaginatedList[T any] struct {
	// Total is the total number of items.
	Total int `json:"total" bson:"total" xml:"total" yaml:"total"`
	// PerPage is the number of items per page.
	PerPage int `json:"per_page" bson:"per_page" xml:"per_page" yaml:"per_page"`
	// CurrentPage is the current page number.
	CurrentPage int `json:"current_page" bson:"current_page" xml:"current_page" yaml:"current_page"`
	// LastPage is the last page number.
	LastPage int `json:"last_page" bson:"last_page" xml:"last_page" yaml:"last_page"`
	// FirstPageUrl is the URL of the first page.
	FirstPageUrl string `json:"first_page_url" bson:"first_page_url" xml:"first_page_url" yaml:"first_page_url"`
	// LastPageUrl is the URL of the last page. It is an empty string if there
	// is only one page.
	LastPageUrl string `json:"last_page_url" bson:"last_page_url" xml:"last_page_url" yaml:"last_page_url"`
	// NextPageUrl is the URL of the next page. It is an empty string if the
	// current page is the last page.
	NextPageUrl string `json:"next_page_url" bson:"next_page_url" xml:"next_page_url" yaml:"next_page_url"`
	// PrevPageUrl is the URL of the previous page. It is an empty string if
	// the current page is the first page.
	PrevPageUrl string `json:"prev_page_url" bson:"prev_page_url" xml:"prev_page_url" yaml:"prev_page_url"`
	// Path is the fully qualified URL without query string.
	Path string `json:"path" bson:"path" xml:"path" yaml:"path"`
	// From is the starting 1-based index of the items.
	From int `json:"from" bson:"from" xml:"from" yaml:"from"`
	// To is the ending 1-based index of the items.
	To int `json:"to" bson:"to" xml:"to" yaml:"to"`
	// Data is the list of items.
	Data []*T `json:"data" bson:"data" xml:"data" yaml:"data"`
}

// PQ is an interface that defines the methods for queries to be paginated.
type PQ[I any, Q any] interface {
	Offset(int) *Q
	Limit(int) *Q
	Count(context.Context) (int, error)
	All(context.Context) ([]*I, error)
}

func AttachTo(op *ogen.Operation, description string, itemRef string) {
	FixParamNames(op.Parameters)
	SetResponse(op, description, itemRef)
}

func AttachAs(
	op *ogen.Operation, description string, itemRef string, pageParam string,
	perPageParam string,
) {
	FixParamNamesWith(op.Parameters, pageParam, perPageParam)
	SetResponse(op, "Paginated list of items", itemRef)
}

// FixParamNames fixes the parameter names to be `per_page` and `page`, to be
// used with oapi-codegen.
func FixParamNames(params []*ogen.Parameter) {
	FixParamNamesWith(params, "page", "itemsPerPage")
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
