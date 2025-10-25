package utils

import "strconv"

type Pagination struct {
	Page         int32 `json:"page"`
	Limit        int32 `json:"limit"`
	TotalRecords int32 `json:"total"`
	TotalPage    int32 `json:"total_page"`
	HasNext      bool  `json:"has_next"`
	HasPrev      bool  `json:"has_prev"`
}

func NewPagination(page, limit, totalRecords int32) *Pagination {
	//totalPage:=math.Ceil(float64(totalRecordss) / float64(limit))
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		envLimit := GetEnv("LIMIT_ITEM_PER_PAGE", "10")
		limitInt, err := strconv.Atoi(envLimit)
		if err != nil && limitInt <= 0 {
			limitInt = 10
		}
		limit = int32(limitInt)
	}
	totalPage := (totalRecords + limit - 1) / limit

	return &Pagination{
		Page:         page,
		Limit:        limit,
		TotalRecords: totalRecords,
		TotalPage:    totalPage,
		HasNext:      page < totalPage,
		HasPrev:      page > 1,
	}
}

func NewPaginationResponse(data any, page, limit, totalRecords int32) map[string]interface{} {

	return map[string]interface{}{
		"pagination": NewPagination(page, limit, totalRecords),
		"data":       data,
	}
}
