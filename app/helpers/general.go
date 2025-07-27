package helpers

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/Kagami/go-avif"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	xdraw "golang.org/x/image/draw"
)

func ConvertAndSaveAsWebp(src io.Reader, filename string, saveDir string) error {
	// Decode gambar (jpeg/png)
	img, _, err := image.Decode(src)
	if err != nil {
		return err
	}

	img = ResizeIfAboveFHD(img)

	// Pilih mode encoding: Lossy dengan kualitas 75
	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 75)
	if err != nil {
		return err
	}

	// Encode ke webp
	var buf bytes.Buffer
	if err := webp.Encode(&buf, img, options); err != nil {
		return err
	}

	// Path penyimpanan final
	fullPath := filepath.Join("storage/app/", saveDir, filename)

	// Buat folder jika belum ada
	if err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm); err != nil {
		return err
	}

	// Tulis hasil ke file
	return os.WriteFile(fullPath, buf.Bytes(), 0644)
}

func CapitalizeEachWords(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			r := []rune(word)
			r[0] = unicode.ToUpper(r[0])
			words[i] = string(r)
		}
	}
	return strings.Join(words, " ")
}

func ConvertAndSaveAsAvif(src io.Reader, filename string, saveDir string) error {
	img, _, err := image.Decode(src)
	if err != nil {
		return err
	}

	img = ResizeIfAboveFHD(img)

	var buf bytes.Buffer
	err = avif.Encode(&buf, img, nil) // Gunakan opsi default
	if err != nil {
		return err
	}

	fullPath := filepath.Join("storage/app/", saveDir, filename)

	if err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm); err != nil {
		return err
	}

	return os.WriteFile(fullPath, buf.Bytes(), 0644)
}

func TanggalStringToIndo(tgl string) string {
	tanggalAsli := tgl // "2025-07-22T00:00:00+07:00"

	t, _ := time.Parse(time.RFC3339, tanggalAsli)

	tanggal := ""
	if results, err := TanggalIndo(t.Format("2006-01-02"), false); err == nil {
		tanggal = results
	}

	return tanggal
}

func BuildDynamicList(items []string) string {
	var builder strings.Builder

	builder.WriteString(`<ul class="ps-2 mb-0">`)

	for _, item := range items {
		item = strings.TrimSpace(item)
		if item != "" {
			builder.WriteString("<li>" + item + "</li>")
		}
	}

	builder.WriteString("</ul>")
	return builder.String()
}

func GenerateCutWords(text string, max int) string {
	// Jika panjangnya sudah di bawah batas, langsung kembalikan
	if utf8.RuneCountInString(text) <= max {
		return text
	}

	// Potong sementara ke panjang maksimum
	trimmed := []rune(text)[:max]
	result := string(trimmed)

	// Cari spasi terakhir untuk memotong di akhir kata
	lastSpace := strings.LastIndex(result, " ")
	if lastSpace != -1 {
		result = result[:lastSpace]
	}

	return result
}

func TanggalIndo(tgl string, tampilHari bool) (string, error) {
	// Nama-nama hari dan bulan
	namaHari := []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
	namaBulan := []string{
		"", "Januari", "Februari", "Maret", "April", "Mei",
		"Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}

	// Parse tanggal string
	parsed, err := time.Parse("2006-01-02", tgl)
	if err != nil {
		return "", fmt.Errorf("format tanggal tidak valid: %v", err)
	}

	tanggal := strconv.Itoa(parsed.Day())
	bulan := namaBulan[int(parsed.Month())]
	tahun := strconv.Itoa(parsed.Year())

	var text strings.Builder
	if tampilHari {
		urutanHari := int(parsed.Weekday()) // 0 = Minggu, 1 = Senin, dst.
		hari := namaHari[urutanHari]
		text.WriteString(hari + ", ")
	}

	text.WriteString(tanggal + " " + bulan + " " + tahun)
	return text.String(), nil
}

func ResizeIfAboveFHD(img image.Image) image.Image {
	origBounds := img.Bounds()
	width := origBounds.Dx()
	height := origBounds.Dy()

	const maxWidth = 1920
	const maxHeight = 1080

	// Skip resizing if image already smaller than FHD
	if width <= maxWidth && height <= maxHeight {
		return img
	}

	// Calculate scaling ratio
	ratioW := float64(maxWidth) / float64(width)
	ratioH := float64(maxHeight) / float64(height)
	ratio := math.Min(ratioW, ratioH)

	newWidth := int(float64(width) * ratio)
	newHeight := int(float64(height) * ratio)

	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	xdraw.CatmullRom.Scale(dst, dst.Bounds(), img, origBounds, draw.Over, nil)

	return dst
}

func StringToInt(v any) (int, error) {
	s, ok := v.(string)
	if !ok {
		return 0, fmt.Errorf("nilai bukan string (actual: %T)", v)
	}
	return strconv.Atoi(s)
}

func StringToUint(v any) (uint, error) {
	s, ok := v.(string)
	if !ok {
		return 0, fmt.Errorf("nilai bukan string (actual: %T)", v)
	}
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(n), nil
}
