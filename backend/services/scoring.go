package services

import (
	"compatiblah/backend/models"
	"math"
	"math/rand"
	"strings"
	"time"
)

type compatibilityScoreSet struct {
	Friend   int
	Coworker int
	Partner  int
}

func calculateCompatibilityScores(person1, person2 models.PersonData) compatibilityScoreSet {
	profile1, ok1 := parseMBTIProfile(person1.MBTI)
	profile2, ok2 := parseMBTIProfile(person2.MBTI)

	seed := time.Now().UnixNano()

	if !ok1 || !ok2 {
		// fallback to neutral scores with slight variation
		return compatibilityScoreSet{
			Friend:   clampScore(applyNoise(3.0, seed+1)),
			Coworker: clampScore(applyNoise(3.0, seed+2)),
			Partner:  clampScore(applyNoise(3.0, seed+3)),
		}
	}

	return compatibilityScoreSet{
		Friend:   computeCategoryScore(profile1, profile2, "friend", seed+1),
		Coworker: computeCategoryScore(profile1, profile2, "coworker", seed+2),
		Partner:  computeCategoryScore(profile1, profile2, "partner", seed+3),
	}
}

func calculateCategoryScore(person1, person2 models.PersonData, category string) int {
	profile1, ok1 := parseMBTIProfile(person1.MBTI)
	profile2, ok2 := parseMBTIProfile(person2.MBTI)
	seed := time.Now().UnixNano()

	if !ok1 || !ok2 {
		return clampScore(applyNoise(3.0, seed))
	}

	return computeCategoryScore(profile1, profile2, category, seed)
}

func blendScores(geminiScore, heuristicScore int) int {
	base := 0.35*float64(geminiScore) + 0.65*float64(heuristicScore)
	return clampScore(base)
}

type mbtiProfile struct {
	energy    rune
	intuition rune
	decision  rune
	lifestyle rune
}

func parseMBTIProfile(raw string) (mbtiProfile, bool) {
	value := strings.ToUpper(strings.TrimSpace(raw))
	if len(value) != 4 {
		return mbtiProfile{}, false
	}

	profile := mbtiProfile{
		energy:    rune(value[0]),
		intuition: rune(value[1]),
		decision:  rune(value[2]),
		lifestyle: rune(value[3]),
	}

	valid := (profile.energy == 'I' || profile.energy == 'E') &&
		(profile.intuition == 'N' || profile.intuition == 'S') &&
		(profile.decision == 'T' || profile.decision == 'F') &&
		(profile.lifestyle == 'J' || profile.lifestyle == 'P')

	return profile, valid
}

func computeCategoryScore(profile1, profile2 mbtiProfile, category string, seed int64) int {
	base := 3.0

	switch category {
	case "friend":
		base += friendEnergyAdjustment(profile1.energy, profile2.energy)
		base += friendIntuitionAdjustment(profile1.intuition, profile2.intuition)
		base += friendDecisionAdjustment(profile1.decision, profile2.decision)
		base += friendLifestyleAdjustment(profile1.lifestyle, profile2.lifestyle)
	case "coworker":
		base += coworkerEnergyAdjustment(profile1.energy, profile2.energy)
		base += coworkerIntuitionAdjustment(profile1.intuition, profile2.intuition)
		base += coworkerDecisionAdjustment(profile1.decision, profile2.decision)
		base += coworkerLifestyleAdjustment(profile1.lifestyle, profile2.lifestyle)
	case "partner":
		base += partnerEnergyAdjustment(profile1.energy, profile2.energy)
		base += partnerIntuitionAdjustment(profile1.intuition, profile2.intuition)
		base += partnerDecisionAdjustment(profile1.decision, profile2.decision)
		base += partnerLifestyleAdjustment(profile1.lifestyle, profile2.lifestyle)
	default:
		base += friendEnergyAdjustment(profile1.energy, profile2.energy)
	}

	return clampScore(applyNoise(base, seed))
}

func friendEnergyAdjustment(a, b rune) float64 {
	if a == b {
		return 0.4
	}
	return -0.2
}

func friendIntuitionAdjustment(a, b rune) float64 {
	if a == b {
		if a == 'N' {
			return 0.3
		}
		return 0.2
	}
	return -0.1
}

func friendDecisionAdjustment(a, b rune) float64 {
	if a == b {
		if a == 'F' {
			return 0.4
		}
		return 0.2
	}
	return 0.1
}

func friendLifestyleAdjustment(a, b rune) float64 {
	if a == b {
		return 0.2
	}
	return 0.1
}

func coworkerEnergyAdjustment(a, b rune) float64 {
	if a == b {
		return 0.1
	}
	return 0.2
}

func coworkerIntuitionAdjustment(a, b rune) float64 {
	if a == b {
		return 0.1
	}
	return 0.2
}

func coworkerDecisionAdjustment(a, b rune) float64 {
	if a == b {
		if a == 'T' {
			return 0.4
		}
		return 0.2
	}
	return 0.2
}

func coworkerLifestyleAdjustment(a, b rune) float64 {
	if a == b {
		if a == 'J' {
			return 0.4
		}
		return 0.1
	}
	return 0.1
}

func partnerEnergyAdjustment(a, b rune) float64 {
	if a == b {
		return -0.1
	}
	return 0.4
}

func partnerIntuitionAdjustment(a, b rune) float64 {
	if a == b {
		if a == 'N' {
			return 0.3
		}
		return 0.2
	}
	return 0.1
}

func partnerDecisionAdjustment(a, b rune) float64 {
	if a == b {
		if a == 'F' {
			return 0.5
		}
		return 0.1
	}
	return 0.2
}

func partnerLifestyleAdjustment(a, b rune) float64 {
	if a == b {
		return 0.1
	}
	return 0.2
}

func applyNoise(value float64, seed int64) float64 {
	rng := rand.New(rand.NewSource(seed))
	noise := (rng.Float64() * 0.6) - 0.3 // [-0.3, 0.3]
	return value + noise
}

func clampScore(value float64) int {
	rounded := int(math.Round(value))
	if rounded < 1 {
		return 1
	}
	if rounded > 5 {
		return 5
	}
	return rounded
}
