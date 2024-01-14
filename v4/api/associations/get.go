//nolint:gomnd,gosmopolitan
package associations

import (
	"net/http"
	"time"

	"bsam-server/v4/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/maps"
)

//nolint:gochecknoglobals
var fetchSample = map[string]models.Association{
	"abcdefgh-1234-5678-90ab-cdefghijklmn": {
		ID:       "sailing-ise",
		Name:     "セーリング伊勢",
		TokenIAT: time.Now(),
		TokenEXP: time.Now().Add(time.Hour * 24),
		Lat:      35.353535,
		Lng:      120.120120,
		RaceName: "全国ハンザクラスブラインドセーリング大会",
	},
	"bacdefgh-1234-5678-90ab-cdefghijklmn": {
		ID:       "hogehoge",
		Name:     "ホゲホゲマリンビレッジ",
		TokenIAT: time.Now(),
		TokenEXP: time.Now().Add(time.Hour * 24),
		Lat:      35.353535,
		Lng:      120.120120,
		RaceName: "ホゲセーリング2023",
	},
	"nmlkjihgfedc-ba09-8765-4321-hgfedcba": {
		ID:       "piyopiyo",
		Name:     "ピヨピヨマリンビレッジ",
		TokenIAT: time.Now(),
		TokenEXP: time.Now().Add(time.Hour * 24),
		Lat:      38.383838,
		Lng:      130.130130,
		RaceName: "ピヨセーリング2023",
	},
}

func AssociationGETAll(c *gin.Context) {
	assocs := getAssociations()

	res := models.AssociationsGETAllRes{
		Assocs: maps.Values(assocs),
	}

	c.JSON(http.StatusOK, res)
}

func AssociationGET(c *gin.Context) {
	id := c.Param("id")

	docID, exist := findAssociationID(id)

	if !exist {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "This association is not found",
		})

		return
	}

	c.JSON(http.StatusOK, getAssociation(docID))
}

func getAssociations() map[string]models.Association {
	return fetchSample
}

func getAssociation(docID string) models.Association {
	return fetchSample[docID]
}

func findAssociationID(id string) (string, bool) {
	assocs := getAssociations()

	for docID, assoc := range assocs {
		if assoc.ID == id {
			return docID, true
		}
	}

	return "", false
}
