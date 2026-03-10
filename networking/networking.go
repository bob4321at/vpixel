package networking

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

var InitialNetworkStarted bool
var JoinServerAddress string = "localhost:8181"

func StartServer() {
	r := gin.Default()

	r.GET("/AddUser", AddUser)
	r.POST("/SetUsersModel", SetUsersModel)
	r.POST("/SetUserFaceTrackingData", SetUserFaceTrackingData)
	r.POST("/GetOtherUsersModels", GetOtherUsersModels)

	go func() {
		InitialNetworkStarted = true
		UploadThisUser()
	}()

	go func() {
		for true {
			UserUpdated.Range(func(key, value any) bool {
				fmt.Println()
				fmt.Println(key.(int), value.(bool))
				fmt.Println()
				return true
			})
			time.Sleep(time.Millisecond * 10)
		}
	}()

	r.Run(":8181")
}

func JoinServer() {
	go func() {
		InitialNetworkStarted = true
		UploadThisUser()
	}()
}
