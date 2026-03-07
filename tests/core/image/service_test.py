from datetime import UTC, datetime
from unittest.mock import Mock

import pytest
from PIL import Image
from pytest_mock import MockerFixture

from src.core.calendar.exceptions import AuthenticationError
from src.core.calendar.service import CalendarEvent, CalendarProcessor
from src.core.exceptions import UpstreamAPIError
from src.core.image.service import CalendarImageRequest, CombinedImageRequest, ImageProcessor, WeatherImageRequest
from src.core.weather.exceptions import CityNotFoundError
from src.core.weather.service import CurrentWeather, WeatherForecast, WeatherForecastRequest, WeatherProcessor


class TestImageService:
    @pytest.fixture
    def mock_weather_service(self, mocker: MockerFixture) -> Mock:
        mock = mocker.Mock(spec=WeatherProcessor)
        mock.get_weather_forecast = Mock()
        return mock

    @pytest.fixture
    def mock_calendar_service(self, mocker: MockerFixture) -> Mock:
        mock = mocker.Mock(spec=CalendarProcessor)
        mock.get_upcoming_events = Mock()
        return mock

    @pytest.fixture
    def sut(self, mock_weather_service: Mock, mock_calendar_service: Mock) -> ImageProcessor:
        return ImageProcessor(weather_service=mock_weather_service, calendar_service=mock_calendar_service)

    def test_render_weather_image_success(self, sut: ImageProcessor, mock_weather_service: Mock) -> None:
        # Given
        forecast = WeatherForecast(city_name="Zurich", current_weather=CurrentWeather(temperature=15.0, windspeed=10.0, weathercode=1))
        mock_weather_service.get_weather_forecast.return_value = forecast

        # When
        img = sut.generate_weather_image(WeatherImageRequest(city_name="Zurich", width=400, height=300))

        # Then
        mock_weather_service.get_weather_forecast.assert_called_once_with(WeatherForecastRequest(city_name="Zurich"))
        assert isinstance(img, Image.Image)
        assert img.width == 400
        assert img.height == 300

    def test_render_calendar_image_success(self, sut: ImageProcessor, mock_calendar_service: Mock) -> None:
        # Given
        events = [CalendarEvent(id="1", summary="Test", start_time=datetime.now(UTC), end_time=datetime.now(UTC))]
        mock_calendar_service.get_upcoming_events.return_value = events

        # When
        img = sut.generate_calendar_image(CalendarImageRequest(width=400, height=300))

        # Then
        mock_calendar_service.get_upcoming_events.assert_called_once()
        assert isinstance(img, Image.Image)
        assert img.width == 400
        assert img.height == 300

    def test_render_combined_image_success(self, sut: ImageProcessor, mock_weather_service: Mock, mock_calendar_service: Mock) -> None:
        # Given
        forecast = WeatherForecast(city_name="Zurich", current_weather=CurrentWeather(temperature=15.0, windspeed=10.0, weathercode=1))
        events = [CalendarEvent(id="1", summary="Test", start_time=datetime.now(UTC), end_time=datetime.now(UTC))]
        mock_weather_service.get_weather_forecast.return_value = forecast
        mock_calendar_service.get_upcoming_events.return_value = events

        # When
        img = sut.generate_combined_image(CombinedImageRequest(city_name="Zurich", width=400, height=300))

        # Then
        mock_weather_service.get_weather_forecast.assert_called_once_with(WeatherForecastRequest(city_name="Zurich"))
        mock_calendar_service.get_upcoming_events.assert_called_once()
        assert isinstance(img, Image.Image)
        assert img.width == 400
        assert img.height == 300

    @pytest.mark.parametrize("service_to_fail, exception_type", [("weather", CityNotFoundError), ("weather", UpstreamAPIError), ("calendar", AuthenticationError), ("calendar", UpstreamAPIError)])
    def test_upstream_failure_handling(self, sut: ImageProcessor, mock_weather_service: Mock, mock_calendar_service: Mock, service_to_fail: str, exception_type: type[Exception]) -> None:
        # Given
        if service_to_fail == "weather":
            mock_weather_service.get_weather_forecast.side_effect = exception_type("Weather failed")
            # Mock calendar to avoid gather failing due to calendar if weather fails first (gather raises the first exception)
            mock_calendar_service.get_upcoming_events.return_value = []
        else:
            mock_weather_service.get_weather_forecast.return_value = WeatherForecast(city_name="Zurich", current_weather=CurrentWeather(temperature=15.0, windspeed=10.0, weathercode=1))
            mock_calendar_service.get_upcoming_events.side_effect = exception_type("Calendar failed")

        # When / Then
        with pytest.raises(exception_type):
            # We also test single functions if appropriate, but combined image covers both
            sut.generate_combined_image(CombinedImageRequest(city_name="Zurich", width=400, height=300))
