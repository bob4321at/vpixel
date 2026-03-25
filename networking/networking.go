package networking

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var InitialNetworkStarted bool
var JoinServerAddress string = "localhost:8181"

func StartServer() {
	r := gin.Default()

	r.GET("/AddUser", AddUser)
	r.GET("/GetUserFaceTrackingData", GetUserFaceTrackingData)
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

		for true {
			time.Sleep(time.Millisecond * 10)

			resp, err := http.Get("http://" + JoinServerAddress + "/GetUserFaceTrackingData")
			if err != nil {
				panic(err)
			}
			resp_bytes, err := io.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}

			face_tracking_data_temp := CollectionFaceTrackingDataNetworked{}
			json.Unmarshal(resp_bytes, &face_tracking_data_temp)

			for _, data := range face_tracking_data_temp.FaceTrackingDatas {
				UsersFaceTrackingData.Store(data.ID, data)
			}
		}
	}()
}
