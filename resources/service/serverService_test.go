package service_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"

	proto "serverGRPC/resources/proto"
	"serverGRPC/resources/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func dial() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()
	proto.RegisterDepositServiceServer(server, &service.DepositService{})

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()
	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestDepositService_Deposit(t *testing.T) {
	test := []struct {
		name    string
		amount  float32
		res     *proto.DepositResponse
		errCode codes.Code
		errMsg  string
	}{
		{
			"Invalid Request with Invalid Amount",
			5000,
			nil,
			codes.InvalidArgument,
			fmt.Sprintf("cannot deposit %v", 5000),
		},
		{
			"Valid Request",
			10000,
			&proto.DepositResponse{Ok: true},
			codes.OK,
			"",
		},
	}

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dial()))

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := proto.NewDepositServiceClient(conn)

	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			request := &proto.DepositRequest{Amount: tc.amount}

			response, err := client.Deposit(ctx, request)

			if response != nil {
				if response.GetOk() != tc.res.Ok {
					t.Error("response: expected", tc.res.GetOk(), "received", response.GetOk())
				}
			}

			if err != nil {
				if er, ok := status.FromError(err); ok {
					if er.Code() != tc.errCode {
						t.Error("error code: expected", codes.InvalidArgument, "received", er.Code())
					}
					if er.Message() != tc.errMsg {
						t.Error("error message: expected", tc.errMsg, "received", er.Message())
					}
				}
			}

		})
	}
}

func TestDepositService_GetDeposit(t *testing.T) {
	test := struct {
		name string
		res  *proto.GetDepositResponse
		err  error
	}{
		"Testing getting total depost",
		&proto.GetDepositResponse{TotalAmount: 10000},
		nil,
	}

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dial()))

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := proto.NewDepositServiceClient(conn)

	t.Run(test.name, func(t *testing.T) {
		request := &proto.GetDepositRequest{}

		response, err := client.GetDeposit(ctx, request)

		if response.TotalAmount != test.res.TotalAmount {
			t.Error("Response expected", test.res, "received", response)
		}

		if err != nil {
			t.Error("Response expected", test.err, "received", err)
		}
	})
}
