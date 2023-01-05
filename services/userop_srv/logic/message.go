package logic

import (
	"context"
	"database/sql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"luke544187758/userop-srv/dao/mysql"
	"luke544187758/userop-srv/models"
	"luke544187758/userop-srv/pkg/snowflake"
	"luke544187758/userop-srv/proto"
	"time"
)

type MessageService struct {
}

func NewMessageService() *MessageService {
	return &MessageService{}
}

func (m *MessageService) MessageList(ctx context.Context, req *proto.MessageRequest) (*proto.MessageListResponse, error) {
	resp := new(proto.MessageListResponse)
	messages, err := mysql.GetMessages(req.UserId)
	if err != nil {
		return nil, err
	}
	resp.Total = int32(len(messages))
	for _, item := range messages {
		msg := &proto.MessageResponse{
			Id:          item.Id,
			UserId:      item.UserId,
			MessageType: item.MessageType,
			Subject:     item.Subject.String,
			Message:     item.Subject.String,
			File:        item.Subject.String,
		}
		resp.Data = append(resp.Data, msg)
	}

	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}

	return resp, nil
}

func (m *MessageService) CreateMessage(ctx context.Context, req *proto.MessageRequest) (*proto.MessageResponse, error) {
	msg := new(models.Message)
	msg.Id = snowflake.GenID()
	msg.UserId = req.UserId
	msg.Subject = sql.NullString{String: req.Subject}
	msg.MessageType = req.MessageType
	msg.Message = sql.NullString{String: req.Message}
	msg.File = sql.NullString{String: req.File}
	msg.IsDeleted = false
	msg.AddTime = sql.NullTime{Time: time.Now()}
	msg.UpdateTime = sql.NullTime{Time: time.Now()}

	if err := mysql.InsertMessage(msg); err != nil {
		return nil, err
	}

	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}

	resp := &proto.MessageResponse{
		Id:          msg.Id,
		UserId:      msg.UserId,
		MessageType: msg.MessageType,
		Subject:     msg.Subject.String,
		Message:     msg.Message.String,
		File:        msg.File.String,
	}
	return resp, nil
}
