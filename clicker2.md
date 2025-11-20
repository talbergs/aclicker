Dark genre game dynamics, a clicker.
You click the friendly rock and thus mine it's health.
The result of efficient clicking allows to upgrade stats making minig more efficient.

Game story:
- main state is 10000000 units of health shown in healthbar
- a click on rock given the accumulated compute gives approriate amount of damge and dust
- use the dust to buy various attack tools and skill coeificients improving mining process efficiency

Game inputs:
- q - quit
- s - save
- l - load
- click-event - mouse-down
- click-position - x,y

Game view:
- 800x600px
- scroll on screen will adjust audio volume with visual feedback

Game technology:
- using golang and popular 2d graphics libraries
- rules of the game are structured in separate package that will encapsulate boundaries
- GPU is used for running various shaders, app decides them dynamically

Game assets:
- sprites.png file used for scene creation
-- capture hardcoded regions from sprite as variables "stone", "marketplace", "clouds" etc..
- game assets package also provides music player api and can be used in shaders

Game design - events driven.
- actions (state transitions)
-- user stats update (i.e. increment damage by 10)
- events (facts)
-- user clicked stone
-- user stats update (an action example, were event is sourced from 2 past events)
--- user clicked stone for 10 times in past second
- audio stream is also passed inbetween shaders
