package dto

type ReportCreateDto struct {
	TargetID   uint   `json:"target_id"   validate:"required"`
	TargetType string `json:"target_type" validate:"required,oneof=USER TEAM CHALLENGE MESSAGE"`
	Reason     string `json:"reason"      validate:"required,min=3"`
	Comment    string `json:"comment"`
}
