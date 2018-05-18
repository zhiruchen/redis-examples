package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"

	"github.com/zhiruchen/redis-examples/db"
)

const (
	redisUserKey     = "user:%s"
	invalidParamCode = 10002
	dbErrCode        = 10003
	redisErrCode     = 10004
)

type user struct {
	ID   string `gorm:"column:id;type:char(20);primary_key" json:"id"`
	Name string `gorm:"column:name;type:varchar(255);not null" json:"name"`
}

func (user) TableName() string {
	return "user"
}

func createUserHandler(c *gin.Context) {
	u := user{}
	err := c.BindJSON(&u)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": invalidParamCode, "message": err.Error()})
		return
	}

	if u.Name == "" {
		c.JSON(http.StatusOK, gin.H{"code": invalidParamCode, "message": "invalid parameter Name"})
		return
	}

	u.ID = xid.New().String()
	if err := db.ORM.Create(&u).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": dbErrCode, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "id": u.ID})
}

func getUserHandler(c *gin.Context) {
	uid := c.Param("id")
	if uid == "" {
		c.JSON(http.StatusOK, gin.H{"code": invalidParamCode, "message": "user id is empty"})
		return
	}

	userkey := fmt.Sprintf(redisUserKey, uid)
	v, err := db.RedisClient.Exists(userkey).Result()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": redisErrCode, "message": err.Error()})
		return
	}

	var result map[string]interface{}
	if v == 1 {
		m := db.RedisClient.HGetAll(userkey).Val()
		result = make(map[string]interface{}, len(m))
		for k, v := range m {
			result[k] = v
		}
	} else {
		u, err := getUserFromDB(uid)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": dbErrCode, "message": err.Error()})
			return
		}

		result = structs.Map(u)
		if err := db.RedisClient.HMSet(userkey, result).Err(); err != nil {
			log.Println("cache user err: ", err)
		}
	}

	c.JSON(http.StatusOK, result)
}

func getUserFromDB(uid string) (*user, error) {
	u := user{}
	if err := db.ORM.Where("id=?", uid).First(&u).Error; err != nil {
		return nil, err
	}

	return &u, nil
}

func updateUserHandler(c *gin.Context) {
	uid := c.Param("id")
	if uid == "" {
		c.Error(errors.New("user id is empty"))
		return
	}

	u := user{}
	err := c.BindJSON(&u)
	if err != nil {
		c.Error(err)
		return
	}

	userkey := fmt.Sprintf(redisUserKey, uid)
	if err := db.RedisClient.Del(userkey).Err(); err != nil {
		c.Error(err)
		return
	}

	if err := updateUser(&user{ID: uid, Name: u.Name}); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK})
}

func updateUser(u *user) error {
	if err := db.ORM.Model(&u).Update("name", u.Name).Error; err != nil {
		return err
	}

	return nil
}

func getRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/users", createUserHandler)
	r.GET("/users/:id", getUserHandler)
	r.PUT("/users/:id", updateUserHandler)
	return r
}

func main() {
	if err := db.InitMysql(); err != nil {
		log.Fatalln(err)
	}

	if err := db.NewRedisClient(); err != nil {
		log.Fatalln(err)
	}

	getRouter().Run(":8989")
}
