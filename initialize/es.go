package initialize

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
	"user/global"
	"user/model"
)

func InitElasticSearch() {
	host := global.ServerConfig.ESConfig.Host
	port := global.ServerConfig.ESConfig.Port

	url := elastic.SetURL(fmt.Sprintf("http://%s:%d", host, port))
	sniff := elastic.SetSniff(false) // 不将本地地址转换
	var err error
	// 输出日志模式
	//logger := log.New(os.Stdout, "elasticsearch", log.LstdFlags) // 设置日志输出位置
	//global.ES, err = elastic.NewClient(url, sniff, elastic.SetTraceLog(logger))

	// 不输出日志
	global.ES, err = elastic.NewClient(url, sniff)
	if err != nil {
		zap.S().Panicf("连接es异常: %s", err.Error())
	}

	// 创建 video mapping
	createVideoIndex()
}

func createVideoIndex() {
	exists, err := global.ES.IndexExists(model.UserES{}.GetIndexName()).Do(context.Background())
	if err != nil {
		zap.S().Panicf("用户索引异常: %s", err)
	}
	if !exists { // 创建索引
		createIndex, err := global.ES.
			CreateIndex(model.UserES{}.GetIndexName()).
			BodyString(model.UserES{}.GetMapping()).
			Do(context.Background())
		if err != nil {
			zap.S().Panicf("创建用户索引异常: %s", err)
		}

		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}
}
