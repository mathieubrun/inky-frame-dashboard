import pytest

from src.core.version import VersionProcessor


class TestVersionProcessor:
    @pytest.fixture
    def sut(self) -> VersionProcessor:
        return VersionProcessor()

    def test_get_version(self, sut: VersionProcessor) -> None:
        # Given
        # Assuming package is installed or falls back

        # When
        result = sut.get_version()

        # Then
        assert isinstance(result, str)
        assert len(result) > 0
