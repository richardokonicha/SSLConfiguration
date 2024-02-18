# SSL Configuration Analysis Web Application

The SSL Configuration Analysis web application is a powerful tool designed to help users evaluate and understand the SSL/TLS security configurations of websites.

Leveraging the comprehensive SSL Labs API, this application performs in-depth analyses of SSL certificates and configurations, offering users a detailed report that includes various security aspects such as certificate validity, encryption strength, and potential vulnerabilities.

## Core Features

SSL Analysis: Users can enter the hostname of any website to analyze its SSL/TLS configuration. The application then queries the SSL Labs API to fetch detailed data about the site's SSL setup.

PDF Report Generation: After the analysis, the application generates a PDF report summarizing the SSL configuration, including grades, warnings, and other critical details. This report provides actionable insights into the website's security posture.

User-Friendly Interface: With a clean and intuitive web interface built using Tailwind CSS, users can easily navigate the application. The interface includes a simple form for inputting the hostname and buttons to view or download the generated PDF report.

## Features

- Analyze SSL configurations for specified hosts.
- Generate PDF reports summarizing the SSL configuration, including grades, warnings, and other pertinent details.
- Download or view the generated PDF reports through a user-friendly web interface.

## Technology Stack

- **Backend**: Written in Go (Golang), leveraging the `gofpdf` package for PDF generation.
- **Frontend**: Utilizes Tailwind CSS for styling and HTMX for dynamic content loading without full page refreshes.
- **Deployment**: Containerized with Docker for easy deployment and scalability.

## Prerequisites

- Go (Golang) 1.20
- Docker

## Local Setup and Running

**Clone the Repository**

```bash
git clone repository-url
cd repository-folder
```

### Run the application

go run .

The web server starts on the default port 8080

Access the application by navigating to http://localhost:8080 in your web browser

### Build the Docker image

docker build -t sslgo-web .

### Run the container from the image

docker run --name sslgo-web-container -p 8080:8080 sslgo-web

Access the application by navigating to http://localhost:8080 in your web browser
