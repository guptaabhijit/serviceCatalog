package handlers

import "serviceCatalog/internal/models"

type ListServicesResponse struct {
	Services    []models.ServiceResponse `json:"services"`
	TotalCount  int64                    `json:"total_count"`
	CurrentPage int                      `json:"current_page"`
	PageSize    int                      `json:"page_size"`
}

type QueryParams struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"pageSize,default=10"`
	Search   string `form:"search"`
	SortBy   string `form:"sortBy,default=id"`
	SortDir  string `form:"sortDir,default=asc"`
}
