package service

import (
	"uniS/internal/model"
	"uniS/internal/repository"
	"uniS/pkg/logger"
)

type SubmitTestReq struct {
	Scores         model.ScoresMap  `json:"scores"         binding:"required"`
	AnswerDetails  model.AnswerList `json:"answer_details" binding:"required"`
	TotalQuestions int              `json:"total_questions"`
	AnsweredCount  int              `json:"answered_count"`
	Timestamp      int64            `json:"timestamp"`
}

type SubmitTestResp struct {
	ID        uint            `json:"id"`
	CreatedAt string          `json:"created_at"`
	Scores    model.ScoresMap `json:"scores"`
	Summary   map[string]int  `json:"summary"`
}

type TestService interface {
	Submit(userID uint, req *SubmitTestReq) (*SubmitTestResp, error)
	GetHistory(userID uint) ([]model.TestResult, error)
	GetByID(id uint, userID uint) (*model.TestResult, error)
}

type testService struct {
	testRepo repository.TestResultRepository
}

func NewTestService(testRepo repository.TestResultRepository) TestService {
	return &testService{testRepo: testRepo}
}

func (s *testService) Submit(userID uint, req *SubmitTestReq) (*SubmitTestResp, error) {
	result := &model.TestResult{
		UserID:         userID,
		Scores:         req.Scores,
		AnswerDetails:  req.AnswerDetails,
		TotalQuestions: req.TotalQuestions,
		AnsweredCount:  req.AnsweredCount,
		Timestamp:      req.Timestamp,
	}

	if err := s.testRepo.Create(result); err != nil {
		logger.Error("test_service", "保存测试结果失败", map[string]any{
			"user_id": userID,
			"error":   err.Error(),
		})
		return nil, err
	}

	logger.Info("test_service", "测试结果提交成功", map[string]any{
		"user_id":         userID,
		"result_id":       result.ID,
		"total_questions": result.TotalQuestions,
		"answered_count":  result.AnsweredCount,
	})

	// 计算各类型分数汇总
	summary := make(map[string]int)
	for k, v := range req.Scores {
		summary[k] = v
	}

	return &SubmitTestResp{
		ID:        result.ID,
		CreatedAt: result.CreatedAt.Format("2006-01-02 15:04:05"),
		Scores:    result.Scores,
		Summary:   summary,
	}, nil
}

func (s *testService) GetHistory(userID uint) ([]model.TestResult, error) {
	results, err := s.testRepo.FindByUserID(userID)
	if err != nil {
		logger.Error("test_service", "查询历史记录失败", map[string]any{
			"user_id": userID,
			"error":   err.Error(),
		})
		return nil, err
	}
	return results, nil
}

func (s *testService) GetByID(id uint, userID uint) (*model.TestResult, error) {
	result, err := s.testRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	// 只能查看自己的记录
	if result.UserID != userID {
		return nil, ErrNotAuthorized
	}
	return result, nil
}

var ErrNotAuthorized = &CustomError{Code: 403, Message: "not authorized"}

type CustomError struct {
	Code    int
	Message string
}

func (e *CustomError) Error() string {
	return e.Message
}
