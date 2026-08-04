package storage

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/port"
	"github.com/google/uuid"
)

var _ port.IAvatarStorage = (*AvatarFS)(nil)

var extByContentType = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

var contentTypeByExt = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".webp": "image/webp",
}

type AvatarFS struct {
	dir string
}

func NewAvatarFS(dir string) (*AvatarFS, error) {
	if dir == "" {
		dir = "./uploads/avatars"
	}

	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(abs, 0o755); err != nil {
		return nil, err
	}

	return &AvatarFS{dir: abs}, nil
}

func (s *AvatarFS) Save(_ context.Context, userID uuid.UUID, contentType string, data []byte) (string, error) {
	ext, ok := extByContentType[contentType]
	if !ok {
		return "", entity.NewBadRequest("unsupported content type")
	}

	filename := userID.String() + ext
	fullPath := filepath.Join(s.dir, filename)

	tmp, err := os.CreateTemp(s.dir, userID.String()+".*.tmp")
	if err != nil {
		return "", entity.NewInternal(err)
	}
	tmpName := tmp.Name()
	defer func() {
		_ = os.Remove(tmpName)
	}()

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return "", entity.NewInternal(err)
	}
	if err := tmp.Close(); err != nil {
		return "", entity.NewInternal(err)
	}

	if err := os.Rename(tmpName, fullPath); err != nil {
		return "", entity.NewInternal(err)
	}

	return filename, nil
}

func (s *AvatarFS) Load(_ context.Context, path string) ([]byte, string, error) {
	filename, err := s.safeFilename(path)
	if err != nil {
		return nil, "", err
	}

	fullPath := filepath.Join(s.dir, filename)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, "", entity.NewNotFound("avatar not found")
		}
		return nil, "", entity.NewInternal(err)
	}

	contentType, ok := contentTypeByExt[strings.ToLower(filepath.Ext(filename))]
	if !ok {
		return nil, "", entity.NewInternal(fmt.Errorf("unknown avatar extension: %s", filepath.Ext(filename)))
	}

	return data, contentType, nil
}

func (s *AvatarFS) Delete(_ context.Context, path string) error {
	if path == "" {
		return nil
	}

	filename, err := s.safeFilename(path)
	if err != nil {
		return err
	}

	fullPath := filepath.Join(s.dir, filename)
	if err := os.Remove(fullPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return entity.NewInternal(err)
	}

	return nil
}

func (s *AvatarFS) safeFilename(path string) (string, error) {
	filename := filepath.Base(path)
	if filename == "." || filename == ".." || filename == string(filepath.Separator) {
		return "", entity.NewBadRequest("invalid avatar path")
	}
	if strings.Contains(filename, string(filepath.Separator)) || strings.Contains(filename, "..") {
		return "", entity.NewBadRequest("invalid avatar path")
	}
	return filename, nil
}
