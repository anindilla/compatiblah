package services

import (
	"bytes"
	"compatiblah/backend/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func AssessCompatibility(person1, person2 models.PersonData) (*models.GeminiResponse, error) {
	prompt := buildPrompt(person1, person2)

	body, err := callGeminiAPI(prompt)
	if err != nil {
		return nil, err
	}

	// Parse Gemini response
	var geminiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no content in response")
	}

	text := geminiResp.Candidates[0].Content.Parts[0].Text

	// Extract JSON from the response text (Gemini might wrap it in markdown)
	jsonText := extractJSON(text)

	// Try to parse as new structured format first
	var result models.GeminiResponse
	err = json.Unmarshal([]byte(jsonText), &result)

	// If parsing fails, try old format (backward compatibility)
	if err != nil {
		// Try intermediate format with sections but content field (intermediate format)
		var intermediateFormat struct {
			FriendScore       int `json:"friend_score"`
			CoworkerScore     int `json:"coworker_score"`
			PartnerScore      int `json:"partner_score"`
			OverallScore      int `json:"overall_score"`
			FriendExplanation struct {
				Sections []struct {
					Heading string `json:"heading"`
					Content string `json:"content"`
				} `json:"sections"`
			} `json:"friend_explanation"`
			CoworkerExplanation struct {
				Sections []struct {
					Heading string `json:"heading"`
					Content string `json:"content"`
				} `json:"sections"`
			} `json:"coworker_explanation"`
			PartnerExplanation struct {
				Sections []struct {
					Heading string `json:"heading"`
					Content string `json:"content"`
				} `json:"sections"`
			} `json:"partner_explanation"`
		}

		// Clean up JSON one more time before trying old format
		jsonText = cleanJSONForParsing(jsonText)

		if oldErr := json.Unmarshal([]byte(jsonText), &intermediateFormat); oldErr == nil {
			// Convert intermediate format (sections with content) to new format
			result = models.GeminiResponse{
				FriendScore:         intermediateFormat.FriendScore,
				CoworkerScore:       intermediateFormat.CoworkerScore,
				PartnerScore:        intermediateFormat.PartnerScore,
				OverallScore:        intermediateFormat.OverallScore,
				FriendExplanation:   convertSectionsToSubcategories(intermediateFormat.FriendExplanation.Sections, "friendship"),
				CoworkerExplanation: convertSectionsToSubcategories(intermediateFormat.CoworkerExplanation.Sections, "workplace"),
				PartnerExplanation:  convertSectionsToSubcategories(intermediateFormat.PartnerExplanation.Sections, "romance"),
			}
		} else {
			// Try oldest format (just strings)
			var stringFormat struct {
				FriendScore         int    `json:"friend_score"`
				CoworkerScore       int    `json:"coworker_score"`
				PartnerScore        int    `json:"partner_score"`
				OverallScore        int    `json:"overall_score"`
				FriendExplanation   string `json:"friend_explanation"`
				CoworkerExplanation string `json:"coworker_explanation"`
				PartnerExplanation  string `json:"partner_explanation"`
			}

			if stringErr := json.Unmarshal([]byte(jsonText), &stringFormat); stringErr == nil {
				result = models.GeminiResponse{
					FriendScore:         stringFormat.FriendScore,
					CoworkerScore:       stringFormat.CoworkerScore,
					PartnerScore:        stringFormat.PartnerScore,
					OverallScore:        stringFormat.OverallScore,
					FriendExplanation:   convertStringToStructured(stringFormat.FriendExplanation, "friendship"),
					CoworkerExplanation: convertStringToStructured(stringFormat.CoworkerExplanation, "workplace"),
					PartnerExplanation:  convertStringToStructured(stringFormat.PartnerExplanation, "romance"),
				}
			} else {
				return nil, fmt.Errorf("failed to parse assessment JSON (all formats): new format error: %w, intermediate format error: %v, string format error: %v, cleaned text: %s", err, oldErr, stringErr, jsonText)
			}
		}
	}

	// Validate scores
	if result.FriendScore < 1 || result.FriendScore > 5 {
		result.FriendScore = 3
	}
	if result.CoworkerScore < 1 || result.CoworkerScore > 5 {
		result.CoworkerScore = 3
	}
	if result.PartnerScore < 1 || result.PartnerScore > 5 {
		result.PartnerScore = 3
	}
	if result.OverallScore < 1 || result.OverallScore > 5 {
		result.OverallScore = (result.FriendScore + result.CoworkerScore + result.PartnerScore) / 3
		if result.OverallScore < 1 {
			result.OverallScore = 3
		}
	}

	heuristicScores := calculateCompatibilityScores(person1, person2)
	result.FriendScore = blendScores(result.FriendScore, heuristicScores.Friend)
	result.CoworkerScore = blendScores(result.CoworkerScore, heuristicScores.Coworker)
	result.PartnerScore = blendScores(result.PartnerScore, heuristicScores.Partner)
	result.OverallScore = clampScore(float64(result.FriendScore+result.CoworkerScore+result.PartnerScore) / 3.0)

	return &result, nil
}

// CategoryResponse represents a single category assessment result
type CategoryResponse struct {
	Score       int                        `json:"score"`
	Explanation models.CategoryExplanation `json:"explanation"`
}

// AssessCategoryCompatibility generates compatibility assessment for a single category
func AssessCategoryCompatibility(person1, person2 models.PersonData, category string) (*CategoryResponse, error) {
	prompt := buildCategoryPrompt(person1, person2, category)

	body, err := callGeminiAPI(prompt)
	if err != nil {
		return nil, err
	}

	// Parse Gemini response
	var geminiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no content in response")
	}

	text := geminiResp.Candidates[0].Content.Parts[0].Text

	// Extract JSON from the response text
	jsonText := extractJSON(text)
	jsonText = cleanJSONForParsing(jsonText)

	// Try to parse as structured format
	var result struct {
		Score       int                        `json:"score"`
		Explanation models.CategoryExplanation `json:"explanation"`
	}

	err = json.Unmarshal([]byte(jsonText), &result)

	// If parsing fails, try old format
	if err != nil {
		var oldFormat struct {
			Score       int    `json:"score"`
			Explanation string `json:"explanation"`
		}

		if oldErr := json.Unmarshal([]byte(jsonText), &oldFormat); oldErr == nil {
			result = struct {
				Score       int                        `json:"score"`
				Explanation models.CategoryExplanation `json:"explanation"`
			}{
				Score:       oldFormat.Score,
				Explanation: convertStringToStructured(oldFormat.Explanation, category),
			}
		} else {
			return nil, fmt.Errorf("failed to parse category assessment JSON: %w", err)
		}
	}

	// Validate score
	if result.Score < 1 || result.Score > 5 {
		result.Score = 3
	}

	heuristic := calculateCategoryScore(person1, person2, category)
	finalScore := blendScores(result.Score, heuristic)

	return &CategoryResponse{
		Score:       finalScore,
		Explanation: result.Explanation,
	}, nil
}

// AssessCategoryCompatibilityWithBase assesses compatibility with optional base explanation for augmentation
func AssessCategoryCompatibilityWithBase(person1, person2 models.PersonData, category string, baseExplanation *models.CategoryExplanation) (*CategoryResponse, error) {
	prompt := buildCategoryPromptWithBase(person1, person2, category, baseExplanation)

	body, err := callGeminiAPI(prompt)
	if err != nil {
		return nil, err
	}

	// Parse Gemini response
	var geminiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no content in response")
	}

	text := geminiResp.Candidates[0].Content.Parts[0].Text

	// Extract JSON from the response text
	jsonText := extractJSON(text)
	jsonText = cleanJSONForParsing(jsonText)

	// Try to parse as structured format
	var result struct {
		Score       int                        `json:"score"`
		Explanation models.CategoryExplanation `json:"explanation"`
	}

	err = json.Unmarshal([]byte(jsonText), &result)

	// If parsing fails, try old format
	if err != nil {
		var oldFormat struct {
			Score       int    `json:"score"`
			Explanation string `json:"explanation"`
		}

		if oldErr := json.Unmarshal([]byte(jsonText), &oldFormat); oldErr == nil {
			result = struct {
				Score       int                        `json:"score"`
				Explanation models.CategoryExplanation `json:"explanation"`
			}{
				Score:       oldFormat.Score,
				Explanation: convertStringToStructured(oldFormat.Explanation, category),
			}
		} else {
			return nil, fmt.Errorf("failed to parse category assessment JSON: %w", err)
		}
	}

	// Validate score
	if result.Score < 1 || result.Score > 5 {
		result.Score = 3
	}

	heuristic := calculateCategoryScore(person1, person2, category)
	finalScore := blendScores(result.Score, heuristic)

	return &CategoryResponse{
		Score:       finalScore,
		Explanation: result.Explanation,
	}, nil
}

func buildPrompt(person1, person2 models.PersonData) string {
	prompt := fmt.Sprintf(`You are a compatibility assessment expert. Analyze the compatibility between two people based on ALL the information provided below. You MUST consider and reference their names and MBTI types when making your assessment.

PERSON 1:
- Name: %s
- MBTI Type: %s`, person1.Name, person1.MBTI)
	prompt += fmt.Sprintf(`

PERSON 2:
- Name: %s
- MBTI Type: %s`, person2.Name, person2.MBTI)

	prompt += `

You are an expert compatibility analyst with deep knowledge of personality psychology, relationship dynamics, and interpersonal communication. Drawing from frameworks like MBTI, cognitive functions, and relationship psychology, provide a comprehensive, insightful assessment based on the MBTI information above.

CRITICAL INSTRUCTIONS:
- You MUST reference and incorporate MBTI types and names in your analysis
- Focus entirely on the provided MBTI information to deliver the most accurate assessment

For each compatibility context (friendship, workplace, romance), provide a structured analysis with AT LEAST 3 distinct sections. Each section should have:
- A clear heading
- 2-3 sub-categories with descriptive titles
- Each sub-category should contain 2-3 bullet points (detailed, not just one word)

Use a mix of consistent labels (like "Strengths", "Challenges") and context-specific labels (like "What Draws Them Together" for romance, "Communication Styles" for friendship, "Collaboration Tips" for workplace) where it makes sense.

Return your response as a JSON object with the following EXACT structure (no markdown, no code blocks):
{
  "friend_score": <integer 1-5>,
  "coworker_score": <integer 1-5>,
  "partner_score": <integer 1-5>,
  "overall_score": <integer 1-5>,
  "friend_explanation": {
    "sections": [
      {
        "heading": "<Section 1 heading, e.g., 'Cognitive Compatibility & Communication'>",
        "subcategories": [
          {
            "title": "<Sub-category title, e.g., 'Communication Styles' or 'Strengths'>",
            "bullets": [
              {"text": "<Detailed bullet point (2-3 sentences worth of information per bullet). Reference specific MBTI traits.>"},
              {"text": "<Another detailed bullet point (2-3 bullets total per sub-category).>"}
            ]
          },
          {
            "title": "<Sub-category 2 title, e.g., 'Potential Misunderstandings' or 'Challenges'>",
            "bullets": [
              {"text": "<Detailed bullet point.>"},
              {"text": "<Another detailed bullet point.>"},
              {"text": "<Third detailed bullet point.>"}
            ]
          },
          {
            "title": "<Sub-category 3 title, e.g., 'Tips for Better Communication' or 'Growth Opportunities'>",
            "bullets": [
              {"text": "<Detailed bullet point.>"},
              {"text": "<Another detailed bullet point.>"},
              {"text": "<Third detailed bullet point.>"}
            ]
          }
        ]
      },
      {
        "heading": "<Section 2 heading, e.g., 'Strengths & Synergies'>",
        "subcategories": [
          {
            "title": "<Sub-category title, e.g., 'What Makes Them Great Together'>",
            "bullets": [
              {"text": "<2-3 detailed bullets describing complementary traits, shared values, etc.>"}
            ]
          }
        ]
      },
      {
        "heading": "<Section 3 heading, e.g., 'Growth Opportunities & Challenges'>",
        "subcategories": [
          {
            "title": "<Sub-category title>",
            "bullets": [
              {"text": "<2-3 detailed bullets>"}
            ]
          }
        ]
      }
    ]
  },
  "coworker_explanation": {
    "sections": [
      {
        "heading": "<Section 1 heading, e.g., 'Work Style Compatibility'>",
        "subcategories": [
          {
            "title": "<Sub-category title, e.g., 'Complementary Skills' or 'Strengths'>",
            "bullets": [
              {"text": "<3-5 detailed bullets on work styles, problem-solving approaches, professional strengths.>"}
            ]
          },
          {
            "title": "<Sub-category 2 title, e.g., 'Potential Friction Points' or 'Challenges'>",
            "bullets": [
              {"text": "<2-3 detailed bullets>"}
            ]
          },
          {
            "title": "<Sub-category 3 title, e.g., 'Collaboration Tips' or 'Professional Growth'>",
            "bullets": [
              {"text": "<2-3 detailed bullets>"}
            ]
          }
        ]
      },
      {
        "heading": "<Section 2 heading, e.g., 'Collaboration Potential'>",
        "subcategories": [
          {
            "title": "<Sub-category title>",
            "bullets": [
              {"text": "<2-3 detailed bullets>"}
            ]
          }
        ]
      },
      {
        "heading": "<Section 3 heading, e.g., 'Professional Development & Considerations'>",
        "subcategories": [
          {
            "title": "<Sub-category title>",
            "bullets": [
              {"text": "<2-3 detailed bullets>"}
            ]
          }
        ]
      }
    ]
  },
  "partner_explanation": {
    "sections": [
      {
        "heading": "<Section 1 heading, e.g., 'Romantic Chemistry & Emotional Connection'>",
        "subcategories": [
          {
            "title": "<Sub-category title, e.g., 'What Draws Them Together' (context-specific)>",
            "bullets": [
              {"text": "<3-5 detailed bullets on romantic compatibility, emotional intimacy, what attracts them to each other.>"}
            ]
          },
          {
            "title": "<Sub-category 2 title, e.g., 'Communication Needs' (context-specific)>",
            "bullets": [
              {"text": "<2-3 detailed bullets>"}
            ]
          },
          {
            "title": "<Sub-category 3 title, e.g., 'Success Strategies' (context-specific)>",
            "bullets": [
              {"text": "<2-3 detailed bullets>"}
            ]
          }
        ]
      },
      {
        "heading": "<Section 2 heading, e.g., 'Relationship Strengths & Values Alignment'>",
        "subcategories": [
          {
            "title": "<Sub-category title>",
            "bullets": [
              {"text": "<2-3 detailed bullets>"}
            ]
          }
        ]
      },
      {
        "heading": "<Section 3 heading, e.g., 'Long-term Potential & Growth Together'>",
        "subcategories": [
          {
            "title": "<Sub-category title>",
            "bullets": [
              {"text": "<2-3 detailed bullets>"}
            ]
          }
        ]
      }
    ]
  }
}

Scoring guidelines:
- 5: Exceptional compatibility with natural synergy and minimal friction
- 4: Strong compatibility with minor areas requiring attention
- 3: Moderate compatibility with some differences that need conscious effort
- 2: Challenging compatibility requiring significant compromise and understanding
- 1: Poor compatibility with fundamental conflicts that are difficult to overcome

Be honest, insightful, and provide genuine value. Reference specific MBTI personality traits when relevant. Make each section detailed and actionable.

Return ONLY the raw JSON object, nothing else.`

	return prompt
}

// buildCategoryPrompt generates a prompt for a single compatibility category
func buildCategoryPrompt(person1, person2 models.PersonData, category string) string {
	categoryContext := ""

	switch category {
	case "friend":
		categoryContext = "as friends"
	case "coworker":
		categoryContext = "as coworkers"
	case "partner":
		categoryContext = "as partners in a romantic relationship"
	default:
		categoryContext = "in general"
	}

	prompt := fmt.Sprintf(`You are a compatibility assessment expert. Analyze the compatibility between two people %s based on ALL the information provided below. You MUST consider and reference their names and MBTI types when making your assessment.

PERSON 1:
- Name: %s
- MBTI Type: %s`, categoryContext, person1.Name, person1.MBTI)

	prompt += fmt.Sprintf(`

PERSON 2:
- Name: %s
- MBTI Type: %s`, person2.Name, person2.MBTI)

	// Category-specific instructions
	categoryInstructions := ""
	switch category {
	case "friend":
		categoryInstructions = `
You are an expert compatibility analyst with deep knowledge of personality psychology, friendship dynamics, and interpersonal communication. Focus specifically on how these two people would interact as FRIENDS. Consider:
- Communication styles and preferences
- Shared interests and activities
- Emotional support and understanding
- Potential conflicts and how they might resolve them
- Complementary personality traits that make them great friends
- Challenges they might face in the friendship`
	case "coworker":
		categoryInstructions = `
You are an expert compatibility analyst with deep knowledge of personality psychology, workplace dynamics, and professional collaboration. Focus specifically on how these two people would interact as COWORKERS. Consider:
- Work styles and approaches to tasks
- Communication in professional settings
- Collaboration and teamwork potential
- Problem-solving approaches
- Complementary professional skills
- Potential workplace conflicts and how they might handle them`
	case "partner":
		categoryInstructions = `
You are an expert compatibility analyst with deep knowledge of personality psychology, romantic relationship dynamics, and emotional intimacy. Focus specifically on how these two people would interact as ROMANTIC PARTNERS. Consider:
- Romantic chemistry and emotional connection
- Communication needs and styles in relationships
- Shared values and life goals
- Intimacy and emotional support
- Conflict resolution in romantic relationships
- Long-term relationship potential`
	}

	prompt += categoryInstructions

	prompt += `

CRITICAL INSTRUCTIONS:
- You MUST reference and incorporate MBTI types and names in your analysis
- Focus entirely on the provided MBTI information to deliver the most accurate assessment

Provide a structured analysis with AT LEAST 3 distinct sections. Each section should have:
- A clear heading
- 2-3 sub-categories with descriptive titles
- Each sub-category should contain 2-3 bullet points (detailed, not just one word)

Return your response as a JSON object with the following EXACT structure (no markdown, no code blocks):
{
  "score": <integer 1-5>,
  "explanation": {
    "sections": [
      {
        "heading": "<Section 1 heading, e.g., 'Cognitive Compatibility & Communication'>",
        "subcategories": [
          {
            "title": "<Sub-category title, e.g., 'Communication Styles' or 'Strengths'>",
            "bullets": [
              {"text": "<Detailed bullet point (2-3 sentences worth of information per bullet). Reference specific MBTI traits.>"},
              {"text": "<Another detailed bullet point (2-3 bullets total per sub-category).>"}
            ]
          },
          {
            "title": "<Sub-category 2 title, e.g., 'Potential Misunderstandings' or 'Challenges'>",
            "bullets": [
              {"text": "<Detailed bullet point.>"},
              {"text": "<Another detailed bullet point.>"},
              {"text": "<Third detailed bullet point.>"}
            ]
          },
          {
            "title": "<Sub-category 3 title, e.g., 'Tips for Better Communication' or 'Growth Opportunities'>",
            "bullets": [
              {"text": "<Detailed bullet point.>"},
              {"text": "<Another detailed bullet point.>"},
              {"text": "<Third detailed bullet point.>"}
            ]
          }
        ]
      },
      {
        "heading": "<Section 2 heading, e.g., 'Strengths & Synergies'>",
        "subcategories": [
          {
            "title": "<Sub-category title, e.g., 'What Makes Them Great Together'>",
            "bullets": [
              {"text": "<2-3 detailed bullets describing complementary traits, shared values, etc.>"}
            ]
          }
        ]
      },
      {
        "heading": "<Section 3 heading, e.g., 'Growth Opportunities & Challenges'>",
        "subcategories": [
          {
            "title": "<Sub-category title>",
            "bullets": [
              {"text": "<2-3 detailed bullets>"}
            ]
          }
        ]
      }
    ]
  }
}

Scoring guidelines:
- 5: Exceptional compatibility with natural synergy and minimal friction
- 4: Strong compatibility with minor areas requiring attention
- 3: Moderate compatibility with some differences that need conscious effort
- 2: Challenging compatibility requiring significant compromise and understanding
- 1: Poor compatibility with fundamental conflicts that are difficult to overcome

Be honest, insightful, and provide genuine value. Reference specific MBTI personality traits when relevant. Make each section detailed and actionable.

Return ONLY the raw JSON object, nothing else.`

	return prompt
}

// buildCategoryPromptWithBase generates a prompt for a single compatibility category with base explanation for augmentation
func buildCategoryPromptWithBase(person1, person2 models.PersonData, category string, baseExplanation *models.CategoryExplanation) string {
	categoryContext := ""

	switch category {
	case "friend":
		categoryContext = "as friends"
	case "coworker":
		categoryContext = "as coworkers"
	case "partner":
		categoryContext = "as partners in a romantic relationship"
	default:
		categoryContext = "in general"
	}

	prompt := fmt.Sprintf(`You are a compatibility assessment expert. Analyze the compatibility between two people %s based on ALL the information provided below. You MUST consider and reference their names and MBTI types when making your assessment.`, categoryContext)

	// If base explanation is provided, add augmentation instructions
	if baseExplanation != nil {
		baseJSON, err := json.MarshalIndent(baseExplanation, "", "  ")
		if err == nil {
			prompt += fmt.Sprintf(`

IMPORTANT: Below is the BASE compatibility assessment for these MBTI types (based on MBTI compatibility alone):

%s

Your task is to ENHANCE and AUGMENT this base assessment to reflect the specific MBTI pairing of these two individuals.
- Use the base assessment as a foundation
- Add new insights and details grounded in the provided MBTI information
- Maintain a similar structure but expand with context-specific details
- If the MBTI information doesn't suggest new insights, enhance the base assessment with more depth
`, string(baseJSON))
		}
	}

	prompt += fmt.Sprintf(`

PERSON 1:
- Name: %s
- MBTI Type: %s`, person1.Name, person1.MBTI)

	prompt += fmt.Sprintf(`

PERSON 2:
- Name: %s
- MBTI Type: %s`, person2.Name, person2.MBTI)

	// Category-specific instructions
	categoryInstructions := ""
	switch category {
	case "friend":
		categoryInstructions = `
You are an expert compatibility analyst with deep knowledge of personality psychology, friendship dynamics, and interpersonal communication. Focus specifically on how these two people would interact as FRIENDS. Consider:
- Communication styles and preferences
- Shared interests and activities
- Emotional support and understanding
- Potential conflicts and how they might resolve them
- Complementary personality traits that make them great friends
- Challenges they might face in the friendship`
	case "coworker":
		categoryInstructions = `
You are an expert compatibility analyst with deep knowledge of personality psychology, workplace dynamics, and professional collaboration. Focus specifically on how these two people would interact as COWORKERS. Consider:
- Work styles and approaches to tasks
- Communication in professional settings
- Collaboration and teamwork potential
- Problem-solving approaches
- Complementary professional skills
- Potential workplace conflicts and how they might handle them`
	case "partner":
		categoryInstructions = `
You are an expert compatibility analyst with deep knowledge of personality psychology, romantic relationship dynamics, and emotional intimacy. Focus specifically on how these two people would interact as ROMANTIC PARTNERS. Consider:
- Romantic chemistry and emotional connection
- Communication needs and styles in relationships
- Shared values and life goals
- Intimacy and emotional support
- Conflict resolution in romantic relationships
- Long-term relationship potential`
	}

	prompt += categoryInstructions

	if baseExplanation != nil {
		prompt += `

CRITICAL AUGMENTATION INSTRUCTIONS:
- You have been provided with a base MBTI compatibility assessment above
- Your task is to ENHANCE this assessment by incorporating the MBTI information provided for each person
- Add depth and detail that contextualizes the MBTI pairing for this specific request
- Maintain the structure of the base assessment but expand it with context-specific information
`
	} else {
		prompt += `

CRITICAL INSTRUCTIONS:
- You MUST reference and incorporate MBTI types and names in your analysis
- Focus entirely on the provided MBTI information to deliver the most accurate assessment
`
	}

	prompt += `
Provide a structured analysis with AT LEAST 3 distinct sections. Each section should have:
- A clear heading
- 2-3 sub-categories with descriptive titles
- Each sub-category should contain 2-3 bullet points (detailed, not just one word)

Return your response as a JSON object with the following EXACT structure (no markdown, no code blocks):
{
  "score": <integer 1-5>,
  "explanation": {
    "sections": [
      {
        "heading": "<Section 1 heading, e.g., 'Cognitive Compatibility & Communication'>",
        "subcategories": [
          {
            "title": "<Sub-category title, e.g., 'Communication Styles' or 'Strengths'>",
            "bullets": [
              {"text": "<Detailed bullet point (2-3 sentences worth of information per bullet). Reference specific MBTI traits.>"},
              {"text": "<Another detailed bullet point (2-3 bullets total per sub-category).>"}
            ]
          },
          {
            "title": "<Sub-category 2 title, e.g., 'Potential Misunderstandings' or 'Challenges'>",
            "bullets": [
              {"text": "<Detailed bullet point.>"},
              {"text": "<Another detailed bullet point.>"},
              {"text": "<Third detailed bullet point.>"}
            ]
          },
          {
            "title": "<Sub-category 3 title, e.g., 'Tips for Better Communication' or 'Growth Opportunities'>",
            "bullets": [
              {"text": "<Detailed bullet point.>"},
              {"text": "<Another detailed bullet point.>"},
              {"text": "<Third detailed bullet point.>"}
            ]
          }
        ]
      },
      {
        "heading": "<Section 2 heading, e.g., 'Strengths & Synergies'>",
        "subcategories": [
          {
            "title": "<Sub-category title, e.g., 'What Makes Them Great Together'>",
            "bullets": [
              {"text": "<2-3 detailed bullets describing complementary traits, shared values, etc.>"}
            ]
          }
        ]
      },
      {
        "heading": "<Section 3 heading, e.g., 'Growth Opportunities & Challenges'>",
        "subcategories": [
          {
            "title": "<Sub-category title>",
            "bullets": [
              {"text": "<2-3 detailed bullets>"}
            ]
          }
        ]
      }
    ]
  }
}

Scoring guidelines:
- 5: Exceptional compatibility with natural synergy and minimal friction
- 4: Strong compatibility with minor areas requiring attention
- 3: Moderate compatibility with some differences that need conscious effort
- 2: Challenging compatibility requiring significant compromise and understanding
- 1: Poor compatibility with fundamental conflicts that are difficult to overcome

Be honest, insightful, and provide genuine value. Reference specific MBTI personality traits when relevant. Make each section detailed and actionable.

Return ONLY the raw JSON object, nothing else.`

	return prompt
}

func callGeminiAPI(prompt string) ([]byte, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{
						"text": prompt,
					},
				},
			},
		},
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1/models/gemini-2.0-flash:generateContent?key=%s", apiKey)
	client := &http.Client{}

	maxAttempts := 3
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadJSON))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to make request: %w", err)
		}

		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			return nil, fmt.Errorf("failed to read response: %w", readErr)
		}

		if resp.StatusCode == http.StatusOK {
			return body, nil
		}

		if (resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable) && attempt < maxAttempts {
			wait := time.Duration(1<<uint(attempt-1)) * time.Second
			time.Sleep(wait)
			continue
		}

		return nil, fmt.Errorf("API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	return nil, fmt.Errorf("API error: exhausted retries")
}

func extractJSON(text string) string {
	// Try to find JSON in the text
	startIdx := -1
	endIdx := -1
	braceCount := 0

	for i, char := range text {
		if char == '{' {
			if startIdx == -1 {
				startIdx = i
			}
			braceCount++
		} else if char == '}' {
			braceCount--
			if braceCount == 0 && startIdx != -1 {
				endIdx = i + 1
				break
			}
		}
	}

	var jsonText string
	if startIdx != -1 && endIdx != -1 {
		jsonText = text[startIdx:endIdx]
	} else {
		// If no JSON found, return the whole text (might be just JSON)
		jsonText = text
	}

	// Fix common JSON issues: trailing commas before closing braces/brackets
	// Use regex-like approach to remove trailing commas more comprehensively

	// Remove trailing comma before } (handle various whitespace patterns)
	jsonText = strings.ReplaceAll(jsonText, ",}", "}")
	jsonText = strings.ReplaceAll(jsonText, ", }", " }")
	jsonText = strings.ReplaceAll(jsonText, ",\n}", "\n}")
	jsonText = strings.ReplaceAll(jsonText, ",\r\n}", "\r\n}")
	jsonText = strings.ReplaceAll(jsonText, ",\r}", "\r}")
	// Handle cases with spaces before comma
	jsonText = strings.ReplaceAll(jsonText, " ,}", "}")
	jsonText = strings.ReplaceAll(jsonText, " , }", " }")

	// Remove trailing comma before ] (handle various whitespace patterns)
	jsonText = strings.ReplaceAll(jsonText, ",]", "]")
	jsonText = strings.ReplaceAll(jsonText, ", ]", " ]")
	jsonText = strings.ReplaceAll(jsonText, ",\n]", "\n]")
	jsonText = strings.ReplaceAll(jsonText, ",\r\n]", "\r\n]")
	jsonText = strings.ReplaceAll(jsonText, ",\r]", "\r]")
	// Handle cases with spaces before comma
	jsonText = strings.ReplaceAll(jsonText, " ,]", "]")
	jsonText = strings.ReplaceAll(jsonText, " , ]", " ]")

	// More aggressive: remove trailing comma after last quote before closing brace
	// This handles cases like: "value",\n}
	jsonBytes := []byte(jsonText)
	result := []byte{}
	inString := false
	escapeNext := false

	for i := 0; i < len(jsonBytes); i++ {
		char := jsonBytes[i]

		if escapeNext {
			result = append(result, char)
			escapeNext = false
			continue
		}

		if char == '\\' {
			escapeNext = true
			result = append(result, char)
			continue
		}

		if char == '"' {
			inString = !inString
			result = append(result, char)
			continue
		}

		// If we're outside a string and find ",}" or ",\n}" or similar, skip the comma
		if !inString && char == ',' {
			// Look ahead to see if next non-whitespace is } or ]
			j := i + 1
			for j < len(jsonBytes) {
				nextChar := jsonBytes[j]
				if nextChar == ' ' || nextChar == '\n' || nextChar == '\r' || nextChar == '\t' {
					j++
					continue
				}
				if nextChar == '}' || nextChar == ']' {
					// Skip this comma, don't append it
					i = j - 1 // Will be incremented by loop
					break
				}
				// Not a closing brace/bracket, keep the comma
				result = append(result, char)
				break
			}
			if j >= len(jsonBytes) {
				result = append(result, char)
			}
		} else {
			result = append(result, char)
		}
	}

	return string(result)
}

// cleanJSONForParsing applies additional cleaning passes to ensure valid JSON
func cleanJSONForParsing(jsonText string) string {
	// Apply multiple cleaning passes
	cleaned := jsonText

	// Remove trailing commas more aggressively
	// Pattern: "value",\n} -> "value"\n}
	cleaned = strings.ReplaceAll(cleaned, "\",\n}", "\"\n}")
	cleaned = strings.ReplaceAll(cleaned, "\",\r\n}", "\"\r\n}")
	cleaned = strings.ReplaceAll(cleaned, "\", }", "\" }")
	cleaned = strings.ReplaceAll(cleaned, "\",}", "\"}")

	// Remove trailing commas after numbers
	cleaned = strings.ReplaceAll(cleaned, ",\n}", "\n}")
	cleaned = strings.ReplaceAll(cleaned, ",\r\n}", "\r\n}")
	cleaned = strings.ReplaceAll(cleaned, ", }", " }")
	cleaned = strings.ReplaceAll(cleaned, ",}", "}")

	// Remove trailing commas before closing bracket
	cleaned = strings.ReplaceAll(cleaned, ",\n]", "\n]")
	cleaned = strings.ReplaceAll(cleaned, ",\r\n]", "\r\n]")
	cleaned = strings.ReplaceAll(cleaned, ", ]", " ]")
	cleaned = strings.ReplaceAll(cleaned, ",]", "]")

	return cleaned
}

// convertStringToStructured converts old string format to new structured format with subcategories and bullets
func convertStringToStructured(text string, category string) models.CategoryExplanation {
	sections := []models.ExplanationSection{}
	headings := getHeadingsForCategory(category)

	// Split text into paragraphs
	paragraphs := splitIntoParagraphs(text)

	if len(paragraphs) >= 3 {
		// Use paragraphs as sections, convert each to subcategories with bullets
		for i := 0; i < 3 && i < len(headings); i++ {
			para := paragraphs[i]
			subcategories := convertParagraphToSubcategories(para, category, i)
			sections = append(sections, models.ExplanationSection{
				Heading:       headings[i],
				Subcategories: subcategories,
			})
		}
	} else if len(paragraphs) > 0 {
		// Fewer paragraphs - split content into 3 sections with subcategories
		words := splitIntoWords(text)
		wordsPerSection := len(words) / 3

		for i := 0; i < 3 && i < len(headings); i++ {
			start := i * wordsPerSection
			end := start + wordsPerSection
			if i == 2 {
				end = len(words)
			}

			if start < len(words) {
				sectionText := strings.Join(words[start:end], " ")
				subcategories := convertParagraphToSubcategories(sectionText, category, i)
				sections = append(sections, models.ExplanationSection{
					Heading:       headings[i],
					Subcategories: subcategories,
				})
			}
		}
	} else {
		// Fallback: create basic structure from full text
		subcategories := convertParagraphToSubcategories(text, category, 0)
		sections = append(sections, models.ExplanationSection{
			Heading:       headings[0],
			Subcategories: subcategories,
		})
	}

	// Ensure at least 3 sections
	if len(sections) < 3 {
		for len(sections) < 3 {
			idx := len(sections)
			if idx < len(headings) {
				subcategories := []models.SubCategory{
					{
						Title: "Additional Insights",
						Bullets: []models.BulletPoint{
							{Text: "Continue reading for more detailed analysis."},
						},
					},
				}
				sections = append(sections, models.ExplanationSection{
					Heading:       headings[idx],
					Subcategories: subcategories,
				})
			} else {
				break
			}
		}
	}

	return models.CategoryExplanation{Sections: sections}
}

// convertParagraphToSubcategories converts a paragraph into subcategories with bullet points
func convertParagraphToSubcategories(text string, category string, sectionIndex int) []models.SubCategory {
	subcategories := []models.SubCategory{}

	// Split paragraph into sentences
	sentences := splitIntoSentences(text)

	if len(sentences) >= 4 {
		// Enough sentences - create 2-3 subcategories with 2-3 bullets each
		subcatTitles := getSubcategoryTitles(category, sectionIndex)
		sentencesPerSubcat := len(sentences) / len(subcatTitles)
		if sentencesPerSubcat < 2 {
			sentencesPerSubcat = 2
		}

		for i, title := range subcatTitles {
			start := i * sentencesPerSubcat
			end := start + sentencesPerSubcat
			if i == len(subcatTitles)-1 {
				end = len(sentences)
			}

			if start < len(sentences) {
				bullets := []models.BulletPoint{}
				for _, sent := range sentences[start:end] {
					if strings.TrimSpace(sent) != "" {
						bullets = append(bullets, models.BulletPoint{Text: strings.TrimSpace(sent)})
					}
				}

				// Group bullets into 2-3 per subcategory
				if len(bullets) > 3 {
					grouped := []models.BulletPoint{}
					for i := 0; i < len(bullets); i += 2 {
						endIdx := i + 3
						if endIdx > len(bullets) {
							endIdx = len(bullets)
						}
						combined := strings.Join(func() []string {
							var texts []string
							for _, b := range bullets[i:endIdx] {
								texts = append(texts, b.Text)
							}
							return texts
						}(), " ")
						grouped = append(grouped, models.BulletPoint{Text: combined})
					}
					bullets = grouped
				}

				if len(bullets) > 0 {
					subcategories = append(subcategories, models.SubCategory{
						Title:   title,
						Bullets: bullets,
					})
				}
			}
		}
	} else {
		// Fewer sentences - create a single subcategory
		bullets := []models.BulletPoint{}
		for _, sent := range sentences {
			sent = strings.TrimSpace(sent)
			if sent != "" {
				bullets = append(bullets, models.BulletPoint{Text: sent})
			}
		}

		// Split into groups of 2-3 bullets
		if len(bullets) > 3 {
			grouped := []models.BulletPoint{}
			for i := 0; i < len(bullets); i += 2 {
				end := i + 3
				if end > len(bullets) {
					end = len(bullets)
				}
				combined := strings.Join(func() []string {
					var texts []string
					for _, b := range bullets[i:end] {
						texts = append(texts, b.Text)
					}
					return texts
				}(), " ")
				grouped = append(grouped, models.BulletPoint{Text: combined})
			}
			bullets = grouped
		}

		if len(bullets) > 0 {
			subcategories = append(subcategories, models.SubCategory{
				Title:   getDefaultSubcategoryTitle(category, sectionIndex),
				Bullets: bullets,
			})
		}
	}

	// Ensure at least one subcategory
	if len(subcategories) == 0 {
		subcategories = append(subcategories, models.SubCategory{
			Title: "Compatibility Analysis",
			Bullets: []models.BulletPoint{
				{Text: text},
			},
		})
	}

	return subcategories
}

func splitIntoSentences(text string) []string {
	// Split by periods, exclamation, question marks followed by space
	sentences := []string{}

	// Replace common sentence endings
	text = strings.ReplaceAll(text, ". ", ".<SPLIT>")
	text = strings.ReplaceAll(text, "! ", "!<SPLIT>")
	text = strings.ReplaceAll(text, "? ", "?<SPLIT>")

	parts := strings.Split(text, "<SPLIT>")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			// Ensure sentence ends with punctuation
			if !strings.HasSuffix(part, ".") && !strings.HasSuffix(part, "!") && !strings.HasSuffix(part, "?") {
				part += "."
			}
			sentences = append(sentences, part)
		}
	}

	return sentences
}

func getSubcategoryTitles(category string, sectionIndex int) []string {
	titles := map[string]map[int][]string{
		"friendship": {
			0: {"Communication Styles", "Potential Misunderstandings", "Tips for Better Communication"},
			1: {"What Makes Them Great Together", "Complementary Strengths"},
			2: {"Growth Opportunities", "Challenges to Navigate"},
		},
		"workplace": {
			0: {"Complementary Skills", "Potential Friction Points", "Collaboration Tips"},
			1: {"Team Dynamics", "Problem-Solving Approaches"},
			2: {"Professional Growth", "Considerations"},
		},
		"romance": {
			0: {"What Draws Them Together", "Communication Needs", "Success Strategies"},
			1: {"Relationship Strengths", "Values Alignment"},
			2: {"Long-term Potential", "Growth Together"},
		},
	}

	if catMap, ok := titles[category]; ok {
		if sectionTitles, ok := catMap[sectionIndex]; ok {
			return sectionTitles
		}
	}

	// Default titles
	return []string{"Strengths", "Challenges", "Growth Opportunities"}
}

func getDefaultSubcategoryTitle(category string, sectionIndex int) string {
	titles := getSubcategoryTitles(category, sectionIndex)
	if len(titles) > 0 {
		return titles[0]
	}
	return "Compatibility Analysis"
}

func splitIntoParagraphs(text string) []string {
	// Split by double newlines or periods followed by newlines
	paragraphs := []string{}

	// First, try splitting by \n\n
	parts := strings.Split(text, "\n\n")
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			paragraphs = append(paragraphs, trimmed)
		}
	}

	// If no double newlines, try splitting by periods + newline
	if len(paragraphs) <= 1 {
		paragraphs = []string{}
		parts = strings.Split(text, ".\n")
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed != "" {
				if !strings.HasSuffix(trimmed, ".") {
					trimmed += "."
				}
				paragraphs = append(paragraphs, trimmed)
			}
		}
	}

	// If still only one, split by sentences (period + space)
	if len(paragraphs) <= 1 {
		paragraphs = []string{}
		sentences := strings.Split(text, ". ")
		currentPara := ""
		for i, sent := range sentences {
			sent = strings.TrimSpace(sent)
			if sent != "" {
				if currentPara != "" {
					currentPara += " "
				}
				currentPara += sent
				if !strings.HasSuffix(sent, ".") {
					currentPara += "."
				}

				// Group 2-3 sentences per paragraph
				if (i+1)%2 == 0 || i == len(sentences)-1 {
					paragraphs = append(paragraphs, currentPara)
					currentPara = ""
				}
			}
		}
		if currentPara != "" {
			paragraphs = append(paragraphs, currentPara)
		}
	}

	return paragraphs
}

func splitIntoWords(text string) []string {
	return strings.Fields(text)
}

// convertSectionsToSubcategories converts old section format (with content string) to new format (with subcategories)
func convertSectionsToSubcategories(oldSections []struct {
	Heading string `json:"heading"`
	Content string `json:"content"`
}, category string) models.CategoryExplanation {
	sections := []models.ExplanationSection{}

	for i, oldSection := range oldSections {
		subcategories := convertParagraphToSubcategories(oldSection.Content, category, i)
		sections = append(sections, models.ExplanationSection{
			Heading:       oldSection.Heading,
			Subcategories: subcategories,
		})
	}

	// Ensure at least 3 sections
	if len(sections) < 3 {
		headings := getHeadingsForCategory(category)
		for len(sections) < 3 {
			idx := len(sections)
			if idx < len(headings) {
				subcategories := []models.SubCategory{
					{
						Title: "Additional Insights",
						Bullets: []models.BulletPoint{
							{Text: "Continue reading for more detailed analysis."},
						},
					},
				}
				sections = append(sections, models.ExplanationSection{
					Heading:       headings[idx],
					Subcategories: subcategories,
				})
			} else {
				break
			}
		}
	}

	return models.CategoryExplanation{Sections: sections}
}

func getHeadingsForCategory(category string) []string {
	switch category {
	case "friendship":
		return []string{
			"Cognitive Compatibility & Communication",
			"Strengths & Synergies",
			"Growth Opportunities & Challenges",
		}
	case "workplace":
		return []string{
			"Work Style Compatibility",
			"Collaboration Potential",
			"Professional Development & Considerations",
		}
	case "romance":
		return []string{
			"Romantic Chemistry & Emotional Connection",
			"Relationship Strengths & Values Alignment",
			"Long-term Potential & Growth Together",
		}
	default:
		return []string{
			"Compatibility Analysis",
			"Key Strengths",
			"Areas for Growth",
		}
	}
}
