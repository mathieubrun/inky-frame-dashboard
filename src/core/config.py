from pathlib import Path
from typing import Any

from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    data_dir: Path = Path(".inky")
    open_meteo_geocoding_url: str = "https://geocoding-api.open-meteo.com/v1/search"
    open_meteo_forecast_url: str = "https://api.open-meteo.com/v1/forecast"
    google_calendar_api_url: str = "https://www.googleapis.com/calendar/v3"
    google_calendar_id: str = "primary"
    google_calendar_api_key: str = ""
    weather_city_name: str = "London"
    weather_cache_ttl_seconds: int = 3600
    calendar_cache_ttl_seconds: int = 3600

    def __init__(self, **data: Any):
        super().__init__(**data)
        self.data_dir.mkdir(parents=True, exist_ok=True)
