package event

import (
	"backend/config"
	"backend/models"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateEvent(c *gin.Context) {
	destinationIDStr := c.PostForm("destinationId")
	nameEvent := c.PostForm("nameevent")
	startDate := c.PostForm("start_date")
	endDate := c.PostForm("end_date")
	description := c.PostForm("description")
	startTime := c.PostForm("start_time")
	endTime := c.PostForm("end_time")
	priceStr := c.PostForm("price")
	maps := c.PostForm("maps")
	doField := c.PostForm("do")
	dontField := c.PostForm("dont")
	safety := c.PostForm("safety")
	location := c.PostForm("location")
	longitudeStr := c.PostForm("longitude")
	latitudeStr := c.PostForm("latitude")

	if nameEvent == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "nameevent wajib diisi",
		})
		return
	}

	var destIDPointer *uint = nil
	if destinationIDStr != "" {
		dID, convErr := strconv.Atoi(destinationIDStr)
		if convErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "destinationId harus angka",
			})
			return
		}

		var destCheck models.Destination
		if err := config.DB.First(&destCheck, "id_destination = ?", dID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "destinationId tidak ditemukan di database",
			})
			return
		}

		tmp := uint(dID)
		destIDPointer = &tmp
	}

	var price float64 = 0
	if priceStr != "" {
		price, _ = strconv.ParseFloat(priceStr, 64)
	}

	var longitude *float64
	if longitudeStr != "" {
		parsedLongitude, convErr := strconv.ParseFloat(longitudeStr, 64)
		if convErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "longitude harus angka",
			})
			return
		}
		longitude = &parsedLongitude
	}

	var latitude *float64
	if latitudeStr != "" {
		parsedLatitude, convErr := strconv.ParseFloat(latitudeStr, 64)
		if convErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "latitude harus angka",
			})
			return
		}
		latitude = &parsedLatitude
	}

	imagesBase64 := []string{}

	if form, formErr := c.MultipartForm(); formErr == nil && form != nil {
		if files, ok := form.File["image"]; ok {
			for idx, file := range files {
				openedFile, openErr := file.Open()
				if openErr != nil {
					c.JSON(http.StatusBadRequest, gin.H{"message": "Gagal membuka file ke-" + strconv.Itoa(idx+1)})
					return
				}

				imageBytes, readErr := io.ReadAll(openedFile)
				openedFile.Close()
				if readErr != nil {
					c.JSON(http.StatusBadRequest, gin.H{"message": "Gagal membaca file ke-" + strconv.Itoa(idx+1)})
					return
				}

				imagesBase64 = append(imagesBase64, base64.StdEncoding.EncodeToString(imageBytes))
			}
		}
	}

	var imageBase64 string
	if len(imagesBase64) > 0 {
		jsonBytes, marshalErr := json.Marshal(imagesBase64)
		if marshalErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal memproses image"})
			return
		}
		imageBase64 = string(jsonBytes)
	}

	event := models.Event{
		Name:        nameEvent,
		StartDate:   startDate,
		EndDate:     endDate,
		Description: description,
		StartTime:   startTime,
		EndTime:     endTime,
		Price:       price,
		Maps:        maps,
		Do:          doField,
		Dont:        dontField,
		Safety:      safety,
		Location:    location,
		Longitude:   longitude,
		Latitude:    latitude,
		ImageEvent:  imageBase64,
	}

	if destIDPointer != nil {
		event.DestinationID = destIDPointer
	}

	if err := config.DB.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal menambahkan event",
		})
		return
	}

	if event.DestinationID != nil {
		if err := config.DB.Preload("Destination").First(&event, event.ID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Gagal memuat data destinasi",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Event berhasil ditambahkan",
		"data":    event,
	})
}
