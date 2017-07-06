package methods

import (
	"context"

	types "github.com/lox/govmomi-vmfork/internal/types"
	"github.com/vmware/govmomi/vim25/soap"
)

type EnableForkParent_TaskBody struct {
	Req    *types.EnableForkParent_Task         `xml:"urn:vim25 EnableForkParent_Task,omitempty"`
	Res    *types.EnableForkParent_TaskResponse `xml:"urn:vim25 EnableForkParent_TaskResponse,omitempty"`
	Fault_ *soap.Fault                          `xml:"http://schemas.xmlsoap.org/soap/envelope/ Fault,omitempty"`
}

func (b *EnableForkParent_TaskBody) Fault() *soap.Fault { return b.Fault_ }

func NewEnableForkParent_Task(ctx context.Context, r soap.RoundTripper, req *types.EnableForkParent_Task) (*types.EnableForkParent_TaskResponse, error) {
	var reqBody, resBody EnableForkParent_TaskBody

	reqBody.Req = req

	if err := r.RoundTrip(ctx, &reqBody, &resBody); err != nil {
		return nil, err
	}

	return resBody.Res, nil
}

type CreateForkChild_TaskBody struct {
	Req    *types.CreateForkChild_Task         `xml:"urn:vim25 CreateForkChild_Task,omitempty"`
	Res    *types.CreateForkChild_TaskResponse `xml:"urn:vim25 CreateForkChild_TaskResponse,omitempty"`
	Fault_ *soap.Fault                         `xml:"http://schemas.xmlsoap.org/soap/envelope/ Fault,omitempty"`
}

func (b *CreateForkChild_TaskBody) Fault() *soap.Fault { return b.Fault_ }

func NewCreateForkChild_Task(ctx context.Context, r soap.RoundTripper, req *types.CreateForkChild_Task) (*types.CreateForkChild_TaskResponse, error) {
	var reqBody, resBody CreateForkChild_TaskBody

	reqBody.Req = req

	if err := r.RoundTrip(ctx, &reqBody, &resBody); err != nil {
		return nil, err
	}

	return resBody.Res, nil
}
