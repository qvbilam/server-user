package api

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"testing"
	proto "user/api/qvbilam/user/v1"
)

func AccountClient() proto.AccountClient {
	host := "127.0.0.1"
	port := 9801
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", host, port),
		grpc.WithInsecure(),
	)
	if err != nil {
		panic(any(err))
	}
	client := proto.NewAccountClient(conn)
	return client
}

// 测试注册
func TestAccountService_Create(t *testing.T) {
	mobile := "13501294164"
	password := "123456"
	email := ""
	c := AccountClient()
	res, err := c.Create(context.Background(), &proto.UpdateAccountRequest{
		Mobile:   mobile,
		Email:    email,
		Password: password,
		Ip:       "",
	})
	if err != nil {
		fmt.Printf("server err: %s\n", err)
		return
	}
	fmt.Println(res)
}

// 测试登陆
func TestAccountService_LoginPassword(t *testing.T) {
	mobile := "13501294164"
	password := "123456"

	c := AccountClient()
	res, err := c.LoginPassword(context.Background(), &proto.LoginPasswordRequest{
		Method:   "mobile",
		Username: "",
		Mobile:   mobile,
		Email:    "",
		Password: password,
		Ip:       "",
	})
	if err != nil {
		fmt.Printf("server err: %s\n", err)
		return
	}

	fmt.Println(res)
}
