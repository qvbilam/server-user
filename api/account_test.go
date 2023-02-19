package api

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
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

func TestIP(t *testing.T) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, value := range addrs {
		if ipnet, ok := value.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println(ipnet.IP.String())
			}
		}
	}
}

// 测试注册
func TestAccountService_Create(t *testing.T) {
	mobile := "13501294111"
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
	mobile := "13501294174"
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
