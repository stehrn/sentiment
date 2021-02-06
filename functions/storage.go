package functions

import (
	"fmt"
	"path/filepath"
	"strings"
)

// GCSEvent is the payload of a Google Cloud Storage (GCS) event.
type GCSEvent struct {
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
}

// StorageImage holds user ID and image name, its derived from GCSEvent.Name
type StorageImage struct {
	UserID    string
	ImageName string
}

// ToStorageImage converts a GCSEvent.Name to a StorageImage
// assumes format of GCSEvent.Name is [UserID]/[ImageName]
func ToStorageImage(event GCSEvent) (*StorageImage, error) {
	// make sure object is of expected format
	split := strings.Split(event.Name, "/")
	if len(split) != 2 {
		return nil, fmt.Errorf("Expected object of form: [userId]/[imageName], but got: %v", event.Name)
	}
	storageImage := &StorageImage{
		UserID:    split[0],
		ImageName: split[1],
	}
	return storageImage, nil
}

// IsThumbNail whether image is thumbnail (starts with 'thumb_')
func (s *StorageImage) IsThumbNail() bool {
	return strings.HasPrefix(s.ImageName, "thumb_")
}

// ToThumbNail get thumbnail path
func (s *StorageImage) ToThumbNail() string {
	return fmt.Sprintf("%s/thumb_%s", s.UserID, s.ImageName)
}

// PhotoID get ID of photo (name of image without extension or 'thumb_')
func (s *StorageImage) PhotoID() string {
	name := s.ImageName
	if strings.HasPrefix(name, "thumb_") {
		name = name[len("thumb_"):]
	}
	return strings.TrimSuffix(name, filepath.Ext(name))
}
