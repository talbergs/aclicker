```
+-------------------+
|       main        |
| (main.go)         |
|-------------------|
| - Game struct     |
| - Update()        |
| - Draw()          |
| - handleInput()   |
+---------^---------+
          |
          | Uses
          |
+---------+---------+         +-------------------+
|       game        |         |       shaders     |
| (game/game.go)    |         | (shaders/*.go)    |
|-------------------|         |-------------------|
| - Rock struct     |         | - ShaderEffect    |
| - Player struct   |         | - Game struct     |
| - EventDispatcher |         | - Grayscale()     |
| - NewGame()       |         | - Invert()        |
| - Click()         |         | - Warp()          |
| - UpgradeDamage() |         | - Apply()         |
| - ApplyDamageUpEvent()|     +---------^---------+
| - Save()          |                   |
| - Load()          |                   | Uses
+---------^---------+                   |
          |                             |
          | Uses                        |
          |                             |
+---------+-----------------------------+
|       game/hud                      |
| (game/hud/hud.go)                   |
|-------------------------------------|
| - HUD struct                        |
| - NewHUD()                          |
| - Draw()                            |
| - DrawHealthBar()                   |
+-------------------------------------+
          ^
          |
          | Uses
          |
+---------+---------+
|       game/events |
| (game/events.go)  |
|-------------------|
| - Event interface |
| - DamageUpgradedEvent |
| - EventHandler    |
| - EventDispatcher |
| - NewEventDispatcher()|
| - Register()      |
| - Dispatch()      |
+-------------------+
```