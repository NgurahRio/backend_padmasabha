package restaurant

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

func GetAllRestaurantsUser(c *gin.Context) {
	var restaurants []models.Restaurant
	if err := config.DB.Find(&restaurants).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengambil data restaurant",
		})
		return
	}

	var response []gin.H

	for _, d := range restaurants {
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
			"id_restaurant":  d.ID,
			"subcategoryId":  d.SubcategoryID,
			"subcategory":    subResp,
			"namerestaurant": d.Name,
			"description":    d.Description,
			"images":         images,
			"operational":    d.Operational,
			"facilityId":     d.FacilityID,
			"facilities":     facilityResp,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil semua restaurant",
		"data":    response,
	})
}

func GetRestaurantByIDUser(c *gin.Context) {
	id := c.Param("id")

	var d models.Restaurant
	if err := config.DB.First(&d, "id_restaurant = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Restaurant tidak ditemukan",
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
		"message": "Berhasil mengambil data restaurant",
		"data": gin.H{
			"id_restaurant":  d.ID,
			"subcategoryId":  d.SubcategoryID,
			"subcategory":    subResp,
			"namerestaurant": d.Name,
			"description":    d.Description,
			"images":         images,
			"operational":    d.Operational,
			"facilityId":     d.FacilityID,
			"facilities":     facilityResp,
		},
	})
}
