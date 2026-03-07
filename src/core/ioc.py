from dishka import AsyncContainer, Provider, Scope, make_async_container, provide

from src.core.battery.persistence import BatteryPersistence
from src.core.battery.service import BatteryProcessor
from src.core.calendar.client import GoogleCalendarClient
from src.core.calendar.service import CalendarCache, CalendarProcessor
from src.core.config import Settings
from src.core.image.service import ImageProcessor
from src.core.version import VersionProcessor
from src.core.weather.client import OpenMeteoClient
from src.core.weather.service import WeatherCache, WeatherProcessor


class AppProvider(Provider):
    @classmethod
    def create_container(cls, scope: Scope) -> AsyncContainer:
        return make_async_container(cls(scope=scope))

    @provide(scope=Scope.APP)
    def get_settings(self) -> Settings:
        return Settings()

    @provide(scope=Scope.APP)
    def get_weather_cache(self, settings: Settings) -> WeatherCache:
        return WeatherCache(maxsize=100, ttl=settings.weather_cache_ttl_seconds)

    @provide(scope=Scope.APP)
    def get_calendar_cache(self, settings: Settings) -> CalendarCache:
        return CalendarCache(maxsize=1, ttl=settings.calendar_cache_ttl_seconds)

    @provide(scope=Scope.REQUEST)
    def get_version_processor(self) -> VersionProcessor:
        return VersionProcessor()

    @provide(scope=Scope.REQUEST)
    def get_battery_persistence(self, settings: Settings) -> BatteryPersistence:
        return BatteryPersistence(settings=settings)

    @provide(scope=Scope.REQUEST)
    def get_battery_processor(self, persistence: BatteryPersistence) -> BatteryProcessor:
        return BatteryProcessor(persistence=persistence)

    @provide(scope=Scope.REQUEST)
    def get_open_meteo_client(self, settings: Settings) -> OpenMeteoClient:
        return OpenMeteoClient(settings)

    @provide(scope=Scope.REQUEST)
    def get_google_calendar_client(self, settings: Settings) -> GoogleCalendarClient:
        return GoogleCalendarClient(settings)

    @provide(scope=Scope.REQUEST)
    def get_weather_processor(self, client: OpenMeteoClient, cache: WeatherCache) -> WeatherProcessor:
        return WeatherProcessor(client, cache)

    @provide(scope=Scope.REQUEST)
    def get_calendar_processor(self, client: GoogleCalendarClient, cache: CalendarCache) -> CalendarProcessor:
        return CalendarProcessor(client, cache)

    @provide(scope=Scope.REQUEST)
    def get_image_processor(self, weather_service: WeatherProcessor, calendar_service: CalendarProcessor) -> ImageProcessor:
        return ImageProcessor(weather_service, calendar_service)
