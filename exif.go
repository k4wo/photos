package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/h2non/filetype"
	"github.com/rwcarlsen/goexif/exif"
	"strconv"
	"strings"
	"time"
)

type exifFields struct {
	ApertureValue                    []string `json:"ApertureValue"`
	BrightnessValue                  []string `json:"BrightnessValue"`
	ColorSpace                       []int    `json:"ColorSpace"`
	ComponentsConfiguration          string   `json:"ComponentsConfiguration"`
	CompressedBitsPerPixel           []string `json:"CompressedBitsPerPixel"`
	Copyright                        string   `json:"Copyright"`
	CustomRendered                   []int    `json:"CustomRendered"`
	DateTime                         string   `json:"DateTime"`
	DateTimeDigitized                string   `json:"DateTimeDigitized"`
	DateTimeOriginal                 string   `json:"DateTimeOriginal"`
	ExifIFDPointer                   []int    `json:"ExifIFDPointer"`
	ExifVersion                      string   `json:"ExifVersion"`
	ExposureBiasValue                []string `json:"ExposureBiasValue"`
	ExposureMode                     []int    `json:"ExposureMode"`
	ExposureProgram                  []int    `json:"ExposureProgram"`
	ExposureTime                     []string `json:"ExposureTime"`
	FNumber                          []string `json:"FNumber"`
	FileSource                       string   `json:"FileSource"`
	Flash                            []int    `json:"Flash"`
	FlashpixVersion                  string   `json:"FlashpixVersion"`
	FocalLength                      []string `json:"FocalLength"`
	FocalPlaneResolutionUnit         []int    `json:"FocalPlaneResolutionUnit"`
	FocalPlaneXResolution            []string `json:"FocalPlaneXResolution"`
	FocalPlaneYResolution            []string `json:"FocalPlaneYResolution"`
	ISOSpeedRatings                  []int    `json:"ISOSpeedRatings"`
	InteroperabilityIFDPointer       []int    `json:"InteroperabilityIFDPointer"`
	InteroperabilityIndex            string   `json:"InteroperabilityIndex"`
	LightSource                      []int    `json:"LightSource"`
	Make                             string   `json:"Make"`
	MakerNote                        string   `json:"MakerNote"`
	MaxApertureValue                 []string `json:"MaxApertureValue"`
	MeteringMode                     []int    `json:"MeteringMode"`
	Model                            string   `json:"Model"`
	Orientation                      []int    `json:"Orientation"`
	PixelXDimension                  []int    `json:"PixelXDimension"`
	PixelYDimension                  []int    `json:"PixelYDimension"`
	ResolutionUnit                   []int    `json:"ResolutionUnit"`
	SceneCaptureType                 []int    `json:"SceneCaptureType"`
	SceneType                        string   `json:"SceneType"`
	SensingMethod                    []int    `json:"SensingMethod"`
	Sharpness                        []int    `json:"Sharpness"`
	ShutterSpeedValue                []string `json:"ShutterSpeedValue"`
	Software                         string   `json:"Software"`
	SubjectDistanceRange             []int    `json:"SubjectDistanceRange"`
	ThumbJPEGInterchangeFormat       []int    `json:"ThumbJPEGInterchangeFormat"`
	ThumbJPEGInterchangeFormatLength []int    `json:"ThumbJPEGInterchangeFormatLength"`
	WhiteBalance                     []int    `json:"WhiteBalance"`
	XResolution                      []string `json:"XResolution"`
	YCbCrPositioning                 []int    `json:"YCbCrPositioning"`
	YResolution                      []string `json:"YResolution"`
}

// ImageInfo image data (should have the same info like in db)
type ImageInfo struct {
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

func parseValues(val []string) (float64, float64) {
	values := strings.Split(val[0], "/")
	first, err := strconv.ParseFloat(values[0], 32)
	second, err := strconv.ParseFloat(values[1], 32)
	if err != nil {
		fmt.Println("Can't parse values: ", val, err)
		return 0, 0
	}

	return first, second
}

func convertToFloat(val []string) float32 {
	first, second := parseValues(val)

	return float32(first / second)
}

func parseExposureTime(val []string) string {
	x, y := parseValues(val)
	if x == 0 || y == 0 {
		return ""
	}

	var second uint
	if int(y) == 1 {
		second = uint(x)
	} else {
		second = uint(y / x)
	}

	return fmt.Sprintf("%d/%d", int(x/x), second)
}

func extractExif(data []byte) (ImageInfo, error) {
	kind, err := filetype.Match(data)
	if kind == filetype.Unknown || err != nil {
		fmt.Println("Unknown file type", err)
	}

	file := bytes.NewReader(data)
	fileExif, err := exif.Decode(file)
	if err != nil {
		return ImageInfo{}, err
	}

	var jsonExif exifFields
	byteJSON, err := fileExif.MarshalJSON()
	err = json.Unmarshal(byteJSON, &jsonExif)
	if err != nil {
		fmt.Println("Can't parse JSON EXIF. ", err)
	}

	dateTime, err := fileExif.DateTime()
	long, lat, err := fileExif.LatLong()

	return ImageInfo{
		dateTime,
		uint(jsonExif.PixelXDimension[0]),
		uint(jsonExif.PixelYDimension[0]),
		convertToFloat(jsonExif.FNumber),
		parseExposureTime(jsonExif.ExposureTime),
		convertToFloat(jsonExif.FocalLength),
		jsonExif.ISOSpeedRatings[0],
		jsonExif.Make,
		jsonExif.Model,
		jsonExif.Orientation[0],
		long,
		lat,
		"",
		"",
		kind.Extension,
		kind.MIME.Value,
		len(data),
		1,
	}, nil
}
