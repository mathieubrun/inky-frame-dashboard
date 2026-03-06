# Quickstart: Google Agenda Display

This guide helps you set up and test the Google Agenda integration.

## 1. Google Cloud Setup

1. Create a project in the [Google Cloud Console](https://console.cloud.google.com/).
2. Enable the **Google Calendar API**.
3. Create a **Service Account** and download the JSON key file.
4. Open your Google Calendar, go to **Settings and sharing**, and share the calendar with the Service Account email (e.g., `dashboard@my-project.iam.gserviceaccount.com`).

## 2. Server Configuration

Add the following to your configuration or environment:
- `GOOGLE_CREDENTIALS_PATH`: Path to your JSON key file.
- `AGENDA_ID`: Your Google Calendar ID (e.g., `your-email@gmail.com`).

## 3. Testing with CLI

1. **List events**:
```bash
inky agenda list --count 5
```

2. **Generate Combined Dashboard**:
```bash
inky dashboard image --location "Zurich" --output dashboard.png
```

## 4. Testing with API

1. **Start the server**:
```bash
inky serve
```

2. **Retrieve JSON data**:
```bash
curl "http://localhost:8080/api/v1/agenda"
```

3. **Fetch full image**:
```bash
curl -o dash.png "http://localhost:8080/api/v1/dashboard/image?location=Zurich"
```
