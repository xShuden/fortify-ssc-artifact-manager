package models

import "time"

type ProjectVersion struct {
	ID                 int64       `json:"id"`
	Name               string      `json:"name"`
	Project            *Project    `json:"project,omitempty"`
	ProjectVersionName string      `json:"projectVersionName,omitempty"`
}

type Artifact struct {
	ID                  int64           `json:"id"`
	FileName            string          `json:"fileName"`
	FileSize            int64           `json:"fileSize"`
	Status              string          `json:"status"`
	UploadDate          *time.Time      `json:"uploadDate"`
	UploadIP            string          `json:"uploadIP"`
	UserName            string          `json:"userName"`
	ArtifactType        string          `json:"artifactType"`
	Messages            interface{}     `json:"messages,omitempty"` // Can be string or []Message
	ProcessingMessages  string          `json:"processingMessages,omitempty"`
	ApprovalRequired    bool            `json:"approvalRequired"`
	ApprovalComment     string          `json:"approvalComment,omitempty"`
	ProjectVersionID    int64           `json:"projectVersionId,omitempty"`
	ProjectVersion      *ProjectVersion `json:"_embed.projectVersion,omitempty"`
	ProjectVersionName  string          `json:"projectVersionName,omitempty"`
}

type Message struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type ApiResponse struct {
	Data       []interface{} `json:"data"`
	Count      int           `json:"count"`
	TotalCount int           `json:"totalCount"`
}

type ProjectResponse struct {
	Data       []Project `json:"data"`
	Count      int       `json:"count"`
	TotalCount int       `json:"totalCount"`
}

type Project struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ProjectVersionResponse struct {
	Data       []ProjectVersion `json:"data"`
	Count      int              `json:"count"`
	TotalCount int              `json:"totalCount"`
}

type ArtifactResponse struct {
	Data       []Artifact `json:"data"`
	Count      int        `json:"count"`
	TotalCount int        `json:"totalCount"`
}