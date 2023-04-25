package routes

import (
	"net/http"
	"strings"

	"github.com/SolProj3ct/Back-end/integrations/algolia"
	"github.com/SolProj3ct/Back-end/models"
	"github.com/SolProj3ct/Back-end/utils"
	"github.com/sirupsen/logrus"
)

func GlobalSearch(w http.ResponseWriter, r *http.Request) {
	//query params
	q := r.URL.Query().Get("q")
	if len(strings.TrimSpace(q)) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// search on items
	results, err := algolia.SearchObject(models.BaseItem{}.IndexAlgolia(), q)
	if err != nil {
		logrus.Errorln("algolia search error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, results)
}
