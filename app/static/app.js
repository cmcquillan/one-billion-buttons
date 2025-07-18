/* GLOBAL BOOTSTRAPPING */

window.CSS.registerProperty({
    name: '--cursor-url',
    inherits: true,
    initialValue: 'url("/cursor/fefefe/cursor.png"',
});

class Api {
    constructor() {
        this.gets = {};
        this.posts = {};
    }

    async _getButtonsAtCoordinates(x, y) {
        const resp = await fetch(`/api/${x}/${y}`);
        if (resp.status === 200) {
            const state = resp.json();
            return state;
        }

        return null;
    };

    async _pressButton(x, y, id, hex) {
        id = parseInt(id);

        const resp = await fetch(`/api/${x}/${y}`, {
            method: 'POST',
            body: JSON.stringify({ id, hex }),
        });

        return resp.json().then((data) => {
            if (resp.status === 200) {
                data.success = true;
            } else {
                data.success = false;
            }

            return data;
        });
    };

    getButtons(x, y) {
        const key = `${x}_${y}`;

        const cache = this.gets[key];
        if (cache) {
            return cache;
        }

        this.gets[key] = this._getButtonsAtCoordinates(x, y)
            .catch(e => {
                this.gets[key] = null;
                return e;
            }).then(r => {
                this.gets[key] = null;
                return r;
            });

        return this.gets[key];
    };

    pressButton(x, y, id, hex) {
        const key = `${x}_${y}`;

        const cache = this.posts[key];
        if (cache) {
            return cache;
        }

        this.posts[key] = this._pressButton(x, y, id, hex)
            .catch(e => {
                this.posts[key] = null;
                return e;
            }).then(r => {
                this.posts[key] = null;
                return r;
            });
        return this.posts[key];
    };
}

const dragReset = 'drag-reset';
const dragStart = 'drag-start';
const dragMove = 'drag-move';
const dragEnd = 'drag-end';

class PanelTracker {
    constructor(window, changeCallback) {
        this.touchBuffer = 0;
        this.posX = 0;
        this.posY = 0;
        this.dragging = false;
        this.current = null;
        this.changeCallback = changeCallback;
        this.window = window;
        this.trackedTouch = null;

        this.window.document.addEventListener('mousemove', (e) => this._onMouseMove(e), false);
        this.window.document.addEventListener('touchmove', (e) => this._onTouchMove(e), false);
        this.window.document.addEventListener('mousedown', (e) => this._onMouseDown(e), false);
        this.window.document.addEventListener('touchstart', (e) => this._onTouchStart(e), false);
        this.window.document.addEventListener('mouseleave', (e) => this._onMouseUp(e), false);
        this.window.document.addEventListener('touchcancel', (e) => this._onTouchEnd(e), false);
        this.window.document.addEventListener('mouseup', (e) => this._onMouseUp(e), false);
        this.window.document.addEventListener('touchend', (e) => this._onTouchEnd(e), false);
        this.window.document.addEventListener('hashchange', (e) => this._onReset(e), false);
        this.window.document.addEventListener('load', (e) => this._onReset(e), false);
    }

    _onReset(evt) {
        this.posX = 0;
        this.posY = 0;

        this.changeCallback(this._createEvent({
            type: dragReset,
            event: evt,
        }));
    }

    _onTouchStart(evt) {
        // Ensure it's a one-finger touch
        if (evt.touches && evt.touches.length === 1) {
            this.trackedTouch = evt.touches[0];
            this.dragging = true;

            this.changeCallback(this._createEvent({
                type: dragStart,
                event: evt,
            }));
        }
    }

    _onTouchMove(evt) {
        this.touchBuffer++;

        if (evt.changedTouches && evt.changedTouches.length > 0 && this.touchBuffer > 5) {
            this.touchBuffer = 0;
            const changes = [...evt.changedTouches];
            const newTouchState = changes.find((t) => t.identifier === this.trackedTouch.identifier);

            const { changeX, changeY } = this._updatePositionFromTouch(newTouchState);
            this.trackedTouch = newTouchState;

            this.changeCallback(this._createEvent({
                type: dragMove,
                event: evt,
                movementX: changeX,
                movementY: changeY,
            }));
        }
    }

    _updatePositionFromTouch(newTouchState) {
        const changeX = this.trackedTouch.clientX - newTouchState.clientX;
        const changeY = this.trackedTouch.clientY - newTouchState.clientY;

        this.posX += changeX;
        this.posY += changeY;
        return { changeX, changeY };
    }

    _onTouchEnd(evt) {
        // Fire one last move event
        const { changeX, changeY } = this._updatePositionFromTouch(this.trackedTouch);
        this.changeCallback(this._createEvent({
            type: dragMove,
            event: evt,
            movementX: changeX,
            movementY: changeY,
        }));

        this.trackedTouch = null;
        this.dragging = false;

        this.changeCallback(this._createEvent({
            type: dragEnd,
            event: evt,
        }));
    }

    _createEvent(opts) {
        return {
            type: null,
            event: null,
            movementX: 0,
            movementY: 0,
            dragging: this.dragging,
            posX: this.posX,
            posY: this.posY,
            ...opts
        };
    }

    _onMouseDown(evt) {
        this.dragging = true;

        this.changeCallback(this._createEvent({
            type: dragStart,
            event: evt,
        }));
    }

    _onMouseMove(evt) {
        if (this.dragging) {
            this.posX -= evt.movementX;
            this.posY -= evt.movementY;

            this.changeCallback(this._createEvent({
                type: dragMove,
                event: evt,
                movementX: evt.movementX,
                movementY: evt.movementY,
            }));
        }
    }

    _onMouseUp(evt) {
        this.dragging = false;

        this.changeCallback(this._createEvent({
            type: dragEnd,
            event: evt,
        }));
    }
}

/* HELPER FUNCTIONS */

function generateRandomHex() {
    let hex = '';

    for (let i = 0; i < 3; i++) {
        hex += Math.floor((Math.random() * 1000) % 256).toString(16).padStart(2, '0');
    }

    return hex;
}

function parseHash(hash) {
    try {
        const [xStr, yStr] = hash.split(',').map(s => s.replace('#', ''));

        return { x: parseInt(xStr), y: parseInt(yStr) };
    } catch {
        return null;
    }
}

function isInViewport(s, elem) {
    const container = s.appDiv.getBoundingClientRect();
    const rect = elem.getBoundingClientRect();

    return (
        rect.top < container.height + 10 &&
        rect.bottom > -10 &&
        rect.left < container.width + 10 &&
        rect.right > -10
    );
}

function createGridElement(w, s, gridX, gridY) {
    const rowLength = Math.sqrt(s.buttonPageSize);

    const gridXOffset = (gridX - s.gridX) * s.gridSizeX;
    const gridYOffset = (gridY - s.gridY) * s.gridSizeY;

    const gridId = `${gridX}-${gridY}`;
    let div = w.document.getElementById(gridId);

    if (!div) {
        div = w.document.createElement('div');
        div.id = gridId;
        div.classList.add('grid-container');
        div.style.top = `calc(${(gridYOffset)}px + var(--offset-y))`;
        div.style.left = `calc(${(gridXOffset)}px + var(--offset-x))`;
        
        if(s.gridSizeX && s.gridSizeY) {
            div.style.height = `${s.gridSizeY}px`;
            div.style.width = `${s.gridSizeX}px`;
        }
        
        div.setAttribute('data-x', gridX.toString());
        div.setAttribute('data-y', gridY.toString());
        div.style['grid-template-columns'] = `repeat(${rowLength}, auto)`;
        s.buttonContainer.appendChild(div);
        s.observer.observe(div);
    }

    return div;
}

function renderButtons(w, s, buttonState) {
    const buttonCount = buttonState.buttons.length;
    const gridX = buttonState.x;
    const gridY = buttonState.y;

    const div = createGridElement(w, s, gridX, gridY);

    for (let i = 0; i < buttonCount; i++) {
        const id = `b${buttonState.buttons[i].id}`;
        let button = w.document.getElementById(id);

        if (!button) {
            button = w.document.createElement('button');
            button.classList.add('button');

            if (buttonState.seen) {
                button.classList.add('seen');
            }

            button.id = `b${buttonState.buttons[i].id}`;
            button.setAttribute('data-x', gridX.toString());
            button.setAttribute('data-y', gridY.toString());
        }

        if (buttonState.buttons[i].hex) {
            button.classList.remove('pressed');
            button.classList.add('pressed');

            button.style.color = `#${buttonState.buttons[i].hex}`;
        }

        if (!button.parentElement) {
            div.appendChild(button);
        }
    }

    // if (!div.parentElement) {
    //     s.buttonContainer.appendChild(div);
    //     s.observer.observe(div);
    // }

    buttonState.seen = true;

    return div;
}

async function storeButtonState(s, buttonState) {
    const key = `${buttonState.x}_${buttonState.y}`;

    s.buttonStates[key] = buttonState;
}

async function retrieveButtonState(s, x, y) {
    const key = `${x}_${y}`;

    return s.buttonStates[key];
}

async function renderGridPoint(w, s, x, y) {

    if (x < 1 || y < 1) {
        return null;
    }

    const buttonState = await retrieveButtonState(s, x, y) || await s.api.getButtons(x, y);
    await storeButtonState(s, buttonState);
    return renderButtons(w, s, buttonState);
}

async function handleButtonClick(w, s, button) {
    const point = getGridPoint(button);
    const id = /[0-9]+/.exec(button.id)[0];

    button.classList.remove('pressed');
    button.classList.add('pressed');

    const buttonState = await s.api.pressButton(point.x, point.y, id, s.hexCode);

    if (buttonState.success) {
        console.log('You totally clicked the button');
    } else {
        console.log('Someone else clicked the button first');
    }

    await storeButtonState(s, buttonState);
    return renderButtons(w, s, buttonState);
}

function getGridPoint(element) {
    const x = parseInt(element.getAttribute('data-x'));
    const y = parseInt(element.getAttribute('data-y'));
    return { x, y };
}

async function eventLoop(w, s) {

    const debugDataDiv = document.getElementById('data');

    if (debugDataDiv) {
        const buttons = w.document.getElementsByClassName('button');
        const inView = [];

        for (let i = 0; i < buttons.length; i++) {
            const b = buttons[i];

            const rect = b.getBoundingClientRect();

            // Probably off to the top or left.
            if (rect.x < 0 || rect.y < 0) {
                continue;
            }

            // Probably off to the bottom or right.
            if (rect.x > w.innerWidth || rect.y > w.innerHeight) {
                continue;
            }

            inView.push(b);
        }

        debugDataDiv.innerHTML = inView.map((val) => val.innerText).join(',');
    }

    updateDocumentCursor(w, s);

    if (s.isScrollDirty) {
        // Get all .grid-container elements that are visible
        // Schedule render of all their buttons
        // Render .grid-container at boundaries
        // Repeat (next event loop)
        // Set isScrollDirty = false when all visibles have buttons renders.

        const list = [...w.document.querySelectorAll('.grid-container')];
        const visisble = list.filter(elem => isInViewport(s, elem));

        if (visisble.length === 0) {
            s.isScrollDirty = false;
        } else {
            await Promise.all(visisble.map(async (elem) => {
                const { x, y } = getGridPoint(elem);
                await renderGridPoint(w, s, x, y);

                for (let xp = x - 1; xp <= x + 1; xp++) {
                    for (let yp = y - 1; yp <= y + 1; yp++) {
                        createGridElement(w, s, xp, yp);
                    }
                }
            }));
        }
    }
}

function updateDocumentCursor(w, s) {
    const colorSelect = w.document.getElementById('color-select');

    if (colorSelect && colorSelect.value) {
        s.hexCode = colorSelect.value.substring(1);
    }

    const cursor = `url("/cursor/${s.hexCode}/cursor.png")`;
    const cursorProp = s.appDiv.style.getPropertyValue('--cursor-url');
    if (cursorProp !== cursor) {
        w.document.body.style.setProperty('--cursor-url', cursor);
    }
}

async function fixStates(w, s) {
    updateDocumentCursor(w, s);
}

async function startApplication(w, s) {
    const hash = parseHash(w.location.hash);

    if (!hash || !hash.x || !hash.y) {
        const toX = Math.ceil((Math.random() * 1000000) % s.gridMaxX);
        const toY = Math.ceil((Math.random() * 1000000) % s.gridMaxY);

        w.location.hash = `#${toX},${toY}`;
        return;
    }

    s.gridX = hash.x;
    s.gridY = hash.y;

    s.appDiv = document.getElementById('app');
    // We reset ALL state here
    s.buttonContainer = document.getElementById('button-box');
    s.buttonContainer.innerHTML = '';

    s.observer = new IntersectionObserver((entries) => {
        entries.forEach((val) => {
            // Ditch the grid containers that are scrolled away from view.
            if (val.target.tagName === 'DIV') {
                if (!val.isIntersecting) {
                    val.target.innerHTML = '';
                }
            }
        });
    }, {
        root: w.document.getElementById('button-box')
    });

    const centerDiv = await renderGridPoint(w, s, s.gridX, s.gridY);

    s.gridSizeX = centerDiv.clientWidth;
    s.gridSizeY = centerDiv.clientHeight;
    s.isScrollDirty = true;

    /* Set some control initial states */
    w.document.getElementById('color-select').value = `#${s.hexCode}`;

    /* Some browsers are just still terrible */
    s.interval = w.setInterval(async () => fixStates(w, s), 100);

    /* Bad way to eal with screen proportions, find something better */
    s.eventInterval = w.setInterval(async () => await eventLoop(w, s,), 1000);
}

async function onPanelStateChange(w, s, data) {
    if (data.dragging && s.appDiv) {
        s.isScrollDirty = true;
        s.appDiv.style.setProperty('--offset-x', `${-data.posX}px`);
        s.appDiv.style.setProperty('--offset-y', `${-data.posY}px`);
    }

    await eventLoop(w, s);
}

const root = window.getComputedStyle(document.documentElement);

window.state = window.state || {
    root: root,
    gridX: null,
    gridY: null,
    gridSizeX: null,
    gridSizeY: null,
    gridMaxX: parseInt(root.getPropertyValue('--button-grid-count-x')),
    gridMaxY: parseInt(root.getPropertyValue('--button-grid-count-y')),
    buttonPageSize: 100,
    appDiv: null,
    buttonContainer: null,

    buttonStates: {},

    interval: null,
    eventInterval: null,
    observer: null,
    api: new Api(),
    panelTracker: null,

    // User
    hexCode: generateRandomHex(),

    debug: false,
};

window.state.panelTracker = new PanelTracker(window, (data) => onPanelStateChange(window, state, data)),

    window.addEventListener('load', function () { startApplication(window, state); });
window.addEventListener('hashchange', function () {
    startApplication(window, state);
});

window.keys = window.keys || { ctrl: false };

window.document.addEventListener('keydown', (evt) => {
    if (evt.ctrlKey && evt.key.toLowerCase() === 'b') {
        const appDiv = window.document.getElementById('app');
        appDiv.classList.toggle('app-debug');

        const debugDiv = window.document.getElementById('debug');

        state.debug = !state.debug;

        debugDiv.style.display = state.debug
            ? 'block'
            : 'none';
    }
});

window.document.addEventListener('click', async (evt) => {
    if (evt.target.tagName === 'BUTTON' && evt.target.classList.contains('button')) {
        await handleButtonClick(window, state, evt.target);
    }

});
