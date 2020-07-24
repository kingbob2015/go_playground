package main

import (
	"context"
	"fmt"
	"grpc_go_course/calculator/calcpb"
	"io"
	"log"
	"math"
	"net"
	"time"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calcpb.SumRequest) (*calcpb.SumResponse, error) {
	num1 := req.GetNum_1()
	num2 := req.GetNum_2()
	res := &calcpb.SumResponse{
		Result: num1 + num2,
	}
	return res, nil
}

func (*server) PrimeNumberDecomposition(req *calcpb.PrimeNumberDecompositionRequest, stream calcpb.CalculatorService_PrimeNumberDecompositionServer) error {
	num := req.GetNum()
	k := int32(2)
	for num > 1 {
		if num%k == 0 {
			res := &calcpb.PrimeNumberDecompositionResponse{
				Result: k,
			}
			num /= k
			stream.Send(res)
			time.Sleep(1000 * time.Millisecond)
		} else {
			k++
		}
	}
	return nil
}

func (*server) ComputeAverage(stream calcpb.CalculatorService_ComputeAverageServer) error {
	sum, totalNums := 0, 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			if totalNums == 0 {
				return stream.SendAndClose(&calcpb.ComputeAverageResponse{
					Result: 0,
				})
			}
			return stream.SendAndClose(&calcpb.ComputeAverageResponse{
				Result: float32(sum) / float32(totalNums),
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}
		sum += int(req.GetNum())
		totalNums++
	}
}

func (*server) FindMaxmimum(stream calcpb.CalculatorService_FindMaxmimumServer) error {
	max := int64(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
			return err
		}
		num := req.GetNum()
		if num > max {
			max = num
			stream.Send(&calcpb.FindMaxmimumResponse{
				Result: max,
			})
		}
	}
}

func (*server) SquareRoot(ctx context.Context, req *calcpb.SquareRootRequest) (*calcpb.SquareRootResponse, error) {
	num := req.GetNum()

	if num < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number: %v", num),
		)
	}

	return &calcpb.SquareRootResponse{
		Result: math.Sqrt(float64(num)),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calcpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}

}
