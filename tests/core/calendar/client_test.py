from unittest.mock import Mock

import httplib2
import pytest
from googleapiclient.errors import HttpError
from pytest_mock import MockerFixture

from src.core.calendar.client import GoogleCalendarClient
from src.core.calendar.exceptions import AuthenticationError
from src.core.config import Settings
from src.core.exceptions import UpstreamAPIError


class TestCalendarClient:
    @pytest.fixture
    def mock_settings(self, mocker: MockerFixture) -> Mock:
        mock = mocker.Mock(spec=Settings)
        mock.google_calendar_api_url = "https://calendar.example.com"
        mock.google_calendar_id = "test_calendar"
        mock.google_calendar_api_key = "test_creds.json"
        return mock

    @pytest.fixture
    def mock_credentials(self, mocker: MockerFixture) -> Mock:
        mock_creds = mocker.Mock()
        mock_creds.valid = True
        mock_creds.token = "mock_jwt_token"

        mocker.patch("src.core.calendar.client.service_account.Credentials.from_service_account_file", return_value=mock_creds)
        return mock_creds

    @pytest.fixture
    def sut(self, mock_settings: Mock, mock_credentials: Mock) -> GoogleCalendarClient:
        return GoogleCalendarClient(settings=mock_settings)

    def test_fetch_upcoming_events_success(self, sut: GoogleCalendarClient, mocker: MockerFixture) -> None:
        # Given
        mock_build = mocker.patch("src.core.calendar.client.build", new_callable=Mock)
        mock_service = mocker.Mock()
        mock_build.return_value = mock_service
        
        mock_events = mock_service.events.return_value
        mock_list = mock_events.list.return_value
        mock_list.execute.return_value = {"items": [{"id": "1", "summary": "Test"}]}

        # When
        result = sut.fetch_upcoming_events()

        # Then
        assert len(result) == 1
        assert result[0].id == "1"
        assert result[0].summary == "Test"
        mock_events.list.assert_called_once()
        args, kwargs = mock_events.list.call_args
        assert kwargs["calendarId"] == "test_calendar"
        assert "timeMin" in kwargs
        assert kwargs["singleEvents"] is True
        assert kwargs["orderBy"] == "startTime"

    def test_fetch_upcoming_events_request_error(self, sut: GoogleCalendarClient, mocker: MockerFixture) -> None:
        # Given
        mock_build = mocker.patch("src.core.calendar.client.build", new_callable=Mock)
        mock_service = mocker.Mock()
        mock_build.return_value = mock_service
        
        mock_events = mock_service.events.return_value
        mock_list = mock_events.list.return_value
        mock_list.execute.side_effect = Exception("Network Error")

        # When / Then
        with pytest.raises(UpstreamAPIError, match="Google Calendar API request failed"):
            sut.fetch_upcoming_events()

    @pytest.mark.parametrize(
        "status_code, expected_exception",
        [
            (401, AuthenticationError),
            (403, AuthenticationError),
            (500, UpstreamAPIError),
            (404, UpstreamAPIError),
        ],
    )
    def test_fetch_upcoming_events_http_error(self, sut: GoogleCalendarClient, mocker: MockerFixture, status_code: int, expected_exception: type[Exception]) -> None:
        # Given
        mock_build = mocker.patch("src.core.calendar.client.build", new_callable=Mock)
        mock_service = mocker.Mock()
        mock_build.return_value = mock_service
        
        mock_events = mock_service.events.return_value
        mock_list = mock_events.list.return_value
        
        resp = httplib2.Response({"status": str(status_code)})
        mock_list.execute.side_effect = HttpError(resp, b"Error")

        # When / Then
        with pytest.raises(expected_exception):
            sut.fetch_upcoming_events()
