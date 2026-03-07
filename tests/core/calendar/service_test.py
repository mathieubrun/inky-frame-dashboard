from datetime import datetime
from unittest.mock import Mock

import pytest
from pytest_mock import MockerFixture

from src.core.calendar.client import GoogleCalendarClient, GoogleCalendarEventItem, GoogleCalendarEventTime
from src.core.calendar.exceptions import AuthenticationError
from src.core.calendar.service import CACHE_KEY, CalendarCache, CalendarProcessor
from src.core.exceptions import UpstreamAPIError


class TestCalendarService:
    @pytest.fixture
    def mock_time(self, mocker: MockerFixture) -> Mock:
        return mocker.Mock(return_value=1000.0)

    @pytest.fixture
    def mock_cache(self, mock_time: Mock) -> CalendarCache:
        return CalendarCache(maxsize=1, ttl=3600, timer=mock_time)

    @pytest.fixture
    def mock_client(self, mocker: MockerFixture) -> Mock:
        mock = mocker.Mock(spec=GoogleCalendarClient)
        mock.fetch_upcoming_events = Mock()
        return mock

    @pytest.fixture
    def sut(self, mock_client: Mock, mock_cache: CalendarCache) -> CalendarProcessor:
        return CalendarProcessor(client=mock_client, cache=mock_cache)

    def test_get_upcoming_events_success(self, sut: CalendarProcessor, mock_client: Mock) -> None:
        # Given
        mock_events = [
            GoogleCalendarEventItem(id="123", summary="Team Meeting", start=GoogleCalendarEventTime(dateTime="2026-06-06T10:00:00Z"), end=GoogleCalendarEventTime(dateTime="2026-06-06T11:00:00Z")),
            GoogleCalendarEventItem(id="124", summary="All Day Event", start=GoogleCalendarEventTime(date="2026-06-07"), end=GoogleCalendarEventTime(date="2026-06-08")),
        ]
        mock_client.fetch_upcoming_events.return_value = mock_events

        # When
        result = sut.get_upcoming_events()

        # Then
        assert len(result) == 2
        assert result[0].id == "123"
        assert result[0].summary == "Team Meeting"
        assert result[0].start_time == datetime.fromisoformat("2026-06-06T10:00:00Z")
        assert result[1].id == "124"
        assert result[1].summary == "All Day Event"
        assert result[1].start_time == datetime.fromisoformat("2026-06-07")
        mock_client.fetch_upcoming_events.assert_called_once()

    def test_in_memory_caching_success(self, sut: CalendarProcessor, mock_client: Mock) -> None:
        # Given
        mock_events = [GoogleCalendarEventItem(id="123", summary="Team Meeting", start=GoogleCalendarEventTime(dateTime="2026-06-06T10:00:00Z"), end=GoogleCalendarEventTime(dateTime="2026-06-06T11:00:00Z"))]
        mock_client.fetch_upcoming_events.return_value = mock_events

        # When
        result1 = sut.get_upcoming_events()
        result2 = sut.get_upcoming_events()
        result3 = sut.get_upcoming_events()

        # Then
        assert len(result1) == 1
        assert result2 == result1
        assert result3 == result1
        mock_client.fetch_upcoming_events.assert_called_once()

    def test_cache_expiration(self, sut: CalendarProcessor, mock_client: Mock, mock_time: Mock) -> None:
        # Given
        mock_events_1 = [GoogleCalendarEventItem(id="123", summary="Team Meeting 1")]
        mock_client.fetch_upcoming_events.return_value = mock_events_1

        # When
        result1 = sut.get_upcoming_events()

        # Advance time beyond TTL
        mock_time.return_value = 1000.0 + 3600 + 1.0

        mock_events_2 = [GoogleCalendarEventItem(id="124", summary="Team Meeting 2")]
        mock_client.fetch_upcoming_events.return_value = mock_events_2

        result2 = sut.get_upcoming_events()

        # Then
        assert result1[0].summary == "Team Meeting 1"
        assert result2[0].summary == "Team Meeting 2"
        assert mock_client.fetch_upcoming_events.call_count == 2

    def test_get_upcoming_events_empty(self, sut: CalendarProcessor, mock_client: Mock, mock_cache: CalendarCache) -> None:
        # Given
        mock_client.fetch_upcoming_events.return_value = []

        # When
        result = sut.get_upcoming_events()

        # Then
        assert result == []
        mock_client.fetch_upcoming_events.assert_called_once()
        assert CACHE_KEY in mock_cache

    @pytest.mark.parametrize(
        "exception_type, exception_msg", [(AuthenticationError, "Google Calendar API authentication failed"), (UpstreamAPIError, "Google Calendar API returned error status: 500"), (UpstreamAPIError, "Google Calendar API request failed")]
    )
    def test_get_upcoming_events_upstream_failure(self, sut: CalendarProcessor, mock_client: Mock, mock_cache: CalendarCache, exception_type: type[Exception], exception_msg: str) -> None:
        # Given
        mock_client.fetch_upcoming_events.side_effect = exception_type(exception_msg)

        # When / Then
        with pytest.raises(exception_type, match=exception_msg):
            sut.get_upcoming_events()

        mock_client.fetch_upcoming_events.assert_called_once()
        assert CACHE_KEY not in mock_cache
