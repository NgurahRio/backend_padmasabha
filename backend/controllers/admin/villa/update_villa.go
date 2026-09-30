package villa

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

func UpdateVilla(c *gin.Context) {
	id := c.Param("id")

	var villa models.Villa

	if err := config.DB.First(&villa, "id_villa = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Villa tidak ditemukan"})
		return
	}

	changed := make(map[string]interface{})
	imageChanged := false

	if v := c.PostForm("subcategoryId"); v != "" {

		rawIDs := strings.Split(v, ",")
		for _, idStr := range rawIDs {
			idInt, _ := strconv.Atoi(strings.TrimSpace(idStr))

			var s models.Subcategory
			if err := config.DB.First(&s, "id_subcategories = ?", idInt).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": fmt.Sprintf("Subcategory ID %d tidak ditemukan", idInt),
				})
				return
			}
		}

		villa.SubcategoryID = v
		changed["subcategoryId"] = v
	}

	updateField := func(field *string, key string) {
		if v := c.PostForm(key); v != "" {
			*field = v
			changed[key] = v
		}
	}

	updateField(&villa.Name, "namevilla")
	updateField(&villa.Description, "description")
	updateField(&villa.Operational, "operational")

	if v := c.PostForm("facilityId"); v != "" {

		rawIDs := strings.Split(v, ",")
		for _, idStr := range rawIDs {
			idInt, _ := strconv.Atoi(strings.TrimSpace(idStr))

			var f models.Facility
			if err := config.DB.First(&f, "id_facility = ?", idInt).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": fmt.Sprintf("Facility ID %d tidak ditemukan", idInt),
				})
				return
			}
		}

		villa.FacilityID = v
		changed["facilityId"] = v
	}

	form, formErr := c.MultipartForm()
	if formErr == nil {

		var oldImages []string
		if err := json.Unmarshal([]byte(villa.Imagedata), &oldImages); err != nil {
			oldImages = []string{}
		}

		for key, files := range form.File {

			if strings.HasPrefix(key, "image[") && strings.HasSuffix(key, "]") {

				indexStr := key[6 : len(key)-1]
				index, err := strconv.Atoi(indexStr)
				if err != nil {
					continue
				}

				if index < len(oldImages) {

					file := files[0]
					openFile, _ := file.Open()
					rawBytes, _ := io.ReadAll(openFile)
					openFile.Close()

					newBase64 := base64.StdEncoding.EncodeToString(rawBytes)

					if newBase64 != oldImages[index] {
						oldImages[index] = newBase64
						changed[key] = "replaced"
						imageChanged = true
					}
				}
			}
		}

		if imageChanged {
			jsonBytes, _ := json.Marshal(oldImages)
			villa.Imagedata = string(jsonBytes)
			changed["images"] = "indexed image updated"
		}
	}

	if err := config.DB.Save(&villa).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengupdate villa",
		})
		return
	}

	if len(changed) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Tidak ada perubahan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil diperbarui",
		"changed": changed,
	})
}
