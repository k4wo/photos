package model

import (
	"gopkg.in/guregu/null.v3"
)

// File file descriptor
type File struct {
	ID           null.Int    `json:"id,omitempty"`
	Date         null.Time   `json:"date,omitempty"`         // DateTime
	Width        null.Int    `json:"width,omitempty"`        // PixelXDimension
	Height       null.Int    `json:"height,omitempty"`       // PixelYDimension
	FNumber      null.Float  `json:"fNumber,omitempty"`      // FNumber
	ExposureTime null.String `json:"exposureTime,omitempty"` // ExposureTime
	FocalLength  null.Float  `json:"focalLength,omitempty"`  // FocalLength (mm)
	Iso          null.Int    `json:"iso,omitempty"`          // ISOSpeedRatings
	Camera       null.String `json:"camera,omitempty"`       // Make
	Model        null.String `json:"model,omitempty"`        // Model
	Orientation  null.Int    `json:"orientation,omitempty"`  // Orientation
	Longitude    null.Float  `json:"longitude,omitempty"`
	Latitude     null.Float  `json:"latitude,omitempty"`
	Name         null.String `json:"name,omitempty"`
	Hash         null.String `json:"hash,omitempty"`
	Extension    null.String `json:"extension,omitempty"`
	MimeType     null.String `json:"mimeType,omitempty"`
	Size         null.Int    `json:"size,omitempty"`
	Owner        null.Int    `json:"owner,omitempty"`
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
	File      Cover    `json:"file,omitempty"`
}
