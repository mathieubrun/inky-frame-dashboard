from datetime import datetime
from typing import Annotated

from dishka.integrations.fastapi import FromDishka, inject
from fastapi import APIRouter, HTTPException, Query, status
from pydantic import BaseModel, Field

from src.core.battery.persistence import BatteryHistoryQuery
from src.core.battery.service import BatteryProcessor, BatteryReportInput

router = APIRouter(prefix="/api/battery", tags=["battery"])


class BatteryReportRequest(BaseModel):
    voltage: float


class BatteryReportResponse(BaseModel):
    voltage: float
    percentage: float
    timestamp: datetime


class BatteryHistoryRequest(BaseModel):
    limit: int = Field(10, ge=1, le=100)
    offset: int = Field(0, ge=0)


class BatteryStatusResponse(BaseModel):
    latest: BatteryReportResponse
    is_low: bool
    last_reported: datetime


@router.post("", response_model=BatteryReportResponse, status_code=status.HTTP_201_CREATED)
@inject
def report_battery(request: BatteryReportRequest, service: FromDishka[BatteryProcessor]):
    report = service.add_report(BatteryReportInput(voltage=request.voltage))
    return BatteryReportResponse(voltage=report.voltage, percentage=report.percentage, timestamp=report.timestamp)


@router.get("/status", response_model=BatteryStatusResponse)
@inject
def get_battery_status(service: FromDishka[BatteryProcessor]):
    status_data = service.get_latest_status()
    if not status_data:
        raise HTTPException(status_code=404, detail="No history available")

    return BatteryStatusResponse(
        latest=BatteryReportResponse(voltage=status_data.latest.voltage, percentage=status_data.latest.percentage, timestamp=status_data.latest.timestamp), is_low=status_data.is_low, last_reported=status_data.last_reported
    )


@router.get("/history", response_model=list[BatteryReportResponse])
@inject
def get_battery_history(request: Annotated[BatteryHistoryRequest, Query()], service: FromDishka[BatteryProcessor] = None):
    history = service.get_history(BatteryHistoryQuery(limit=request.limit, offset=request.offset))
    return [BatteryReportResponse(voltage=report.voltage, percentage=report.percentage, timestamp=report.timestamp) for report in history]
