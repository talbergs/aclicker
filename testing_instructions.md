**Testing Instructions (with emphasis on boundary conditions and assertions):**

To perform comprehensive testing of all mechanics, upgrades, and endings, follow these steps:

1.  **Build the game:**
    ```bash
    go build -o clicker2 main.go
    ```
2.  **Run the game:**
    ```bash
    ./clicker2
    ```
3.  **Test Core Mechanics:**
    *   **Click the rock:**
        *   Verify `TheRock.Health` decreases by `ThePlayer.Damage` on each click.
        *   Verify `ThePlayer.Dust` increases by 1 on each click.
        *   **Assertion:** Ensure health does not go below 0.
    *   **Observe rock messages:**
        *   Verify messages appear after clicks.
        *   Verify messages disappear after approximately 3 seconds.
        *   **Assertion:** Check that `g.state.CurrentRockMessage` becomes empty after the timer expires.
    *   **Observe rock cracking (Visual Feedback):**
        *   **Boundary Condition:** Reduce rock health to just below 75%, 50%, and 25% of `InitialRockHealth`.
        *   **Assertion:** Verify the rock sprite visually changes to `RockSpriteCracked1`, `RockSpriteCracked2`, and `RockSpriteShattered` respectively at these thresholds.
    *   **Observe environmental decay (Visual Feedback):**
        *   **Boundary Condition:** Observe the background aesthetics as rock health crosses the 75%, 50%, and 25% thresholds.
        *   **Assertion:** Verify the background visually changes (e.g., colors become more decayed) in correlation with `HealthPercentage` passed to the shader.

4.  **Test Upgrades:**
    *   **Purchase "Stronger Pickaxe":**
        *   **Boundary Condition:** Attempt to purchase when `ThePlayer.Dust` is exactly `UpgradeCost` (for level 0).
        *   **Assertion:** Verify `ThePlayer.Damage` increases by 1. Verify `ThePlayer.Dust` decreases by `UpgradeCost`.
        *   **Boundary Condition:** Attempt to purchase when `ThePlayer.Dust` is less than `UpgradeCost`.
        *   **Assertion:** Verify purchase fails and `ErrorSFXPlayer` plays.
        *   **Boundary Condition:** Purchase "Stronger Pickaxe" until `MaxLevel` (5). Attempt to purchase again.
        *   **Assertion:** Verify purchase fails and `ErrorSFXPlayer` plays.
    *   **Purchase "Auto-Clicker v0.1":**
        *   **Assertion:** Verify `g.state.AutoClickerActive` becomes `true` and `g.state.AutoClickerRate` is 1.
        *   Observe the rock being clicked automatically.
    *   **Purchase "Auto-Clicker v1.0":**
        *   **Assertion:** Verify `g.state.AutoClickerActive` remains `true` and `g.state.AutoClickerRate` increases to 5.
        *   Observe the rock being clicked faster automatically.
    *   **Verify upgrade buttons:**
        *   **Assertion:** All defined upgrades (`stronger_pickaxe`, `auto_clicker_v0_1`, `heart_of_the_mountain`) are listed in the UI.
        *   **Assertion:** Their names, current levels, and costs are displayed correctly.

5.  **Test Audio:**
    *   **Background Music:**
        *   **Boundary Condition:** Observe music when `TheRock.Health` is near `InitialRockHealth` (healthy).
        *   **Assertion:** `HealthyMusicPlayer` is dominant.
        *   **Boundary Condition:** Observe music when `TheRock.Health` is near 0 (decayed).
        *   **Assertion:** `MelancholicMusicPlayer` is dominant.
        *   **Assertion:** Verify smooth crossfade between tracks as health changes.
    *   **Sound Effects:**
        *   **Assertion:** `ClickSFXPlayer` plays on each rock click.
        *   **Assertion:** `UpgradeSFXPlayer` plays on successful upgrade purchases.
        *   **Assertion:** `ErrorSFXPlayer` plays on failed upgrade purchases.
    *   **Volume Control:**
        *   Scroll mouse wheel up/down.
        *   **Assertion:** Verify music volume increases/decreases.
        *   **Boundary Condition:** Scroll to max/min volume.
        *   **Assertion:** Volume does not exceed 1.0 or go below 0.0.

6.  **Test Save/Load:**
    *   Play for a bit, making some progress (e.g., purchase an upgrade, reduce rock health).
    *   Press 'S' to save.
    *   Quit the game ('Q').
    *   Run the game again, then press 'L' to load.
    *   **Assertion:** Verify that `TheRock.Health`, `ThePlayer.Dust`, `ThePlayer.Damage`, `UpgradeManager.PlayerUpgrades` (levels of purchased upgrades), `CurrentRockMessage`, and `RockMessageTimer` are restored correctly.

7.  **Test Endings:**
    *   **"The Heart of the Mountain" purchase:**
        *   Play until you can afford "heart_of_the_mountain".
        *   Purchase it.
        *   **Assertion:** Verify `g.state.EndGameChoicePending` becomes `true`.
        *   **Assertion:** Verify the end-game choice buttons ("Take the Heart", "Let it Rest") appear and the specific message is displayed.
    *   **"Take the Heart" (Bad Ending):**
        *   Click "Take the Heart".
        *   **Assertion:** Verify the game exits.
        *   Upon re-running the game (without loading a previous save), if the game state persists the ending (e.g., by checking a flag in a separate file or if the save file itself reflects the ending), **Assertion:** Verify the rock is shattered (`TheRock.Health` is 0) and the "bad ending" message is displayed.
    *   **"Let it Rest" (Good Ending):**
        *   Restart the game and reach the end-game choice again.
        *   Click "Let it Rest".
        *   **Assertion:** Verify the game exits.
        *   Upon re-running the game, **Assertion:** Verify the rock is at peace (e.g., `g.state.GameWon` is true, and the "good ending" message is displayed, potentially with a visual change like a flower).

This comprehensive testing approach, focusing on specific assertions and boundary conditions, will ensure the game mechanics are working as intended.