package dtos

import (
	"errors"
	"fmt"
	"strconv"
)

type Order string

func (o Order) IsValid() bool {
	return o == OrderAsc || o == OrderDesc
}

const (
	OrderAsc  Order = "asc"
	OrderDesc Order = "desc"
)

type PaginationOptions struct {
	Page              int                      `json:"page"`
	PerPage           int                      `json:"per_page"`
	SortBy            string                   `json:"sort_by"`
	Order             Order                    `json:"orderBy"`
	Filters           map[string]string        `json:"filters"`
	AllowedSortFields map[string]bool          `json:"-"`
	AllowedFilters    map[string]AllowedFilter `json:"-"`
}

type AllowedFilter struct {
	Type FilterType
}

type FilterType string

const (
	IntegerFilter FilterType = "integer"
	BooleanFilter FilterType = "boolean"
	StringFilter  FilterType = "string"
)

func (p *PaginationOptions) ApplyDefaults() {
	if p.Page <= 0 {
		p.Page = 1
	}

	if p.PerPage <= 0 {
		p.PerPage = 20
	}

	if p.SortBy == "" {
		p.SortBy = "id"
	}

	if !p.Order.IsValid() {
		p.Order = OrderAsc
	}

	if p.Filters == nil {
		p.Filters = make(map[string]string)
	}
}

func (p *PaginationOptions) ValidateSortField() {
	if !p.AllowedSortFields[p.SortBy] {
		p.SortBy = "id"
	}
}

func (p *PaginationOptions) ValidateFilters() {
	for _, filter := range p.Filters {
		if _, ok := p.AllowedFilters[filter]; !ok {
			delete(p.Filters, filter)
		}
	}
}

func (p *PaginationOptions) Configure() {
	p.ApplyDefaults()
	p.ValidateSortField()
	p.ValidateFilters()
}

func (p *PaginationOptions) Offset() int {
	return (p.Page - 1) * p.PerPage
}

func (p *PaginationOptions) ConvertFilter(column string) (interface{}, error) {
	filter, ok := p.Filters[column]
	if !ok {
		return nil, errors.New(fmt.Sprintf("invalid filter column: %s", column))
	}

	allowedFilter, ok := p.AllowedFilters[column]
	if !ok {
		return nil, errors.New(fmt.Sprintf("filter %s is not allowed", filter))
	}

	switch allowedFilter.Type {
	case IntegerFilter:
		convertedInt, err := strconv.ParseInt(filter, 10, 32)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("unable to convert filter column: %s value: %v to int", column, filter))
		}
		return convertedInt, nil
	case BooleanFilter:
		if filter == "false" || filter == "0" {
			return false, nil
		}
		return true, nil
	case StringFilter:
		return filter, nil
	default:
		return filter, nil
	}
}

type PaginatedResponse[T any] struct {
	Data []T                   `json:"data"`
	Meta PaginatedResponseMeta `json:"meta"`
}

type PaginatedResponseMeta struct {
	TotalCount  int `json:"total_count"`
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	LastPage    int `json:"last_page"`
	FirstPage   int `json:"first_page"`
}
