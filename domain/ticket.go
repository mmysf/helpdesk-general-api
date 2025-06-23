package domain

import "app/domain/model"

type TicketRequest struct {
	Subject    string   `form:"subject"`
	Content    string   `form:"content"`
	Priority   string   `form:"priority"`
	AttachIds  []string `form:"attachIds"`
	ParentId   string   `form:"parentId"`
	ProjectId  string   `form:"projectId"`
	CategoryId string   `form:"categoryId"`
}

type TicketCommentRequest struct {
	TicketId  string             `json:"ticketId"`
	Content   string             `json:"content"`
	AttachIds []string           `json:"attachIds"`
	Status    model.TicketStatus `json:"status"`
}

type SuperadminTicketCommentRequest struct {
	AgentId   string             `json:"agentId"`
	TicketId  string             `json:"ticketId"`
	Content   string             `json:"content"`
	AttachIds []string           `json:"attachIds"`
	Status    model.TicketStatus `json:"status"`
}

type CloseTicketRequest struct {
	TicketId string `json:"id"`
}
type CloseTicketbyEmailRequest struct {
	Token string `json:"token"`
}

type ReopenTicketRequest struct {
	TicketId string `json:"id"`
}

type LoggingTicketRequest struct {
	TicketId string `json:"id"`
}

type UploadAttachment struct {
	Title string `form:"title"`
}

type TimeTrackRequest struct {
	Hour   int8 `json:"hour"`
	Minute int  `json:"minute"`
	Second int  `json:"second"`
}

type CancelTicketRequest struct {
	TicketId string `json:"id"`
}

var AllowedMimeTypes = []string{
	"application/pdf", // pdf
	// "application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document", // doc & docx
	"image/jpeg", "image/png", "image/gif", "image/svg+xml", "image/webp", "image/vnd.microsoft.icon", "image/x-icon", // images
	"video/mp4", "video/mpeg", "video/ogg", "video/webm", // videos
	// "audio/mpeg", "audio/ogg", "audio/wav", "audio/webm", // auidos
	// "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", // xlsx
}

var AllowedImgMimeTypes = []string{
	"image/jpeg", "image/png", "image/gif", "image/svg+xml", "image/webp", "image/vnd.microsoft.icon", "image/x-icon", // images
}
