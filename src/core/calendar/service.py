import logging
from datetime import datetime

from cachetools import TTLCache
from pydantic import BaseModel

from src.core.calendar.client import GoogleCalendarClient

logger = logging.getLogger(__name__)


class CalendarEvent(BaseModel):
    id: str
    summary: str
    start_time: datetime
    end_time: datetime


class CalendarCache(TTLCache):
    pass


CACHE_KEY = "calendar_events"


class CalendarProcessor:
    def __init__(self, client: GoogleCalendarClient, cache: CalendarCache) -> None:
        self.client = client
        self.cache = cache

    def get_upcoming_events(self) -> list[CalendarEvent]:
        if CACHE_KEY in self.cache:
            return self.cache[CACHE_KEY]

        logger.info("Fetching upcoming events from Google Calendar API")
        raw_events = self.client.fetch_upcoming_events()

        events: list[CalendarEvent] = []
        for item in raw_events:
            event_id = item.id
            summary = item.summary

            start_str = item.start.datetime or item.start.date
            end_str = item.end.datetime or item.end.date

            try:
                start_time = datetime.fromisoformat(str(start_str)) if start_str else datetime.min
                end_time = datetime.fromisoformat(str(end_str)) if end_str else datetime.max
            except (ValueError, TypeError):
                continue

            events.append(CalendarEvent(id=event_id, summary=summary, start_time=start_time, end_time=end_time))

        self.cache[CACHE_KEY] = events
        return events
