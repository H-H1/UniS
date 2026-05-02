package handler

import (
	"net/http"
	"path/filepath"
	"strings"
	"time"
	"uniS/internal/service"
	"uniS/pkg/logger"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userSvc service.UserService
}

func NewUserHandler(userSvc service.UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

// UpdateProfile godoc
// POST /api/user/profile
// 更新昵称和头像 URL（JSON body，字段均可选，只更新传入的非空字段）
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req service.UpdateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("user_handler", "UpdateProfile 参数绑定失败", map[string]any{
			"user_id": userID,
			"error":   err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	user, err := h.userSvc.UpdateProfile(userID, &req)
	if err != nil {
		logger.Error("user_handler", "UpdateProfile 失败", map[string]any{
			"user_id": userID,
			"error":   err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	logger.Info("user_handler", "UpdateProfile 成功", map[string]any{
		"user_id":    userID,
		"nick_name":  user.NickName,
		"avatar_url": user.AvatarURL,
	})
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok", "data": user})
}

// UploadAvatar godoc
// POST /api/user/avatar
// 上传头像文件（multipart/form-data，字段名 file），保存后更新用户 avatar_url。
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	userID := c.GetUint("user_id")

	fh, err := c.FormFile("file")
	if err != nil {
		logger.Error("user_handler", "UploadAvatar 获取文件失败", map[string]any{
			"user_id": userID,
			"error":   err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "file required"})
		return
	}

	// 只允许图片格式
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowed[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "only jpg/png/webp allowed"})
		return
	}

	// 限制 5MB
	if fh.Size > 5<<20 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "file too large (max 5MB)"})
		return
	}

	// 生成唯一文件名存到 uploads/avatars/
	filename := strings.ReplaceAll(time.Now().Format("20060102150405.000"), ".", "") + ext
	savePath := filepath.Join("uploads", "avatars", filename)

	if err := c.SaveUploadedFile(fh, savePath); err != nil {
		logger.Error("user_handler", "UploadAvatar 保存文件失败", map[string]any{
			"user_id":   userID,
			"save_path": savePath,
			"error":     err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "save file failed"})
		return
	}

	// 拼接对外可访问的 URL
	avatarURL := "/static/avatars/" + filename

	user, err := h.userSvc.UpdateProfile(userID, &service.UpdateProfileReq{AvatarURL: avatarURL})
	if err != nil {
		logger.Error("user_handler", "UploadAvatar 更新头像 URL 失败", map[string]any{
			"user_id":    userID,
			"avatar_url": avatarURL,
			"error":      err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	logger.Info("user_handler", "UploadAvatar 成功", map[string]any{
		"user_id":    userID,
		"avatar_url": avatarURL,
	})
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok", "data": gin.H{
		"avatar_url": avatarURL,
		"user":       user,
	}})
}
