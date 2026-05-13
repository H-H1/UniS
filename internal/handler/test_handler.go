package handler

import (
	"net/http"
	"uniS/internal/service"
	"uniS/pkg/logger"

	"github.com/gin-gonic/gin"
)

type TestHandler struct {
	testSvc service.TestService
}

func NewTestHandler(testSvc service.TestService) *TestHandler {
	return &TestHandler{testSvc: testSvc}
}

// SubmitTest godoc
// POST /api/test/submit
// 提交测试结果（需登录）
func (h *TestHandler) SubmitTest(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req service.SubmitTestReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("test_handler", "SubmitTest 参数绑定失败", map[string]any{
			"user_id": userID,
			"error":   err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "invalid request body"})
		return
	}

	resp, err := h.testSvc.Submit(userID, &req)
	if err != nil {
		logger.Error("test_handler", "SubmitTest 失败", map[string]any{
			"user_id": userID,
			"error":   err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	logger.Info("test_handler", "SubmitTest 成功", map[string]any{
		"user_id":   userID,
		"result_id": resp.ID,
	})

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": resp,
	})
}

// GetTestHistory godoc
// GET /api/test/history
// 获取当前用户的测试历史记录
func (h *TestHandler) GetTestHistory(c *gin.Context) {
	userID := c.GetUint("user_id")

	results, err := h.testSvc.GetHistory(userID)
	if err != nil {
		logger.Error("test_handler", "GetTestHistory 失败", map[string]any{
			"user_id": userID,
			"error":   err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": results,
	})
}

// GetTestResult godoc
// GET /api/test/:id
// 获取单条测试结果详情
func (h *TestHandler) GetTestResult(c *gin.Context) {
	userID := c.GetUint("user_id")

	var uri struct {
		ID uint `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "invalid id"})
		return
	}

	result, err := h.testSvc.GetByID(uri.ID, userID)
	if err != nil {
		if err == service.ErrNotAuthorized {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "not authorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": result,
	})
}
