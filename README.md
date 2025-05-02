# Electronic Signature Label Generator

A simple HTTP service that generates electronic signature label images on demand. This service creates professional-looking signatures with rounded borders that can be used in electronic documents.

## Features

- Dynamically generates signature images with professional appearance
- Customizable name and job title parameters
- Rounded border design with consistent formatting
- Returns JPEG images ready for embedding in documents
- Detailed logging for all requests

## Requirements

- Go 1.24 or higher
- Port 8080 available for the HTTP server

## Installation

1. Clone this repository:
   ```
   git clone https://github.com/yohangdev/esign-generator.git
   ```

2. Build and run the application:
   ```
   cd src/
   go build -o esign-generator
   ./esign-generator
   ```
