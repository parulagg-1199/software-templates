package services

import (
	"context"
	"log"
	"os"

	"git.xenonstack.com/lib/golang-boilerplate/auth"

	"git.xenonstack.com/util/golang-boilerplate/grpc-example/micro-service1/config"
	pb "git.xenonstack.com/util/golang-boilerplate/grpc-example/micro-service1/pb"
	"git.xenonstack.com/util/golang-boilerplate/grpc-example/micro-service1/src/dbtypes"
	"git.xenonstack.com/util/golang-boilerplate/grpc-example/micro-service1/src/details"
	s2pb "git.xenonstack.com/util/golang-boilerplate/grpc-example/micro-service2/pb"

	ot "github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Server is a structure used for calling all the rpc's methods
type Server struct {
	Service2Client s2pb.UserInformationClient
}

// GetS2Client is a function used to connect to micro-service2
func GetS2Client() s2pb.UserInformationClient {
	// create connection with another service
	conn, err := grpc.Dial(config.Conf.OtherServices.MicroService2, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to start gRPC connection: %v", err)
	}
	//defer conn.Close()

	return s2pb.NewUserInformationClient(conn)
}

// GetUserDetails is a rpc method to fetch user details from database and user information from another grpc service
func (s *Server) GetUserDetails(ctx context.Context, req *pb.EmptyRequest) (*pb.UserDetailInfo, error) {
	// start span from context
	zap.S().Info("get user details and information")
	span, _ := ot.StartSpanFromContext(ctx, "add user details")
	defer span.Finish()
	span.SetTag("event", "fetch user details and information")

	span.LogKV("task", "extract jwt claims and token")
	//extracting jwt claims
	os.Setenv("PRIVATE_KEY", config.Conf.JWT.PrivateKey)
	claims, token, err := auth.ValidateToken(ctx)
	zap.S().Info(claims)
	if err != nil {
		zap.S().Error("error in fetching token....", err)
		span.LogKV("task", "send final output after fetching token and claims")
		return &pb.UserDetailInfo{}, status.Error(codes.Unauthenticated, err.Error())
	}
	span.LogKV("task", "extract email data from jwt claims")
	// extract email from jwt claims and before assigning check email exists in claims map
	val, ok := claims["email"]
	zap.S().Info(val)
	if !ok {
		// if email is not there
		span.LogKV("task", "send final output after extract email from claims")
		zap.S().Error("email claim is not set")
		return &pb.UserDetailInfo{}, status.Error(codes.Internal, "email claim is not set")
	}

	// fetch user information from second micro-service
	span.LogKV("task", "fetch user information from other service")
	s2ctx := context.Background()
	// set authorization
	s2ctx = metadata.AppendToOutgoingContext(s2ctx, "Authorization", "Bearer "+token)
	// call rpc method to fetch user information
	reply, err := s.Service2Client.GetUserInformation(s2ctx, &s2pb.EmptyRequest{})
	zap.S().Info(reply, "===", err)

	if err != nil {
		span.LogKV("task", "send final output after fetching user information")
		zap.S().Error(err)
		return &pb.UserDetailInfo{}, err
	}
	span.LogKV("task", "fetch user details")
	// pass email in baggage
	span.SetBaggageItem("email", val.(string))
	// fetch user details from mysql database
	data, err := details.Fetch(span)
	zap.S().Info(data)
	if err != nil {
		zap.S().Error(err)
		return &pb.UserDetailInfo{}, err
	}
	// send final details
	result := convertData(data)
	return &pb.UserDetailInfo{
		UserDetail: &result,
		UserInfo:   reply.GetUserInfo(),
	}, status.Error(codes.OK, "user detail and information successfully fetched")
}

func convertData(data dbtypes.UserDetail) pb.UserDetail {
	return pb.UserDetail{
		Name:  data.Name,
		Phone: data.Contact,
	}
}

// AddUserDetails is a rpc method to save details in mysql database
func (s *Server) AddUserDetails(ctx context.Context, req *pb.UserDetail) (*pb.AddReply, error) {
	// start span using context
	zap.S().Info("add user details")
	zap.S().Info(req.Name, "===", req.Phone, "====", req.GetName(), "======", req.GetPhone())
	span, _ := ot.StartSpanFromContext(ctx, "add user details")
	defer span.Finish()
	span.SetTag("event", "insert user details in db")

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
	if ok {
		// if email is there save data in db
		span.LogKV("task", "save data in db")
		data := dbtypes.UserDetail{
			Name:    req.Name,
			Contact: req.Phone,
			Email:   val.(string),
		}
		msg, err := details.Save(span, data)
		span.LogKV("task", "send final output after saving data in db")
		if err != nil {
			zap.S().Error(err)
			return &pb.AddReply{Error: true, Message: msg}, status.Error(codes.Internal, msg)
		}
		zap.S().Info(msg)
		return &pb.AddReply{Error: false, Message: msg}, status.Error(codes.OK, msg)
	}
	// if email is not there
	span.LogKV("task", "send final output after extract email from claims")
	zap.S().Error("email claim is not set")
	return &pb.AddReply{Error: true, Message: "email claim is not set"}, status.Error(codes.Internal, "email claim is not set")
}
