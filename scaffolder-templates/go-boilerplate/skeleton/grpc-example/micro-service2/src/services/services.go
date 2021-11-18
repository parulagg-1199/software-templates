package services

import (
	"context"
	"os"

	"git.xenonstack.com/lib/golang-boilerplate/auth"

	"git.xenonstack.com/util/golang-boilerplate/grpc-example/micro-service2/config"
	pb "git.xenonstack.com/util/golang-boilerplate/grpc-example/micro-service2/pb"
	"git.xenonstack.com/util/golang-boilerplate/grpc-example/micro-service2/src/info"
	"git.xenonstack.com/util/golang-boilerplate/grpc-example/micro-service2/src/types"

	ot "github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server is a structure used for calling all the rpc's methods
type Server struct {
}

// AddUserInformation is a rpc method to save user information in mongo database
func (s *Server) AddUserInformation(ctx context.Context, req *pb.UserInfo) (*pb.AddReply, error) {
	zap.S().Info(req)
	// start span from context
	span, _ := ot.StartSpanFromContext(ctx, "add user information")
	defer span.Finish()
	span.SetTag("event", "insert user information in mongodb")
	span.LogKV("task", "extract jwt claims")
	//extracting jwt claims
	os.Setenv("PRIVATE_KEY", config.Conf.JWT.PrivateKey)
	claims, _, err := auth.ValidateToken(ctx)
	zap.S().Info(claims)
	if err != nil {
		span.LogKV("task", "send final output after extracting jwt claims")
		zap.S().Error(err)
		return &pb.AddReply{Error: true, Message: "error in extracting jwt claims"}, status.Error(codes.Unauthenticated, err.Error())
	}
	span.LogKV("task", "extract email data from jwt claims")
	// extract email from jwt claims and before assigning check email exists in claims map
	val, ok := claims["email"]
	zap.S().Info(val, " ", ok)
	if ok {
		// if email is there save data in mongodb
		// convert data according to service structure from proto structure
		span.LogKV("task", "save data in a structure")
		data := convertData(req, val.(string))
		// save data in db
		span.LogKV("task", "save data in mongodb")
		err := info.Save(span, data)
		zap.S().Error(err)
		span.LogKV("task", "send final output after saving data in db")
		if err != nil {
			return &pb.AddReply{Error: true, Message: err.Error()}, status.Error(codes.Internal, err.Error())
		}
		return &pb.AddReply{Error: false, Message: "information successfully saved"}, status.Error(codes.OK, "information successfully saved")
	}
	// if email is not there
	span.LogKV("task", "send final output after extract email from claims")
	zap.S().Error("email claim is not set")
	return &pb.AddReply{Error: true, Message: "email claim is not set"}, status.Error(codes.Internal, "email claim is not set")
}

// convertData is a function to convert proto data structure into service specific structure
func convertData(req *pb.UserInfo, email string) types.Info {
	address := types.Address{
		Country: req.Address.GetCountry(),
		State:   req.Address.State,
		City:    req.Address.City,
		Postal:  req.Address.GetPostal(),
	}

	interests := make([]types.Interests, 0)
	for i := 0; i < len(req.Interests); i++ {
		interests = append(interests, types.Interests{
			Interest: req.Interests[i].GetInterest(),
			Priority: req.Interests[i].Priority,
		})
	}
	return types.Info{
		Email:     email,
		Address:   address,
		Interests: interests,
	}
}

// GetUserInformation is a rpc method to fetch user information from mongo database
func (s *Server) GetUserInformation(ctx context.Context, req *pb.EmptyRequest) (*pb.UserInfoReply, error) {
	// start span from context
	span, _ := ot.StartSpanFromContext(ctx, "add user information")
	defer span.Finish()
	span.SetTag("event", "fetch user information from mongodb")
	span.LogKV("task", "extract jwt claims")
	//extracting jwt claims
	os.Setenv("PRIVATE_KEY", config.Conf.JWT.PrivateKey)
	claims, _, err := auth.ValidateToken(ctx)
	zap.S().Info(claims)
	if err != nil {
		span.LogKV("task", "send final output after extracting jwt claims")
		zap.S().Error(err)
		return &pb.UserInfoReply{}, status.Error(codes.Unauthenticated, err.Error())
	}
	span.LogKV("task", "extract email data from jwt claims")
	// extract email from jwt claims and before assigning check email exists in claims map
	val, ok := claims["email"]
	zap.S().Info(val, " ", ok)
	if ok {
		// fetch data from db
		span.LogKV("task", "Fetch data from db")
		data, err := info.Fetch(span, val.(string))
		span.LogKV("task", "send final output fetching data from db")
		if err != nil {
			return &pb.UserInfoReply{Error: true}, status.Error(codes.Internal, err.Error())
		}
		// convert data to proto structure
		reply := reConvertData(data)
		return &pb.UserInfoReply{Error: false, UserInfo: reply}, status.Error(codes.OK, "information fetched succesfully")

	}
	// if email is not there
	span.LogKV("task", "send final output after extract email from claims")
	zap.S().Error("email claim is not set")
	return &pb.UserInfoReply{}, status.Error(codes.Internal, "email claim is not set")
}

// reConvertData is a function to convert service specific structure into proto data structure
func reConvertData(data types.Info) *pb.UserInfo {
	address := pb.Address{
		Country: data.Address.Country,
		State:   data.Address.State,
		City:    data.Address.City,
		Postal:  data.Address.Postal,
	}

	interests := make([]*pb.Interests, 0)
	for i := 0; i < len(data.Interests); i++ {
		interests = append(interests, &pb.Interests{
			Interest: data.Interests[i].Interest,
			Priority: data.Interests[i].Priority,
		})
	}

	return &pb.UserInfo{
		Address:   &address,
		Interests: interests,
	}
}
