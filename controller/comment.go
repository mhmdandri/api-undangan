package controller

import (
	"api-undangan/config"
	"api-undangan/database"
	"api-undangan/models"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const commentCookie = "comment_fp"
type CommentRequest struct {
	Name		string `json:"name" binding:"required"`
	Message		string `json:"message" binding:"required"`
}

func GetComments(c *gin.Context){
	var comments []models.Comment
	if err := database.DB.Order("created_at desc").Find(&comments).Error; err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch comments",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": comments,
		"message": "Get comments successfully!",
	})
}

func clientFingerprint(c *gin.Context) string {
	if fp, err := c.Cookie(commentCookie); err == nil && fp != "" {
		// c.JSON(http.StatusBadRequest, gin.H{
		// 	"error": "Tunggu 1 menit sebelum ucapan lagi",
		// })
		return fp
	}
	fp := uuid.NewString()
	secure := c.Request.TLS != nil || strings.EqualFold(c.GetHeader("X-Forwarded-Proto"), "https") || strings.EqualFold(c.GetHeader("X-Forwarded-Scheme"), "https")
	sameSite := http.SameSiteNoneMode
	if !secure {
		sameSite = http.SameSiteLaxMode
	}
	c.SetSameSite(sameSite)
	c.SetCookie(commentCookie, fp, 86400 * 30, "/", "", secure, true)
	return fp
}

func PostComment(c *gin.Context){
	var req CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Name dan message wajib di isi",
		})
		return
	}
	if config.ContainsBadword(req.Message) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Jangan berkata kasar ya kak",
		})
		return
	}
	fingerprint := clientFingerprint(c)
	ip := c.ClientIP()
	now := time.Now()
	lastMinute := now.Add(-time.Minute)
	var recentFp models.Comment
	errFp := database.DB.Where("fingerprint = ? AND created_at >= ?", fingerprint, lastMinute).First(&recentFp).Error
	if errFp == nil {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": "Tunggu 1 menit sebelum ucapan lagi",
			"code": "RATE_LIMIT",
		})
		return
	}
	if !errors.Is(errFp, gorm.ErrRecordNotFound){
		 c.JSON(http.StatusInternalServerError, gin.H{
        "error": "Cek spam gagal",
        "code":  "CHECK_SPAM_FAILED",
    })
    return
	}
	var countIP int64
	if err := database.DB.Model(&models.Comment{}).Where("ip_address = ? AND created_at >= ?", ip, lastMinute).Count(&countIP).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
        "error": "Cek spam gagal",
        "code":  "CHECK_SPAM_FAILED",
    })
    return
	}
	if countIP >= 5 {
    c.JSON(http.StatusTooManyRequests, gin.H{
        "error": "Terlalu banyak komentar dari jaringan ini, coba sebentar lagi",
        "code":  "RATE_LIMIT_IP",
    })
    return
	}
	comment := models.Comment{
		Name: req.Name,
		Message: req.Message,
		IPAddress: c.ClientIP(),
		Fingerprint: fingerprint,
	}
	if err := database.DB.Create(&comment).Error; err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to post comment",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"data": comment,
		"message": "Post comment successfully!",
	})
}