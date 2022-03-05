package db

import "secret-scanner/models"

func (db *database) AddRepository(repos []*models.Repository) error {
	result := db.conn.Table("repository").Create(&repos)
	return result.Error
}

func (db *database) GetRepositories() ([]*models.Repository, error) {
	var repos []*models.Repository
	result := db.conn.Table("repository").Find(&repos)
	return repos, result.Error
}

func (db *database) GetRepositoryById(repoId string) (*models.Repository, error) {
	var repo models.Repository
	result := db.conn.Table("repository").First(&repo, repoId)
	return &repo, result.Error
}

func (db *database) UpdateRepository(repo *models.Repository) error {
	result := db.conn.Table("repository").Save(repo)
	return result.Error
}

func (db *database) DeleteRepository(repoId string) error {
	result := db.conn.Table("repository").Delete(&models.Repository{}, repoId)
	return result.Error
}
