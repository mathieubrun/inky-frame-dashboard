import logging

import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient

from src.api.logger_config import RequestIdFilter, RequestIdMiddleware, request_id_context_var


class TestRequestIdFilter:
    @pytest.fixture
    def sut(self) -> RequestIdFilter:
        return RequestIdFilter()

    def test_filter_injects_request_id(self, sut: RequestIdFilter) -> None:
        # Given
        record = logging.LogRecord(name="test", level=logging.INFO, pathname="", lineno=0, msg="test", args=(), exc_info=None)
        token = request_id_context_var.set("test-123")

        try:
            # When
            result = sut.filter(record)

            # Then
            assert result is True
            assert getattr(record, "request_id", None) == "test-123"
        finally:
            request_id_context_var.reset(token)


class TestRequestIdMiddleware:
    @pytest.fixture
    def inner_app(self) -> FastAPI:
        app = FastAPI()

        @app.get("/test")
        def test_endpoint() -> dict[str, str]:
            return {"request_id": request_id_context_var.get()}

        return app

    @pytest.fixture
    def sut(self, inner_app: FastAPI) -> RequestIdMiddleware:
        return RequestIdMiddleware(app=inner_app)

    def test_middleware_injects_and_returns_request_id(self, sut: RequestIdMiddleware) -> None:
        # Given
        with TestClient(app=sut, base_url="http://test") as client:
            # When
            response = client.get("/test")

            # Then
            assert response.status_code == 200
            data = response.json()
            request_id = data["request_id"]

            assert len(request_id) == 8
            assert response.headers["X-Request-ID"] == request_id
