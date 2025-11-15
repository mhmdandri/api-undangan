package controller

import (
	"api-undangan/config"
	"api-undangan/database"
	"api-undangan/email"
	"api-undangan/models"
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ReservationRequest struct {
	Name			string `json:"name" binding:"required"`
	IsPresent		*bool   `json:"is_present" binding:"required"`
	Email			string `json:"email" binding:"required,email"`
    TotalGuests    int    `json:"total_guests" binding:"omitempty,min=1,max=3"`
}

type ConfirmReservationRequest struct {
    Code  string `json:"code" binding:"required"`
}

func GetReservations(c *gin.Context){
	var reservations []models.Reservation
	if err := database.DB.Order("created_at desc").Find(&reservations).Error; err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch reservations",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": reservations,
	})
}

func FindReservationByCode(c *gin.Context){
    code := c.Param("code")
    var reservation models.Reservation
    if err := database.DB.Where("code = ?", code).First(&reservation).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": "Reservation not found",
        })
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "data": reservation,
    })
}

func ConfirmReservation(c *gin.Context){
    var req ConfirmReservationRequest
		if err := c.ShouldBindJSON(&req); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "code wajib di isi",
			})
			return
		}
		code := req.Code
    var reservation models.Reservation
    if err := database.DB.Where("code = ?", code).First(&reservation).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": "Reservation not found",
        })
        return
    }
    if strings.EqualFold(reservation.Status, "hadir") {
        c.JSON(http.StatusOK, gin.H{
            "message": "Reservation already confirmed",
            "data":    reservation,
        })
        return
    }
    reservation.Status = "hadir"
		if err := database.DB.Save(&reservation).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
						"error": "Failed to confirm reservation",
				})
				return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": reservation,
			"message": "Reservation confirmed successfully",
		})
	}

func CreateReservation(c *gin.Context){
	var req ReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "nama, kehadiran, email wajib di isi",
		})
		return
	}
	code, err := generateUniqueReservationCode()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate reservation code",
		})
		return
	}
	reservation := models.Reservation{
		Name: req.Name,
		IsPresent: *req.IsPresent,
		Email: req.Email,
        Code: code,
        TotalGuests: req.TotalGuests,
	}
	if err := database.DB.Create(&reservation).Error; err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create reservation",
		})
		return
	}
	go sendReservationEmailByProvider(reservation)
	c.JSON(http.StatusCreated, gin.H{
		"data": reservation,
		"message": "Reservation created successfully",
	})
}

func sendReservationEmailByProvider(r models.Reservation) {
	emailLower := strings.ToLower(strings.TrimSpace(r.Email))
	if strings.HasSuffix(emailLower, "@icloud.com") || strings.HasSuffix(emailLower, "@me.com") || strings.HasSuffix(emailLower, "@mac.com") {
		sendReservationEmailIcloud(r)
		return
	}
	sendReservationEmail(r)
}
func sendReservationEmailIcloud(r models.Reservation){
	cfg := config.Cfg

    if cfg.MailtrapToken == "" {
        fmt.Println("MAILTRAP_TOKEN is empty, skip sending email")
        return
    }
   
    data := email.WeddingEmailData{
		Name:           r.Name,       
		Email:          r.Email,      
		EventDate:           "02 January 2006", 
		EventTime:           "09:00",           
		VenueName:           "Gedung A",
		VenueAddress:        "Jakarta, Indonesia",
		ReservationCode:     r.Code, // misal kode unik
		ReservationDetailURL: "https://wedding.mohaproject.dev",
		BrideName:           "Cica Purwanti", // bisa dari config
		GroomName:           "Muhamad Andriyansyah",   // bisa dari config
		Year:                time.Now().Year(),
	}
    htmlBody, err := email.BuildWeddingReservationEmailIcloud(data)
	if err != nil {
		fmt.Println("failed build email template:", err)
		return
	}

    url := "https://send.api.mailtrap.io/api/send"

	payload := map[string]interface{}{
		"from": map[string]string{
			"email": cfg.MailtrapFromEmail,
			"name":  cfg.MailtrapFromName,
		},
		"to": []map[string]string{
			{
				"email": r.Email,
			},
		},
		"subject": "Konfirmasi Reservasi Wedding",
		"text": "Halo {{.Name}}, terima kasih sudah mengonfirmasi kehadiran. Tanggal: ...",
		"html":    htmlBody, // PENTING: pakai html, bukan text
		"category": "Wedding Reservation",
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Failed to marshal Mailtrap payload:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		fmt.Println("Failed to create Mailtrap request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.MailtrapToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Failed to send Mailtrap request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Println("Mailtrap returned non-2xx status:", resp.Status)
		return
	}

	fmt.Println("Reservation email sent to:", r.Email)

    // url := "https://send.api.mailtrap.io/api/send"

    // // payload mengikuti contoh kamu, tapi dinamis
    // payload := map[string]interface{}{
    //     "from": map[string]string{
    //         "email": cfg.MailtrapFromEmail,
    //         "name":  cfg.MailtrapFromName,
    //     },
    //     "to": []map[string]string{
    //         {
    //             "email": r.Email,
    //         },
    //     },
    //     "subject": "Konfirmasi Reservasi",
    //     "text": fmt.Sprintf(
    //         "Hallo %s,\n\nTerima kasih sudah melakukan reservasi.\n\nNama: %s\nEmail: %s\n\nSampai jumpa!",
    //         r.Name,
    //         r.Name,
    //         r.Email,
    //     ),
    //     "category": "Reservation",
    // }

    // bodyBytes, err := json.Marshal(payload)
    // if err != nil {
    //     fmt.Println("Failed to marshal Mailtrap payload:", err)
    //     return
    // }

    // req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
    // if err != nil {
    //     fmt.Println("Failed to create Mailtrap request:", err)
    //     return
    // }

    // req.Header.Add("Authorization", "Bearer "+cfg.MailtrapToken)
    // req.Header.Add("Content-Type", "application/json")

    // client := &http.Client{}
    // res, err := client.Do(req)
    // if err != nil {
    //     fmt.Println("Failed to send Mailtrap request:", err)
    //     return
    // }
    // defer res.Body.Close()

    // if res.StatusCode < 200 || res.StatusCode >= 300 {
    //     fmt.Println("Mailtrap API returned non-2xx status:", res.Status)
    //     return
    // }

    // fmt.Println("Reservation email sent via Mailtrap to", r.Email)
}
func sendReservationEmail(r models.Reservation){
	cfg := config.Cfg

    if cfg.MailtrapToken == "" {
        fmt.Println("MAILTRAP_TOKEN is empty, skip sending email")
        return
    }
   
    data := email.WeddingEmailData{
		Name:           r.Name,       
		Email:          r.Email,      
		EventDate:           "02 January 2006", 
		EventTime:           "09:00",           
		VenueName:           "Gedung A",
		VenueAddress:        "Jakarta, Indonesia",
		ReservationCode:     r.Code, // misal kode unik
		ReservationDetailURL: "https://wedding.mohaproject.dev",
		BrideName:           "Cica Purwanti", // bisa dari config
		GroomName:           "Muhamad Andriyansyah",   // bisa dari config
		Year:                time.Now().Year(),
	}
    htmlBody, err := email.BuildWeddingReservationEmail(data)
	if err != nil {
		fmt.Println("failed build email template:", err)
		return
	}

    url := "https://send.api.mailtrap.io/api/send"

	payload := map[string]interface{}{
		"from": map[string]string{
			"email": cfg.MailtrapFromEmail,
			"name":  cfg.MailtrapFromName,
		},
		"to": []map[string]string{
			{
				"email": r.Email,
			},
		},
		"subject": "Konfirmasi Reservasi Wedding",
		"html":    htmlBody, // PENTING: pakai html, bukan text
		"category": "Wedding Reservation",
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Failed to marshal Mailtrap payload:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		fmt.Println("Failed to create Mailtrap request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.MailtrapToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Failed to send Mailtrap request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Println("Mailtrap returned non-2xx status:", resp.Status)
		return
	}

	fmt.Println("Reservation email sent to:", r.Email)
}

func generateUniqueReservationCode() (string, error) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		code := fmt.Sprintf("%05d", rnd.Intn(100000))
		var count int64
		if err := database.DB.Model(&models.Reservation{}).Where("code = ?", code).Count(&count).Error; err != nil {
			return "", err
		}
		if count == 0 {
			return code, nil
		}
	}
}
