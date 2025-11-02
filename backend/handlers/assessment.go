package handlers

import (
	"net/http"
	"compatiblah/backend/db"
	"compatiblah/backend/models"
	"compatiblah/backend/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AssessCompatibility(c *gin.Context) {
	var req models.AssessmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Validate required fields
	if req.Person1.Name == "" || req.Person1.MBTI == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Person 1 must have a name and MBTI type"})
		return
	}

	if req.Person2.Name == "" || req.Person2.MBTI == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Person 2 must have a name and MBTI type"})
		return
	}

	// Call Gemini API
	geminiResp, err := services.AssessCompatibility(req.Person1, req.Person2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assess compatibility: " + err.Error()})
		return
	}

	// Create assessment record (privacy-first: only save results, NOT personal data)
	assessment := &models.Assessment{
		ID:                  uuid.New().String(),
		FriendScore:         geminiResp.FriendScore,
		CoworkerScore:       geminiResp.CoworkerScore,
		PartnerScore:        geminiResp.PartnerScore,
		OverallScore:        geminiResp.OverallScore,
		FriendExplanation:   geminiResp.FriendExplanation,
		CoworkerExplanation: geminiResp.CoworkerExplanation,
		PartnerExplanation:  geminiResp.PartnerExplanation,
	}

	// Save to database (only assessment results, no personal data stored)
	if err := db.SaveAssessment(assessment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save assessment: " + err.Error()})
		return
	}

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"id":                    assessment.ID,
		"friend_score":          assessment.FriendScore,
		"coworker_score":        assessment.CoworkerScore,
		"partner_score":         assessment.PartnerScore,
		"overall_score":         assessment.OverallScore,
		"friend_explanation":    assessment.FriendExplanation,
		"coworker_explanation":  assessment.CoworkerExplanation,
		"partner_explanation":   assessment.PartnerExplanation,
	})
}

func GetAssessment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Assessment ID is required"})
		return
	}

	assessment, err := db.GetAssessment(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assessment not found"})
		return
	}

	// Return only assessment results, NOT personal data (privacy-first)
	c.JSON(http.StatusOK, gin.H{
		"id":                    assessment.ID,
		"friend_score":          assessment.FriendScore,
		"coworker_score":        assessment.CoworkerScore,
		"partner_score":         assessment.PartnerScore,
		"overall_score":         assessment.OverallScore,
		"friend_explanation":    assessment.FriendExplanation,
		"coworker_explanation":  assessment.CoworkerExplanation,
		"partner_explanation":   assessment.PartnerExplanation,
		"created_at":            assessment.CreatedAt,
	})
}

func GetAllAssessments(c *gin.Context) {
	assessments, err := db.GetAllAssessments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve assessments: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, assessments)
}

