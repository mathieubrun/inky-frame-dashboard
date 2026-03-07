from datetime import datetime

from pydantic import BaseModel

from src.core.config import Settings


class BatteryReport(BaseModel):
    voltage: float
    percentage: float
    timestamp: datetime


class BatteryHistoryQuery(BaseModel):
    limit: int
    offset: int


class BatteryPersistence:
    def __init__(self, settings: Settings) -> None:
        self.file_path = settings.data_dir / "battery_history.csv"

    def save_report(self, report: BatteryReport) -> None:
        file_exists = self.file_path.exists()
        with self.file_path.open(mode="a", newline="") as f:
            if not file_exists:
                f.write("voltage,percentage,timestamp\n")
            f.write(f"{report.voltage},{report.percentage},{report.timestamp.isoformat()}\n")

    def get_latest_report(self) -> BatteryReport | None:
        if not self.file_path.exists():
            return None

        last_line = ""
        with self.file_path.open(newline="") as f:
            lines = f.readlines()
            if len(lines) <= 1:
                return None
            last_line = lines[-1]

        parts = last_line.strip().split(",")
        if len(parts) == 3:
            return BatteryReport(voltage=float(parts[0]), percentage=float(parts[1]), timestamp=datetime.fromisoformat(parts[2]))
        return None

    def get_history(self, query: BatteryHistoryQuery) -> list[BatteryReport]:
        if not self.file_path.exists() or query.limit <= 0 or query.offset < 0:
            return []

        reports: list[BatteryReport] = []
        with self.file_path.open(newline="") as f:
            lines = f.readlines()
            if len(lines) <= 1:
                return []
            data_lines = lines[1:]

        data_lines.reverse()

        start = query.offset
        end = query.offset + query.limit
        paginated_lines = data_lines[start:end]

        for line in paginated_lines:
            parts = line.strip().split(",")
            if len(parts) == 3:
                reports.append(BatteryReport(voltage=float(parts[0]), percentage=float(parts[1]), timestamp=datetime.fromisoformat(parts[2])))
        return reports
