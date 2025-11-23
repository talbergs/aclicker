1.  **Navigability Improvements**
    *   [ ] Move `upgrades.go` into its own `game/upgrades` subpackage.
    *   [ ] Create a `game/state` subpackage to encapsulate the `Game`, `Rock`, and `Player` structs.

2.  **Readability Improvements**
    *   [ ] Extract debug functionality into a separate `debug` package/file.
    *   [ ] Refactor `handleInput()` into smaller, more focused functions.

3.  **System Stability and Resilience Improvements**
    *   [ ] Replace direct `os.Exit(0)` calls with a mechanism to signal game termination to the main loop.
    *   [ ] Implement robust error handling for file operations (save/load) with user feedback.
    *   [ ] Add bounds checking or validation for configuration values.