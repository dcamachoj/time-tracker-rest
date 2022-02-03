package common

import (
	"context"
	"math"
)

// DefaultPage
const DefaultPage int = 20

type Resthandler func(ctx context.Context) *Response

// RequestPage
type RequestPage struct {
	Index int `form:"index"`
	Size  int `form:"size"`
}

// RequestId
type RequestId struct {
	ID int64 `uri:"id" binding:"required"`
}

// RequestChildId
type RequestChildId struct {
	ID      int64 `uri:"id" binding:"required"`
	ChildID int64 `uri:"idChild" binding:"required"`
}

// RequestPageId
type RequestPageId struct {
	ID    int64 `uri:"id" binding:"required"`
	Index int   `json:"index,omitempty"`
	Size  int   `json:"size,omitempty"`
}

// ResponsePage
type ResponsePage struct {
	Index int `json:"index"`
	Size  int `json:"size"`
	Count int `json:"count"`
	Total int `json:"total"`
}

// ResponseMap
type ResponseMap map[string]interface{}

func (p *ResponsePage) SetTotal(value int64) *ResponsePage {
	p.Total = int(value)
	p.Count = 1
	if p.Size > 0 {
		p.Count = int(math.Ceil(float64(value) / float64(p.Size)))
	}
	return p
}
func (p *ResponsePage) Assign(src *RequestPage) *ResponsePage {
	p.Index = src.Index
	p.Size = src.Size
	p.Count = 1
	p.Total = 0
	return p
}
func (p *ResponsePage) AssignId(src *RequestPageId) *ResponsePage {
	p.Index = src.Index
	p.Size = src.Size
	p.Count = 1
	p.Total = 0
	return p
}
