package image

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	model "photos/model"

	"github.com/h2non/filetype"
	"github.com/rwcarlsen/goexif/exif"
	"gopkg.in/guregu/null.v3"
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

func convertToFloat(val []string) null.Float {
	first, second := parseValues(val)

	return null.FloatFrom(float64(first / second))
}

func convertToInt(val int) null.Int {
	return null.IntFrom(int64(val))
}

func parseExposureTime(val []string) null.String {
	x, y := parseValues(val)
	if x == 0 || y == 0 {
		return null.StringFrom("")
	}

	var second uint
	if int(y) == 1 {
		second = uint(x)
	} else {
		second = uint(y / x)
	}

	return null.StringFrom(fmt.Sprintf("%d/%d", int(x/x), second))
}

// ExtractExif extracts exif from file and returns File
func ExtractExif(data []byte) (model.File, error) {
	kind, err := filetype.Match(data)
	if kind == filetype.Unknown || err != nil {
		fmt.Println("Unknown file type", err)
	}

	file := bytes.NewReader(data)
	fileExif, err := exif.Decode(file)
	if err != nil {
		return model.File{}, err
	}

	var jsonExif exifFields
	byteJSON, err := fileExif.MarshalJSON()
	err = json.Unmarshal(byteJSON, &jsonExif)
	if err != nil {
		fmt.Println("Can't parse JSON EXIF. ", err)
	}

	dateTime, err := fileExif.DateTime()
	long, lat, err := fileExif.LatLong()

	return model.File{
		Camera:       null.StringFrom(jsonExif.Make),
		Date:         null.TimeFrom(dateTime),
		ExposureTime: parseExposureTime(jsonExif.ExposureTime),
		Extension:    null.StringFrom(kind.Extension),
		FNumber:      convertToFloat(jsonExif.FNumber),
		FocalLength:  convertToFloat(jsonExif.FocalLength),
		Height:       convertToInt(jsonExif.PixelYDimension[0]),
		Iso:          convertToInt(jsonExif.ISOSpeedRatings[0]),
		Latitude:     null.FloatFrom(lat),
		Longitude:    null.FloatFrom(long),
		MimeType:     null.StringFrom(kind.MIME.Value),
		Model:        null.StringFrom(jsonExif.Model),
		Orientation:  convertToInt(jsonExif.Orientation[0]),
		Owner:        convertToInt(1),
		Size:         convertToInt(len(data)),
		Width:        convertToInt(jsonExif.PixelXDimension[0]),
	}, nil
}
