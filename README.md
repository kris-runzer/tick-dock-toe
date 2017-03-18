# tick-dock-toe
Simple tic-tac-toe game implemented with an AngularJS frontend and Go JSON API for game state.

All static assets are embedded as binary data into the executables using [go-bindata](https://github.com/jteeuwen/go-bindata).

## Starting the Game
#### 1. Clone repo then build and run Dockerfile

```Bash
git clone https://github.com/kris-runzer/tick-dock-toe.git

cd tick-dock-toe

docker build -t tick-dock-toe:latest .
docker run --net=host tick-dock-toe:latest
```

#### 2. Download binaries
Please see the [releases](https://github.com/kris-runzer/tick-dock-toe/releases) page.

## Playing the Game
After starting the game, goto http://localhost:3000 in your browser.  If you changed the `-bind` flag please adjust accordingly.

## Configuration
```Bash
Usage of tick-dock-toe:
  -bind string
    	the http binding port (default ":3000")
```

## Notes
1. Only the core game logic (game.go) is tested
