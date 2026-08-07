package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// CustomValidator โครงสร้างหลักสำหรับจัดการ Validator Instance
type CustomValidator struct {
	validator *validator.Validate
}

// ErrorResponse ตัวแปรสำหรับจัด Format การตอบกลับเมื่อเกิด Validation Error
type ErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// New สร้าง Instance ใหม่ของ CustomValidator
func New() *CustomValidator {
	v := validator.New()

	// สามารถลงทะเบียน Custom Validation Rules เพิ่มเติมตรงนี้ได้ (ถ้ามี)
	// v.RegisterValidation("custom_tag", customFunc)

	return &CustomValidator{
		validator: v,
	}
}

// ValidateStruct ตรวจสอบข้อมูลใน Struct และคืนค่ารายชื่อ Error ทั้งหมดถ้าพบข้อผิดพลาด
func (cv *CustomValidator) ValidateStruct(s interface{}) []ErrorResponse {
	var errors []ErrorResponse

	err := cv.validator.Struct(s)
	if err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range validationErrs {
				errors = append(errors, ErrorResponse{
					Field:   fieldErr.Field(),
					Message: msgForTag(fieldErr),
				})
			}
		}
	}

	return errors
}

// msgForTag แปลง Validation Tag เป็นข้อความอ่านง่าย
func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Must be at least %s characters/value", fe.Param())
	case "max":
		return fmt.Sprintf("Must be at most %s characters/value", fe.Param())
	case "url":
		return "Must be a valid URL"
	default:
		return fmt.Sprintf("Failed validation on tag '%s'", fe.Tag())
	}
}