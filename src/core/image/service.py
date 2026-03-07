from PIL import Image, ImageDraw
from pydantic import BaseModel

from src.core.calendar.service import CalendarProcessor
from src.core.weather.service import WeatherForecastRequest, WeatherProcessor


class WeatherImageRequest(BaseModel):
    city_name: str
    width: int
    height: int


class CalendarImageRequest(BaseModel):
    width: int
    height: int


class CombinedImageRequest(BaseModel):
    city_name: str
    width: int
    height: int


class ImageProcessor:
    def __init__(self, weather_service: WeatherProcessor, calendar_service: CalendarProcessor) -> None:
        self.weather_service = weather_service
        self.calendar_service = calendar_service

    def _draw_weather_image(self, forecast, width: int, height: int) -> Image.Image:
        img = Image.new("RGB", (width, height), color="white")
        draw = ImageDraw.Draw(img)
        text = f"Weather in {forecast.city_name}\nTemp: {forecast.current_weather.temperature}°C\nWind: {forecast.current_weather.windspeed} km/h"
        draw.text((10, 10), text, fill="black")
        return img

    def generate_weather_image(self, request: WeatherImageRequest) -> Image.Image:
        forecast = self.weather_service.get_weather_forecast(WeatherForecastRequest(city_name=request.city_name))
        return self._draw_weather_image(forecast, request.width, request.height)

    def _draw_calendar_image(self, events, width: int, height: int) -> Image.Image:
        img = Image.new("RGB", (width, height), color="white")
        draw = ImageDraw.Draw(img)
        text = "Upcoming Events:\n"
        for event in events[:5]:
            text += f"- {event.summary}\n"
        draw.text((10, 10), text, fill="black")
        return img

    def generate_calendar_image(self, request: CalendarImageRequest) -> Image.Image:
        events = self.calendar_service.get_upcoming_events()
        return self._draw_calendar_image(events, request.width, request.height)

    def _draw_combined_image(self, weather_img: Image.Image, calendar_img: Image.Image, width: int, height: int) -> Image.Image:
        half_width = width // 2
        img = Image.new("RGB", (width, height), color="white")
        img.paste(weather_img, (0, 0))
        img.paste(calendar_img, (half_width, 0))
        draw = ImageDraw.Draw(img)
        draw.line([(half_width, 0), (half_width, height)], fill="black", width=2)
        return img

    def generate_combined_image(self, request: CombinedImageRequest) -> Image.Image:
        half_width = request.width // 2
        weather_img = self.generate_weather_image(WeatherImageRequest(city_name=request.city_name, width=half_width, height=request.height))
        calendar_img = self.generate_calendar_image(CalendarImageRequest(width=half_width, height=request.height))
        return self._draw_combined_image(weather_img, calendar_img, request.width, request.height)
