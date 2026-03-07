from datetime import UTC, datetime
from unittest.mock import Mock

import pytest
from pytest_mock import MockerFixture

from src.core.battery.persistence import BatteryHistoryQuery, BatteryPersistence, BatteryReport
from src.core.battery.service import BatteryProcessor, BatteryReportInput
from src.core.config import Settings


class TestBatteryService:
    @pytest.fixture
    def mock_settings(self, mocker: MockerFixture) -> Mock:
        mock = mocker.Mock(spec=Settings)
        return mock

    @pytest.fixture
    def mock_persistence(self, mocker: MockerFixture) -> Mock:
        mock = mocker.Mock(spec=BatteryPersistence)
        mock.save_report = Mock()
        mock.get_latest_report = Mock()
        mock.get_history = Mock()
        return mock

    @pytest.fixture
    def sut(self, mock_persistence: Mock) -> BatteryProcessor:
        return BatteryProcessor(persistence=mock_persistence)

    @pytest.mark.parametrize(
        "voltage, expected",
        [
            (2.5, 0.0),
            (3.0, 0.0),
            (3.75, 50.0),
            (4.5, 100.0),
            (5.0, 100.0),
        ],
    )
    def test_calculate_percentage(self, sut: BatteryProcessor, voltage: float, expected: float) -> None:
        # Given
        # Parameterized inputs
        # When
        result = sut.calculate_percentage(voltage)

        # Then
        assert result == expected

    def test_add_report(self, sut: BatteryProcessor, mock_persistence: Mock, mock_settings: Mock) -> None:
        # Given
        # A mocked BatteryPersistence layer are provided, with the Configuration module mocked to provide a temporary data_dir.
        voltage = 4.5

        # When
        result = sut.add_report(BatteryReportInput(voltage=voltage))

        # Then
        assert result.voltage == 4.5
        assert result.percentage == 100.0
        mock_persistence.save_report.assert_called_once_with(result)

    @pytest.mark.parametrize(
        "mocked_report, expected_is_low",
        [
            (None, None),
            (BatteryReport(voltage=3.3, percentage=19.9, timestamp=datetime.now(UTC)), True),
            (BatteryReport(voltage=3.31, percentage=20.0, timestamp=datetime.now(UTC)), False),
        ],
    )
    def test_get_latest_status(self, sut: BatteryProcessor, mock_persistence: Mock, mock_settings: Mock, mocked_report: BatteryReport | None, expected_is_low: bool | None) -> None:
        # Given
        # A parameterized set of returns from the mocked BatteryPersistence.get_latest_report() is provided, and the Configuration module is mocked
        mock_persistence.get_latest_report.return_value = mocked_report

        # When
        result = sut.get_latest_status()

        # Then
        if mocked_report is None:
            assert result is None
        else:
            assert result is not None
            assert result.is_low is expected_is_low
            assert result.latest == mocked_report
            assert result.last_reported == mocked_report.timestamp

    @pytest.mark.parametrize("limit, offset", [(10, 0), (0, 5), (100, 150)])
    def test_get_history(self, sut: BatteryProcessor, mock_persistence: Mock, mock_settings: Mock, limit: int, offset: int) -> None:
        # Given
        # A mocked BatteryPersistence layer containing historical report data is provided, with the Configuration module mocked
        mock_reports = [BatteryReport(voltage=4.0, percentage=66.6, timestamp=datetime.now(UTC))]
        mock_persistence.get_history.return_value = mock_reports

        # When
        result = sut.get_history(BatteryHistoryQuery(limit=limit, offset=offset))

        # Then
        assert result == mock_reports
        mock_persistence.get_history.assert_called_once_with(BatteryHistoryQuery(limit=limit, offset=offset))
