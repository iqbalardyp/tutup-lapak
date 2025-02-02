package repository

import (
	"context"

	"tutup-lapak/internal/file/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FileRepository struct {
	pool *pgxpool.Pool
}

func NewFileRepository(pool *pgxpool.Pool) *FileRepository {
	return &FileRepository{pool: pool}
}

const insertFileQuery = `-- name: InsertFile :one
INSERT INTO files (uri, thumbnail_uri) VALUES ($1, $2) RETURNING id, uri, thumbnail_uri
`

type InsertFileParams struct {
	URI          string
	ThumbnailURI string
}

func (r *FileRepository) InsertFile(ctx context.Context, arg InsertFileParams) (model.File, error) {
	row := r.pool.QueryRow(ctx, insertFileQuery, arg.URI, arg.ThumbnailURI)

	var file model.File
	err := row.Scan(
		&file.ID,
		&file.URI,
		&file.ThumbnailURI,
	)
	return file, err
}
