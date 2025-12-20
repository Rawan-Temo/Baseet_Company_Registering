package utils

import (
	"strconv"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type QueryBuilder struct {
	DB          *gorm.DB
	Query       map[string][]string
	AllowedCols map[string]bool // whitelist
}

func NewQueryBuilder(db *gorm.DB, query map[string][]string, allowedCols []string) *QueryBuilder {
	cols := make(map[string]bool)
	for _, col := range allowedCols {
		cols[col] = true
	}
	return &QueryBuilder{DB: db, Query: query, AllowedCols: cols}
}

func (q *QueryBuilder) Filter() *QueryBuilder {
	db := q.DB
	for key, values := range q.Query {

		if key == "page" || key == "sort" || key == "limit" || key == "fields" {
			continue
		}
		if len(values) == 0 {
			continue
		}

		value := values[0]
		field := key

		switch {
		case strings.HasSuffix(key, "_gte"):
			field = strings.TrimSuffix(key, "_gte")
			if q.AllowedCols[field] {
				db = db.Where(field+" >= ?", value)
			}
		case strings.HasSuffix(key, "_gt"):
			field = strings.TrimSuffix(key, "_gt")
			if q.AllowedCols[field] {
				db = db.Where(field+" > ?", value)
			}
		case strings.HasSuffix(key, "_lte"):
			field = strings.TrimSuffix(key, "_lte")
			if q.AllowedCols[field] {
				db = db.Where(field+" <= ?", value)
			}
		case strings.HasSuffix(key, "_lt"):
			field = strings.TrimSuffix(key, "_lt")
			if q.AllowedCols[field] {
				db = db.Where(field+" < ?", value)
			}
		case strings.HasSuffix(key, "_contains"):
			field = strings.TrimSuffix(key, "_contains")
			if q.AllowedCols[field] {
				db = db.Where(field+" ILIKE ?", "%"+value+"%")
			}
		default:
			if q.AllowedCols[field] {
				db = db.Where(field+" = ?", value)
			}
		}
	}
	q.DB = db
	return q
}

func (q *QueryBuilder) Sort() *QueryBuilder {
	sort, ok := q.Query["sort"]
	if ok && len(sort) > 0 {
		fields := strings.Split(sort[0], ",")
		for _, field := range fields {
			desc := false
			if strings.HasPrefix(field, "-") {
				desc = true
				field = strings.TrimPrefix(field, "-")
			}
			if q.AllowedCols[field] {
				q.DB = q.DB.Order(clause.OrderByColumn{Column: clause.Column{Name: field}, Desc: desc})
			}
		}
	}
	return q
}

func (q *QueryBuilder) LimitFields() *QueryBuilder {
	fields, ok := q.Query["fields"]
	if ok && len(fields) > 0 {
		fieldList := strings.Split(fields[0], ",")
		validFields := make([]string, 0)
		for _, f := range fieldList {
			if q.AllowedCols[f] {
				validFields = append(validFields, f)
			}
		}
		if len(validFields) > 0 {
			q.DB = q.DB.Select(validFields)
		}
	}
	return q
}

func (q *QueryBuilder) Paginate() *QueryBuilder {
	page := 1
	limit := 100
	if p, ok := q.Query["page"]; ok && len(p) > 0 {
		if v, err := strconv.Atoi(p[0]); err == nil && v > 0 {
			page = v
		}
	}
	if l, ok := q.Query["limit"]; ok && len(l) > 0 {
		if v, err := strconv.Atoi(l[0]); err == nil && v > 0 && v <= 1000 {
			limit = v
		}
	}
	skip := (page - 1) * limit
	q.DB = q.DB.Offset(skip).Limit(limit)
	return q
}

func (q *QueryBuilder) Apply() *gorm.DB {
	return q.DB
}
