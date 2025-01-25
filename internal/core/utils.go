package core

//func parseQueryParams(r *http.Request) handlers.QueryParams {
//	params := handlers.QueryParams{
//		Page:     1,
//		PageSize: 10,
//		SortBy:   "name",
//		SortDir:  "asc",
//	}
//
//	if page := r.URL.Query().Get("page"); page != "" {
//		if pageNum, err := strconv.Atoi(page); err == nil && pageNum > 0 {
//			params.Page = pageNum
//		}
//	}
//
//	if pageSize := r.URL.Query().Get("page_size"); pageSize != "" {
//		if size, err := strconv.Atoi(pageSize); err == nil && size > 0 {
//			params.PageSize = size
//		}
//	}
//
//	params.Search = r.URL.Query().Get("search")
//
//	if sortBy := r.URL.Query().Get("sort_by"); sortBy != "" {
//		params.SortBy = sortBy
//	}
//
//	if sortDir := r.URL.Query().Get("sort_dir"); sortDir != "" {
//		params.SortDir = sortDir
//	}
//
//	return params
//}
