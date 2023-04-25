package utils

import (
	"context"
	"net/http"
	"strings"

	magic "github.com/SolProj3ct/Back-end/integrations/magic"
	"github.com/SolProj3ct/Back-end/models"
	"github.com/golang-jwt/jwt"
	"github.com/magiclabs/magic-admin-go/token"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type JwtType string

// singleton keyfunc
var keyFunc func(token *jwt.Token) (interface{}, error) = func(token *jwt.Token) (interface{}, error) {
	return []byte(viper.GetString("auth.jwtPrivateKey")), nil
}

const (
	JwtTypeMagic  JwtType = "magic"
	JwtTypeSimple JwtType = "simple"
)

func InjectUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		jwtString := strings.TrimPrefix(h, "Bearer ")
		var u *models.User
		switch viper.GetString("auth.jwtType") {
		case string(JwtTypeMagic):
			tk, err := magic.ValidateDidToken(jwtString)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			// pass db in order to use users collection
			u, err = getUserInfoTokenMagic(tk.GetClaim(), r.Context().Value(DB).(*mongo.Database).Collection(models.User{}.DBCollectionName()))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		case string(JwtTypeSimple):
			claims := &models.JWTData{}
			_, err := jwt.ParseWithClaims(jwtString, claims, keyFunc)
			if err != nil {
				logrus.Errorln(err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			uid, ok := claims.CustomClaims["uid"]
			if !ok {
				logrus.Errorln("user id not present")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			uidObjectID, err := primitive.ObjectIDFromHex(uid)
			if err != nil {
				logrus.Errorln(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// add user
			if err := r.Context().Value(DB).(*mongo.Database).Collection(models.User{}.DBCollectionName()).FindOne(context.Background(), bson.M{"_id": uidObjectID}).Decode(&u); err != nil {
				logrus.Errorln(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		default:
			logrus.Errorln("jwt type %s is not supported", viper.GetString("auth.jwtType"))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if u == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), User, u))
		next.ServeHTTP(w, r)
	})
}

func getUserInfoTokenMagic(tk token.Claim, userDbCollection *mongo.Collection) (*models.User, error) {
	u := models.User{}
	uidObjectID, err := primitive.ObjectIDFromHex(tk.Tid)
	if err != nil {
		return nil, err
	}
	if err := userDbCollection.FindOne(context.Background(), bson.M{"_id": uidObjectID}).Decode(&u); err != nil {
		return nil, err
	}
	return &u, nil
}
