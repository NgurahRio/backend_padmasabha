package destination

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func normalizeImages(data string) []string {
	if data == "" {
		return []string{}
	}

	var images []string

	if err := json.Unmarshal([]byte(data), &images); err == nil {
		return images
	}

	raw := strings.Split(data, ",")
	for _, item := range raw {
		item = strings.TrimSpace(item)
		if item != "" {
			images = append(images, item)
		}
	}

	return images
}

func GetAllDestinations(c *gin.Context) {

	var destinations []models.Destination
	if err := config.DB.Preload("Sos").Find(&destinations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengambil data destinasi",
		})
		return
	}

	var response []gin.H

	for _, d := range destinations {
		images := normalizeImages(d.Imagedata)

		var subcategories []models.Subcategory
		var subResp []gin.H

		if d.SubcategoryID != "" {
			rawSubs := strings.Split(d.SubcategoryID, ",")
			intIDs := []int{}
			for _, idStr := range rawSubs {
				if n, err := strconv.Atoi(strings.TrimSpace(idStr)); err == nil {
					intIDs = append(intIDs, n)
				}
			}

			config.DB.Preload("Category").Where("id_subcategories IN ?", intIDs).Find(&subcategories)

			for _, s := range subcategories {
				subResp = append(subResp, gin.H{
					"id_subcategories":  s.ID,
					"namesubcategories": s.Name,
					"categoriesId":      s.CategoryID,
					"category": gin.H{
						"id_categories": s.Category.ID,
						"name":          s.Category.Name,
					},
				})
			}
		}

		var facilities []models.Facility
		var facilityResp []gin.H

		if d.FacilityID != "" {
			rawIDs := strings.Split(d.FacilityID, ",")
			intIDs := []int{}
			for _, idStr := range rawIDs {
				if n, err := strconv.Atoi(strings.TrimSpace(idStr)); err == nil {
					intIDs = append(intIDs, n)
				}
			}

			config.DB.Where("id_facility IN ?", intIDs).Find(&facilities)

			for _, f := range facilities {
				facilityResp = append(facilityResp, gin.H{
					"id_facility":  f.IDFacility,
					"namefacility": f.NameFacility,
					"icon":         utils.ToBase64(f.Icon),
				})
			}
		}

		response = append(response, gin.H{
			"id_destination":  d.ID,
			"subcategoryId":   d.SubcategoryID,
			"subcategory":     subResp,
			"namedestination": d.Name,
			"location":        d.Location,
			"description":     d.Description,
			"images":          images,
			"do":              d.Do,
			"dont":            d.Dont,
			"safety":          d.Safety,
			"operational":     d.Operational,
			"maps":            d.Maps,
			"longitude":       d.Longitude,
			"latitude":        d.Latitude,
			"facilityId":      d.FacilityID,
			"facilities":      facilityResp,
			"sosId":           d.SosID,
			"sos": gin.H{
				"id_sos":     d.Sos.ID,
				"name_sos":   d.Sos.Name,
				"alamat_sos": d.Sos.Alamat,
				"telepon":    d.Sos.Telepon,
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil semua destinasi",
		"data":    response,
	})
}

func GetDestinationByID(c *gin.Context) {
	id := c.Param("id")

	var d models.Destination
	if err := config.DB.Preload("Sos").First(&d, "id_destination = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Destinasi tidak ditemukan",
		})
		return
	}

	images := normalizeImages(d.Imagedata)

	var subcategories []models.Subcategory
	var subResp []gin.H

	if d.SubcategoryID != "" {
		rawSubs := strings.Split(d.SubcategoryID, ",")
		intIDs := []int{}
		for _, idStr := range rawSubs {
			if n, err := strconv.Atoi(strings.TrimSpace(idStr)); err == nil {
				intIDs = append(intIDs, n)
			}
		}

		config.DB.Preload("Category").Where("id_subcategories IN ?", intIDs).Find(&subcategories)

		for _, s := range subcategories {
			subResp = append(subResp, gin.H{
				"id_subcategories":  s.ID,
				"namesubcategories": s.Name,
				"categoriesId":      s.CategoryID,
				"category": gin.H{
					"id_categories": s.Category.ID,
					"name":          s.Category.Name,
				},
			})
		}
	}

	var facilities []models.Facility
	var facilityResp []gin.H

	if d.FacilityID != "" {
		rawIDs := strings.Split(d.FacilityID, ",")
		intIDs := []int{}
		for _, idStr := range rawIDs {
			if n, err := strconv.Atoi(strings.TrimSpace(idStr)); err == nil {
				intIDs = append(intIDs, n)
			}
		}

		config.DB.Where("id_facility IN ?", intIDs).Find(&facilities)

		for _, f := range facilities {
			facilityResp = append(facilityResp, gin.H{
				"id_facility":  f.IDFacility,
				"namefacility": f.NameFacility,
				"icon":         utils.ToBase64(f.Icon),
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data destinasi",
		"data": gin.H{
			"id_destination":  d.ID,
			"subcategoryId":   d.SubcategoryID,
			"subcategory":     subResp,
			"namedestination": d.Name,
			"location":        d.Location,
			"description":     d.Description,
			"images":          images,
			"do":              d.Do,
			"dont":            d.Dont,
			"safety":          d.Safety,
			"operational":     d.Operational,
			"maps":            d.Maps,
			"longitude":       d.Longitude,
			"latitude":        d.Latitude,
			"facilityId":      d.FacilityID,
			"facilities":      facilityResp,
			"sosId":           d.SosID,
			"sos": gin.H{
				"id_sos":     d.Sos.ID,
				"name_sos":   d.Sos.Name,
				"alamat_sos": d.Sos.Alamat,
				"telepon":    d.Sos.Telepon,
			},
		},
	})
}
