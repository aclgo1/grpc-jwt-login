package service

import (
	"context"
	"fmt"

	"github.com/aclgo/grpc-jwt/internal/user"
	"github.com/aclgo/grpc-jwt/pkg/grpc_errors"
	"github.com/aclgo/grpc-jwt/proto"
	"github.com/pkg/errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// func (us *UserService) mustEmbedUnimplementedUserServiceServer() {}

func (us *UserService) Register(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreatedUserResponse, error) {
	params := user.ParamsCreateUser{
		Name:     req.Name,
		Lastname: req.LastName,
		Password: req.Password,
		Email:    req.Email,
	}

	created, err := us.userUC.Register(ctx, &params)
	if err != nil {
		return nil, status.Errorf(grpc_errors.ParseGRPCErrors(err), "useUC.Register: %v", err)
	}

	return &proto.CreatedUserResponse{User: parseModelsToProto(created)}, nil
}

func (us *UserService) Login(ctx context.Context, req *proto.UserLoginRequest) (*proto.UserLoginResponse, error) {
	email, password := req.Email, req.Password
	if email == "" || password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Login: %v", grpc_errors.EmptyCredentials{})
	}

	tokens, err := us.userUC.Login(ctx, email, password)
	if err != nil {
		return nil, status.Errorf(grpc_errors.ParseGRPCErrors(err), "Login: %v", err)
	}

	// fmt.Println("tokens generateds", tokens)

	return &proto.UserLoginResponse{
		Tokens: &proto.Tokens{
			AccessToken:  tokens.Access,
			RefreshToken: tokens.Refresh,
		},
	}, nil
}

func (us *UserService) Logout(ctx context.Context, req *proto.UserLogoutRequest) (*proto.UserLogoutResponse, error) {
	accessTK, refreshTK := req.AccessToken, req.RefreshToken

	paramLogout := user.ParamLogoutInput{
		AccessToken:  accessTK,
		RefreshToken: refreshTK,
	}

	if err := paramLogout.Validate(); err != nil {
		return nil, fmt.Errorf("paramLogout.Validate: %w", err)
	}

	if err := us.userUC.Logout(ctx, &paramLogout); err != nil {
		return nil, status.Errorf(grpc_errors.ParseGRPCErrors(err), "Logout: %v", err)
	}

	return &proto.UserLogoutResponse{}, nil
}

func (us *UserService) RefreshTokens(ctx context.Context, req *proto.RefreshTokensRequest) (*proto.RefreshTokensResponse, error) {

	paramRefresh := user.ParamsRefreshTokens{
		AccessToken:  req.AccessToken,
		RefreshToken: req.RefreshToken,
	}

	if err := paramRefresh.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	tokens, err := us.userUC.RefreshTokens(ctx, &paramRefresh)

	//
	// if tokens == nil {
	// 	return nil, fmt.Errorf("us.userUC.RefreshTokens: tokens is nil")
	// }

	if err != nil {
		return nil, fmt.Errorf("us.userUC.RefreshTokens: %w", err)
	}

	out := proto.RefreshTokensResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	return &out, nil
}

func (us *UserService) FindById(ctx context.Context, req *proto.FindByIdRequest) (*proto.FindByIdResponse, error) {
	// tokenString, err := us.getToken(ctx, KeyAccessToken)
	// if err != nil {
	// 	return nil, err
	// }

	// if err := us.userUC.ValidToken(ctx, tokenString); err != nil {
	// 	return nil, status.Errorf(grpc_errors.ParseGRPCErrors(err), "FindById: %v", err)
	// }

	id := req.Id

	found, err := us.userUC.FindByID(ctx, id)
	if err != nil {
		return nil, status.Errorf(grpc_errors.ParseGRPCErrors(err), "userUC.FindByID: %vs", err)
	}

	return &proto.FindByIdResponse{
		User: parseModelsToProto(found),
	}, nil
}

func (us *UserService) FindByEmail(ctx context.Context, req *proto.FindByEmailRequest) (*proto.FindByEmailResponse, error) {
	// tokenString, err := us.getToken(ctx, KeyAccessToken)
	// if err != nil {
	// 	return nil, err
	// }

	// if err := us.userUC.ValidToken(ctx, tokenString); err != nil {
	// 	return nil, status.Errorf(grpc_errors.ParseGRPCErrors(err), "FindByEmail: %v", err)
	// }

	email := req.Email

	found, err := us.userUC.FindByEmail(ctx, email)
	if err != nil {
		return nil, status.Errorf(grpc_errors.ParseGRPCErrors(err), "userUC.FindByEmail: %v", err)
	}

	return &proto.FindByEmailResponse{
		User: parseModelsToProto(found),
	}, nil
}

func (us *UserService) Update(ctx context.Context, req *proto.UpdateRequest) (*proto.UpdateResponse, error) {

	params := user.ParamsUpdateUser{
		UserID:   req.Id,
		Name:     req.Name,
		Lastname: req.Lastname,
		Password: req.Password,
		Email:    req.Email,
		Verified: req.Verified,
	}

	if err := params.Validate(); err != nil {
		return nil, errors.Wrap(err, "Update.Validate")
	}

	updatedUser, err := us.userUC.Update(
		ctx,
		&params,
	)

	if err != nil {
		return nil, status.Errorf(grpc_errors.ParseGRPCErrors(err), "Update.Update: %v", err)
	}

	return &proto.UpdateResponse{
		User: parseModelsToProto(updatedUser),
	}, nil
}

func (us *UserService) Delete(ctx context.Context, req *proto.DeleteRequest) (*proto.DeleteResponse, error) {
	err := us.userUC.Delete(ctx, &user.ParamsDeleteUser{
		UserID: req.Id,
	})

	if err != nil {
		return nil, status.Errorf(grpc_errors.ParseGRPCErrors(err), "Delete.Delete: %v", err)
	}

	return &proto.DeleteResponse{}, nil
}

func (us *UserService) ValidateToken(ctx context.Context, req *proto.ValidateTokenRequest) (*proto.ValidateTokenResponse, error) {
	resp, err := us.userUC.ValidToken(ctx, &user.ParamsValidToken{AccessToken: req.Token})
	if err != nil {
		return nil, err
	}

	return &proto.ValidateTokenResponse{
		UserId:   resp.UserID,
		UserRole: resp.Role,
	}, nil
}

// func (us *UserService) getToken(ctx context.Context, key string) (string, error) {
// 	md, ok := metadata.FromIncomingContext(ctx)
// 	if !ok {
// 		return "", status.Errorf(codes.Unauthenticated, "metadata.FromIncomingContext: %v", grpc_errors.ErrNoCtxMetaData)
// 	}

// 	token := md.Get(key)

// 	if len(token) == 0 {
// 		return "", status.Errorf(codes.PermissionDenied, "md.Get access_token: %v", grpc_errors.ErrInvalidToken)
// 	}

// 	if token[0] == "" {
// 		return "", status.Errorf(codes.PermissionDenied, "md.Get access_token: %v", grpc_errors.ErrInvalidToken)
// 	}

// 	return token[0], nil
// }

func parseModelsToProto(user *user.ParamsOutputUser) *proto.User {
	return &proto.User{
		Id:        user.Id,
		Name:      user.Name,
		LastName:  user.Lastname,
		Password:  user.Password,
		Email:     user.Email,
		Role:      user.Role,
		Verified:  user.Verified,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}
