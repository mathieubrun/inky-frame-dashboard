from datetime import UTC, datetime

from google.oauth2 import service_account
from googleapiclient.discovery import build
from googleapiclient.errors import HttpError
from pydantic import BaseModel, Field

from src.core.calendar.exceptions import AuthenticationError
from src.core.config import Settings
from src.core.exceptions import UpstreamAPIError


class GoogleCalendarEventTime(BaseModel):
    datetime: str | None = Field(None, alias="dateTime")
    date: str | None = None


class GoogleCalendarEventItem(BaseModel):
    id: str = ""
    summary: str = "No Title"
    start: GoogleCalendarEventTime = Field(default_factory=GoogleCalendarEventTime)
    end: GoogleCalendarEventTime = Field(default_factory=GoogleCalendarEventTime)


class GoogleCalendarEventsResponse(BaseModel):
    items: list[GoogleCalendarEventItem] = Field(default_factory=list)


class GoogleCalendarClient:
    def __init__(self, settings: Settings) -> None:
        self.settings = settings
        try:
            self.credentials = service_account.Credentials.from_service_account_file(self.settings.google_calendar_api_key, scopes=["https://www.googleapis.com/auth/calendar.readonly"])
        except Exception as e:
            # If the file is missing or invalid, we catch it when making the request or log it
            self.credentials = None
            self._init_error = e

    def fetch_upcoming_events(self) -> list[GoogleCalendarEventItem]:
        now = datetime.now(UTC).isoformat()

        if not self.credentials:
            raise AuthenticationError(f"Google Calendar credentials invalid: {self._init_error}")

        try:
            service = build("calendar", "v3", credentials=self.credentials)
            events_result = service.events().list(calendarId=self.settings.google_calendar_id, timeMin=now, singleEvents=True, orderBy="startTime").execute()
        except Exception as e:
            if isinstance(e, HttpError):
                if e.resp.status in (401, 403):
                    raise AuthenticationError(f"Google Calendar API authentication failed: {e}") from e
                raise UpstreamAPIError(f"Google Calendar API returned error status: {e}") from e
            raise UpstreamAPIError(f"Google Calendar API request failed: {e}") from e

        return GoogleCalendarEventsResponse.model_validate(events_result).items
