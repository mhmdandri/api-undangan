package controller

import (
	"api-undangan/config"
	"api-undangan/database"
	"api-undangan/email"
	"api-undangan/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)
const maxReservationsCodeAttempts = 5
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
		"message": "Get reservations successfully!",
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
				"message": "Get reservation successfully!",
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
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Reservation already confirmed",
        })
        return
    }
		if !reservation.IsPresent{
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Reservasi tamu memilih tidak hadir",
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
	var reservation models.Reservation
	for attempt := 0; attempt < maxReservationsCodeAttempts; attempt++ {
		code, err := generateUniqueReservationCode()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to generate reservation code",
			})
			return
		}
		reservation = models.Reservation{
			Name: req.Name,
			IsPresent: *req.IsPresent,
			Email: req.Email,
			Code: code,
			TotalGuests: req.TotalGuests,
		}
		if err := database.DB.Create(&reservation).Error; err != nil{
			if errors.Is(err, gorm.ErrDuplicatedKey){
				continue
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to create reservation",
			})
			return
		}
		if reservation.IsPresent {
			go sendReservationEmailByProvider(reservation)
		}
		c.JSON(http.StatusCreated, gin.H{
			"data": reservation,
			"message": "Reservation created successfully",
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "failed to create reservation, please try again",
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

func sendReservationEmailIcloud(r models.Reservation) {
	cfg := config.Cfg

	if cfg.MailtrapToken == "" {
		fmt.Println("MAILTRAP_TOKEN is empty, skip sending email")
		return
	}
	url := "https://send.api.mailtrap.io/api/send"

	data := email.DefaultWeddingEmailData(r.Name, r.Email, r.Code, r.TotalGuests)
	htmlBody, err := email.BuildWeddingReservationEmailIcloud(data)
	if err != nil {
		fmt.Println("failed build icloud email template:", err)
		return
	}
	textBody := email.BuildWeddingReservationPlainText(data)

	payload := map[string]interface{}{
		"from": map[string]string{
			"email": cfg.MailtrapFromEmail,
			"name":  cfg.MailtrapFromName,
		},
		"to": []map[string]string{
			{
				"email": r.Email,
				"name":  r.Name,
			},
		},
		"subject":  fmt.Sprintf("Konfirmasi Kehadiran: %s - The Wedding of Andri & Cica", r.Code),
		"html":     htmlBody,
		"text":     textBody,
		"category": "Wedding Reservation",
		"headers": map[string]string{
			"X-Entity-Ref-ID": r.Code,
		},
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
		fmt.Println("Mailtrap returned non-2xx status for iCloud:", resp.Status)
		return
	}

	fmt.Println("Reservation email (iCloud) sent successfully to:", r.Email)
}

func sendReservationEmail(r models.Reservation) {
	cfg := config.Cfg

	if cfg.MailtrapToken == "" {
		fmt.Println("MAILTRAP_TOKEN is empty, skip sending email")
		return
	}

	data := email.DefaultWeddingEmailData(r.Name, r.Email, r.Code, r.TotalGuests)
	htmlBody, err := email.BuildWeddingReservationEmail(data)
	if err != nil {
		fmt.Println("failed build email template:", err)
		return
	}
	textBody := email.BuildWeddingReservationPlainText(data)

	url := "https://send.api.mailtrap.io/api/send"

	payload := map[string]interface{}{
		"from": map[string]string{
			"email": cfg.MailtrapFromEmail,
			"name":  cfg.MailtrapFromName,
		},
		"to": []map[string]string{
			{
				"email": r.Email,
				"name":  r.Name,
			},
		},
		"subject":  fmt.Sprintf("Konfirmasi Kehadiran: %s - The Wedding of Andri & Cica", r.Code),
		"html":     htmlBody,
		"text":     textBody,
		"category": "Wedding Reservation",
		"headers": map[string]string{
			"X-Entity-Ref-ID": r.Code,
		},
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

	fmt.Println("Reservation email sent successfully to:", r.Email)
}

// PreviewReservationEmail lets anyone preview the rendered email (HTML or plain text) in browser
func PreviewReservationEmail(c *gin.Context) {
	name := c.DefaultQuery("name", "Tamu Undangan")
	code := c.DefaultQuery("code", "74921")
	provider := strings.ToLower(c.DefaultQuery("provider", "gmail"))
	format := strings.ToLower(c.DefaultQuery("format", "html"))

	data := email.DefaultWeddingEmailData(name, "preview@example.com", code, 2)

	if format == "text" {
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(http.StatusOK, email.BuildWeddingReservationPlainText(data))
		return
	}

	var htmlContent string
	var err error
	if provider == "icloud" {
		htmlContent, err = email.BuildWeddingReservationEmailIcloud(data)
	} else {
		htmlContent, err = email.BuildWeddingReservationEmail(data)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, htmlContent)
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
