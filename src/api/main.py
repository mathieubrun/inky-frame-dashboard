import logging
from contextlib import asynccontextmanager

from dishka import AsyncContainer, Scope
from dishka.integrations.fastapi import setup_dishka
from fastapi import FastAPI, Request
from fastapi.responses import JSONResponse

from src.api.logger_config import RequestIdMiddleware
from src.api.routers import battery, dashboard, version
from src.core.calendar.exceptions import AuthenticationError
from src.core.exceptions import DomainError, UpstreamAPIError
from src.core.ioc import AppProvider
from src.core.weather.exceptions import CityNotFoundError


@asynccontextmanager
async def lifespan(app: FastAPI):
    logging.debug("app : start")
    yield
    logging.debug("app : close container")
    app.state.dishka_container.close()
    logging.debug("app : stop")


def create_app(container: AsyncContainer) -> FastAPI:

    app = FastAPI(title="Inky Frame Dashboard API", version="3.0.0", lifespan=lifespan)
    setup_dishka(container, app)

    app.include_router(version.router)
    app.include_router(battery.router)
    app.include_router(dashboard.router)

    app.add_middleware(RequestIdMiddleware)

    @app.exception_handler(CityNotFoundError)
    def city_not_found_error_handler(request: Request, exc: CityNotFoundError) -> JSONResponse:
        return JSONResponse(status_code=404, content={"detail": str(exc)})

    @app.exception_handler(AuthenticationError)
    def auth_error_handler(request: Request, exc: AuthenticationError) -> JSONResponse:
        return JSONResponse(status_code=401, content={"detail": str(exc)})

    @app.exception_handler(UpstreamAPIError)
    def upstream_error_handler(request: Request, exc: UpstreamAPIError) -> JSONResponse:
        return JSONResponse(status_code=502, content={"detail": str(exc)})

    @app.exception_handler(DomainError)
    def domain_error_handler(request: Request, exc: DomainError) -> JSONResponse:
        return JSONResponse(
            status_code=400,
            content={"detail": str(exc)},
        )

    return app


app = create_app(AppProvider.create_container(Scope.REQUEST))
