package routes

import (
	"context"
	"fmt"
	"net/http"
	"net/mail"
	"time"

	"github.com/SolProj3ct/Back-end/models"
	"github.com/SolProj3ct/Back-end/utils"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func SignIn(w http.ResponseWriter, r *http.Request) {
	thedb := r.Context().Value(utils.DB).(*mongo.Database)
	// get user payload
	u := models.User{}
	err := utils.UnmarshalObject(r.Body, &u)
	if err != nil {
		logrus.Errorln(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// check if email has correct format
	if _, err := mail.ParseAddress(u.Email); err != nil {
		logrus.Errorln(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// check if email is already used
	if count, err := thedb.Collection(u.DBCollectionName()).CountDocuments(r.Context(), bson.M{"email": u.Email}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if count > 0 {
		logrus.Errorln(fmt.Errorf("email %s already used", u.Email))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// hash password
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), utils.HashCost)
	if err != nil {
		logrus.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	u.ID = primitive.NewObjectID()
	u.Password = string(bytes)
	// persist in db
	if _, err := thedb.Collection(u.DBCollectionName()).InsertOne(r.Context(), u); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, u)
}

func LogIn(w http.ResponseWriter, r *http.Request) {
	thedb := r.Context().Value(utils.DB).(*mongo.Database)
	// get user payload
	u := models.User{}
	err := utils.UnmarshalObject(r.Body, &u)
	if err != nil {
		logrus.Errorln(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// check password
	userFromDb := models.User{}
	if err := thedb.Collection(u.DBCollectionName()).FindOne(context.Background(), bson.M{"_id": u.ID}).Decode(&userFromDb); err != nil {
		logrus.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(userFromDb.Password), []byte(u.Password)); err != nil {
		logrus.Errorln(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// prepare claims for token
	claims := models.JWTData{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(viper.GetDuration("auth.jwtDurationInHours") * time.Hour).Unix(),
		},
		CustomClaims: map[string]string{
			"uid": userFromDb.ID.Hex(),
		},
	}
	// generate token
	tokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenString.SignedString([]byte(viper.GetString("auth.jwtPrivateKey")))
	if err != nil {
		logrus.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, map[string]string{"jwt": token})
}
