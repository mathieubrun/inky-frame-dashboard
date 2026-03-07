from dishka.integrations.fastapi import FromDishka, inject
from fastapi import APIRouter
from pydantic import BaseModel

from src.core.version import VersionProcessor

router = APIRouter(prefix="/api/version", tags=["version"])


class VersionResponse(BaseModel):
    version: str


@router.get("", response_model=VersionResponse)
@inject
def get_version(processor: FromDishka[VersionProcessor]) -> VersionResponse:
    version = processor.get_version()
    return VersionResponse(version=version)
