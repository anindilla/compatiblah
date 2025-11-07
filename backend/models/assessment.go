package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type PersonData struct {
	Name string `json:"name"`
	MBTI string `json:"mbti"`
}

type Assessment struct {
	ID                  string              `json:"id" db:"id"`
	Person1Name         string              `json:"person1_name" db:"person1_name"`
	Person1Data         PersonData          `json:"person1_data" db:"person1_data"`
	Person2Name         string              `json:"person2_name" db:"person2_name"`
	Person2Data         PersonData          `json:"person2_data" db:"person2_data"`
	FriendScore         int                 `json:"friend_score" db:"friend_score"`
	CoworkerScore       int                 `json:"coworker_score" db:"coworker_score"`
	PartnerScore        int                 `json:"partner_score" db:"partner_score"`
	OverallScore        int                 `json:"overall_score" db:"overall_score"`
	FriendExplanation   CategoryExplanation `json:"friend_explanation" db:"friend_explanation"`
	CoworkerExplanation CategoryExplanation `json:"coworker_explanation" db:"coworker_explanation"`
	PartnerExplanation  CategoryExplanation `json:"partner_explanation" db:"partner_explanation"`
	CreatedAt           time.Time           `json:"created_at" db:"created_at"`
}

type AssessmentRequest struct {
	Person1 PersonData `json:"person1"`
	Person2 PersonData `json:"person2"`
}

type BulletPoint struct {
	Text string `json:"text"`
}

type SubCategory struct {
	Title   string        `json:"title"`
	Bullets []BulletPoint `json:"bullets"`
}

type ExplanationSection struct {
	Heading       string        `json:"heading"`
	Subcategories []SubCategory `json:"subcategories"`
}

type CategoryExplanation struct {
	Sections []ExplanationSection `json:"sections"`
}

type GeminiResponse struct {
	FriendScore         int                 `json:"friend_score"`
	CoworkerScore       int                 `json:"coworker_score"`
	PartnerScore        int                 `json:"partner_score"`
	OverallScore        int                 `json:"overall_score"`
	FriendExplanation   CategoryExplanation `json:"friend_explanation"`
	CoworkerExplanation CategoryExplanation `json:"coworker_explanation"`
	PartnerExplanation  CategoryExplanation `json:"partner_explanation"`
}

// Value implements driver.Valuer for PersonData
func (p PersonData) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan implements sql.Scanner for PersonData
func (p *PersonData) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, p)
}

// Value implements driver.Valuer for CategoryExplanation
func (c CategoryExplanation) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// Scan implements sql.Scanner for CategoryExplanation
func (c *CategoryExplanation) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, c)
}
