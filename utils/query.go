package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type QueryParams struct {
	Filters map[string]string
	Limit   int
	Offset  int
	Sort    string
}

func ParseQueryParams(c *gin.Context) QueryParams {
	filters := map[string]string{}
	query := c.Request.URL.Query()

	for key, values := range query {
		if key == "limit" || key == "offset" || key == "sort" {
			continue
		}
		if len(values) > 0 {
			filters[key] = values[0]
		}
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	sort := c.DefaultQuery("sort", "created_at desc")

	return QueryParams{
		Filters: filters,
		Limit:   limit,
		Offset:  offset,
		Sort:    sort,
	}
}
