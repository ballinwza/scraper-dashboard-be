package helper

import "time"

// NowUTC คืนค่าเวลาปัจจุบันในรูปแบบ UTC
func NowUTC() time.Time {
	return time.Now().UTC()
}

// ToBangkokTime แปลงเวลา UTC เป็นเวลาไทย (Asia/Bangkok) เผื่อใช้ในกรณีต้องการ Log หรือแสดงผลฝั่ง Backend
func ToBangkokTime(t time.Time) time.Time {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return t // ถ้าดึง location ไม่ได้ให้คืนค่าเดิม
	}
	return t.In(loc)
}

// FormatISO8601 แปลงเวลาเป็น String รูปแบบ RFC3339/ISO8601
func FormatISO8601(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}
