### Core Game Mechanic: The Friendly Rock

The central element of the game is a "friendly rock." Players click on the rock to mine "dust," the primary in-game currency. However, each click also depletes the rock's health. The player's goal is to gather enough resources to achieve a state where they can choose to stop mining and let the rock "live," effectively winning the game by closing it.

### Game Loop

1.  **Click:** The player clicks the friendly rock.
2.  **Gather:** Each click generates "dust" and slightly decreases the rock's health.
3.  **Upgrade:** The player spends dust on upgrades.
4.  **Decide:** The player's progress and upgrades lead them towards a final decision: continue mining to oblivion or stop and "save" the rock.

### Currency: Dust

*   **Dust:** The primary resource, gathered by clicking the rock. Used to purchase all upgrades.

### Upgrade Paths

Upgrades are designed to create a conflict between the desire for progression and the well-being of the rock.

#### Tier 1: The Grind (Early Game)

These initial upgrades seem like standard clicker game improvements, encouraging the player to mine more efficiently.

*   **Stronger Pickaxe (Levels 1-5):** Increases dust per click.
    *   *Mechanic:* A straightforward upgrade to accelerate early-game resource gathering.
*   **Dust Goggles (Passive):** A small chance to find extra dust on each click.
    *   *Mechanic:* Introduces a bit of randomness and reward, keeping the player engaged.
*   **Auto-Clicker v0.1 (Toggleable):** Clicks the rock automatically at a slow pace.
    *   *Mechanic:* The first introduction to automation. At this stage, it's a convenience that can be turned off.

#### Tier 2: The Awakening (Mid Game)

This tier introduces upgrades that reveal the rock's "sentience" and the negative consequences of mining.

*   **Geode Sonar (Level 1):** A "bleep" sound is added to each click. After this upgrade, the game's music becomes slightly more melancholic.
    *   *Mechanic:* The first hint that the rock is more than just an inanimate object. The change in music sets a new tone.
*   **Rock Empathy (Level 1):** Text snippets occasionally appear on screen, expressing the rock's "thoughts" or "feelings" of discomfort. (e.g., "that tickles... a little too much," "feeling a bit crumbly today"). The rock's sprite may show small, temporary cracks after a series of rapid clicks.
    *   *Mechanic:* Directly communicates the rock's "pain" to the player, creating a sense of guilt.
*   **Auto-Clicker v1.0 (Permanent):** The auto-clicker is now always on and cannot be disabled. The click rate increases.
    *   *Mechanic:* This is a crucial turning point. The player loses some control, and the rock's health will now constantly decrease while the game is open. This forces the player to consider closing the game to prevent the rock's destruction.

#### Tier 3: The Consequence (Late Game)

The final set of upgrades forces the player to confront the ultimate fate of the rock.

*   **Earth-Shattering Pickaxe (Level 1):** A massive boost to dust per click, but each click now takes a significant chunk of the rock's health. The rock's sprite now shows permanent, growing cracks. The environment/background of the game starts to wither and decay.
    *   *Mechanic:* High risk, high reward. This accelerates the end-game but also dramatically speeds up the rock's demise, making the player's impact on the world undeniable.
*   **The Heart of the Mountain (1-time purchase):** An extremely expensive upgrade. Upon purchase, all mining stops. A message appears:
    > "You have reached the Heart of the Mountain. The rock is now still. It has given all it can. You have gathered enough. Will you take the final piece, or will you let it rest?"
    *   A choice is presented:
        1.  **"Take the Heart"**: The rock shatters, the screen goes white, and the game closes. A final message appears: "The mountain is no more. You are alone with your dust." (This is the "bad" ending).
        2.  **"Let it Rest"**: The game saves and then closes. (This is the "win" condition).

### Winning the Game

The true "win" is to purchase "The Heart of the Mountain" and choose to "Let it Rest." When the player re-opens the game, they will see the still, cracked rock in a recovering environment, with a small, permanent flower growing beside it. They can no longer click the rock. The game is, for all intents and purposes, over. The player has "won" by choosing preservation over consumption.