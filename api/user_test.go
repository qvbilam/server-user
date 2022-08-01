package api

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"testing"
	proto "user/api/v1"
)

func TestUserService_CreateUser(t *testing.T) {
	mobile := "11000000000"
	signInRequest := proto.SignInRequest{
		Nickname: "test",
		Mobile:   mobile,
		Password: "123456",
		Gender:   "male",
	}

	client := client()
	// 创建用户
	user, err := client.CreateUser(context.Background(), &signInRequest)
	if err != nil {
		t.Errorf("create user failed: %v", err)
	}
	// 获取用户
	getUser, err := client.GetUser(context.Background(), &proto.GetUserRequest{Id: user.Id})
	if err != nil {
		t.Errorf("get user failed: %v", err)
	}
	if getUser.Id != user.Id {
		t.Errorf("user id error: %v", err)
	}
}

func client() proto.UserClient {
	host := "127.0.0.1"
	port := 9801
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", host, port),
		grpc.WithInsecure(),
	)
	if err != nil {
		panic(any(err))
	}
	client := proto.NewUserClient(conn)
	return client
}
