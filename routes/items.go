package routes

import (
	"net/http"

	"github.com/SolProj3ct/Back-end/integrations/algolia"
	"github.com/SolProj3ct/Back-end/models"
	"github.com/SolProj3ct/Back-end/utils"
	"github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewItem(w http.ResponseWriter, r *http.Request) {
	thedb := r.Context().Value(utils.DB).(*mongo.Database)
	newItem := models.BaseItem{}
	err := utils.UnmarshalObject(r.Body, &newItem)
	if err != nil {
		logrus.Errorln(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	newItem.ID = primitive.NewObjectID()
	// validate item
	/*if err := newItem.ValidateItem(thedb); err != nil {
		logrus.Errorln(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}*/
	// persist
	if _, err := thedb.Collection(newItem.DBCollectionName()).InsertOne(r.Context(), newItem); err != nil {
		logrus.Errorln("error inserting item in db", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// insert object in algolia
	if err := algolia.InsertObject(newItem.IndexAlgolia(), newItem); err != nil {
		logrus.Errorln("failed to save item in algolia", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, newItem)
}

func GetPopularItems(w http.ResponseWriter, r *http.Request) {
	thedb := r.Context().Value(utils.DB).(*mongo.Database)
	baseItems := []models.BaseItem{}
	cur, err := thedb.Collection(models.BaseItem{}.DBCollectionName()).Find(r.Context(), bson.M{})
	if err != nil {
		logrus.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := cur.All(r.Context(), &baseItems); err != nil {
		logrus.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, baseItems)
}

func GetRecommendedItems(w http.ResponseWriter, r *http.Request) {
	thedb := r.Context().Value(utils.DB).(*mongo.Database)
	baseItems := []models.BaseItem{}
	cur, err := thedb.Collection(models.BaseItem{}.DBCollectionName()).Find(r.Context(), bson.M{})
	if err != nil {
		logrus.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := cur.All(r.Context(), &baseItems); err != nil {
		logrus.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, baseItems)
}

func GetItemBrands(w http.ResponseWriter, r *http.Request) {
	thedb := r.Context().Value(utils.DB).(*mongo.Database)
	brands, err := thedb.Collection(models.BaseItem{}.DBCollectionName()).Distinct(r.Context(), "brand", bson.M{})
	if err != nil {
		logrus.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, brands)
}
