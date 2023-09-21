package models

type Student struct {
	Stundent_ID      uint64 `json:"student_id" binding:"required"`
	Stundent_name    string `json:"student_name" binding:"required"`
	Stundent_age     uint64 `json:"student_age" binding:"required"`
	Stundent_address string `json:"student_address" binding:"required"`
	Stundent_phone   string `json:"student_phone" binding:"required"`
}
