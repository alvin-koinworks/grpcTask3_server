package service

import (
	"context"
	proto "serverGRPC/resources/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Deposit struct {
	Amount float32
}

var Dep Deposit

type DepositService struct {
	proto.UnimplementedDepositServiceServer
}

func NewDepositServiceServer() *DepositService {
	return &DepositService{}
}

func (d *DepositService) Deposit(ctx context.Context, in *proto.DepositRequest) (*proto.DepositResponse, error) {
	if in.GetAmount() < 10000 {
		return nil, status.Errorf(codes.InvalidArgument, "cannot deposit %v", in.GetAmount())
	}

	Dep.Amount = Dep.Amount + in.GetAmount()

	return &proto.DepositResponse{Ok: true}, nil
}

func (ds *DepositService) GetDeposit(ctx context.Context, in *proto.GetDepositRequest) (*proto.GetDepositResponse, error) {
	getResponse := proto.GetDepositResponse{
		TotalAmount: Dep.Amount,
	}

	return &getResponse, nil
}
