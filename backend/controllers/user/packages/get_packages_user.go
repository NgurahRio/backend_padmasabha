package packages

import (
	"backend/config"
	"backend/models"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type subPackageInclude struct {
	Image string `json:"image"`
	Name  string `json:"name"`
}

type subPackageDetail struct {
	Price   float64             `json:"price"`
	Include []subPackageInclude `json:"include"`
}

func normalizeSubPackageData(data string) map[string]subPackageDetail {
	result := map[string]subPackageDetail{}
	if data == "" {
		return result
	}

	raw := map[string]interface{}{}
	if err := json.Unmarshal([]byte(data), &raw); err != nil {
		return result
	}

	for key, val := range raw {
		valMap, ok := val.(map[string]interface{})
		if !ok {
			continue
		}

		detail := subPackageDetail{}

		if price, ok := valMap["price"].(float64); ok {
			detail.Price = price
		}

		if includeArr, ok := valMap["include"].([]interface{}); ok {
			for _, item := range includeArr {
				if includeRaw, ok := item.(map[string]interface{}); ok {
					entry := subPackageInclude{}
					if img, ok := includeRaw["image"].(string); ok {
						entry.Image = img
					}
					if name, ok := includeRaw["name"].(string); ok {
						entry.Name = name
					}
					detail.Include = append(detail.Include, entry)
				}
			}
		} else if includeRaw, ok := valMap["include"].(map[string]interface{}); ok {
			entry := subPackageInclude{}
			if img, ok := includeRaw["image"].(string); ok {
				entry.Image = img
			}
			if name, ok := includeRaw["name"].(string); ok {
				entry.Name = name
			}
			detail.Include = append(detail.Include, entry)
		}

		if len(detail.Include) == 0 {
			entry := subPackageInclude{}
			if img, ok := valMap["image"].(string); ok {
				entry.Image = img
			}
			if name, ok := valMap["packagetype"].(string); ok {
				entry.Name = name
			}
			if entry.Image != "" || entry.Name != "" {
				detail.Include = append(detail.Include, entry)
			}
		}

		result[key] = detail
	}

	return result
}

func GetAllPackagesUser(c *gin.Context) {
	var pkgs []models.Packages

	if err := config.DB.Preload("Destination").Find(&pkgs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengambil semua packages",
		})
		return
	}

	result := []gin.H{}

	for _, pkg := range pkgs {
		data := normalizeSubPackageData(pkg.SubPackageData)

		subIDs := []int{}
		for k := range data {
			if n, err := strconv.Atoi(strings.TrimSpace(k)); err == nil {
				subIDs = append(subIDs, n)
			}
		}

		var subpackages []models.SubPackage
		if len(subIDs) > 0 {
			if err := config.DB.Where("id_subpackage IN ?", subIDs).Find(&subpackages).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Gagal memuat data subpackage",
				})
				return
			}
		}

		subResp := []gin.H{}
		for _, sp := range subpackages {
			subResp = append(subResp, gin.H{
				"id_subpackage": sp.ID,
				"jenispackage":  sp.Packagetype,
				"image":         sp.Image,
			})
		}

		result = append(result, gin.H{
			"id_packages":     pkg.ID,
			"destinationId":   pkg.DestinationID,
			"destination": gin.H{
				"id_destination":  pkg.Destination.ID,
				"namedestination": pkg.Destination.Name,
				"location":        pkg.Destination.Location,
				"description":     pkg.Destination.Description,
			},
			"subpackage_data": data,
			"subpackages":     subResp,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil semua packages",
		"data":    result,
	})
}

func GetPackageByDestinationIDUser(c *gin.Context) {
	destIDStr := c.Param("destinationId")
	if destIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "destinationId wajib diisi"})
		return
	}

	destID, convErr := strconv.Atoi(destIDStr)
	if convErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "destinationId harus berupa angka"})
		return
	}

	var dest models.Destination
	if err := config.DB.First(&dest, "id_destination = ?", destID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "destinationId tidak ditemukan di database"})
		return
	}

	var pkg models.Packages
	if err := config.DB.Preload("Destination").First(&pkg, "destinationId = ?", destID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Package untuk destinationId ini belum ada",
		})
		return
	}

	subData := normalizeSubPackageData(pkg.SubPackageData)

	subIDs := []int{}
	for k := range subData {
		if n, err := strconv.Atoi(strings.TrimSpace(k)); err == nil {
			subIDs = append(subIDs, n)
		}
	}

	var subpackages []models.SubPackage
	if len(subIDs) > 0 {
		if err := config.DB.Where("id_subpackage IN ?", subIDs).Find(&subpackages).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Gagal memuat data subpackage",
			})
			return
		}
	}

	subResp := []gin.H{}
	for _, sp := range subpackages {
		subResp = append(subResp, gin.H{
			"id_subpackage": sp.ID,
			"jenispackage":  sp.Packagetype,
			"image":         sp.Image,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data package",
		"data": gin.H{
			"id_packages":     pkg.ID,
			"destinationId":   pkg.DestinationID,
			"destination": gin.H{
				"id_destination":  pkg.Destination.ID,
				"namedestination": pkg.Destination.Name,
				"location":        pkg.Destination.Location,
				"description":     pkg.Destination.Description,
			},
			"subpackage_data": subData,
			"subpackages":     subResp,
		},
	})
}
