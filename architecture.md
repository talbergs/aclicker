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
| - Player struct   |         | - Grayscale()     |
| - Game struct     |         | - Invert()        |
| - Warp()          |         | - Apply()         |
| - NewGame()       |         +---------^---------+
| - Click()         |                   |
| - UpgradeDamage() |                   | Uses
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
| - NewHUD()                          |
| - Draw()                            |
| - DrawHealthBar()                   |
+-------------------------------------+
```