package grpc

import (
	"context"
	"strconv"

	"github.com/RedWood011/ServiceURL/internal/apperror"
	"github.com/RedWood011/ServiceURL/internal/entities"
	"github.com/RedWood011/ServiceURL/internal/service"
	"github.com/RedWood011/ServiceURL/internal/transport/grpc/pb"
	"github.com/pkg/errors"
)

func NewGRPCServer(service service.Translation) *URLServer {
	return &URLServer{
		service: service,
	}
}

type URLServer struct {
	pb.UnimplementedURLServer
	service service.Translation
}

func (us *URLServer) GetURLByID(ctx context.Context, in *pb.RetrieveRequest) (*pb.RetrieveResponse, error) {
	var appErr *apperror.AppError
	long, err := us.service.GetURLByID(ctx, in.ShortUrlId)
	if err != nil {

		switch errors.As(err, &appErr) {
		case errors.Is(err, apperror.ErrGone):
			return &pb.RetrieveResponse{
				Status: "gone",
			}, nil
		case errors.Is(err, apperror.ErrNotFound):
			return &pb.RetrieveResponse{
				Status: "not found",
			}, nil
		default:
			return &pb.RetrieveResponse{
				Status: "internal server error",
			}, nil
		}
	}
	return &pb.RetrieveResponse{
		RedirectUrl: long,
		Status:      "ok",
	}, nil
}

func (us *URLServer) PostOneURL(ctx context.Context, in *pb.CreateRequest) (*pb.CreateResponse, error) {
	var appErr *apperror.AppError
	responseURL, err := us.service.CreateShortURL(ctx, entities.URL{FullURL: in.OriginalUrl,
		UserID: in.UserId})
	if err != nil {
		switch errors.As(err, &appErr) {
		case errors.Is(err, apperror.ErrConflict):
			return &pb.CreateResponse{
				Status: "conflict",
			}, nil
		default:
			return &pb.CreateResponse{
				Status: "internal server error",
			}, nil
		}
	}
	return &pb.CreateResponse{
		Status:      "ok",
		ResponseUrl: responseURL,
	}, nil
}

func (us *URLServer) GetUserURLs(ctx context.Context, in *pb.GetUserURLsRequest) (*pb.GetUserURLsResponse, error) {
	var appErr *apperror.AppError
	urls, err := us.service.GetAllURLsByUserID(ctx, in.UserId)
	if err != nil {

		switch errors.As(err, &appErr) {
		case errors.Is(err, apperror.ErrNoContent):
			return &pb.GetUserURLsResponse{
				Status: "no content",
			}, nil
		default:
			return &pb.GetUserURLsResponse{
				Status: "internal server error",
			}, nil
		}
	}
	var result []*pb.GetUserURLsResponse_URL
	for i := 0; i < len(urls); i++ {
		result = append(result, &pb.GetUserURLsResponse_URL{
			OriginalUrl: urls[0].FullURL,
			ShortUrl:    urls[0].ShortURL,
		})
	}
	return &pb.GetUserURLsResponse{
		Status: "ok",
		URLs:   result,
	}, nil
}

func (us *URLServer) PostBatchURLs(ctx context.Context, in *pb.CreateBatchRequest) (*pb.CreateBatchResponse, error) {
	var data []entities.URL
	for i := 0; i < len(in.URLs); i++ {
		data = append(data, entities.URL{
			FullURL: in.URLs[i].OriginalUrl,
			UserID:  in.UserId,
		})
	}
	urls, err := us.service.CreateShortURLs(ctx, data)
	if err != nil {
		return &pb.CreateBatchResponse{
			Status: "internal server error",
		}, nil
	}
	var response []*pb.CreateBatchResponse_URL
	for i := 0; i < len(urls); i++ {
		id, _ := strconv.ParseInt(urls[i].CorrelationID, 10, 32)
		response = append(response, &pb.CreateBatchResponse_URL{
			CorrelationId: int32(id),
			ShortUrl:      urls[i].ShortURL,
		})
	}
	return &pb.CreateBatchResponse{
		Status: "ok",
		Urls:   response,
	}, nil
}

func (us *URLServer) DeleteBatchURLs(ctx context.Context, in *pb.DeleteBatchRequest) (*pb.DeleteBatchResponse, error) {
	us.service.DeleteShortURLs(ctx, in.Urls, in.UserId)
	return &pb.DeleteBatchResponse{
		Status: "accepted",
	}, nil
}
