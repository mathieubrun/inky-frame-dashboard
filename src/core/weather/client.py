from typing import Any

import httpx
from pydantic import BaseModel, Field

from src.core.config import Settings
from src.core.exceptions import UpstreamAPIError


class Coordinates(BaseModel):
    latitude: float
    longitude: float


class GeocodingRequest(BaseModel):
    city_name: str


class GeocodingResponse(BaseModel):
    results: list[Coordinates] = Field(default_factory=list)


class OpenMeteoCurrentWeather(BaseModel):
    temperature: float = 0.0
    windspeed: float = 0.0
    weathercode: int = 0


class OpenMeteoForecastResponse(BaseModel):
    current_weather: OpenMeteoCurrentWeather = Field(default_factory=OpenMeteoCurrentWeather)


class OpenMeteoClient:
    def __init__(self, settings: Settings) -> None:
        self.settings = settings

    def get_coordinates(self, request: GeocodingRequest) -> Coordinates | None:
        with httpx.Client() as client:
            try:
                response = client.get(self.settings.open_meteo_geocoding_url, params={"name": request.city_name, "count": 1, "language": "en", "format": "json"})
                response.raise_for_status()
            except httpx.RequestError as e:
                raise UpstreamAPIError(f"Geocoding API request failed: {e}") from e
            except httpx.HTTPStatusError as e:
                raise UpstreamAPIError(f"Geocoding API returned error status: {e}") from e

            data = response.json()
            geocoding_response = GeocodingResponse.model_validate(data)
            if not geocoding_response.results:
                return None
            return geocoding_response.results[0]

    def get_forecast(self, coords: Coordinates) -> OpenMeteoForecastResponse:
        with httpx.Client() as client:
            try:
                response = client.get(self.settings.open_meteo_forecast_url, params={"latitude": coords.latitude, "longitude": coords.longitude, "models": "icon_seamless", "current_weather": "true"})
                response.raise_for_status()
            except httpx.RequestError as e:
                raise UpstreamAPIError(f"Forecast API request failed: {e}") from e
            except httpx.HTTPStatusError as e:
                raise UpstreamAPIError(f"Forecast API returned error status: {e}") from e

            data = response.json()
            return OpenMeteoForecastResponse.model_validate(data)
