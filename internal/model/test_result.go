package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// TestResult 测试结果表
type TestResult struct {
	ID             uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt      time.Time  `                                json:"created_at"`
	UserID         uint       `gorm:"index;not null"           json:"user_id"`
	Scores         ScoresMap  `gorm:"type:json;not null"       json:"scores"`
	AnswerDetails  AnswerList `gorm:"type:json;not null"       json:"answer_details"`
	TotalQuestions int        `gorm:"not null"                 json:"total_questions"`
	AnsweredCount  int        `gorm:"not null"                 json:"answered_count"`
	Timestamp      int64      `gorm:"not null"                 json:"timestamp"`
}

// ScoresMap 存储各类型分数，key 是类型 ID，value 是分数
type ScoresMap map[string]int

// AnswerDetail 单道题的答题详情
type AnswerDetail struct {
	QuestionID int    `json:"question_id"`
	Type       int    `json:"type"`
	Answer     int    `json:"answer"`
	Score      int    `json:"score"`
	Text       string `json:"text"`
}

// AnswerList 答题详情列表
type AnswerList []AnswerDetail

// ============ GORM JSON 序列化实现 ============

func (s ScoresMap) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *ScoresMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, s)
}

func (a AnswerList) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *AnswerList) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, a)
}
