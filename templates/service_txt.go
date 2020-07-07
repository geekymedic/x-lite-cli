package templates

const (
	serviceMainTxt = `package services

func Main() error {
	return nil
}`

	serviceServerTxt = `package cal

type Server struct {
	star int
}

var s *Server

func init() {
	s = &Server{}
}

func CalServer() *Server {
	return s
}
`

	serviceImlTxt = `package cal

import (
	"context"

	validate "github.com/geekymedic/x-lite/pkg/xvalidator"
)

// @type: s.i.rt
// @desc:
type SumRequest struct {
	N int
}

// @type: s.i.re
// @desc:
type SumResponse struct {
	Ret int
}

// @type: s.i
// @desc:
func (s *Server) SumHandler(ctx context.Context, req *SumRequest) (*SumResponse, error) {
	if err := validate.ValidateStruct(req); err != nil {
		return nil, err
	}
	return &SumResponse{
		Ret: req.N + s.star,
	}, nil
}`
)
