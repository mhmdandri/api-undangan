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
	"io"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
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

type ManualReservationRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"omitempty,email"`
	Phone string `json:"phone"`
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
	// Parse pagination parameters (default 25 per page)
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limitStr := c.Query("limit")
	if limitStr == "" {
		limitStr = c.Query("per_page")
	}

	fetchAll := strings.EqualFold(c.Query("all"), "true") || limitStr == "-1" || strings.EqualFold(limitStr, "all")

	limit := 25
	if !fetchAll && limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	queryDB := database.DB.Model(&models.Reservation{})

	// Optional search query parameter
	q := strings.TrimSpace(c.Query("q"))
	if q == "" {
		q = strings.TrimSpace(c.Query("search"))
	}
	if q != "" {
		pattern := "%" + strings.ToLower(q) + "%"
		queryDB = queryDB.Where(
			"LOWER(name) LIKE ? OR LOWER(code) LIKE ? OR LOWER(email) LIKE ? OR phone LIKE ?",
			pattern, pattern, pattern, pattern,
		)
	}

	// Optional status filter
	status := strings.TrimSpace(c.Query("status"))
	if status != "" && !strings.EqualFold(status, "all") {
		queryDB = queryDB.Where("LOWER(status) = ?", strings.ToLower(status))
	}

	// Optional source filter
	source := strings.TrimSpace(c.Query("source"))
	if source != "" && !strings.EqualFold(source, "all") {
		queryDB = queryDB.Where("LOWER(source) = ?", strings.ToLower(source))
	}

	var totalRecords int64
	if err := queryDB.Count(&totalRecords).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to count reservations",
		})
		return
	}

	totalPages := 0
	if fetchAll {
		limit = int(totalRecords)
		if limit <= 0 {
			limit = 25
		}
		totalPages = 1
	} else if limit > 0 {
		totalPages = int(math.Ceil(float64(totalRecords) / float64(limit)))
	}

	reservations := make([]models.Reservation, 0)
	execDB := queryDB.Order("created_at desc")

	if !fetchAll {
		offset := (page - 1) * limit
		execDB = execDB.Offset(offset).Limit(limit)
	}

	if err := execDB.Find(&reservations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch reservations",
		})
		return
	}

	var hadirCount int64
	var preRegCount int64
	var selfRegCount int64
	_ = database.DB.Model(&models.Reservation{}).Where("LOWER(status) = ? OR is_present = ?", "hadir", true).Count(&hadirCount).Error
	_ = database.DB.Model(&models.Reservation{}).Where("LOWER(source) = ?", "pre_registered").Count(&preRegCount).Error
	_ = database.DB.Model(&models.Reservation{}).Where("LOWER(source) != ? OR source IS NULL", "pre_registered").Count(&selfRegCount).Error

	c.JSON(http.StatusOK, gin.H{
		"data":    reservations,
		"message": "Get reservations successfully!",
		"pagination": gin.H{
			"page":          page,
			"limit":         limit,
			"total_records": totalRecords,
			"total_pages":   totalPages,
		},
		"total":       totalRecords,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
		"stats": gin.H{
			"total":           totalRecords,
			"hadir":           hadirCount,
			"pre_registered":  preRegCount,
			"self_registered": selfRegCount,
		},
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
	emailClean := strings.ToLower(strings.TrimSpace(req.Email))
	phoneClean := strings.TrimSpace(req.Phone)

	// JALUR 1: Jika request membawa kode reservasi (tamu pre-registered atau update RSVP)
	if reqCode != "" {
		var existing models.Reservation
		if err := database.DB.Where("code = ?", reqCode).First(&existing).Error; err == nil {
			// Cek apakah tamu sudah melakukan reservasi / konfirmasi kehadiran sebelumnya
			if strings.EqualFold(existing.Status, "konfirmasi_hadir") || strings.EqualFold(existing.Status, "hadir") {
				guestName := strings.TrimSpace(existing.Name)
				msg := "Anda sudah melakukan konfirmasi kehadiran sebelumnya"
				if guestName != "" {
					msg = fmt.Sprintf("Halo %s, Anda sudah melakukan konfirmasi kehadiran sebelumnya", guestName)
				}

				c.JSON(http.StatusOK, gin.H{
					"status":               "already_confirmed",
					"message":              msg,
					"data":                 existing,
					"is_already_confirmed": true,
				})
				return
			}

			existing.Name = req.Name
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

			if *req.IsPresent && strings.TrimSpace(existing.Email) != "" {
				go sendReservationEmailByProvider(existing)
			}

			c.JSON(http.StatusOK, gin.H{
				"status":  "success",
				"data":    existing,
				"message": "Konfirmasi kehadiran berhasil disimpan",
			})
			return
		}
	}

	// JALUR 2: Cek berdasarkan Email atau Nomor Telepon jika tamu tidak memasukkan kode
	var existingByContact models.Reservation
	foundByContact := false

	if emailClean != "" {
		if err := database.DB.Where("LOWER(email) = ? AND email != ''", emailClean).First(&existingByContact).Error; err == nil {
			foundByContact = true
		}
	}
	if !foundByContact && phoneClean != "" {
		if err := database.DB.Where("phone = ? AND phone != ''", phoneClean).First(&existingByContact).Error; err == nil {
			foundByContact = true
		}
	}

	if foundByContact {
		// Jika kontak (email/telepon) sudah pernah melakukan konfirmasi kehadiran sebelumnya
		if strings.EqualFold(existingByContact.Status, "konfirmasi_hadir") || strings.EqualFold(existingByContact.Status, "hadir") {
			guestName := strings.TrimSpace(existingByContact.Name)
			msg := "Anda sudah melakukan konfirmasi kehadiran sebelumnya"
			if guestName != "" {
				msg = fmt.Sprintf("Halo %s, Anda sudah melakukan konfirmasi kehadiran sebelumnya", guestName)
			}

			c.JSON(http.StatusOK, gin.H{
				"status":               "already_confirmed",
				"message":              msg,
				"data":                 existingByContact,
				"is_already_confirmed": true,
			})
			return
		}

		// Jika ada data tapi statusnya belum_konfirmasi (misal data pre-registered dari panitia)
		existingByContact.Name = req.Name
		if req.Email != "" {
			existingByContact.Email = strings.TrimSpace(req.Email)
		}
		if req.Phone != "" {
			existingByContact.Phone = strings.TrimSpace(req.Phone)
		}
		if req.TotalGuests > 0 {
			existingByContact.TotalGuests = req.TotalGuests
		}
		if *req.IsPresent {
			existingByContact.Status = "konfirmasi_hadir"
		} else {
			existingByContact.Status = "tidak_datang"
		}

		if err := database.DB.Save(&existingByContact).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gagal memperbarui konfirmasi reservasi",
			})
			return
		}

		if *req.IsPresent && strings.TrimSpace(existingByContact.Email) != "" {
			go sendReservationEmailByProvider(existingByContact)
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"data":    existingByContact,
			"message": "Konfirmasi kehadiran berhasil disimpan",
		})
		return
	}

	// JALUR 3: Tamu Mandiri Baru (Self-registered tanpa riwayat data sebelumnya)
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
			IsPresent:   false,
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

		if *req.IsPresent && strings.TrimSpace(reservation.Email) != "" {
			go sendReservationEmailByProvider(reservation)
		}

		c.JSON(http.StatusCreated, gin.H{
			"status":  "success",
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
	sendReservationEmail(r)
}

func sendReservationEmailIcloud(r models.Reservation) {
	sendReservationEmail(r)
}

// EmailAPIPayload represents the payload required by the custom Email API (SMTP iCloud)
type EmailAPIPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
	Html    string `json:"html"`
}

func normalizeEmailAPIURL(rawURL string) string {
	targetURL := strings.TrimSpace(rawURL)
	if targetURL == "" {
		targetURL = "http://127.0.0.1:8082/api/send-email"
	}
	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		targetURL = "http://" + targetURL
	}
	if !strings.HasSuffix(targetURL, "/api/send-email") {
		targetURL = strings.TrimRight(targetURL, "/") + "/api/send-email"
	}
	return targetURL
}

func sendEmailViaAPI(to, subject, text, htmlContent string) error {
	cfg := config.Cfg
	if cfg == nil || strings.TrimSpace(cfg.EmailAPIKey) == "" {
		return errors.New("EMAIL_API_KEY is empty")
	}

	targetURL := normalizeEmailAPIURL(cfg.EmailAPIURL)

	payload := EmailAPIPayload{
		To:      to,
		Subject: subject,
		Text:    text,
		Html:    htmlContent,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal email payload: %w", err)
	}

	req, err := http.NewRequest("POST", targetURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create email request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", strings.TrimSpace(cfg.EmailAPIKey))

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call email API at %s: %w", targetURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("email API returned status %s: %s", resp.Status, strings.TrimSpace(string(respBody)))
	}

	return nil
}

func sendReservationEmail(r models.Reservation) {
	cfg := config.Cfg
	if cfg == nil {
		fmt.Println("Config is nil, skip sending email")
		return
	}

	data := email.DefaultWeddingEmailData(r.Name, r.Email, r.Code, r.TotalGuests)
	htmlBody, err := email.BuildWeddingReservationEmail(data)
	if err != nil {
		fmt.Println("failed build email template:", err)
		return
	}
	textBody := email.BuildWeddingReservationPlainText(data)
	subject := fmt.Sprintf("Konfirmasi Kehadiran: %s - The Wedding of Andri & Cica", r.Code)

	// Primary: Use new Email API (SMTP iCloud microservice)
	if strings.TrimSpace(cfg.EmailAPIKey) != "" {
		if err := sendEmailViaAPI(r.Email, subject, textBody, htmlBody); err != nil {
			fmt.Printf("Failed to send reservation email to %s via Email API: %v\n", r.Email, err)
			return
		}
		fmt.Printf("Reservation email sent successfully via Email API to: %s (Code: %s)\n", r.Email, r.Code)
		return
	}

	// Fallback: Mailtrap
	if strings.TrimSpace(cfg.MailtrapToken) != "" {
		sendMailtrapEmail(r.Email, r.Name, subject, textBody, htmlBody, r.Code)
		return
	}

	fmt.Println("EMAIL_API_KEY is empty, skip sending email")
}

func sendMailtrapEmail(toEmail, toName, subject, textBody, htmlBody, code string) {
	cfg := config.Cfg
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
				"email": toEmail,
				"name":  toName,
			},
		},
		"subject":  subject,
		"html":     htmlBody,
		"text":     textBody,
		"category": "Wedding Reservation",
		"headers": map[string]string{
			"X-Entity-Ref-ID": code,
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

	fmt.Println("Reservation email sent successfully via Mailtrap to:", toEmail)
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

// CreateManualReservation adds one guest without requiring a spreadsheet import.
func CreateManualReservation(c *gin.Context) {
	var req ManualReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Nama tamu wajib diisi dan format data tidak valid",
			"detail": err.Error(),
		})
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nama tamu wajib diisi",
		})
		return
	}

	var reservation models.Reservation
	for attempt := 0; attempt < maxReservationsCodeAttempts; attempt++ {
		code, err := generateUniqueReservationCode()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gagal membuat kode reservasi",
			})
			return
		}

		reservation = models.Reservation{
			Name:        name,
			Email:       strings.TrimSpace(req.Email),
			Phone:       strings.TrimSpace(req.Phone),
			Code:        code,
			QuotaGuests: 1,
			TotalGuests: 1,
			IsPresent:   false,
			Status:      "belum_konfirmasi",
			Source:      "manual",
		}

		if err := database.DB.Create(&reservation).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				continue
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gagal menyimpan data tamu",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"data":    reservation,
			"message": "Tamu berhasil ditambahkan",
		})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "Gagal membuat reservasi, silakan coba lagi",
	})
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
