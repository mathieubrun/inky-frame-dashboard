from unittest.mock import Mock

import pytest
from pytest_mock import MockerFixture

from src.core.exceptions import UpstreamAPIError
from src.core.weather.client import Coordinates, GeocodingRequest, OpenMeteoClient, OpenMeteoCurrentWeather, OpenMeteoForecastResponse
from src.core.weather.exceptions import CityNotFoundError
from src.core.weather.service import WeatherCache, WeatherForecastRequest, WeatherProcessor


class TestWeatherService:
    @pytest.fixture
    def mock_time(self, mocker: MockerFixture) -> Mock:
        return mocker.Mock(return_value=1000.0)

    @pytest.fixture
    def mock_cache(self, mock_time: Mock) -> WeatherCache:
        return WeatherCache(maxsize=100, ttl=3600, timer=mock_time)

    @pytest.fixture
    def mock_client(self, mocker: MockerFixture) -> Mock:
        mock = mocker.Mock(spec=OpenMeteoClient)
        mock.get_coordinates = Mock()
        mock.get_forecast = Mock()
        return mock

    @pytest.fixture
    def sut(self, mock_client: Mock, mock_cache: WeatherCache) -> WeatherProcessor:
        return WeatherProcessor(client=mock_client, cache=mock_cache)

    def test_get_forecast_success(self, sut: WeatherProcessor, mock_client: Mock) -> None:
        # Given
        city_name = "Zurich"
        mock_coords = Coordinates(latitude=47.37, longitude=8.54)
        mock_client.get_coordinates.return_value = mock_coords

        mock_forecast_data = OpenMeteoForecastResponse(current_weather=OpenMeteoCurrentWeather(temperature=15.5, windspeed=10.2, weathercode=3))
        mock_client.get_forecast.return_value = mock_forecast_data

        # When
        result = sut.get_weather_forecast(WeatherForecastRequest(city_name=city_name))

        # Then
        assert result.city_name == city_name
        assert result.current_weather.temperature == 15.5
        assert result.current_weather.windspeed == 10.2
        assert result.current_weather.weathercode == 3
        mock_client.get_coordinates.assert_called_once_with(GeocodingRequest(city_name=city_name))
        mock_client.get_forecast.assert_called_once_with(mock_coords)

    def test_in_memory_caching_success(self, sut: WeatherProcessor, mock_client: Mock) -> None:
        # Given
        city_name = "Zurich"
        mock_coords = Coordinates(latitude=47.37, longitude=8.54)
        mock_client.get_coordinates.return_value = mock_coords
        mock_forecast_data = OpenMeteoForecastResponse(current_weather=OpenMeteoCurrentWeather(temperature=15.5, windspeed=10.2, weathercode=3))
        mock_client.get_forecast.return_value = mock_forecast_data

        # When
        result1 = sut.get_weather_forecast(WeatherForecastRequest(city_name=city_name))
        result2 = sut.get_weather_forecast(WeatherForecastRequest(city_name=city_name))
        result3 = sut.get_weather_forecast(WeatherForecastRequest(city_name=city_name))

        # Then
        assert result1.city_name == city_name
        assert result2 == result1
        assert result3 == result1
        mock_client.get_coordinates.assert_called_once_with(GeocodingRequest(city_name=city_name))
        mock_client.get_forecast.assert_called_once_with(mock_coords)

    def test_cache_expiration(self, sut: WeatherProcessor, mock_client: Mock, mock_time: Mock) -> None:
        # Given
        city_name = "Zurich"
        mock_coords = Coordinates(latitude=47.37, longitude=8.54)
        mock_client.get_coordinates.return_value = mock_coords
        mock_forecast_data = OpenMeteoForecastResponse(current_weather=OpenMeteoCurrentWeather(temperature=15.5))
        mock_client.get_forecast.return_value = mock_forecast_data

        # When
        result1 = sut.get_weather_forecast(WeatherForecastRequest(city_name=city_name))

        # Advance time beyond TTL
        mock_time.return_value = 1000.0 + 3600 + 1.0

        mock_forecast_data_new = OpenMeteoForecastResponse(current_weather=OpenMeteoCurrentWeather(temperature=18.0))
        mock_client.get_forecast.return_value = mock_forecast_data_new
        result2 = sut.get_weather_forecast(WeatherForecastRequest(city_name=city_name))

        # Then
        assert result1.current_weather.temperature == 15.5
        assert result2.current_weather.temperature == 18.0
        assert mock_client.get_coordinates.call_count == 2
        assert mock_client.get_forecast.call_count == 2

    @pytest.mark.parametrize("invalid_city", ["", "NonExistentCity123", "!@#$%"])
    def test_get_forecast_unknown_city(self, sut: WeatherProcessor, mock_client: Mock, mock_cache: WeatherCache, invalid_city: str) -> None:
        # Given
        mock_client.get_coordinates.return_value = None

        # When / Then
        with pytest.raises(CityNotFoundError):
            sut.get_weather_forecast(WeatherForecastRequest(city_name=invalid_city))

        mock_client.get_coordinates.assert_called_once_with(GeocodingRequest(city_name=invalid_city))
        mock_client.get_forecast.assert_not_called()
        assert invalid_city not in mock_cache

    def test_get_forecast_upstream_failure(self, sut: WeatherProcessor, mock_client: Mock, mock_cache: WeatherCache) -> None:
        # Given
        city_name = "Zurich"
        mock_coords = Coordinates(latitude=47.37, longitude=8.54)
        mock_client.get_coordinates.return_value = mock_coords
        mock_client.get_forecast.side_effect = UpstreamAPIError("Forecast API returned error status: 500")

        # When / Then
        with pytest.raises(UpstreamAPIError):
            sut.get_weather_forecast(WeatherForecastRequest(city_name=city_name))

        mock_client.get_coordinates.assert_called_once_with(GeocodingRequest(city_name=city_name))
        mock_client.get_forecast.assert_called_once_with(mock_coords)
        assert city_name not in mock_cache
