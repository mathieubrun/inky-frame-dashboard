from datetime import UTC, datetime

from pydantic import BaseModel

from src.core.battery.persistence import BatteryHistoryQuery, BatteryPersistence, BatteryReport


class BatteryStatus(BaseModel):
    latest: BatteryReport
    is_low: bool
    last_reported: datetime


class BatteryReportInput(BaseModel):
    voltage: float


class BatteryProcessor:
    def __init__(self, persistence: BatteryPersistence) -> None:
        self.persistence = persistence

    def calculate_percentage(self, voltage: float) -> float:
        percentage = (voltage - 3.0) / (4.5 - 3.0) * 100.0
        return max(0.0, min(100.0, percentage))

    def add_report(self, input_data: BatteryReportInput) -> BatteryReport:
        percentage = self.calculate_percentage(input_data.voltage)
        report = BatteryReport(voltage=input_data.voltage, percentage=percentage, timestamp=datetime.now(UTC))
        self.persistence.save_report(report)
        return report

    def get_latest_status(self) -> BatteryStatus | None:
        latest = self.persistence.get_latest_report()
        if not latest:
            return None

        is_low = latest.percentage < 20.0
        return BatteryStatus(latest=latest, is_low=is_low, last_reported=latest.timestamp)

    def get_history(self, query: BatteryHistoryQuery) -> list[BatteryReport]:
        return self.persistence.get_history(query)
