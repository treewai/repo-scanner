package db

import (
	"secret-scanner/models"
)

// AddResult adds the given result to database
func (db *database) AddResult(res *models.Result) (*models.Result, error) {
	result := db.conn.Table("result").Create(&res)
	return res, result.Error
}

// GetResult gets all results from database
func (db *database) GetResults() ([]*models.Result, error) {
	var results []*models.Result
	result := db.conn.Table("result").Select(
		"result.id", "result.status", "repository.name as repo_name", "repository.url as repo_url", "result.findings",
		"result.queued_at", "result.scanning_at", "result.finished_at",
	).Joins("left join repository on repository.id = result.repo_id").Scan(&results)
	return results, result.Error
}

// UpdateResult updates the given result to database
func (db *database) UpdateResult(res *models.Result) error {
	result := db.conn.Table("result").Save(res)
	return result.Error
}
