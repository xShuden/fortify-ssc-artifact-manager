package main

import (
	"encoding/json"
	"fmt"
	"os"
	"ssc-approver/client"
	"ssc-approver/config"
	"ssc-approver/models"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	projectFilter string
	showDetails   bool
	outputFormat  string
	sscURL       string
	sscToken     string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "ssc-approver",
		Short: "SSC Approval Tool - List and manage artifacts requiring approval",
		Long: `A CLI tool to interact with Fortify Software Security Center (SSC) API 
to list artifacts that require approval and display their processing messages.`,
	}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List artifacts requiring approval",
		Long:  `List all artifacts that have "Requires Approval" status across all projects`,
		Run:   listArtifacts,
	}

	var projectsCmd = &cobra.Command{
		Use:   "projects",
		Short: "List all projects",
		Long:  `List all projects available in SSC`,
		Run:   listProjects,
	}

	var artifactsCmd = &cobra.Command{
		Use:   "artifacts [project-version-id]",
		Short: "List artifacts for a specific project version",
		Long:  `List all artifacts for a specific project version ID`,
		Args:  cobra.ExactArgs(1),
		Run:   listProjectArtifacts,
	}

	// Add global flags
	rootCmd.PersistentFlags().StringVarP(&sscURL, "url", "u", "", "SSC URL (e.g., https://sast.example.com)")
	rootCmd.PersistentFlags().StringVarP(&sscToken, "token", "t", "", "SSC API Token")

	// Add command-specific flags
	listCmd.Flags().StringVarP(&projectFilter, "project", "p", "", "Filter by project name (partial match)")
	listCmd.Flags().BoolVarP(&showDetails, "details", "d", false, "Show detailed processing messages")
	listCmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Output format: table, json, csv")

	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(projectsCmd)
	rootCmd.AddCommand(artifactsCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func listArtifacts(cmd *cobra.Command, args []string) {
	cfg, err := config.LoadConfigWithOverrides(sscURL, sscToken)
	if err != nil {
		color.Red("Error loading configuration: %v", err)
		os.Exit(1)
	}

	sscClient := client.NewSSCClient(cfg)
	
	artifacts, err := sscClient.GetAllArtifactsRequiringApproval()
	if err != nil {
		color.Red("Error fetching artifacts: %v", err)
		os.Exit(1)
	}

	// Filter by project if specified
	if projectFilter != "" {
		filtered := []models.Artifact{}
		for _, artifact := range artifacts {
			if strings.Contains(strings.ToLower(artifact.ProjectVersionName), strings.ToLower(projectFilter)) {
				filtered = append(filtered, artifact)
			}
		}
		artifacts = filtered
	}

	if len(artifacts) == 0 {
		color.Green("No artifacts requiring approval found.")
		return
	}

	displayArtifacts(artifacts)
}

func displayArtifacts(artifacts []models.Artifact) {
	color.Yellow("\nFound %d artifacts requiring approval:\n", len(artifacts))

	switch outputFormat {
	case "json":
		// JSON output
		type JSONArtifact struct {
			Project         string `json:"project"`
			UploadDate      string `json:"upload_date"`
			FileName        string `json:"file_name"`
			FileSizeBytes   int64  `json:"file_size_bytes"`
			FileSizeMB      string `json:"file_size_mb"`
			UploadIP        string `json:"upload_ip"`
			Status          string `json:"status"`
			Messages        string `json:"messages,omitempty"`
		}
		
		var jsonArtifacts []JSONArtifact
		for _, artifact := range artifacts {
			uploadDate := ""
			if artifact.UploadDate != nil {
				uploadDate = artifact.UploadDate.Format("2006-01-02 15:04:05")
			}
			
			jsonArtifact := JSONArtifact{
				Project:       artifact.ProjectVersionName,
				UploadDate:    uploadDate,
				FileName:      artifact.FileName,
				FileSizeBytes: artifact.FileSize,
				FileSizeMB:    fmt.Sprintf("%.2f", float64(artifact.FileSize)/(1024*1024)),
				UploadIP:      artifact.UploadIP,
				Status:        "Requires Approval",
			}
			
			if showDetails {
				jsonArtifact.Messages = extractMessages(artifact)
			}
			
			jsonArtifacts = append(jsonArtifacts, jsonArtifact)
		}
		
		output, err := json.MarshalIndent(jsonArtifacts, "", "  ")
		if err != nil {
			color.Red("Error formatting JSON: %v", err)
			return
		}
		fmt.Println(string(output))
		return
		
	case "csv":
		fmt.Println("Project,Upload Date,File Name,Size,Upload IP,Messages")
		for _, artifact := range artifacts {
			uploadDate := ""
			if artifact.UploadDate != nil {
				uploadDate = artifact.UploadDate.Format("2006-01-02 15:04:05")
			}
			messages := extractMessages(artifact)
			fmt.Printf("%s,%s,%s,%d,%s,\"%s\"\n",
				artifact.ProjectVersionName,
				uploadDate,
				artifact.FileName,
				artifact.FileSize,
				artifact.UploadIP,
				strings.ReplaceAll(messages, "\"", "\"\""))
		}
		return
		
	default:
		// Table format (default)
		table := tablewriter.NewWriter(os.Stdout)
	
	if showDetails {
		table.SetHeader([]string{"Project", "Upload Date", "File Name", "Size (MB)", "Upload IP", "Processing Messages"})
		table.SetRowLine(true)
		table.SetAutoWrapText(true)
		table.SetColWidth(50)
	} else {
		table.SetHeader([]string{"Project", "Upload Date", "File Name", "Size (MB)", "Upload IP", "Status"})
	}

	for _, artifact := range artifacts {
		uploadDate := ""
		if artifact.UploadDate != nil {
			uploadDate = artifact.UploadDate.Format("2006-01-02 15:04:05")
		}
		
		sizeInMB := fmt.Sprintf("%.2f", float64(artifact.FileSize)/(1024*1024))

		if showDetails {
			messages := extractMessages(artifact)
			table.Append([]string{
				artifact.ProjectVersionName,
				uploadDate,
				artifact.FileName,
				sizeInMB,
				artifact.UploadIP,
				messages,
			})
		} else {
			table.Append([]string{
				artifact.ProjectVersionName,
				uploadDate,
				artifact.FileName,
				sizeInMB,
				artifact.UploadIP,
				"Requires Approval",
			})
		}
	}

		table.Render()

		if !showDetails {
			color.Cyan("\nTip: Use -d or --details flag to see processing messages for each artifact")
		}
	}
}

func extractMessages(artifact models.Artifact) string {
	var messages []string
	
	// Add processing messages if any
	if artifact.ProcessingMessages != "" {
		messages = append(messages, artifact.ProcessingMessages)
	}
	
	// Handle Messages field which can be string or []Message
	if artifact.Messages != nil {
		switch v := artifact.Messages.(type) {
		case string:
			if v != "" {
				messages = append(messages, v)
			}
		case []interface{}:
			for _, item := range v {
				if msgMap, ok := item.(map[string]interface{}); ok {
					msg := ""
					code := ""
					if m, exists := msgMap["message"]; exists {
						msg = fmt.Sprintf("%v", m)
					}
					if c, exists := msgMap["code"]; exists {
						code = fmt.Sprintf("%v", c)
					}
					if msg != "" {
						if code != "" {
							messages = append(messages, fmt.Sprintf("[%s] %s", code, msg))
						} else {
							messages = append(messages, msg)
						}
					}
				}
			}
		}
	}
	
	if len(messages) == 0 {
		return "No processing messages available"
	}
	
	return strings.Join(messages, "\n")
}

func listProjects(cmd *cobra.Command, args []string) {
	cfg, err := config.LoadConfigWithOverrides(sscURL, sscToken)
	if err != nil {
		color.Red("Error loading configuration: %v", err)
		os.Exit(1)
	}

	sscClient := client.NewSSCClient(cfg)
	
	fmt.Println("Fetching projects...")
	projects, err := sscClient.GetProjects()
	if err != nil {
		color.Red("Error fetching projects: %v", err)
		os.Exit(1)
	}

	if len(projects) == 0 {
		color.Yellow("No projects found.")
		return
	}

	color.Green("\nFound %d projects:\n", len(projects))

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Description"})

	for _, project := range projects {
		description := project.Description
		if len(description) > 50 {
			description = description[:47] + "..."
		}
		table.Append([]string{
			fmt.Sprintf("%d", project.ID),
			project.Name,
			description,
		})
	}

	table.Render()
}

func listProjectArtifacts(cmd *cobra.Command, args []string) {
	projectVersionID := args[0]
	
	cfg, err := config.LoadConfigWithOverrides(sscURL, sscToken)
	if err != nil {
		color.Red("Error loading configuration: %v", err)
		os.Exit(1)
	}

	sscClient := client.NewSSCClient(cfg)
	
	var pvID int64
	fmt.Sscanf(projectVersionID, "%d", &pvID)
	
	fmt.Printf("Fetching artifacts for project version %d...\n", pvID)
	artifacts, err := sscClient.GetArtifacts(pvID)
	if err != nil {
		color.Red("Error fetching artifacts: %v", err)
		os.Exit(1)
	}

	if len(artifacts) == 0 {
		color.Yellow("No artifacts found for this project version.")
		return
	}

	color.Green("\nFound %d artifacts:\n", len(artifacts))

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "File Name", "Status", "Upload Date", "Size (MB)", "Upload IP"})

	for _, artifact := range artifacts {
		uploadDate := ""
		if artifact.UploadDate != nil {
			uploadDate = artifact.UploadDate.Format("2006-01-02 15:04:05")
		}
		
		sizeInMB := fmt.Sprintf("%.2f", float64(artifact.FileSize)/(1024*1024))
		
		// Color status based on value
		status := artifact.Status
		if status == "REQUIRE_AUTH" || status == "Requires Approval" {
			status = color.RedString(status)
		} else if status == "PROCESSED" || status == "Complete" {
			status = color.GreenString(status)
		}

		table.Append([]string{
			fmt.Sprintf("%d", artifact.ID),
			artifact.FileName,
			status,
			uploadDate,
			sizeInMB,
			artifact.UploadIP,
		})
	}

	table.Render()
}