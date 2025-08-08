package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"ssc-approver/config"
	"ssc-approver/models"
	"time"
)

type SSCClient struct {
	config     *config.Config
	httpClient *http.Client
}

func NewSSCClient(cfg *config.Config) *SSCClient {
	// Create HTTP client with TLS config that accepts self-signed certificates
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	
	return &SSCClient{
		config: cfg,
		httpClient: &http.Client{
			Transport: tr,
			Timeout:   60 * time.Second,
		},
	}
}

func (c *SSCClient) makeRequest(endpoint string, params map[string]string) ([]byte, error) {
	fullURL := fmt.Sprintf("%s/ssc/api/v1/%s", c.config.SSCUrl, endpoint)
	
	// Add query parameters if any
	if len(params) > 0 {
		u, err := url.Parse(fullURL)
		if err != nil {
			return nil, err
		}
		q := u.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
		fullURL = u.String()
	}

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	// Add authorization header
	req.Header.Set("Authorization", fmt.Sprintf("FortifyToken %s", c.config.SSCToken))
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

func (c *SSCClient) GetProjects() ([]models.Project, error) {
	params := map[string]string{
		"limit": "200",
		"fields": "id,name,description",
	}
	
	body, err := c.makeRequest("projects", params)
	if err != nil {
		return nil, err
	}

	var response models.ProjectResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

func (c *SSCClient) GetProjectVersions(projectID int64) ([]models.ProjectVersion, error) {
	endpoint := fmt.Sprintf("projects/%d/versions", projectID)
	params := map[string]string{
		"limit": "200",
		"fields": "id,name",
	}
	
	body, err := c.makeRequest(endpoint, params)
	if err != nil {
		return nil, err
	}

	var response models.ProjectVersionResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

func (c *SSCClient) GetAllProjectVersions() ([]models.ProjectVersion, error) {
	params := map[string]string{
		"limit": "500",
		"fields": "id,name,project",
	}
	
	body, err := c.makeRequest("projectVersions", params)
	if err != nil {
		return nil, err
	}

	var response models.ProjectVersionResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

func (c *SSCClient) GetArtifacts(projectVersionID int64) ([]models.Artifact, error) {
	endpoint := fmt.Sprintf("projectVersions/%d/artifacts", projectVersionID)
	params := map[string]string{
		"limit": "200",
		"embed": "messages",
	}
	
	body, err := c.makeRequest(endpoint, params)
	if err != nil {
		return nil, err
	}

	var response models.ArtifactResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

func (c *SSCClient) GetAllArtifactsRequiringApproval() ([]models.Artifact, error) {
	// First get all project versions
	fmt.Println("Fetching project versions...")
	projectVersions, err := c.GetAllProjectVersions()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Found %d project versions, checking for artifacts requiring approval...\n", len(projectVersions))
	
	var artifactsRequiringApproval []models.Artifact
	checkedCount := 0

	// For each project version, get artifacts
	for _, pv := range projectVersions {
		checkedCount++
		if checkedCount%10 == 0 {
			fmt.Printf("Checked %d/%d project versions...\n", checkedCount, len(projectVersions))
		}
		
		artifacts, err := c.GetArtifacts(pv.ID)
		if err != nil {
			// Continue with other project versions even if one fails
			continue
		}

		for _, artifact := range artifacts {
			// Check if artifact requires approval
			if artifact.Status == "REQUIRE_AUTH" || artifact.Status == "Requires Approval" {
				artifact.ProjectVersionID = pv.ID
				// Create project version name from project and version
				projectName := ""
				if pv.Project != nil {
					projectName = pv.Project.Name
				}
				artifact.ProjectVersionName = fmt.Sprintf("%s - %s", projectName, pv.Name)
				artifactsRequiringApproval = append(artifactsRequiringApproval, artifact)
			}
		}
	}

	return artifactsRequiringApproval, nil
}

func (c *SSCClient) GetArtifactDetails(artifactID int64) (*models.Artifact, error) {
	endpoint := fmt.Sprintf("artifacts/%d", artifactID)
	params := map[string]string{
		"embed": "messages",
	}
	
	body, err := c.makeRequest(endpoint, params)
	if err != nil {
		return nil, err
	}

	var artifact models.Artifact
	if err := json.Unmarshal(body, &artifact); err != nil {
		return nil, err
	}

	return &artifact, nil
}