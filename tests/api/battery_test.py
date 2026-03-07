from datetime import UTC, datetime
from unittest.mock import Mock

import pytest
from dishka import Provider, Scope, make_async_container, provide
from fastapi.testclient import TestClient
from pytest_mock import MockerFixture

from src.api.main import create_app
from src.core.battery.persistence import BatteryHistoryQuery, BatteryReport
from src.core.battery.service import BatteryProcessor, BatteryReportInput, BatteryStatus
from src.core.ioc import AppProvider


class MockBatteryProvider(Provider):
    def __init__(self, mock_service):
        super().__init__()
        self.mock_service = mock_service

    @provide(scope=Scope.APP)
    def get_service(self) -> BatteryProcessor:
        return self.mock_service


class TestApiBattery:
    @pytest.fixture
    def mock_service(self, mocker: MockerFixture) -> Mock:
        mock = mocker.Mock(spec=BatteryProcessor)
        mock.add_report = Mock()
        mock.get_latest_status = Mock()
        mock.get_history = Mock()
        return mock

    @pytest.fixture
    def sut(self, mock_service: Mock) -> TestClient:
        return TestClient(create_app(make_async_container(AppProvider(), MockBatteryProvider(mock_service))))

    def test_post_report(self, sut: TestClient, mock_service: Mock) -> None:
        # Given
        payload = {"voltage": 4.5}
        mock_report = BatteryReport(voltage=4.5, percentage=100.0, timestamp=datetime.now(UTC))
        mock_service.add_report.return_value = mock_report

        # When
        response = sut.post("/api/battery", json=payload)

        # Then
        assert response.status_code == 201
        data = response.json()
        assert data["voltage"] == 4.5
        assert data["percentage"] == 100.0
        mock_service.add_report.assert_called_once_with(BatteryReportInput(voltage=4.5))

    @pytest.mark.parametrize("payload", [{}, {"voltage": "high"}, {"wrong_field": 4.5}])
    def test_post_report_invalid(self, sut: TestClient, mock_service: Mock, payload: dict) -> None:
        # Given
        # (Parameterized invalid payload)

        # When
        response = sut.post("/api/battery", json=payload)

        # Then
        assert response.status_code == 422
        mock_service.add_report.assert_not_called()

    def test_get_status_empty(self, sut: TestClient, mock_service: Mock) -> None:
        # Given
        mock_service.get_latest_status.return_value = None

        # When
        response = sut.get("/api/battery/status")

        # Then
        assert response.status_code == 404

    def test_get_status(self, sut: TestClient, mock_service: Mock) -> None:
        # Given
        mock_report = BatteryReport(voltage=3.75, percentage=50.0, timestamp=datetime.now(UTC))
        mock_status = BatteryStatus(latest=mock_report, is_low=False, last_reported=mock_report.timestamp)
        mock_service.get_latest_status.return_value = mock_status

        # When
        response = sut.get("/api/battery/status")

        # Then
        assert response.status_code == 200
        data = response.json()
        assert data["latest"]["voltage"] == 3.75
        assert data["is_low"] is False

    def test_get_status_invalid_method(self, sut: TestClient, mock_service: Mock) -> None:
        # Given
        # A mocked BatteryService is provided. A malformed request (unsupported HTTP method) is prepared.

        # When
        response = sut.post("/api/battery/status")

        # Then
        assert response.status_code == 405

    @pytest.mark.parametrize("limit, offset", [(10, 5), (100, 0), (1, 1000)])
    def test_get_history(self, sut: TestClient, mock_service: Mock, limit: int, offset: int) -> None:
        # Given
        mock_report = BatteryReport(voltage=4.2, percentage=80.0, timestamp=datetime.now(UTC))
        mock_service.get_history.return_value = [mock_report]

        # When
        response = sut.get(f"/api/battery/history?limit={limit}&offset={offset}")

        # Then
        assert response.status_code == 200
        data = response.json()
        assert isinstance(data, list)
        assert len(data) == 1
        assert data[0]["voltage"] == 4.2
        mock_service.get_history.assert_called_once_with(BatteryHistoryQuery(limit=limit, offset=offset))

    @pytest.mark.parametrize("limit, offset", [(0, 0), (150, 0), (10, -1)])
    def test_get_history_invalid(self, sut: TestClient, mock_service: Mock, limit: int, offset: int) -> None:
        # Given
        # Invalid parameters are provided.

        # When
        response = sut.get(f"/api/battery/history?limit={limit}&offset={offset}")

        # Then
        assert response.status_code == 422
        mock_service.get_history.assert_not_called()
