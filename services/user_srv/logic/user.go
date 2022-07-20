package logic

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"luke544187758/user-srv/dao/mysql"
	"luke544187758/user-srv/proto"
	"luke544187758/user-srv/utils"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (u *UserService) GetUserList(ctx context.Context, req *proto.PageInfoRequest) (rsp *proto.UserListResponse, err error) {
	users, err := mysql.GetUserList(req.Page, req.PerSize)
	if err != nil {
		return nil, err
	}
	rsp = new(proto.UserListResponse)
	rsp.Total = int32(len(users))

	for _, user := range users {
		info := proto.UserInfoResponse{
			Id:       user.ID,
			Role:     user.Role,
			Password: user.Password,
			Mobile:   user.Mobile,
			NickName: user.NickName,
			Gender:   user.Gender.String,
			Address:  user.Address.String,
			Birthday: user.Birthday.Time.String(),
		}
		rsp.Data = append(rsp.Data, &info)
	}
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}
	return rsp, nil
}

func (u *UserService) GetUserByMobile(ctx context.Context, req *proto.UserMobileRequest) (rsp *proto.UserInfoResponse, err error) {
	rsp = new(proto.UserInfoResponse)
	user, err := mysql.GetUserByMobile(req.Mobile)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, status.Errorf(codes.NotFound, "the user is not found")
	}

	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}
	rsp.Id = user.ID
	rsp.Mobile = user.Mobile
	rsp.Birthday = user.Birthday.Time.Format("2006-01-02")
	rsp.NickName = user.NickName
	rsp.Gender = user.Gender.String
	rsp.Password = user.Password
	rsp.Address = user.Address.String
	rsp.Role = user.Role
	return rsp, nil
}

func (u *UserService) GetUserByID(ctx context.Context, req *proto.UserIdRequest) (rsp *proto.UserInfoResponse, err error) {
	rsp = new(proto.UserInfoResponse)
	user, err := mysql.GetUserByID(req.Id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, status.Errorf(codes.NotFound, "the user is not found")
	}

	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}

	rsp.Id = user.ID
	rsp.Mobile = user.Mobile
	rsp.Birthday = user.Birthday.Time.String()
	rsp.NickName = user.NickName
	rsp.Gender = user.Gender.String
	rsp.Password = user.Password
	rsp.Address = user.Address.String
	rsp.Role = user.Role
	return rsp, nil
}

func (u *UserService) CreateUser(ctx context.Context, req *proto.UserInfoRequest) (rsp *proto.CreateUserResponse, err error) {
	rsp = new(proto.CreateUserResponse)
	user, err := mysql.GetUserByMobile(req.Mobile)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return nil, status.Errorf(codes.AlreadyExists, "user is already exist")
	}
	id, err := mysql.CreateUser(req.NickName, req.Mobile, req.Password)
	if err != nil {
		return nil, err
	}
	rsp.Id = id
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}
	return rsp, nil
}

func (u *UserService) ModifyUser(ctx context.Context, req *proto.ModifyUserInfoRequest) (*emptypb.Empty, error) {
	user, err := mysql.GetUserByID(req.Id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, status.Errorf(codes.NotFound, "the user is not found")
	}
	if err = mysql.ModifyUserInfo(req.Id, req.NickName, req.Gender, req.Birthday); err != nil {
		return nil, err
	}
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}
	return nil, nil
}

func (u *UserService) RemoveUser(ctx context.Context, req *proto.UserIdRequest) (*emptypb.Empty, error) {
	user, err := mysql.GetUserByID(req.Id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, status.Errorf(codes.NotFound, "the user is not found")
	}
	if err = mysql.RemoveUserInfo(req.Id); err != nil {
		return nil, err
	}

	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}
	return nil, nil
}

func (u *UserService) CheckPassword(ctx context.Context, req *proto.PasswordCheckRequest) (rsp *proto.CheckPasswordResponse, err error) {
	rsp = new(proto.CheckPasswordResponse)
	rsp.Success = utils.EncryptPassword(req.Password) == req.EncryptPwd
	return rsp, nil
}
