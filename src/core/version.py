import importlib.metadata
import logging

logger = logging.getLogger(__name__)


class VersionProcessor:
    def get_version(self) -> str:
        try:
            ver = importlib.metadata.version("inky-frame-dashboard")
            logger.debug(f"Retrieved version from package metadata: {ver}")
            return ver
        except importlib.metadata.PackageNotFoundError:
            logger.warning("Package not installed, falling back to 3.0.0")
            return "3.0.0"
