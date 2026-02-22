package networking

import (
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

	r.Run(":8181")
}

func JoinServer() {
	go func() {
		InitialNetworkStarted = true
		UploadThisUser()
	}()
}
