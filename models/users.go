package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID            primitive.ObjectID `json:"id" bson:"_id"`
	Name          string             `json:"name" bson:"name"`
	Surname       string             `json:"surname" bson:"surname"`
	Email         string             `json:"email" bson:"email"`
	PictureURL    *image             `json:"picture,omitempty" bson:"picture,omitempty"`
	Bio           string             `json:"bio" bson:"bio"`
	StarRating    float64            `json:"star_rating" bson:"star_rating"`
	NumberReviews float64            `json:"number_reviews" bson:"number_reviews"`
	Password      string             `json:"password" bson:"password"`
}

func (i User) DBCollectionName() string {
	return "users"
}
