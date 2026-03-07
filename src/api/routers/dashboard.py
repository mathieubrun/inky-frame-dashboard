import io
from typing import Annotated

from dishka.integrations.fastapi import FromDishka, inject
from fastapi import APIRouter, Query, Response
from pydantic import BaseModel, Field

from src.core.config import Settings
from src.core.image.service import CalendarImageRequest, CombinedImageRequest, ImageProcessor, WeatherImageRequest


class DashboardImageRequest(BaseModel):
    width: int = Field(600, gt=0)
    height: int = Field(448, gt=0)


router = APIRouter(prefix="/api/dashboard", tags=["dashboard"])


@router.get("/weather", response_class=Response)
@inject
def get_weather_image(request: Annotated[DashboardImageRequest, Query()], service: FromDishka[ImageProcessor] = None, settings: FromDishka[Settings] = None):
    img = service.generate_weather_image(WeatherImageRequest(city_name=settings.weather_city_name, width=request.width, height=request.height))
    buf = io.BytesIO()
    img.save(buf, format="PNG")
    return Response(content=buf.getvalue(), media_type="image/png")


@router.get("/calendar", response_class=Response)
@inject
def get_calendar_image(request: Annotated[DashboardImageRequest, Query()], service: FromDishka[ImageProcessor] = None):
    img = service.generate_calendar_image(CalendarImageRequest(width=request.width, height=request.height))
    buf = io.BytesIO()
    img.save(buf, format="PNG")
    return Response(content=buf.getvalue(), media_type="image/png")


@router.get("/combined", response_class=Response)
@inject
def get_combined_image(request: Annotated[DashboardImageRequest, Query()], service: FromDishka[ImageProcessor] = None, settings: FromDishka[Settings] = None):
    img = service.generate_combined_image(CombinedImageRequest(city_name=settings.weather_city_name, width=request.width, height=request.height))
    buf = io.BytesIO()
    img.save(buf, format="PNG")
    return Response(content=buf.getvalue(), media_type="image/png")
