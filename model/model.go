package model

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

// File file descriptor
type File struct {
	ID           int       `json:"id"`
	Date         time.Time `json:"date"`         // DateTime
	Width        uint      `json:"width"`        // PixelXDimension
	Height       uint      `json:"height"`       // PixelYDimension
	FNumber      float32   `json:"fNumber"`      // FNumber
	ExposureTime string    `json:"exposureTime"` // ExposureTime
	FocalLength  float32   `json:"focalLength"`  // FocalLength (mm)
	Iso          int       `json:"iso"`          // ISOSpeedRatings
	Camera       string    `json:"camera"`       // Make
	Model        string    `json:"model"`        // Model
	Orientation  int       `json:"orientation"`  // Orientation
	Longitude    float64   `json:"longitude"`
	Latitude     float64   `json:"latitude"`
	Name         string    `json:"name"`
	Hash         string    `json:"hash"`
	Extension    string    `json:"extension"`
	MimeType     string    `json:"mimeType"`
	Size         int       `json:"size"`
	Owner        int       `json:"owner"`
}

// Album descriptor
type Album struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Size      int      `json:"size"`
	Owner     int      `json:"owner"`
	Cover     null.Int `json:"cover"`
	UpdatedAt string   `json:"updatedAt"`
	CreatedAt string   `json:"createdAt"`
}
