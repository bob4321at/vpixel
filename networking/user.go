package networking

import (
	"bytes"
	"encoding/json"
	"io"
	"main/models"
	"main/tracking"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type UserJson struct {
	ID uint8
}

type FaceTrackingNetworked struct {
	ID               uint8
	FaceTrackingData map[string]float64
}

var Users []UserJson
var UsersModel sync.Map

var UsersFaceTrackingData sync.Map

var ThisUsersId uint8

var UserUpdated sync.Map

func UploadThisUser() {
	time.Sleep(time.Second)

	resp, err := http.Get("http://" + JoinServerAddress + "/AddUser")
	if err != nil {
		panic(err)
	}

	users_id_in_user_bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	user_id_in_user := UserJson{}
	json.Unmarshal(users_id_in_user_bytes, &user_id_in_user)

	ThisUsersId = user_id_in_user.ID

	go func() {
		for true {
			time.Sleep(time.Second / 10)
			networked_face_tracking_data := FaceTrackingNetworked{ThisUsersId, map[string]float64{}}

			for key, data := range tracking.WeightOptions {
				networked_face_tracking_data.FaceTrackingData[key] = data
			}

			networked_face_tracking_data_bytes, err := json.Marshal(networked_face_tracking_data)
			if err != nil {
				panic(err)
			}

			if resp, err = http.Post("http://"+JoinServerAddress+"/SetUserFaceTrackingData", "application/json", bytes.NewBuffer(networked_face_tracking_data_bytes)); err != nil {
				panic(err)
			}

			resp_bytes, err := io.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}

			IsUserUpdated := UserJson{}
			if err := json.Unmarshal(resp_bytes, &IsUserUpdated); err != nil {
				panic(err)
			}

			if IsUserUpdated.ID == 0 {
				user_to_send := UserJson{ID: ThisUsersId}
				user_to_send_bytes, err := json.Marshal(user_to_send)
				if err != nil {
					panic(err)
				}
				resp, err := http.Post("http://"+JoinServerAddress+"/GetOtherUsersModels", "application/json", bytes.NewBuffer(user_to_send_bytes))
				if err != nil {
					panic(err)
				}

				resp_bytes, err := io.ReadAll(resp.Body)
				if err != nil {
					panic(err)
				}

				other_models := OtherUsersModels{}
				if err := json.Unmarshal(resp_bytes, &other_models); err != nil {
					panic(err)
				}

				for _, model := range other_models.OtherModels {
					new_model := TurnNetworkedModelToModel(model)
					UsersModel.Store(model.UserID, new_model)
				}
			}
		}
	}()
}

type OtherUsersModels struct {
	OtherModels []NetworkedModelJson
}

func GetOtherUsersModels(c *gin.Context) {
	user_ID, err := io.ReadAll(c.Request.Body)
	if err != nil {
		panic(err)
	}

	GetUserIDJson := UserJson{}
	if err := json.Unmarshal(user_ID, &GetUserIDJson); err != nil {
		panic(err)
	}

	models_to_send := OtherUsersModels{}

	UsersModel.Range(func(key, value any) bool {
		if key == GetUserIDJson.ID {
			return true
		}

		model := value.(models.Model)
		new_model := TurnModelToNetworkedModel(&model)

		models_to_send.OtherModels = append(models_to_send.OtherModels, new_model)
		return true
	})

	UserUpdated.Store(GetUserIDJson.ID, true)

	c.JSON(http.StatusOK, models_to_send)
}

func SetUserFaceTrackingData(c *gin.Context) {
	face_tracking_json_bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		panic(err)
	}

	var face_tracking_json FaceTrackingNetworked

	if err = json.Unmarshal(face_tracking_json_bytes, &face_tracking_json); err != nil {
		panic(err)
	}

	UsersFaceTrackingData.Store(face_tracking_json.ID, face_tracking_json)

	is_updated, ok := UserUpdated.Load(face_tracking_json.ID)
	if ok {
		updated := is_updated.(bool)
		if updated {
			c.JSON(http.StatusOK, UserJson{ID: 1})
			return
		}
	}

	if len(Users) == 1 {
		c.JSON(http.StatusOK, UserJson{ID: 1})
	} else {
		c.JSON(http.StatusOK, UserJson{ID: 0})
	}
}

func SetUsersModel(c *gin.Context) {
	user_model_json_bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		panic(err)
	}

	user_model_json := NetworkedModelJson{}
	json.Unmarshal(user_model_json_bytes, &user_model_json)

	new_model := TurnNetworkedModelToModel(user_model_json)

	UsersModel.Store(user_model_json.UserID, new_model)

	UsersModel.Range(func(key, value any) bool {
		UserUpdated.Store(key, false)
		return true
	})
}

func AddUser(c *gin.Context) {
	Usr := UserJson{}
	Usr.ID = uint8(len(Users))

	Users = append(Users, Usr)

	c.JSON(http.StatusOK, Usr)
}
