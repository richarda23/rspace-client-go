package rspace

import (
	"strings"
)

// Boolean combinator - AND or OR
type SearchOperator int

const (
	And SearchOperator = iota
	Or
)

var searchStrings = [2]string{"and", "or"}

//Stringer implementation
func (op SearchOperator) String() string {
	return searchStrings[op]
}

// QueryType restricts search to a particular category
type QueryType int

var qTypeStrings = [9]string{"global", "fullText", "tag", "name", "created", "lastModified", "form", "attachment", "owner"}

const (
	GLOBAL QueryType = iota
	FULL_TEXT
	TAG
	NAME
	CREATED
	LAST_MODIFIED
	FORM
	ATTACHMENT
	OWNER
)

//Stringer implementation
func (op QueryType) String() string {
	return qTypeStrings[op]
}

//SearchTerm is used by SearchQueryBuilder to construct a valid search query
type SearchTerm struct {
	QueryType QueryType
	Term      string
}

//Stringer implementation for SearchTerm
func (op SearchTerm) String() string {
	return "term=" + op.Term + ", queryType=" + op.QueryType.String()
}

type SearchQueryBuilder struct {
	operator SearchOperator
	terms    []SearchTerm
}

//operator sets the boolean type of the search query
func (qb *SearchQueryBuilder) Operator(op SearchOperator) *SearchQueryBuilder {
	qb.operator = op
	return qb
}

func (qb *SearchQueryBuilder) AddGlobalTerm(term string) *SearchQueryBuilder {
	return qb.AddTerm(term, GLOBAL)
}

// addTerm appends a search term in the given category. If the term is empty or nil
// the term is not added
func (qb *SearchQueryBuilder) AddTerm(term string, queryType QueryType) *SearchQueryBuilder {
	if len(term) == 0 {
		return qb
	}
	sterm := SearchTerm{queryType, term}
	if qb.terms == nil {
		qb.terms = make([]SearchTerm, 0)
	}
	qb.terms = append(qb.terms, sterm)
	return qb
}

//build generates a SearchQuery object and returns its pointer
func (qb *SearchQueryBuilder) Build() *SearchQuery {
	rc := SearchQuery{}
	rc.Operator = qb.operator.String()
	rc.Terms = make([]STerm, 0)
	for _, v := range qb.terms {
		t := STerm{v.Term, v.QueryType.String()}
		rc.Terms = append(rc.Terms, t)
	}
	return &rc
}

type SearchQuery struct {
	Operator string  `json:"operator"`
	Terms    []STerm `json:"terms"`
}

type STerm struct {
	Term      string `json:"query"`
	QueryType string `json:"queryType"`
}

func (q *SearchQuery) String() string {
	var b strings.Builder
	for _, v := range q.Terms {
		b.WriteString(v.Term + "=" + v.QueryType)
		b.WriteString(";")
	}
	return q.Operator + " terms:" + b.String()
}
