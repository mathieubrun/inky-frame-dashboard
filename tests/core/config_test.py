from pathlib import Path

import pytest

from src.core.config import Settings


class TestSettingsDefault:
    @pytest.fixture
    def sut(self) -> Settings:
        return Settings()

    def test_default_settings_loading(self, sut: Settings) -> None:
        # Given
        # (Environment has no specific DATA_DIR set)

        # When
        # (sut is initialized)

        # Then
        assert sut.data_dir == Path(".inky")


class TestSettingsEnvironmentOverride:
    @pytest.fixture
    def sut(self, monkeypatch: pytest.MonkeyPatch) -> Settings:
        monkeypatch.setenv("DATA_DIR", "/tmp/custom_inky")
        return Settings()

    def test_env_override(self, sut: Settings) -> None:
        # Given
        # (Environment has DATA_DIR set via fixture)

        # When
        # (sut is initialized via fixture)

        # Then
        assert sut.data_dir == Path("/tmp/custom_inky")


class TestSettingsParameterized:
    @pytest.fixture
    def sut(self, request: pytest.FixtureRequest, tmp_path: Path) -> Settings:
        param_type = getattr(request, "param", "relative")
        if param_type == "relative":
            target = tmp_path / "new_relative"
        elif param_type == "absolute":
            target = Path(tmp_path / "new_absolute").resolve()
        else:
            target = tmp_path / "existing"
            target.mkdir()

        return Settings(data_dir=target)

    @pytest.mark.parametrize("sut", ["relative", "absolute", "existing"], indirect=True)
    def test_data_directory_creation(self, sut: Settings) -> None:
        # Given
        # (sut is parameterized with different paths)

        # When
        # (sut is initialized by fixture, triggering creation)

        # Then
        assert sut.data_dir.exists()
        assert sut.data_dir.is_dir()
