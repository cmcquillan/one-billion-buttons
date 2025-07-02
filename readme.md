# One Billion Buttons

It's a grid of buttons, roughly 1 billion of them. Each coordinate of the grid is 1000 buttons.

* When the client connects, a random color HEX code is chosen for them and stored in the client local storage. 
* When a button is pressed, its color changes to the user's HEX code.
* Buttons remain in pressed state forever.
* There is a minimap that, on some schedule, takes the state of every button and draws a "close enough" distillation of all the buttons states and colors.
* Clicking on the minimap will teleport you to the grid coordinate.


## Design

* Go Backend
  - Redis encoding for button state
* JS / HTML Frontend
  - TBD: htmx? Go templates?
* JSON Api w/ aggressive cacheability
* Game State is initialized with 3 parameters
  - X distance
  - Y distance
  - Buttons per coordinate
  - e.g. 1000 * 1000 * 1000 = 1,000,000,000 buttons
  - e.g. 3163 * 3163 * 100 = 1,000,456,900 buttons
* Redis keys for button state
  - key: `x,y`
  - value: raw byte array. Every 3 bytes is a hex code for a button index w/in the grid coordinate.

### GET Routes

* `/` -- Serve index.html
* `/#{x:int},{y:int}` -- Serve index.html, but URL becomes center point
* `/*.(js|css)` -- Serve static files. Highly cacheable.
* `/minimap.{ext}` -- Serve image of minimap. Highly cacheable.
* `/api/{x:int},{y:int}` -- Serve button state + optional hashed link to most recent state.
  - `x` x coordinate
  - `y` y coordinate
  - `buttons[]`
    + `id` -- id of the button
    + `hex` -- Hex code of the button color, `null` is default and unpressed. All other colors are pressed.
  - `next` -- The hash link to (see below) to poll for more recent state. 
* `/api/{x:int},{y:int}/{hash}` -- Same as above, but aggressively cacheable. 
  - Same return shape as above.
  - Sends `cache-control` that is long-lived, server and client cacheable.
  - Idea is that the `next` link will serve

### POST Routes

* `/api/{x:int},{y,int}` -- Send a button index along with hex code to push the button.