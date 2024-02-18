package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// APIServiceInfo represents information about the SSL Labs API service.
type APIServiceInfo struct {
	EngineVersion        string   `json:"engineVersion"`
	CriteriaVersion      string   `json:"criteriaVersion"`
	ClientMaxAssessments int      `json:"clientMaxAssessments"`
	Messages             []string `json:"messages"`
}

// Endpoint represents an SSL endpoint.
type Endpoint struct {
	IPAddress         string `json:"ipAddress"`
	ServerName        string `json:"serverName"`
	StatusMessage     string `json:"statusMessage"`
	Grade             string `json:"grade"`
	GradeTrustIgnored string `json:"gradeTrustIgnored"`
	HasWarnings       bool   `json:"hasWarnings"`
	IsExceptional     bool   `json:"isExceptional"`
}

// Response represents the response from the SSL Labs API.
type Response struct {
	Host            string     `json:"host"`
	Port            int        `json:"port"`
	Protocol        string     `json:"protocol"`
	IsPublic        bool       `json:"isPublic"`
	Status          string     `json:"status"`
	StartTime       int64      `json:"startTime"`
	TestTime        int64      `json:"testTime"`
	EngineVersion   string     `json:"engineVersion"`
	CriteriaVersion string     `json:"criteriaVersion"`
	Endpoints       []Endpoint `json:"endpoints"`
}

type AnalysisStatus struct {
	Completed bool
	Error     string
	Result    string 
  Name      string
}

const (
    defaultServerPort = ":8080"
    sslLabsAPIBaseURL = "https://api.ssllabs.com/api/v3/analyze"
		filesBaseURL = "files/"
)

func main() {
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
			serverPort = defaultServerPort
	}
	if !checkAPIService() {
		log.Println("SSL Labs API service is not available. Exiting.")
		return
	}
  // go cleanupGeneratedPDFs()
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/analyze", analyzeHandler)
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("files"))))
	fmt.Printf("Server starting on %s\n", serverPort)
  log.Fatal(http.ListenAndServe(serverPort, nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, nil)
}

func analyzeHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    host := r.FormValue("host")
    if host == "" {
        http.Error(w, "Host is required", http.StatusBadRequest)
        return
    }

    response, err := fetchSSLData(host)
    if err != nil {
        log.Printf("Failed to fetch SSL data for %s: %v", host, err)
        http.Error(w, "Failed to fetch SSL data", http.StatusInternalServerError)
        return
    }

    pdfName, err := generatePDF(response)

    if err != nil {
        log.Printf("Failed to generate PDF for %s: %v", host, err)
        http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
        return
    }

		pdfPath := "http://" + r.Host + "/" + pdfName
    w.Header().Set("Content-Type", "application/json")
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.ExecuteTemplate(w, "result", AnalysisStatus{Result: pdfPath, Completed: true, Error: "", Name: pdfName})
}

func checkAPIService() bool {
	resp, err := http.Get("https://api.ssllabs.com/api/v2/info")
	if err != nil {
		log.Fatalf("Failed to fetch API service info: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	var info APIServiceInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		log.Fatalf("Failed to decode API service info: %v\n", err)
		return false
	}

	fmt.Println("SSL Labs API Service Information:")
	fmt.Printf("Engine Version: %s\n", info.EngineVersion)
	fmt.Printf("Criteria Version: %s\n", info.CriteriaVersion)
	for _, message := range info.Messages {
		fmt.Println(message)
	}

	return true
}

func fetchSSLData(host string) (Response, error) {
    var response Response

    url := fmt.Sprintf("%s?host=%s", sslLabsAPIBaseURL, host)
    resp, err := http.Get(url)
    if err != nil {
        log.Printf("Error executing request to SSL Labs API: %v\n", err)
        return response, err
    }
    defer resp.Body.Close()

    err = json.NewDecoder(resp.Body).Decode(&response)
    if err != nil {
        log.Printf("Error decoding response from SSL Labs API: %v\n", err)
        return response, err
    }

    return response, nil
}

func generatePDF(response Response) (string, error) {
    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.AddPage()

    // Add title
    pdf.SetFont("Arial", "B", 16)
    pdf.Cell(0, 10, "SSL Labs API Response")
    pdf.Ln(20) // Increase the line height for better spacing

    // Add host information
    pdf.SetFont("Arial", "", 12)
    pdf.Cell(0, 10, fmt.Sprintf("Host: %s", response.Host))
    pdf.Ln(12) // Adjusted line height for consistency

    // Add timestamp
    pdf.SetFont("Arial", "", 10)
    pdf.Cell(0, 10, "Timestamp: "+time.Now().Format("2006-01-02 15:04:05"))
    pdf.Ln(15) // Adjust for spacing

    // Adding SSL data as a structured table
    pdf.SetFont("Arial", "", 10)
    pdf.SetFillColor(240, 240, 240) // Grey background for header
    pdf.SetTextColor(0, 0, 0)       // Black text
    pdf.SetDrawColor(0, 0, 0)       // Black border
    pdf.SetLineWidth(.3)

    // Table headers
    header := []string{"IP Address", "Server Name", "Status Message", "Grade", "Grade Trust Ignored", "Has Warnings", "Is Exceptional"}
    w := []float64{40, 35, 40, 15, 20, 20, 20} // Column widths
    for i, str := range header {
        pdf.CellFormat(w[i], 7, str, "1", 0, "C", true, 0, "")
    }
    pdf.Ln(-1)

    // Table data
    for _, ep := range response.Endpoints {
        pdf.SetFillColor(224, 235, 255) // Light blue background for rows
        pdf.CellFormat(w[0], 6, ep.IPAddress, "1", 0, "L", false, 0, "")
        pdf.CellFormat(w[1], 6, ep.ServerName, "1", 0, "L", false, 0, "")
        pdf.CellFormat(w[2], 6, ep.StatusMessage, "1", 0, "L", false, 0, "")
        pdf.CellFormat(w[3], 6, ep.Grade, "1", 0, "L", false, 0, "")
        pdf.CellFormat(w[4], 6, ep.GradeTrustIgnored, "1", 0, "L", false, 0, "")
        pdf.CellFormat(w[5], 6, fmt.Sprintf("%t", ep.HasWarnings), "1", 0, "L", false, 0, "")
        pdf.CellFormat(w[6], 6, fmt.Sprintf("%t", ep.IsExceptional), "1", 0, "L", false, 0, "")
        pdf.Ln(-1)
    }

		filename := fmt.Sprintf(filesBaseURL + "ssl_report_%s.pdf", response.Host)
    err := pdf.OutputFileAndClose(filename)
    if err != nil {
        return "", fmt.Errorf("error generating PDF: %v", err)
    }
		
  	pdfURL := filename
    return pdfURL, nil
}
