from unittest.mock import Mock

import pytest
from dishka import Provider, Scope, make_async_container, provide
from fastapi.testclient import TestClient
from PIL import Image
from pytest_mock import MockerFixture

from src.api.main import create_app
from src.core.calendar.exceptions import AuthenticationError
from src.core.exceptions import UpstreamAPIError
from src.core.image.service import CalendarImageRequest, CombinedImageRequest, ImageProcessor, WeatherImageRequest
from src.core.ioc import AppProvider
from src.core.weather.exceptions import CityNotFoundError


class MockDashboardProvider(Provider):
    def __init__(self, mock_service):
        super().__init__()
        self.mock_service = mock_service

    @provide(scope=Scope.APP)
    def get_service(self) -> ImageProcessor:
        return self.mock_service


class TestApiDashboard:
    @pytest.fixture
    def mock_service(self, mocker: MockerFixture) -> Mock:
        mock = mocker.Mock(spec=ImageProcessor)

        # Create a dummy image
        dummy_img = Image.new("RGB", (10, 10))
        mock.generate_weather_image = Mock(return_value=dummy_img)
        mock.generate_calendar_image = Mock(return_value=dummy_img)
        mock.generate_combined_image = Mock(return_value=dummy_img)

        return mock

    @pytest.fixture
    def sut(self, mock_service: Mock) -> TestClient:
        return TestClient(create_app(make_async_container(AppProvider(), MockDashboardProvider(mock_service))))

    def test_get_weather_image(self, sut: TestClient, mock_service: Mock) -> None:
        # Given
        # When
        response = sut.get("/api/dashboard/weather?width=600&height=448")

        # Then
        assert response.status_code == 200
        assert response.headers["content-type"] == "image/png"
        mock_service.generate_weather_image.assert_called_once_with(WeatherImageRequest(city_name="London", width=600, height=448))

    def test_get_calendar_image(self, sut: TestClient, mock_service: Mock) -> None:
        # Given
        # When
        response = sut.get("/api/dashboard/calendar?width=600&height=448")

        # Then
        assert response.status_code == 200
        assert response.headers["content-type"] == "image/png"
        mock_service.generate_calendar_image.assert_called_once_with(CalendarImageRequest(width=600, height=448))

    def test_get_combined_image(self, sut: TestClient, mock_service: Mock) -> None:
        # Given
        # When
        response = sut.get("/api/dashboard/combined?width=600&height=448")

        # Then
        assert response.status_code == 200
        assert response.headers["content-type"] == "image/png"
        mock_service.generate_combined_image.assert_called_once_with(CombinedImageRequest(city_name="London", width=600, height=448))

    @pytest.mark.parametrize(
        "endpoint, params",
        [
            ("/api/dashboard/weather", {"width": -10}),  # invalid width
            ("/api/dashboard/calendar", {"height": 0}),  # invalid height
            ("/api/dashboard/combined", {"width": -100}),  # invalid width
        ],
    )
    def test_invalid_payloads(self, sut: TestClient, endpoint: str, params: dict) -> None:
        # Given
        # When
        response = sut.get(endpoint, params=params)

        # Then
        assert response.status_code == 422  # Unprocessable Entity

    @pytest.mark.parametrize("endpoint", ["/api/dashboard/weather", "/api/dashboard/calendar", "/api/dashboard/combined"])
    def test_invalid_method(self, sut: TestClient, endpoint: str) -> None:
        # Given
        # When
        response = sut.post(endpoint)

        # Then
        assert response.status_code == 405  # Method Not Allowed

    @pytest.mark.parametrize("exception_class, expected_status", [(CityNotFoundError, 404), (AuthenticationError, 401), (UpstreamAPIError, 502)])
    def test_domain_exception_mapping(self, sut: TestClient, mock_service: Mock, exception_class: type[Exception], expected_status: int) -> None:
        # Given
        mock_service.generate_combined_image.side_effect = exception_class("Domain error occurred")

        # When
        response = sut.get("/api/dashboard/combined")

        # Then
        assert response.status_code == expected_status
        assert response.json()["detail"] == "Domain error occurred"
