package api

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"testing"
	proto "user/api/qvbilam/user/v1"
)

func TestUserService_CreateUser(t *testing.T) {
	//mobile := "13501294174"
	//signInRequest := proto.UpdateRequest{
	//	Nickname: "4和9是一个用户",
	//	Mobile:   mobile,
	//	Password: "123456",
	//	Gender:   "male",
	//}

	users := []*proto.UpdateRequest{
		{
			AccountId: 1,
			Nickname:  "QvBiLam806607",
			Mobile:    "13501284164",
			Password:  "123456",
			Gender:    "male",
		},
		{
			AccountId: 2,
			Nickname:  "二滑大魔王",
			Mobile:    "13501284165",
			Password:  "123456",
			Gender:    "male",
		},
		{
			AccountId: 3,
			Nickname:  "星期二要上班",
			Mobile:    "13501284166",
			Password:  "123456",
			Gender:    "male",
		},
		{
			AccountId: 4,
			Nickname:  "进击的二滑",
			Mobile:    "13501284167",
			Password:  "123456",
			Gender:    "male",
		},
		{
			AccountId: 5,
			Nickname:  "4399小游戏",
			Mobile:    "13501284168",
			Password:  "123456",
			Gender:    "male",
		},
		{
			AccountId: 6,
			Nickname:  "牛魔王",
			Mobile:    "13501284169",
			Password:  "123456",
			Gender:    "male",
		},
		{
			AccountId: 7,
			Nickname:  "FATE ZERO",
			Mobile:    "13501284170",
			Password:  "123456",
			Gender:    "male",
		},
		{
			AccountId: 8,
			Nickname:  "fate-耳朵",
			Mobile:    "13501284171",
			Password:  "123456",
			Gender:    "male",
		},
		{
			AccountId: 9,
			Nickname:  "Fate养成",
			Mobile:    "13501284172",
			Password:  "123456",
			Gender:    "male",
		},
		{
			AccountId: 10,
			Nickname:  "4和9是一个用户",
			Mobile:    "13501284173",
			Password:  "123456",
			Gender:    "male",
		},
		{
			AccountId: 11,
			Nickname:  "糊涂的小神仙",
			Mobile:    "13501284174",
			Password:  "123456",
			Gender:    "male",
		},
		{
			AccountId: 12,
			Nickname:  "胡图图",
			Mobile:    "13501284175",
			Password:  "123456",
			Gender:    "male",
		},
	}

	// 创建用户
	for _, u := range users {
		c := client()
		//user, err := c.Create(context.Background(), &signInRequest)
		user, err := c.Create(context.Background(), u)
		if err != nil {
			t.Errorf("create user failed: %v", err)
		}
		// 获取用户
		getUser, err := c.Detail(context.Background(), &proto.GetUserRequest{Id: user.Id})
		if err != nil {
			t.Errorf("get user failed: %v", err)
		}
		if getUser.Nickname != u.Nickname {
			t.Errorf("user id error: %v", err)
		}
	}
}

//func TestUserService_CheckPassword(t *testing.T) {
//	r := proto.CheckPasswordRequest{
//		Password:         "123456",
//		EnctypedPassword: "$pbkdf2-sha512$KFO0YCkULF$ba00053b89c2369dbbf6884d7ae052b2dfe878e92f3384778071e8374b793e76",
//	}
//	c := client()
//	res, err := c.CheckPassword(context.Background(), &r)
//	if err != nil {
//		t.Errorf("create user failed: %v", err)
//	}
//	fmt.Printf("验证结果: %+v\n", res)
//}

func TestUserService_Delete(t *testing.T) {
	var userId int64
	userId = 2

	r := proto.UpdateRequest{
		Id: userId,
	}

	c := client()
	// 删除用户
	_, err := c.Delete(context.Background(), &r)
	if err != nil {
		t.Errorf("create user failed: %v", err)
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
