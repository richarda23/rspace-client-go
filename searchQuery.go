package rspace
import (
	"strings"
)

type SearchOperator int
const (
	AND SearchOperator = iota
	OR
)
var searchStrings =[2]string{"AND", "OR"}

//Stringer implementation
func (op SearchOperator) String () string {
	return searchStrings[op]
}

// QueryType restricts search to a particulr category
type QueryType int 

var qTypeStrings =[9]string{"global", "fullText", "tag", "name", "created", "lastModified", "form", "attachment", "owner"}
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
func (op QueryType) String () string {
	return qTypeStrings[op]
}

type SearchTerm struct {
	QueryType QueryType
	Term	string
}
func (op SearchTerm) String () string {
	return"term=" + op.Term +", queryType=" + op.QueryType.String() 
}

type SearchQueryBuilder struct {
	Operator  SearchOperator
	Terms []SearchTerm
}
func (qb *SearchQueryBuilder) operator (op SearchOperator) *SearchQueryBuilder {
	qb.Operator = op
	return qb
}
func (qb *SearchQueryBuilder) addTerm (term string, queryType QueryType) *SearchQueryBuilder {
	sterm := SearchTerm{queryType, term}
	if qb.Terms == nil {
		qb.Terms = make([]SearchTerm, 0)
	}
	qb.Terms = append(qb.Terms, sterm)
	return qb
}

func (qb *SearchQueryBuilder) build ()  *SearchQuery {
	rc := SearchQuery{}
	rc.Operator = qb.Operator.String()
	rc.Terms = make([]STerm, 0)
	for _,v := range qb.Terms {
		t := STerm{v.Term, v.QueryType.String()}
		rc.Terms=append(rc.Terms, t)
	}
	return &rc
}

type SearchQuery struct {
	Operator  string
	Terms [] STerm 
}

type STerm struct {
	Term string
	QueryType string
}
func (q *SearchQuery) String () string {
	var b strings.Builder
	for _, v := range q.Terms {
		b.WriteString(v.Term + "="+v.QueryType)
		b.WriteString(";")
	}
	return q.Operator+ " terms:" + b.String()
}
