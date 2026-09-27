package packagess

import (
	"backend/config"
	"backend/models"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func UpdatePackages(c *gin.Context) {
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
	if err := config.DB.First(&pkg, "destinationId = ?", destID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Package untuk destinationId ini belum ada",
		})
		return
	}

	// Parse subPackageId dari string yang dipisahkan koma (contoh: "1,2,3,4")
	subIDsStr := c.PostForm("subPackageId")
	if subIDsStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "subPackageId wajib diisi"})
		return
	}

	// Split string menjadi array
	subIDsRaw := strings.Split(subIDsStr, ",")
	subIDs := []string{}
	for _, id := range subIDsRaw {
		trimmed := strings.TrimSpace(id)
		if trimmed != "" {
			subIDs = append(subIDs, trimmed)
		}
	}

	if len(subIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "subPackageId tidak boleh kosong"})
		return
	}

	// Parse price dari string yang dipisahkan koma (contoh: "100000,200000,300000")
	pricesStr := c.PostForm("price")
	if pricesStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "price wajib diisi"})
		return
	}

	pricesRaw := strings.Split(pricesStr, ",")
	prices := []string{}
	for _, p := range pricesRaw {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			prices = append(prices, trimmed)
		}
	}

	// Parse includeName dari string yang dipisahkan koma (opsional)
	includeNamesStr := c.PostForm("includeName")
	includeNames := []string{}
	if includeNamesStr != "" {
		includeNamesRaw := strings.Split(includeNamesStr, ",")
		for _, name := range includeNamesRaw {
			trimmed := strings.TrimSpace(name)
			includeNames = append(includeNames, trimmed)
		}
	}

	form, formErr := c.MultipartForm()
	if formErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Form tidak valid"})
		return
	}

	files := form.File["image"]

	fmt.Println("DEBUG SUB IDS  :", subIDs)
	fmt.Println("DEBUG PRICES  :", prices)
	fmt.Println("DEBUG IMAGES  :", len(files))

	if len(prices) != len(subIDs) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Jumlah price harus sama dengan subPackageId",
			"detail": gin.H{
				"subPackageId": len(subIDs),
				"price":        len(prices),
			},
		})
		return
	}

	if len(files) != len(subIDs) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Jumlah image harus sama dengan subPackageId",
			"detail": gin.H{
				"subPackageId": len(subIDs),
				"image":        len(files),
			},
		})
		return
	}

	// Grouping data berdasarkan subPackageId
	priceMap := map[string]float64{}
	includesMap := map[string][]subPackageInclude{}

	for i := range subIDs {
		subID := subIDs[i]

		if subID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "subPackageId di index " + strconv.Itoa(i) + " kosong",
			})
			return
		}

		subIDInt, errConv := strconv.Atoi(subID)
		if errConv != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "subPackageId harus angka",
				"value":   subID,
			})
			return
		}

		var subCheck models.SubPackage
		if err := config.DB.First(&subCheck, "id_subpackage = ?", subIDInt).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "subPackageId tidak ditemukan di tabel subpackage",
				"value":   subIDInt,
			})
			return
		}

		priceVal, priceErr := strconv.ParseFloat(prices[i], 64)
		if priceErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Harga tidak valid pada index " + strconv.Itoa(i),
				"value":   prices[i],
			})
			return
		}

		// Set price (ambil yang pertama untuk setiap subPackageId)
		if _, exists := priceMap[subID]; !exists {
			priceMap[subID] = priceVal
		}
		
		file := files[i]
		opened, openErr := file.Open()
		if openErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal membuka file image"})
			return
		}

		imgBytes, readErr := io.ReadAll(opened)
		opened.Close()
		if readErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal membaca file image"})
			return
		}

		base64Img := base64.StdEncoding.EncodeToString(imgBytes)

		nameVal := subCheck.Packagetype
		if len(includeNames) > i && includeNames[i] != "" {
			nameVal = includeNames[i]
		}

		// Append include ke array
		includesMap[subID] = append(includesMap[subID], subPackageInclude{
			Image: base64Img,
			Name:  nameVal,
		})
	}

	// Build updatedData dari map yang sudah di-group
	updatedData := map[string]subPackageDetail{}
	for subID, includes := range includesMap {
		updatedData[subID] = subPackageDetail{
			Price:   priceMap[subID],
			Include: includes,
		}
	}

	jsonBytes, _ := json.MarshalIndent(updatedData, "", "  ")
	pkg.SubPackageData = string(jsonBytes)

	if err := config.DB.Save(&pkg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengupdate package",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Package berhasil diupdate",
		"data": gin.H{
			"id_packages":     pkg.ID,
			"destinationId":   pkg.DestinationID,
			"subpackage_data": updatedData,
		},
	})
}
