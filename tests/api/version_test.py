from unittest.mock import Mock

import pytest
from dishka import Provider, Scope, make_async_container, provide
from fastapi.testclient import TestClient
from pytest_mock import MockerFixture

from src.api.main import app
from src.core.ioc import AppProvider
from src.core.version import VersionProcessor


class MockVersionProvider(Provider):
    def __init__(self, mock_processor):
        super().__init__()
        self.mock_processor = mock_processor

    @provide(scope=Scope.APP)
    def get_processor(self) -> VersionProcessor:
        return self.mock_processor


class TestApiVersion:
    @pytest.fixture
    def mock_processor(self, mocker: MockerFixture) -> Mock:
        mock = mocker.Mock(spec=VersionProcessor)
        mock.get_version.return_value = "1.0.0"
        return mock

    @pytest.fixture
    def sut(self, mock_processor: Mock) -> TestClient:
        old_container = app.state.dishka_container
        app.state.dishka_container = make_async_container(AppProvider(), MockVersionProvider(mock_processor))
        yield TestClient(app)
        app.state.dishka_container = old_container

    def test_get_version(self, sut: TestClient) -> None:
        # Given
        # When
        response = sut.get("/api/version")

        # Then
        assert response.status_code == 200
        assert "version" in response.json()

    def test_get_version_invalid_method(self, sut: TestClient) -> None:
        # Given
        # When
        response = sut.post("/api/version")

        # Then
        assert response.status_code == 405
