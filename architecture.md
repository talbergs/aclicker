+-------------------+
|       main        |
| (main.go)         |
|-------------------|
| - EbitenGame struct|
| - Update()        |
| - Draw()          |
| - handleInput()   |
| - main()          |
+---------^---------+
          |
          | Uses
          |
+---------+---------+         +-------------------+
|       game        |         |       shaders     |
| (shaders/*.kage)  |
|-------------------|         |-------------------|
| - Game struct     |         | - DesertShader    |
| - Rock struct     |         | - TimeClickShader |
| - Player struct   |         | - GrayscaleShader |
| - UpgradeManager  |         | - InvertShader    |
| - EventDispatcher |         | - WarpShader      |
| - NewGame()       |         +---------^---------+
| - Click()         |                   |
| - PurchaseUpgrade()|                  | Uses
| - TakeHeart()     |                   |
| - LetRest()       |                   |
| - Save()          |                   |
| - Load()          |                   |
+---------^---------+                   |
          |                             |
          | Uses                        |
          |                             |
+---------+-----------------------------+
|       game/hud                      |
| (game/hud/hud.go)                   |
|-------------------------------------|
| - HUD struct                        |
| - UpgradeButton struct              |
| - NewHUD()                          |
| - Draw()                            |
| - DrawHealthBar()                   |
| - GetClickedUpgradeID()             |
| - GetClickedChoiceID()              |
+-------------------------------------+
          ^
          |
          | Uses
          |
+---------+---------+
|       game/upgrades|
| (game/upgrades.go)|
|-------------------|
| - Upgrade struct  |
| - UpgradeManager  |
| - NewUpgradeManager()|
| - Init()          |
| - GetAllUpgrades()|
| - GetUpgrade()    |
| - GetPlayerUpgradeLevel()|
+-------------------+
          ^
          |
          | Uses
          |
+---------+---------+
|       game/events |
| (game/events.go)  |
|-------------------|
| - Event interface |
| - ClickEvent      |
| - EventHandler    |
| - EventDispatcher |
| - NewEventDispatcher()|
| - Register()      |
| - Dispatch()      |
+-------------------+
          ^
          |
          | Uses
          |
+---------+---------+
|       assets      |
| (assets/assets.go)|
|-------------------|
| - Image assets    |
| - Audio assets    |
| - AudioContext    |
| - Music Players   |
| - SFX Players     |
+-------------------+