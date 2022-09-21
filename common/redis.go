// @Title  redis
// @Description  该文件用于初始化redis数据库，以及包装一个向外提供数据库的功能
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package common

import (
	"github.com/go-redis/redis/v9"
	"github.com/spf13/viper"
)

var Client *redis.Client

// @title    InitRedis
// @description   从配置文件中读取数据库相关信息后，完成数据库初始化
// @auth      MGAronya（张健）             2022-9-16 10:07
// @param     void        void         没有入参
// @return    *redis.Client         将返回一个初始化后的redis数据库指针
func InitRedis() *redis.Client {
	addr := viper.GetString("datasource.addredis")
	dbid := viper.GetInt("datasource.redisid")
	password := viper.GetString("datasource.redispass")
	host := viper.GetString("datasource.host")

	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + addr,
		Password: password,
		DB:       dbid,
	})
	Client = client
	return client
}

// @title    GetRedisClient
// @description   返回数据库的指针
// @auth      MGAronya（张健）             2022-9-16 10:08
// @param     void        void         没有入参
// @return    *redis.Client         将返回一个初始化后的数据库指针
func GetRedisClient() *redis.Client {
	return Client
}
