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
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const maxReservationsCodeAttempts = 5

type ReservationRequest struct {
	Code        string `json:"code"`
	Name        string `json:"name" binding:"required"`
	IsPresent   *bool  `json:"is_present" binding:"required"`
	Email       string `json:"email" binding:"omitempty,email"`
	Phone       string `json:"phone"`
	TotalGuests int    `json:"total_guests" binding:"omitempty,min=1,max=10"`
}

type BatchGuestItem struct {
	Name        string `json:"name" binding:"required"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	QuotaGuests int    `json:"quota_guests" binding:"omitempty,min=1,max=10"`
}

type BatchReservationRequest struct {
	Guests []BatchGuestItem `json:"guests" binding:"required"`
}

type BatchReservationItemResponse struct {
	models.Reservation
	InviteURL string `json:"invite_url"`
	QRCodeURL string `json:"qr_code_url"`
}
type ConfirmReservationRequest struct {
	Code string `json:"code" binding:"required"`
}
type GuestCheckInRequest struct {
	Code         string `json:"code" binding:"required"`
	ActualGuests int    `json:"actual_guests" binding:"omitempty,min=1,max=3"`
}

func GetReservations(c *gin.Context) {
	var reservations []models.Reservation
	if err := database.DB.Order("created_at desc").Find(&reservations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch reservations",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    reservations,
		"message": "Get reservations successfully!",
	})
}
func FindReservationByCode(c *gin.Context) {
	code := c.Param("code")
	var reservation models.Reservation
	if err := database.DB.Where("code = ?", code).First(&reservation).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Reservation not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    reservation,
		"message": "Get reservation successfully!",
	})
}

func ConfirmReservation(c *gin.Context) {
	var req ConfirmReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
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
	if !reservation.IsPresent {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Reservasi tamu memilih tidak hadir",
		})
		return
	}
	now := time.Now()
	reservation.Status = "hadir"
	reservation.AttendedAt = &now
	if err := database.DB.Save(&reservation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to confirm reservation",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    reservation,
		"message": "Reservation confirmed successfully",
	})
}

// CheckInLookup is a public endpoint to check guest details by reservation code
func CheckInLookup(c *gin.Context) {
	code := strings.TrimSpace(c.Param("code"))
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "Kode reservasi wajib diisi",
		})
		return
	}

	var reservation models.Reservation
	if err := database.DB.Where("code = ?", code).First(&reservation).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "not_found",
			"error":  "Kode reservasi tidak valid atau tidak ditemukan. Mohon pastikan kode sesuai dengan email konfirmasi Anda.",
		})
		return
	}

	isAlreadyCheckedIn := strings.EqualFold(reservation.Status, "hadir")

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data reservasi ditemukan",
		"data": gin.H{
			"id":                    reservation.ID,
			"name":                  reservation.Name,
			"email":                 reservation.Email,
			"code":                  reservation.Code,
			"total_guests":          reservation.TotalGuests,
			"status":                reservation.Status,
			"is_present":            reservation.IsPresent,
			"is_already_checked_in": isAlreadyCheckedIn,
			"attended_at":           reservation.AttendedAt,
			"created_at":            reservation.CreatedAt,
		},
	})
}

// GuestCheckIn is a public endpoint allowing guests to submit their reservation code to confirm arrival
func GuestCheckIn(c *gin.Context) {
	var req GuestCheckInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "Kode reservasi wajib diisi",
		})
		return
	}

	code := strings.TrimSpace(req.Code)
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "Kode reservasi tidak boleh kosong",
		})
		return
	}

	var reservation models.Reservation
	if err := database.DB.Where("code = ?", code).First(&reservation).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "not_found",
			"error":  "Kode reservasi tidak valid atau tidak ditemukan. Mohon cek kembali kode di email Anda.",
		})
		return
	}

	// Cek apakah tamu sudah check-in sebelumnya
	if strings.EqualFold(reservation.Status, "hadir") {
		c.JSON(http.StatusOK, gin.H{
			"status":  "already_checked_in",
			"message": fmt.Sprintf("Halo %s, Anda sudah melakukan check-in kehadiran sebelumnya!", reservation.Name),
			"data":    reservation,
		})
		return
	}

	// Update status kehadiran tamu
	now := time.Now()
	reservation.Status = "hadir"
	reservation.IsPresent = true
	reservation.AttendedAt = &now

	if req.ActualGuests >= 1 && req.ActualGuests <= 3 {
		reservation.TotalGuests = req.ActualGuests
	}

	if err := database.DB.Save(&reservation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Gagal menyimpan konfirmasi kehadiran, silakan coba lagi.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": fmt.Sprintf("Selamat datang %s! Kehadiran Anda berhasil dikonfirmasi. Selamat menikmati acara pernikahan Andri & Cica.", reservation.Name),
		"data":    reservation,
	})
}

// GuestCheckInByCode allows a check-in scanner to submit the scanned reservation code directly.
func GuestCheckInByCode(c *gin.Context) {
	code := strings.TrimSpace(c.Param("code"))
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "Kode reservasi tidak boleh kosong",
		})
		return
	}

	var reservation models.Reservation
	if err := database.DB.Where("code = ?", code).First(&reservation).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "not_found",
			"error":  "Kode reservasi tidak valid atau tidak ditemukan.",
		})
		return
	}

	if strings.EqualFold(reservation.Status, "hadir") {
		c.JSON(http.StatusOK, gin.H{
			"status":  "already_checked_in",
			"message": fmt.Sprintf("Halo %s, Anda sudah melakukan check-in kehadiran sebelumnya!", reservation.Name),
			"data":    reservation,
		})
		return
	}

	now := time.Now()
	reservation.Status = "hadir"
	reservation.IsPresent = true
	reservation.AttendedAt = &now

	if err := database.DB.Save(&reservation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Gagal menyimpan konfirmasi kehadiran, silakan coba lagi.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": fmt.Sprintf("Selamat datang %s! Kehadiran Anda berhasil dikonfirmasi. Selamat menikmati acara pernikahan Andri & Cica.", reservation.Name),
		"data":    reservation,
	})
}

func CreateReservation(c *gin.Context) {
	var req ReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Nama dan status kehadiran wajib diisi",
			"detail": err.Error(),
		})
		return
	}

	reqCode := strings.TrimSpace(req.Code)

	// JALUR 1: Jika request membawa kode reservasi (tamu pre-registered atau update RSVP)
	if reqCode != "" {
		var existing models.Reservation
		if err := database.DB.Where("code = ?", reqCode).First(&existing).Error; err == nil {
			existing.Name = req.Name
			existing.IsPresent = *req.IsPresent
			if req.Email != "" {
				existing.Email = strings.TrimSpace(req.Email)
			}
			if req.Phone != "" {
				existing.Phone = strings.TrimSpace(req.Phone)
			}
			if req.TotalGuests > 0 {
				existing.TotalGuests = req.TotalGuests
			}
			if *req.IsPresent {
				existing.Status = "konfirmasi_hadir"
			} else {
				existing.Status = "tidak_datang"
			}

			if err := database.DB.Save(&existing).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Gagal memperbarui konfirmasi reservasi",
				})
				return
			}

			if existing.IsPresent && strings.TrimSpace(existing.Email) != "" {
				go sendReservationEmailByProvider(existing)
			}

			c.JSON(http.StatusOK, gin.H{
				"data":    existing,
				"message": "Konfirmasi kehadiran berhasil disimpan",
			})
			return
		}
	}

	// JALUR 2: Tamu Mandiri (Self-registered tanpa kode)
	var reservation models.Reservation
	for attempt := 0; attempt < maxReservationsCodeAttempts; attempt++ {
		code, err := generateUniqueReservationCode()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to generate reservation code",
			})
			return
		}

		status := "tidak_datang"
		if *req.IsPresent {
			status = "konfirmasi_hadir"
		}

		totalGuests := req.TotalGuests
		if totalGuests <= 0 {
			totalGuests = 1
		}

		reservation = models.Reservation{
			Name:        req.Name,
			IsPresent:   *req.IsPresent,
			Email:       strings.TrimSpace(req.Email),
			Phone:       strings.TrimSpace(req.Phone),
			Code:        code,
			QuotaGuests: totalGuests,
			TotalGuests: totalGuests,
			Status:      status,
			Source:      "self_registered",
		}

		if err := database.DB.Create(&reservation).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				continue
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to create reservation",
			})
			return
		}

		if reservation.IsPresent && strings.TrimSpace(reservation.Email) != "" {
			go sendReservationEmailByProvider(reservation)
		}

		c.JSON(http.StatusCreated, gin.H{
			"data":    reservation,
			"message": "Reservation created successfully",
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "failed to create reservation, please try again",
	})
}
func sendReservationEmailByProvider(r models.Reservation) {
	if strings.TrimSpace(r.Email) == "" {
		return
	}
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
		"reply_to": map[string]string{
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
		"reply_to": map[string]string{
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
	return generateUniqueReservationCodeWithDB(database.DB)
}

func generateUniqueReservationCodeWithDB(db *gorm.DB) (string, error) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for attempt := 0; attempt < 50; attempt++ {
		code := fmt.Sprintf("%05d", rnd.Intn(100000))
		var count int64
		if err := db.Model(&models.Reservation{}).Where("code = ?", code).Count(&count).Error; err != nil {
			return "", err
		}
		if count == 0 {
			return code, nil
		}
	}
	return "", fmt.Errorf("exhausted unique reservation code generation attempts")
}


// BatchCreateReservations handles bulk pre-registration for VIP / shared guests
func BatchCreateReservations(c *gin.Context) {
	var req BatchReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Format data tamu tidak valid. Pastikan array 'guests' dengan nama terisi.",
			"detail": err.Error(),
		})
		return
	}

	if len(req.Guests) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Daftar tamu tidak boleh kosong",
		})
		return
	}

	siteURL := "https://weddingofandricica.me"
	if config.Cfg != nil && config.Cfg.ClientOrigin != "" {
		siteURL = strings.TrimRight(config.Cfg.ClientOrigin, "/")
	}

	var createdList []BatchReservationItemResponse
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, item := range req.Guests {
		name := strings.TrimSpace(item.Name)
		if name == "" {
			continue
		}

		code, err := generateUniqueReservationCodeWithDB(tx)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gagal generate kode reservasi unik",
			})
			return
		}

		quota := item.QuotaGuests
		if quota <= 0 {
			quota = 2
		}

		guest := models.Reservation{
			Name:        name,
			Phone:       strings.TrimSpace(item.Phone),
			Email:       strings.TrimSpace(item.Email),
			Code:        code,
			QuotaGuests: quota,
			TotalGuests: quota,
			IsPresent:   false,
			Status:      "belum_konfirmasi",
			Source:      "pre_registered",
		}

		if err := tx.Create(&guest).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":  "Gagal menyimpan data tamu",
				"detail": err.Error(),
			})
			return
		}

		inviteURL := fmt.Sprintf("%s?code=%s", siteURL, code)
		qrCodeURL := fmt.Sprintf("https://api.qrserver.com/v1/create-qr-code/?size=250x250&margin=10&data=%s", url.QueryEscape(inviteURL))
		createdList = append(createdList, BatchReservationItemResponse{
			Reservation: guest,
			InviteURL:   inviteURL,
			QRCodeURL:   qrCodeURL,
		})
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal commit database",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("Berhasil membuat %d undangan pre-registered", len(createdList)),
		"total":   len(createdList),
		"data":    createdList,
	})
}
