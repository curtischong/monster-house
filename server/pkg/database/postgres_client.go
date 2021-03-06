package database

import (
	"fmt"

	"../common"
	"../config"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type PostgresClient struct {
	config config.DatabaseConfig
	db     *sqlx.DB
}

func NewPostgresClient(
	config *config.Config,
) *PostgresClient {
	client := PostgresClient{
		config: config.DatabaseConfig,
	}
	db, err := client.getDBClient()
	if err != nil {
		log.Fatalf("cannot instantiate DB client, err=%s", err)
	}
	// We want to share one DB client across all queries.
	client.db = db
	return &client
}

func (client *PostgresClient) QueryAllPhotoIDs() (allIDs []uuid.UUID, err error) {
	query := `SELECT id FROM photo`

	tx := client.db.MustBegin()
	rows, err := tx.Query(query)
	if err != nil {
		return
	}
	allIDs = make([]uuid.UUID, 0)
	for rows.Next() {
		var photoID uuid.UUID
		err = rows.Scan(&photoID)
		if err != nil {
			return
		}
		allIDs = append(allIDs, photoID)
	}
	err = tx.Commit()
	return
}

func (client *PostgresClient) QueryAllTagsForPhoto(
	photoID uuid.UUID,
) (allTags []common.TagResponseData, err error) {
	query := `
	WITH tag_ids AS (
		SELECT tag_id, is_auto_generated FROM photo_tag
		WHERE photo_id=$1
	)
	SELECT name, is_auto_generated FROM tag_ids JOIN tag ON (tag_ids.tag_id=tag.id);
	`

	tx := client.db.MustBegin()
	rows, err := tx.Query(query, photoID)
	if err != nil {
		return
	}
	allTags = make([]common.TagResponseData, 0)
	for rows.Next() {
		var tagName string
		var isGenerated bool
		err = rows.Scan(&tagName, &isGenerated)
		if err != nil {
			return nil, err
		}
		tagResponseData := common.TagResponseData{
			Name:        tagName,
			IsGenerated: isGenerated,
		}
		allTags = append(allTags, tagResponseData)
	}
	err = tx.Commit()
	return
}

func (client *PostgresClient) QueryAllPhotosWithTag(
	tag string,
) (photos []uuid.UUID, err error) {
	query := `
	WITH tag_ids AS (
		SELECT id FROM tag
		WHERE name LIKE $1
	)
	SELECT photo_id FROM tag_ids JOIN photo_tag ON (tag_ids.id = photo_tag.tag_id);
	`

	tx := client.db.MustBegin()
	rows, err := tx.Query(query, tag)
	if err != nil {
		return
	}
	photos = make([]uuid.UUID, 0)
	for rows.Next() {
		var photoID uuid.UUID
		err = rows.Scan(&photoID)
		if err != nil {
			return nil, err
		}
		photos = append(photos, photoID)
	}
	err = tx.Commit()
	return
}

func (client *PostgresClient) InsertPhoto(
	photoID uuid.UUID,
	photoName, extensionName string,
) (err error) {
	tx := client.db.MustBegin()
	_, err = tx.Exec("INSERT INTO photo (id, name, extension) VALUES ($1, $2, $3);",
		photoID, photoName, extensionName)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

func (client *PostgresClient) InsertTagIfNotExist(
	tagName string,
) (tagID uuid.UUID, err error) {
	// Query from https://stackoverflow.com/a/62205017/4647924
	query := `WITH e AS(
		INSERT INTO tag (id, "name")
	VALUES ($1, $2)
	ON CONFLICT DO NOTHING
	RETURNING id
	)
	SELECT * FROM e
	UNION
	SELECT id FROM tag WHERE name=$2;
	`

	tagID = uuid.New()

	tx := client.db.MustBegin()
	err = tx.QueryRow(query, tagID, tagName).Scan(&tagID)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

func (client *PostgresClient) InsertPhotoTags(
	photoID uuid.UUID,
	tagIDs []uuid.UUID,
	isAutoGenerated bool,
) (err error) {
	if len(tagIDs) == 0 {
		return nil
	}

	query := "INSERT INTO photo_tag (photo_id, tag_id, is_auto_generated) VALUES"
	params := []interface{}{}
	// Build the query to built insert all the IDs
	for i, tagID := range tagIDs {
		tagIdx := i * 3
		query += fmt.Sprintf("($%d,$%d,$%d),", tagIdx+1, tagIdx+2, tagIdx+3)
		params = append(params, photoID, tagID, isAutoGenerated)
	}
	query = query[:len(query)-1] // remove trailing ","

	tx := client.db.MustBegin()
	_, err = tx.Exec(query, params...)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

func (client *PostgresClient) getDBClient() (*sqlx.DB, error) {
	dataSourceName := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		client.config.Username, client.config.Password, client.config.DBName)
	db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to database. err=%s", err)
	}
	return db, nil
}
