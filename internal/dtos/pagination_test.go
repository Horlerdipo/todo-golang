package dtos

import "testing"

func TestPaginationOptions_ApplyDefaults(t *testing.T) {
	t.Parallel()

	paginationOptions := PaginationOptions{
		Page:    -1,
		PerPage: 0,
		SortBy:  "",
		Order:   "random",
		Filters: nil,
	}

	paginationOptions.ApplyDefaults()
	if paginationOptions.Page != 1 {
		t.Errorf("expected paginationOptions.Page to be 1, got %v", paginationOptions.Page)
	}
	if paginationOptions.PerPage != 20 {
		t.Errorf("expected paginationOptions.Page to be 1, got %v", paginationOptions.PerPage)
	}
	if paginationOptions.SortBy != "id" {
		t.Errorf("expected paginationOptions.Page to be 1, got %v", paginationOptions.SortBy)
	}
	if paginationOptions.Order != OrderAsc {
		t.Errorf("expected paginationOptions.Order to be %v, got %v", OrderAsc, paginationOptions.Order)
	}
	if paginationOptions.Filters == nil {
		t.Error("expected paginationOptions.Filters to not be nil")
	}
}

func TestPaginationOptions_ValidateFilters(t *testing.T) {
	t.Parallel()

	paginationOptions := PaginationOptions{}
	paginationOptions.ApplyDefaults()
	paginationOptions.Filters = map[string]string{
		"id":     "1",
		"status": "1",
		"title":  "1",
		"pinned": "true",
	}
	paginationOptions.AllowedFilters = map[string]AllowedFilter{
		"title": {
			Type: StringFilter,
		},
		"pinned": {
			Type: BooleanFilter,
		},
	}

	paginationOptions.ValidateFilters()
	t.Log(paginationOptions.Filters)
	if _, ok := paginationOptions.Filters["id"]; ok {
		t.Error("expected paginationOptions.Filters['id'] to have been removed")
	}

	if _, ok := paginationOptions.Filters["status"]; ok {
		t.Error("expected paginationOptions.Filters['status'] to have been removed")
	}

	if _, ok := paginationOptions.Filters["title"]; !ok {
		t.Error("expected paginationOptions.Filters['title'] to not be removed")
	}

	if _, ok := paginationOptions.Filters["pinned"]; !ok {
		t.Error("expected paginationOptions.Filters['pinned'] to not be removed")
	}
}

func TestPaginationOptions_ValidateSortField(t *testing.T) {
	t.Parallel()

	paginationOptions := PaginationOptions{
		SortBy: "unallowed-column",
	}
	paginationOptions.ApplyDefaults()
	paginationOptions.AllowedSortFields = map[string]bool{
		"allowed-column": true,
	}
	paginationOptions.ValidateSortField()

	if paginationOptions.SortBy == "unallowed-column" {
		t.Errorf("expected paginationOptions.SortBy to be defaulted to id, found: %v", paginationOptions.SortBy)
	}

	paginationOptions.SortBy = "allowed-column"
	paginationOptions.ValidateSortField()

	if paginationOptions.SortBy != "allowed-column" {
		t.Errorf("expected paginationOptions.SortBy to be allowed-column, found: %v", paginationOptions.SortBy)
	}
}

func TestPaginationOptions_Offset(t *testing.T) {
	t.Parallel()

	paginationOptions := PaginationOptions{
		Page:    5,
		PerPage: 20,
	}
	paginationOptions.ApplyDefaults()
	offset := paginationOptions.Offset()
	if offset != 80 {
		t.Errorf("expected paginationOptions.Offset to be 80, got %v", offset)
	}
}

func TestPaginationOptions_ConvertFilter(t *testing.T) {
	t.Parallel()
	paginationOptions := PaginationOptions{}
	paginationOptions.ApplyDefaults()
	paginationOptions.Filters = map[string]string{
		"title":          "1",
		"pinned":         "true",
		"number":         "150",
		"unknown-column": "nice",
	}
	paginationOptions.AllowedFilters = map[string]AllowedFilter{
		"title": {
			Type: StringFilter,
		},
		"pinned": {
			Type: BooleanFilter,
		},
		"number": {
			Type: IntegerFilter,
		},
		"error": {
			Type: StringFilter,
		},
	}

	number, err := paginationOptions.ConvertFilter("number")
	if err != nil {
		t.Error(err)
	}
	if number != int64(150) {
		t.Errorf("expected paginationOptions.Filters['number'] to be 150, got %v", number)
	}

	title, err := paginationOptions.ConvertFilter("title")
	if err != nil {
		t.Error(err)
	}
	if title != "1" {
		t.Errorf("expected title to be 1, got %v", title)
	}

	_, err = paginationOptions.ConvertFilter("error")
	if err == nil {
		t.Error("paginationOptions.ConvertFilter(\"error\") should have been an error")
	}

	if err != nil {
		if err.Error() != "invalid filter column: error" {
			t.Errorf("unexpected error: %v", err)
		}
	}

	_, err = paginationOptions.ConvertFilter("unknown-column")
	if err == nil {
		t.Error("paginationOptions.ConvertFilter(\"unknown-column\") should have been an error")
	}

	if err != nil {
		if err.Error() != "filter unknown-column is not allowed" {
			t.Errorf("unexpected error: %v", err)
		}
	}
}
