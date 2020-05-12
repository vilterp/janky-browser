# JankyBrowser

A simple browser for pedagogical purposes. It only renders a very small SVG-like
XML dialect; `testdata/` contains some example files.

## Install

1. Install Go 1.9.4 or similar.
2. `make deps`

## Build and Run

1. Set up a web server to serve the files in `testdata` (it won't work on
   real SVG files). I recommend
   `npm install -g http-server && cd testdata && http-server`
2. In this directory, `make`
3. Browse away. The browser is hardcoded to hit
   `http://localhost:8081/circleRectText.svg` first; you may need to type in
   a different port to match your server.
