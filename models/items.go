package models

import (
	"encoding/json"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ItemCategory string

const (
	ItemCategoryClothe ItemCategory = "clothe"
	ItemCategoryHitech ItemCategory = "hitech"
)

type image struct {
	ImgURL string `json:"img_url" bson:"img_url"`
}

type ItemMeta map[string]interface{}

type BaseItem struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Category    ItemCategory       `json:"cat" bson:"cat"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"desc" bson:"desc"`
	Brand       string             `json:"brand" bson:"brand"`
	//Vendors      []Vendor           `json:"-" bson:"vendors"`
	Metas        *ItemMeta `json:"metas" bson:"metas"`
	PrimaryImage *image    `json:"primary_img,omitempty" bson:"primary_img,omitempty"`
	ImageList    []image   `json:"imgs,omitempty" bson:"imgs,omitempty"`
	Price        float64   `json:"price" bson:"price"`
	Likes        float64   `json:"likes" bson:"likes"`
}

type ItemClothe struct {
	Size  string `json:"size" bson:"size"`
	Color string `json:"color" bson:"color"`
}

type ItemHiTech struct {
	Battery float32 `json:"battery" bson:"battery"`
}

func (i BaseItem) DBCollectionName() string {
	return "items"
}

func (i BaseItem) IndexAlgolia() string {
	return "items"
}

func (u *BaseItem) MarshalJSON() ([]byte, error) {
	type Alias BaseItem
	var aux interface{}
	switch u.Category {
	case ItemCategoryClothe:
		var c *ItemClothe
		mapstructure.Decode(u.Metas, &c)
		aux = &struct {
			*Alias
			Metas *ItemClothe `json:"metas" bson:"metas"`
		}{
			Alias: (*Alias)(u),
			Metas: c,
		}
	case ItemCategoryHitech:
		var ht *ItemHiTech
		mapstructure.Decode(u.Metas, &ht)
		aux = &struct {
			*Alias
			Metas *ItemHiTech `json:"metas" bson:"metas"`
		}{
			Alias: (*Alias)(u),
			Metas: ht,
		}
	default:
		return nil, fmt.Errorf("item type %v not supported", u.Category)
	}
	return json.Marshal(aux)
}

func (u *BaseItem) UnmarshalJSON(data []byte) error {
	type Alias BaseItem
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	switch aux.Category {
	case ItemCategoryClothe:
		// decoder
		var c ItemClothe
		msd, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			ErrorUnused: true,
			ErrorUnset:  true,
			Result:      &c,
		})
		if err != nil {
			return err
		}
		if err := msd.Decode(u.Metas); err != nil {
			return err
		}
	case ItemCategoryHitech:
		// decoder
		var ht ItemHiTech
		msd, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			ErrorUnused: true,
			ErrorUnset:  true,
			Result:      &ht,
		})
		if err != nil {
			return err
		}
		if err := msd.Decode(u.Metas); err != nil {
			return err
		}
	default:
		return fmt.Errorf("item type %v not supported", aux.Category)
	}
	return nil
}

/*func (item BaseItem) ValidateItem(thedb *mongo.Database) error {
	// check if all vendors exist in db
	if len(item.Vendors) > 0 {
		var idsVendors = make([]primitive.ObjectID, len(item.Vendors))
		for i := range item.Vendors {
			idsVendors[i] = item.Vendors[i].ID
		}
		// check if already exists analysis with this type on this ticket
		count, err := thedb.Collection(Vendor{}.DBCollectionName()).CountDocuments(context.Background(), bson.M{"_id": bson.M{"$in": idsVendors}})
		if err != nil {
			return err
		}
		if count != int64(len(item.Vendors)) {
			return errors.New("some vendors not exist")
		}
	}
	return nil
}*/
