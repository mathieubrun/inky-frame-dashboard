# Configuration Module Test Specification

## 1. Default Settings Loading
### Given
No environment variables are set (empty environment).
### When
The `Settings` class is instantiated.
### Then
The `data_dir` attribute is defaulted correctly to `.inky`.

## 2. Data Directory Creation (Parameterized)
### Given
A parameterized list of `data_dir` paths is provided, including:
  * A non-existent relative path.
  * A non-existent absolute path.
  * An existing path.
### When
The application is initialized or directory creation logic is explicitly triggered.
### Then
The specified directory is verified to exist on the filesystem post-initialization.

## 3. Environment Variable Parsing Override
### Given
The environment variable for `data_dir` (e.g., `DATA_DIR`) is set to a custom path (e.g., `/tmp/custom_inky`).
### When
The `Settings` class is instantiated.
### Then
The custom path specified in the environment variable is matched by the `data_dir` attribute, successfully overriding the default `.inky` value.
