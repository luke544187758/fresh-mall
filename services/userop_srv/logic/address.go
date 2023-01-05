package logic

import (
	"context"
	"database/sql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"luke544187758/userop-srv/dao/mysql"
	"luke544187758/userop-srv/models"
	"luke544187758/userop-srv/pkg/snowflake"
	"luke544187758/userop-srv/proto"
	"time"
)

type AddressService struct {
}

func NewAddressService() *AddressService {
	return &AddressService{}
}

func (a *AddressService) GetAddressList(ctx context.Context, req *proto.AddressRequest) (*proto.AddressListResponse, error) {
	resp := new(proto.AddressListResponse)
	addresses, err := mysql.GetAddresses(req.UserId)
	if err != nil {
		return nil, err
	}
	resp.Total = int32(len(addresses))
	for _, item := range addresses {
		addr := &proto.AddressResponse{
			Id:           item.Id,
			UserId:       item.UserId,
			Province:     item.Province,
			City:         item.City,
			District:     item.District,
			Address:      item.Address,
			SignerName:   item.SignerName,
			SignerMobile: item.SignerMobile,
		}
		resp.Data = append(resp.Data, addr)
	}
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}
	return resp, nil
}

func (a *AddressService) CreateAddress(ctx context.Context, req *proto.AddressRequest) (*proto.AddressResponse, error) {
	resp := new(proto.AddressResponse)
	addr := &models.Address{
		Id:           snowflake.GenID(),
		UserId:       req.UserId,
		Province:     req.Province,
		City:         req.City,
		District:     req.District,
		Address:      req.Address,
		SignerName:   req.SignerName,
		SignerMobile: req.SignerMobile,
		AddTime:      sql.NullTime{Time: time.Now()},
		UpdateTime:   sql.NullTime{Time: time.Now()},
		IsDeleted:    false,
	}

	if err := mysql.InsertAddress(addr); err != nil {
		return nil, err
	}

	resp.Id = addr.Id

	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}
	return resp, nil
}

func (a *AddressService) DeleteAddress(ctx context.Context, req *proto.AddressRequest) (*emptypb.Empty, error) {

	if err := mysql.DeleteAddress(req.Id); err != nil {
		return &emptypb.Empty{}, err
	}

	if ctx.Err() == context.Canceled {
		return &emptypb.Empty{}, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return &emptypb.Empty{}, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}
	return &emptypb.Empty{}, nil
}

func (a *AddressService) UpdateAddress(ctx context.Context, req *proto.AddressRequest) (*emptypb.Empty, error) {
	addr := &models.Address{
		Province:     req.Province,
		City:         req.City,
		District:     req.District,
		Address:      req.Address,
		SignerName:   req.SignerName,
		SignerMobile: req.SignerMobile,
		UpdateTime:   sql.NullTime{Time: time.Now()},
	}

	if err := mysql.UpdateAddress(addr); err != nil {
		return &emptypb.Empty{}, err
	}

	if ctx.Err() == context.Canceled {
		return &emptypb.Empty{}, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return &emptypb.Empty{}, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}
	return &emptypb.Empty{}, nil
}
