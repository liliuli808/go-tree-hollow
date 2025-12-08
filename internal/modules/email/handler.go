package email

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


type EmailHandler struct {
    emailService EmailService
}

func NewEmailHandler(emailService EmailService) *EmailHandler {
    return &EmailHandler{emailService: emailService}
}


func (h *EmailHandler) SendVerificationCode(c *gin.Context) {
    var req SendCodeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
        return
    }
    
    if err := h.emailService.SendVerificationCode(c.Request.Context(), req.Email); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, SendCodeResponse{Message: "验证码已发送，请查收邮件"})
}