from unittest.mock import Mock

import httpx
import pytest
from pytest_mock import MockerFixture

from src.core.config import Settings
from src.core.exceptions import UpstreamAPIError
from src.core.weather.client import Coordinates, GeocodingRequest, OpenMeteoClient


class TestWeatherClient:
    @pytest.fixture
    def mock_settings(self, mocker: MockerFixture) -> Mock:
        mock = mocker.Mock(spec=Settings)
        mock.open_meteo_geocoding_url = "https://geocoding.example.com"
        mock.open_meteo_forecast_url = "https://forecast.example.com"
        return mock

    @pytest.fixture
    def sut(self, mock_settings: Mock) -> OpenMeteoClient:
        return OpenMeteoClient(settings=mock_settings)

    def test_get_coordinates_success(self, sut: OpenMeteoClient, mocker: MockerFixture) -> None:
        # Given
        mock_response = mocker.Mock()
        mock_response.json.return_value = {"results": [{"latitude": 47.37, "longitude": 8.54}]}
        mock_response.raise_for_status = mocker.Mock()

        mock_get = mocker.patch("httpx.Client.get", new_callable=Mock)
        mock_get.return_value = mock_response

        # When
        result = sut.get_coordinates(GeocodingRequest(city_name="Zurich"))

        # Then
        assert result is not None
        assert result.latitude == 47.37
        assert result.longitude == 8.54

    def test_get_coordinates_not_found(self, sut: OpenMeteoClient, mocker: MockerFixture) -> None:
        # Given
        mock_response = mocker.Mock()
        mock_response.json.return_value = {"results": []}
        mock_response.raise_for_status = mocker.Mock()

        mock_get = mocker.patch("httpx.Client.get", new_callable=Mock)
        mock_get.return_value = mock_response

        # When
        result = sut.get_coordinates(GeocodingRequest(city_name="Unknown"))

        # Then
        assert result is None

    def test_get_coordinates_request_error(self, sut: OpenMeteoClient, mocker: MockerFixture) -> None:
        # Given
        mock_get = mocker.patch("httpx.Client.get", new_callable=Mock)
        mock_get.side_effect = httpx.RequestError("Network Error")

        # When / Then
        with pytest.raises(UpstreamAPIError, match="Geocoding API request failed"):
            sut.get_coordinates(GeocodingRequest(city_name="Zurich"))

    def test_get_coordinates_http_error(self, sut: OpenMeteoClient, mocker: MockerFixture) -> None:
        # Given
        mock_response = mocker.Mock()
        mock_response.raise_for_status.side_effect = httpx.HTTPStatusError("404 Not Found", request=mocker.Mock(), response=mocker.Mock())

        mock_get = mocker.patch("httpx.Client.get", new_callable=Mock)
        mock_get.return_value = mock_response

        # When / Then
        with pytest.raises(UpstreamAPIError, match="Geocoding API returned error status"):
            sut.get_coordinates(GeocodingRequest(city_name="Zurich"))

    def test_get_forecast_success(self, sut: OpenMeteoClient, mocker: MockerFixture) -> None:
        # Given
        mock_response = mocker.Mock()
        mock_response.json.return_value = {"current_weather": {"temperature": 15.0}}
        mock_response.raise_for_status = mocker.Mock()

        mock_get = mocker.patch("httpx.Client.get", new_callable=Mock)
        mock_get.return_value = mock_response

        coords = Coordinates(latitude=47.37, longitude=8.54)

        # When
        result = sut.get_forecast(coords)

        # Then
        assert result.current_weather.temperature == 15.0
        mock_get.assert_called_once()
        args, kwargs = mock_get.call_args
        assert args[0] == "https://forecast.example.com"
        assert kwargs["params"]["latitude"] == 47.37
        assert kwargs["params"]["longitude"] == 8.54
        assert kwargs["params"]["models"] == "icon_seamless"
        assert kwargs["params"]["current_weather"] == "true"

    def test_get_forecast_request_error(self, sut: OpenMeteoClient, mocker: MockerFixture) -> None:
        # Given
        mock_get = mocker.patch("httpx.Client.get", new_callable=Mock)
        mock_get.side_effect = httpx.RequestError("Network Error")
        coords = Coordinates(latitude=47.37, longitude=8.54)

        # When / Then
        with pytest.raises(UpstreamAPIError, match="Forecast API request failed"):
            sut.get_forecast(coords)

    def test_get_forecast_http_error(self, sut: OpenMeteoClient, mocker: MockerFixture) -> None:
        # Given
        mock_response = mocker.Mock()
        mock_response.raise_for_status.side_effect = httpx.HTTPStatusError("500 Internal Error", request=mocker.Mock(), response=mocker.Mock())

        mock_get = mocker.patch("httpx.Client.get", new_callable=Mock)
        mock_get.return_value = mock_response
        coords = Coordinates(latitude=47.37, longitude=8.54)

        # When / Then
        with pytest.raises(UpstreamAPIError, match="Forecast API returned error status"):
            sut.get_forecast(coords)
