package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

var arialFace font.Face

func loadFont() font.Face {
	fontBytes, err := os.ReadFile("arial.ttf")
	if err != nil {
		log.Fatalf("failed to read font: %v", err)
	}

	fnt, err := opentype.Parse(fontBytes)
	if err != nil {
		log.Fatalf("failed to parse font: %v", err)
	}

	face, err := opentype.NewFace(fnt, &opentype.FaceOptions{
		Size:    48,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatalf("failed to create font face: %v", err)
	}

	return face
}

func drawText(img *image.RGBA, x, y int, label string) {
	col := color.Black
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: arialFace,
		Dot:  fixed.P(x, y),
	}
	d.DrawString(label)
}

func drawWrappedText(img *image.RGBA, x, y, maxWidth int, text string) int {
	words := strings.Split(text, " ")
	var line []string
	lineHeight := int(float64(arialFace.Metrics().Height.Ceil()) * 1) // Add some spacing between lines
	startY := y

	// Measure text width
	measureString := func(s string) int {
		return font.MeasureString(arialFace, s).Ceil()
	}

	for _, word := range words {
		newLine := strings.Join(append(line, word), " ")
		if measureString(newLine) > maxWidth {
			if len(line) > 0 {
				// Draw current line
				drawText(img, x, y, strings.Join(line, " "))
				y += lineHeight
				line = []string{word}
			} else {
				// Word is too long for the line, draw it anyway
				drawText(img, x, y, word)
				y += lineHeight
			}
		} else {
			line = append(line, word)
		}
	}

	// Draw the last line
	if len(line) > 0 {
		drawText(img, x, y, strings.Join(line, " "))
		y += lineHeight
	}

	return y - startY // Return total height used
}

func drawWrappedTextFromBottom(img *image.RGBA, x, bottomY, maxWidth int, text string) int {
	words := strings.Split(text, " ")
	var lines []string
	var currentLine []string
	lineHeight := int(float64(arialFace.Metrics().Height.Ceil() * 1)) // Add some spacing between lines

	// Measure text width
	measureString := func(s string) int {
		return font.MeasureString(arialFace, s).Ceil()
	}

	// First, calculate all lines
	for _, word := range words {
		newLine := strings.Join(append(currentLine, word), " ")
		if measureString(newLine) > maxWidth {
			if len(currentLine) > 0 {
				lines = append(lines, strings.Join(currentLine, " "))
				currentLine = []string{word}
			} else {
				lines = append(lines, word)
			}
		} else {
			currentLine = append(currentLine, word)
		}
	}
	if len(currentLine) > 0 {
		lines = append(lines, strings.Join(currentLine, " "))
	}

	// Calculate starting Y position
	startY := bottomY - (lineHeight * (len(lines) - 1))

	// Draw lines from top to bottom
	for i, line := range lines {
		drawText(img, x, startY+(i*lineHeight), line)
	}

	return startY // Return the top Y position where drawing started
}

func drawCornerArc(img *image.RGBA, cx, cy, r int, startAngle, endAngle float64, borderWidth int, col color.Color) {
	for a := startAngle; a <= endAngle; a += 0.01 {
		for b := 0; b < borderWidth; b++ {
			dx := int(float64(r-b) * math.Cos(a))
			dy := int(float64(r-b) * math.Sin(a))
			img.Set(cx+dx, cy+dy, col)
		}
	}
}

func drawRoundedBorder(img *image.RGBA, width, height, borderWidth, radius int, col color.Color) {
	// Draw straight edges
	for i := 0; i < borderWidth; i++ {
		for x := radius; x < width-radius; x++ {
			img.Set(x, i, col)
			img.Set(x, height-1-i, col)
		}
		for y := radius; y < height-radius; y++ {
			img.Set(i, y, col)
			img.Set(width-1-i, y, col)
		}
	}

	// Draw corner arcs
	drawCornerArc(img, radius, radius, radius, math.Pi, 1.5*math.Pi, borderWidth, col)            // top-left
	drawCornerArc(img, width-1-radius, radius, radius, 1.5*math.Pi, 2*math.Pi, borderWidth, col)  // top-right
	drawCornerArc(img, radius, height-1-radius, radius, 0.5*math.Pi, math.Pi, borderWidth, col)   // bottom-left
	drawCornerArc(img, width-1-radius, height-1-radius, radius, 0, 0.5*math.Pi, borderWidth, col) // bottom-right
}

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	jabatan := query.Get("jabatan")
	nama := query.Get("nama")
	jabatan, _ = url.QueryUnescape(jabatan)
	nama, _ = url.QueryUnescape(nama)

	// Check for empty parameters
	missingParams := make([]string, 0)

	if nama == "" {
		missingParams = append(missingParams, "nama")
	}

	// If any required parameter is missing, return error response
	if len(missingParams) > 0 {
		statusCode := http.StatusBadRequest // 400

		// Log the error with missing parameters
		logData := map[string]interface{}{
			"time":          time.Now().Format(time.RFC3339),
			"method":        r.Method,
			"url":           r.RequestURI,
			"remoteIP":      r.RemoteAddr,
			"error":         "Missing required parameters",
			"missingParams": missingParams,
			"statusCode":    statusCode,
		}

		var buf bytes.Buffer
		encoder := json.NewEncoder(&buf)
		encoder.SetEscapeHTML(false)
		encoder.Encode(logData)
		fmt.Print(buf.String())

		// Return error response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		response := map[string]interface{}{
			"error":   "Missing required parameters",
			"missing": missingParams,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	if jabatan != "" {
		jabatan = strings.ToUpper(jabatan)
	}

	statusCode := http.StatusOK // 200

	// Log successful request
	logData := map[string]interface{}{
		"time":       time.Now().Format(time.RFC3339),
		"method":     r.Method,
		"url":        r.RequestURI,
		"jabatan":    jabatan,
		"nama":       nama,
		"remoteIP":   r.RemoteAddr,
		"statusCode": statusCode,
	}
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.Encode(logData)
	fmt.Print(buf.String())

	// Create a blank white canvas
	width := 1489
	height := 485
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	white := color.RGBA{255, 255, 255, 255}
	draw.Draw(rgba, rgba.Bounds(), &image.Uniform{white}, image.Point{}, draw.Src)

	// Draw rounded border
	drawRoundedBorder(rgba, width, height, 4, 70, color.Black)

	// Draw the text
	const maxWidth = 980
	const topMargin = 65
	const bottomMargin = 440

	// Draw from top
	y := topMargin
	drawText(rgba, 400, y, "Ditandatangani secara elektronik oleh:")

	// Draw jabatan with wrapping from top
	y = 125
	if jabatan != "" {
		y += drawWrappedText(rgba, 400, y, maxWidth, jabatan+",")
	}

	// Draw footer text from bottom
	drawText(rgba, 400, bottomMargin, "Penyedia")

	// Draw nama with wrapping from bottom
	// Calculate position above footer text
	namaBottomY := bottomMargin - 55 // Leave some space above "Penyedia"
	drawWrappedTextFromBottom(rgba, 400, namaBottomY, maxWidth, nama)

	w.Header().Set("Content-Type", "image/jpeg")
	jpeg.Encode(w, rgba, &jpeg.Options{Quality: 100})
}

func main() {
	arialFace = loadFont()
	http.HandleFunc("/generate", handler)

	logData := map[string]interface{}{
		"time":    time.Now().Format(time.RFC3339),
		"message": "Starting server",
		"port":    ":8080",
	}
	logJSON, _ := json.Marshal(logData)
	fmt.Println(string(logJSON))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
