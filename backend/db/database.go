package db

import (
	"compatiblah/backend/models"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	// Check if old schema exists and migrate if needed
	if err := migrateIfNeeded(); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS assessments (
		id TEXT PRIMARY KEY,
		friend_score INTEGER NOT NULL,
		coworker_score INTEGER NOT NULL,
		partner_score INTEGER NOT NULL,
		overall_score INTEGER NOT NULL,
		friend_explanation TEXT NOT NULL,
		coworker_explanation TEXT NOT NULL,
		partner_explanation TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = DB.Exec(createTableQuery)
	if err != nil {
		return err
	}

	return nil
}

func migrateIfNeeded() error {
	// Check if old schema exists (has person1_name column)
	var exists bool
	err := DB.QueryRow(`
		SELECT COUNT(*) > 0 
		FROM sqlite_master 
		WHERE type='table' AND name='assessments'
	`).Scan(&exists)

	if err != nil || !exists {
		return nil // Table doesn't exist yet, no migration needed
	}

	// Check if old columns exist
	var hasOldColumns bool
	err = DB.QueryRow(`
		SELECT COUNT(*) > 0 
		FROM pragma_table_info('assessments') 
		WHERE name='person1_name'
	`).Scan(&hasOldColumns)

	if err != nil || !hasOldColumns {
		return nil // Already migrated or new schema
	}

	// Migration: Create new table, migrate data, drop old table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS assessments_new (
			id TEXT PRIMARY KEY,
			friend_score INTEGER NOT NULL,
			coworker_score INTEGER NOT NULL,
			partner_score INTEGER NOT NULL,
			overall_score INTEGER NOT NULL,
			friend_explanation TEXT NOT NULL,
			coworker_explanation TEXT NOT NULL,
			partner_explanation TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return err
	}

	// Migrate only assessment results (not personal data)
	_, err = DB.Exec(`
		INSERT INTO assessments_new 
		(id, friend_score, coworker_score, partner_score, overall_score,
		 friend_explanation, coworker_explanation, partner_explanation, created_at)
		SELECT 
		id, friend_score, coworker_score, partner_score, overall_score,
		friend_explanation, coworker_explanation, partner_explanation, created_at
		FROM assessments;
	`)
	if err != nil {
		return err
	}

	// Replace old table with new one
	_, err = DB.Exec(`DROP TABLE assessments;`)
	if err != nil {
		return err
	}

	_, err = DB.Exec(`ALTER TABLE assessments_new RENAME TO assessments;`)
	return err
}

func SaveAssessment(assessment *models.Assessment) error {
	// Only save assessment results, NOT personal data (privacy-first approach)
	query := `
	INSERT INTO assessments (
		id, friend_score, coworker_score, partner_score, overall_score,
		friend_explanation, coworker_explanation, partner_explanation, created_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := DB.Exec(query,
		assessment.ID,
		assessment.FriendScore,
		assessment.CoworkerScore,
		assessment.PartnerScore,
		assessment.OverallScore,
		assessment.FriendExplanation,
		assessment.CoworkerExplanation,
		assessment.PartnerExplanation,
		assessment.CreatedAt,
	)

	return err
}

func GetAssessment(id string) (*models.Assessment, error) {
	// Only retrieve assessment results, NOT personal data (privacy-first approach)
	query := `
	SELECT id, friend_score, coworker_score, partner_score, overall_score,
		   friend_explanation, coworker_explanation, partner_explanation, created_at
	FROM assessments
	WHERE id = ?
	`

	var assessment models.Assessment
	var createdAt string

	err := DB.QueryRow(query, id).Scan(
		&assessment.ID,
		&assessment.FriendScore,
		&assessment.CoworkerScore,
		&assessment.PartnerScore,
		&assessment.OverallScore,
		&assessment.FriendExplanation,
		&assessment.CoworkerExplanation,
		&assessment.PartnerExplanation,
		&createdAt,
	)

	if err != nil {
		return nil, err
	}

	assessment.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
	if err != nil {
		// Try alternative format
		assessment.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
		if err != nil {
			assessment.CreatedAt = time.Now()
		}
	}

	return &assessment, nil
}

func GetAllAssessments() ([]*models.Assessment, error) {
	// Privacy-first: only return assessment results, NOT personal data
	query := `
	SELECT id, overall_score, created_at
	FROM assessments
	ORDER BY created_at DESC
	LIMIT 100
	`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assessments []*models.Assessment
	for rows.Next() {
		var assessment models.Assessment
		var createdAt string

		err := rows.Scan(
			&assessment.ID,
			&assessment.OverallScore,
			&createdAt,
		)
		if err != nil {
			continue
		}

		var parseErr error
		assessment.CreatedAt, parseErr = time.Parse("2006-01-02 15:04:05", createdAt)
		if parseErr != nil {
			assessment.CreatedAt, parseErr = time.Parse(time.RFC3339, createdAt)
			if parseErr != nil {
				assessment.CreatedAt = time.Now()
			}
		}

		assessments = append(assessments, &assessment)
	}

	return assessments, rows.Err()
}
