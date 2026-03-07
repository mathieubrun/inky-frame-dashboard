from datetime import UTC, datetime
from unittest.mock import Mock

import pytest
from pytest_mock import MockerFixture

from src.core.battery.persistence import BatteryHistoryQuery, BatteryPersistence, BatteryReport
from src.core.config import Settings


class TestBatteryPersistence:
    @pytest.fixture
    def mock_settings(self, mocker: MockerFixture, tmp_path) -> Mock:
        mock = mocker.Mock(spec=Settings)
        mock.data_dir = tmp_path
        return mock

    @pytest.fixture
    def sut(self, mock_settings: Mock) -> BatteryPersistence:
        return BatteryPersistence(settings=mock_settings)

    def test_save_and_get_latest_report(self, sut: BatteryPersistence) -> None:
        # Given
        report = BatteryReport(voltage=4.5, percentage=100.0, timestamp=datetime.now(UTC))

        # When
        sut.save_report(report)
        latest = sut.get_latest_report()

        # Then
        assert latest is not None
        assert latest.voltage == 4.5
        assert latest.percentage == 100.0
        assert latest.timestamp.isoformat() == report.timestamp.isoformat()

    def test_get_history(self, sut: BatteryPersistence) -> None:
        # Given
        report1 = BatteryReport(voltage=3.0, percentage=0.0, timestamp=datetime.now(UTC))
        report2 = BatteryReport(voltage=4.5, percentage=100.0, timestamp=datetime.now(UTC))
        sut.save_report(report1)
        sut.save_report(report2)

        # When
        history = sut.get_history(BatteryHistoryQuery(limit=10, offset=0))

        # Then
        assert len(history) == 2
        # History is returned descending (newest first)
        assert history[0].voltage == 4.5
        assert history[1].voltage == 3.0

    def test_get_latest_report_empty(self, sut: BatteryPersistence) -> None:
        # Given empty persistence
        # When
        latest = sut.get_latest_report()
        # Then
        assert latest is None

    def test_get_history_empty(self, sut: BatteryPersistence) -> None:
        # Given empty persistence
        # When
        history = sut.get_history(BatteryHistoryQuery(limit=10, offset=0))
        # Then
        assert history == []
