package helper

import (
	"regexp"
	"strconv"
	"strings"
)

func ExtractFloat(text string) float64 {
	if text == "" {
		return 0
	}

	// แปลงตัวเลขไทยเป็นอารบิกก่อน (ถ้ามี)
	cleaned := convertThaiDigits(text)

	// Pattern ดักจับตัวเลขที่มีจุดทศนิยม และรองรับเครื่องหมาย comma (เช่น 1,500.50)
	re := regexp.MustCompile(`([\d,]+\.?\d*)`)
	match := re.FindString(cleaned)

	if match == "" {
		return 0
	}

	// ลบเครื่องหมาย comma (,) ออกเพื่อให้ Parse ได้ถูกต้อง
	rawNum := strings.ReplaceAll(match, ",", "")

	val, err := strconv.ParseFloat(rawNum, 64)
	if err != nil {
		return 0
	}

	return val
}

// ExtractInt สกัดตัวเลขจำนวนเต็มตัวแรกที่พบในข้อความ
func ExtractInt(text string) int {
	if text == "" {
		return 0
	}

	cleaned := convertThaiDigits(text)

	// Pattern ดักจับเฉพาะตัวเลขล้วน
	re := regexp.MustCompile(`\d+`)
	match := re.FindString(cleaned)

	if match == "" {
		return 0
	}

	val, err := strconv.Atoi(match)
	if err != nil {
		return 0
	}

	return val
}

// convertThaiDigits แปลงตัวเลขไทย (๐-๙) เป็นตัวเลขอารบิก (0-9)
func convertThaiDigits(text string) string {
	thaiDigits := []string{"๐", "๑", "๒", "๓", "๔", "๕", "๖", "๗", "๘", "๙"}
	arabicDigits := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

	result := text
	for i := 0; i < len(thaiDigits); i++ {
		result = strings.ReplaceAll(result, thaiDigits[i], arabicDigits[i])
	}
	return result
}
