import logging

from cachetools import TTLCache
from pydantic import BaseModel

from src.core.weather.client import GeocodingRequest, OpenMeteoClient
from src.core.weather.exceptions import CityNotFoundError

logger = logging.getLogger(__name__)


class CurrentWeather(BaseModel):
    temperature: float
    windspeed: float
    weathercode: int


class WeatherForecast(BaseModel):
    city_name: str
    current_weather: CurrentWeather


class WeatherForecastRequest(BaseModel):
    city_name: str


class WeatherCache(TTLCache):
    pass


class WeatherProcessor:
    def __init__(self, client: OpenMeteoClient, cache: WeatherCache) -> None:
        self.client = client
        self.cache = cache

    def get_weather_forecast(self, request: WeatherForecastRequest) -> WeatherForecast:
        if request.city_name in self.cache:
            return self.cache[request.city_name]  # type: ignore

        logger.info(f"Fetching coordinates from Open-Meteo API for city: {request.city_name}")
        coords = self.client.get_coordinates(GeocodingRequest(city_name=request.city_name))
        if not coords:
            raise CityNotFoundError(f"City not found: {request.city_name}")

        logger.info(f"Fetching forecast from Open-Meteo API for coordinates: {coords}")
        forecast_data = self.client.get_forecast(coords)
        current = forecast_data.current_weather

        forecast = WeatherForecast(city_name=request.city_name, current_weather=CurrentWeather(temperature=current.temperature, windspeed=current.windspeed, weathercode=current.weathercode))

        self.cache[request.city_name] = forecast
        return forecast
