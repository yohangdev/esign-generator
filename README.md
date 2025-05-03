# Electronic Signature Label Generator

A simple HTTP service that generates electronic signature label images on demand. This service creates professional-looking signatures with rounded borders that can be used in electronic documents.

**Demo/Preview:**

<img src="https://github.com/user-attachments/assets/9fad424d-4bc9-47d5-95d6-291eaafef6b4" height="110">

## Features

- Dynamically generates signature images with professional appearance.
- Customizable name and job title parameters.
- Rounded border design with consistent formatting.
- Returns JPEG or PNG images ready for embedding in documents.
- Detailed logging for all requests.

## Requirements

- Go 1.24 or higher.
- Port 8080 available for the HTTP server.

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

## API Documentation

### Generate Signature Image

Generates a signature image with the specified parameters.

**Endpoint:** `GET /generate`

**Parameters:**

| Parameter | Type   | Required | Description                                              |
|-----------|--------|----------|----------------------------------------------------------|
| nama      | string | Yes      | The name to be displayed in the signature.               |
| jabatan   | string | No       | The job title/position (will be displayed in uppercase). |
| pangkat   | string | No       | The rank/grade to be displayed below the name.           |
| format    | string | No       | Image format: 'jpeg' (default) or 'png'.                 |

**Response:**
- Content-Type: `image/jpeg` or `image/png` depending on the format parameter.
- Image dimensions: 1489x485 pixels.

**Example Request:**
```
GET /generate?nama=John%20Doe&jabatan=Software%20Engineer&pangkat=Penyedia&format=jpeg
```

**Error Responses:**

- 400 Bad Request
   - Returned when required parameters are missing
   - Response format:
     ```json
     {
       "error": "Missing required parameters",
       "missing": ["nama"]
     }
     ```

- 500 Internal Server Error
   - Returned when image encoding fails

## Image Specifications

- Width: 1489 pixels
- Height: 485 pixels
- Background: White
- Border: Black rounded border (4px width, 70px radius)
- Font: Arial
- Format: JPEG (default, 100% quality) or PNG
