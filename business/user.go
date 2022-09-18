package business

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
	"user/global"
	"user/model"
	"user/utils"
)

const defaultUserAvatar = "https://blogupy.qvbilam.xin/bg/6666.JPG"
const defaultUserNicknamePrefixLogin = "QvBiLam"
const defaultNicknamePrefixLogoff = "注销用户"
const defaultNicknameSuffixLen = 6

type UserBusiness struct {
	Id         int64
	Code       int64
	Mobile     string
	Gender     string
	Nickname   string
	Password   string
	Avatar     string
	Ids        []int64
	Keyword    string
	Sort       string
	Level      int64
	IsVisible  *bool
	DeletedAt  *time.Time
	ModelQuery ModelQuery
	Page       int64
	PerPage    int64
}

type ModelQuery struct {
	Fields string
}

// ExistsMobile 验证手机号
func (b *UserBusiness) ExistsMobile() bool {
	fields := ModelQuery{Fields: "mobile"}
	entity := model.User{}
	res := global.DB.Where(&model.User{
		Mobile: b.Mobile,
	}).Select(fields).First(&entity)
	if res.RowsAffected == 0 {
		return false
	}
	return true
}

// Create 创建用户
func (b *UserBusiness) Create() (*model.User, error) {
	if b.ExistsMobile() {
		return nil, status.Errorf(codes.AlreadyExists, "手机号已注册")
	}

	ucB := UserCodeBusiness{}
	userCode, err := ucB.RandomCode(false)
	if err != nil {
		log.Printf("生成用户code失败: %v\n", err)
		return nil, status.Errorf(codes.Internal, "创建用户异常")
	}
	if userCode == 0 {
		return nil, status.Errorf(codes.Internal, "生成用户信息失败")
	}

	if b.Avatar == "" {
		b.Avatar = defaultUserAvatar
	}

	if b.Nickname == "" {
		b.Nickname = defaultUserNicknamePrefixLogin + utils.RandomNumber(defaultNicknameSuffixLen)
	}

	if b.Password != "" {
		b.Password = utils.GeneratePassword(b.Password)
	}

	entity := model.User{
		Code:     userCode,
		Mobile:   b.Mobile,
		Password: b.Password,
		Nickname: b.Nickname,
		Avatar:   b.Avatar,
		Gender:   b.Gender,
		Level:    0,
		Visible: model.Visible{
			IsVisible: true,
		},
	}

	if res := global.DB.Save(&entity); res.RowsAffected == 0 {
		zap.S().Errorf("创建用户失败: %s", res.Error)
		return nil, status.Errorf(codes.Internal, "创建失败")
	}

	return &entity, nil
}

func (b *UserBusiness) GetByMobile() (*model.User, error) {
	fields := b.SelectEntityFields()
	entity := model.User{}
	if res := global.DB.Unscoped().Select(fields).Where(model.User{
		Mobile: b.Mobile,
	}).First(&entity); res.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	if entity.DeletedAt != nil {
		return nil, status.Errorf(codes.NotFound, "用户已注销")
	}

	return &entity, nil
}

// GetById 获取用户 todo 批量获取用户等级
func (b *UserBusiness) GetById() (*model.User, error) {
	fields := b.SelectEntityFields()
	entity := model.User{}
	if res := global.DB.Unscoped().Select(fields).First(&entity, b.Id); res.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if entity.DeletedAt != nil {
		return nil, status.Errorf(codes.NotFound, "用户已注销")
	}

	return &entity, nil
}

// Search 搜索用户
func (b *UserBusiness) Search() (*[]model.User, int64) {

	switch {
	case b.PerPage <= 0:
		b.PerPage = 10
	case b.PerPage > 1000:
		b.PerPage = 1000
	}
	// 分页数据
	if b.Page == 0 {
		b.Page = 1
	}
	b.Page = (b.Page - 1) * b.PerPage

	// 获取 ES query
	q := b.GetUserESQuery()

	// 查询
	result, err := global.ES.
		Search().
		Index(model.UserES{}.GetIndexName()).
		Query(q).
		SortWithInfo(b.GetUserESSort()).
		From(int(b.Page)).
		Size(int(b.PerPage)).
		Do(context.Background())
	if err != nil {
		log.Printf("search user es error: %s\n", err)
		return nil, 0
	}

	fmt.Println("es查询结果:", result)

	// 获取总数
	total := result.Hits.TotalHits.Value

	// 获取视频 ids
	userIds := make([]int64, 0)
	for _, user := range result.Hits.Hits {
		userESModel := model.UserES{}
		_ = json.Unmarshal(user.Source, &userESModel)
		userIds = append(userIds, userESModel.ID)
	}

	b.Ids = userIds

	res, _ := b.GetByIds()
	return res, total
}

// GetByIds todo 批量获取用户等级
func (b *UserBusiness) GetByIds() (*[]model.User, int64) {
	fields := b.SelectEntityFields()
	var entity []model.User

	res := global.DB.Select(fields).Find(&entity, b.Ids)
	return &entity, res.RowsAffected
}

func (b *UserBusiness) Update() error {
	var user model.User
	if result := global.DB.Where(user.DeletedAt, nil).First(&user, b.Id); result.RowsAffected == 0 {
		return status.Errorf(codes.NotFound, "用户不存在")
	}

	b.EntityToUpdateModel(&user)

	result := global.DB.Save(&user)
	if result.RowsAffected == 0 {
		return status.Errorf(codes.NotFound, "修改用户信息失败")
	}
	return nil
}

func (b *UserBusiness) Delete() error {
	b.Avatar = defaultUserAvatar
	b.Nickname = defaultNicknamePrefixLogoff + utils.RandomNumber(defaultNicknameSuffixLen)
	visible := false
	b.IsVisible = &visible
	at := time.Now()
	b.DeletedAt = &at
	return b.Update()
}

func (b *UserBusiness) EntityToUpdateModel(user *model.User) {
	user.ID = b.Id
	if b.Code != 0 {
		user.Code = b.Code
	}
	if b.Mobile != "" {
		user.Mobile = b.Mobile
	}
	if b.Nickname != "" {
		user.Nickname = b.Nickname
	}
	if b.Password != "" {
		user.Password = utils.GeneratePassword(b.Password)
	}
	if b.Avatar != "" {
		user.Avatar = defaultUserAvatar
	}
	if b.Gender != "" {
		user.Gender = b.Gender
	}
	if b.Level > 0 {
		user.Level = b.Level
	}
	if b.IsVisible != nil {
		user.IsVisible = *b.IsVisible
	}
	if b.DeletedAt != nil {
		user.DeletedAt = b.DeletedAt
	}
}

func (b *UserBusiness) SelectEntityFields() string {
	if b.ModelQuery.Fields == "" {
		b.ModelQuery.Fields = "*"
	}
	return b.ModelQuery.Fields
}

func (b *UserBusiness) GetUserESQuery() *elastic.BoolQuery {
	// match bool 复合查询
	q := elastic.NewBoolQuery()

	if b.Keyword != "" { // 搜索 名称, 简介
		q = q.Must(elastic.NewMultiMatchQuery(b.Keyword, "code", "nickname"))
	}
	if b.Gender != "" { // 搜索用户
		q = q.Filter(elastic.NewTermQuery("gender", b.Gender))
	}

	//if b.IsVisible != nil { // 搜索展示状态
	//	q = q.Filter(elastic.NewTermQuery("is_visible", b.IsVisible))
	//}

	return q
}

func (b *UserBusiness) GetUserESSort() elastic.SortInfo {
	sort := elastic.SortInfo{
		Field:     "fans_count",
		Ascending: false,
	}

	if b.Sort != "" {
		if string(b.Sort[0]) == "-" {
			sort.Field = b.Sort[0:]
		} else {
			sort.Field = b.Sort
			sort.Ascending = true
		}
	}

	return sort
}
