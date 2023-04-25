package routes

import (
	"fmt"
	"net/http"
	"net/mail"

	"github.com/SolProj3ct/Back-end/models"
	"github.com/SolProj3ct/Back-end/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewUser(w http.ResponseWriter, r *http.Request) {
	thedb := r.Context().Value(utils.DB).(*mongo.Database)
	// get user payload
	newUser := models.User{}
	err := utils.UnmarshalObject(r.Body, &newUser)
	if err != nil {
		logrus.Errorln(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// check if email has correct format
	if _, err := mail.ParseAddress(newUser.Email); err != nil {
		logrus.Errorln(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// check if email is already used
	if count, err := thedb.Collection(newUser.DBCollectionName()).CountDocuments(r.Context(), bson.M{"email": newUser.Email}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if count > 0 {
		logrus.Errorln(fmt.Errorf("email %s already used", newUser.Email))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// hash password
	bytes, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), utils.HashCost)
	if err != nil {
		logrus.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newUser.ID = primitive.NewObjectID()
	newUser.Password = string(bytes)
	// chck auth token
	/*h := r.Header.Get("Authorization")
	if !strings.HasPrefix(h, "Bearer ") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	jwtString := strings.TrimPrefix(h, "Bearer ")
	switch viper.GetString("auth.jwtType") {
	case string(utils.JwtTypeMagic):
		// validate did token
		tk, err := magic.ValidateDidToken(jwtString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// ge user info from magic
		mgcUserInfo, err := magic.GetUserInfo(jwtString)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// create user db istance
		uid, err := primitive.ObjectIDFromHex(tk.GetClaim().Tid)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		u.ID = uid
		u.Name = newUser.Name
		u.Surname = newUser.Surname
		u.Email = mgcUserInfo.Email
		u.Bio = newUser.Bio
	default:
		logrus.Errorln("jwt type %s is not supported", viper.GetString("auth.jwtType"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}*/
	// persist in db
	if _, err := thedb.Collection(newUser.DBCollectionName()).InsertOne(r.Context(), newUser); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, newUser)
}
