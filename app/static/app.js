/* GLOBAL BOOTSTRAPPING */

/**
 * @typedef GridPoint
 * @property {number} x
 * @property {number} y
 */

/**
 * @typedef ButtonState
 * @property {number} id
 * @property {string?} hex
 */

/**
 * @interface GridState 
 * @implements {GridPoint}
 * @property {ButtonState[]} buttons
 * @property {string?} next
 */

/**
 * @typedef ButtonStat
 * @property {string} stat_key
 * @property {string} stat_name
 * @property {string} stat_desc
 * @property {number} val
 * @property {number} scale
 * @property {number} order
 */

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

    /**
     * 
     * @param {number} x 
     * @param {number} y 
     * @param {string} id 
     * @param {string} hex 
     * @returns {Promise<GridState>}
     */
    pressButton(x, y, id, hex) {
        return this._pressButton(x, y, id, hex);
    };

    /**
     * 
     * @returns {Promise<ButtonStat[]>}
     */
    getStats() {
        return fetch('/api/stats')
            .then(resp => {
                if (resp.status === 200) {
                    return resp.json();
                }
                return {};
            });
    }
}

const dragReset = 'drag-reset';
const dragStart = 'drag-start';
const dragMove = 'drag-move';
const dragEnd = 'drag-end';

class PanelTracker {
    /**
     * 
     * @param {Window} window 
     * @param {Function} changeCallback 
     */
    constructor(window, changeCallback) {
        this.touchBuffer = 0;
        this.posX = 0;
        this.posY = 0;
        this.dragging = false;
        this.current = null;
        this.changeCallback = changeCallback;
        this.window = window;
        this.trackedTouch = null;

    }

    init() {
        const appDiv = this.window.document.getElementById('app');
        appDiv.addEventListener('mousemove', (e) => this._onMouseMove(e), false);
        appDiv.addEventListener('touchmove', (e) => this._onTouchMove(e), false);
        appDiv.addEventListener('mousedown', (e) => this._onMouseDown(e), false);
        appDiv.addEventListener('touchstart', (e) => this._onTouchStart(e), false);
        appDiv.addEventListener('mouseleave', (e) => this._onMouseUp(e), false);
        appDiv.addEventListener('touchcancel', (e) => this._onTouchEnd(e), false);
        appDiv.addEventListener('mouseup', (e) => this._onMouseUp(e), false);
        appDiv.addEventListener('touchend', (e) => this._onTouchEnd(e), false);
        appDiv.addEventListener('hashchange', (e) => this._onReset(e), false);
        appDiv.addEventListener('load', (e) => this._onReset(e), false);
    }

    /**
     * 
     * @param {Event} evt 
     */
    _onReset(evt) {
        this.posX = 0;
        this.posY = 0;

        this.changeCallback(this._createEvent({
            type: dragReset,
            event: evt,
        }));
    }

    /**
     * 
     * @param {TouchEvent} evt 
     */
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

    /**
     * 
     * @param {MouseEvent|TouchEvent} evt 
     */
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

    /**
     * 
     * @param {TouchEvent} evt 
     */
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

    /**
     * 
     * @param {MouseEvent} evt 
     */
    _onMouseDown(evt) {
        this.dragging = true;

        this.changeCallback(this._createEvent({
            type: dragStart,
            event: evt,
        }));
    }

    /**
     * 
     * @param {MouseEvent} evt 
     */
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

    /**
     * 
     * @param {MouseEvent} evt 
     */
    _onMouseUp(evt) {
        this.dragging = false;

        this.changeCallback(this._createEvent({
            type: dragEnd,
            event: evt,
        }));
    }
}

class LocalState {
    /**
     * 
     * @param {Window} window 
     * @param {number} gridMaxX
     * @param {number} gridMaxY  
     */
    constructor(window, gridMaxX, gridMaxY) {
        this.document = window.document;
        this.storage = window.localStorage;

        this.gridMaxX = gridMaxX;
        this.gridMaxY = gridMaxY;
        this.gridX = null;
        this.gridY = null;
        this.gridSizeX = null;
        this.gridSizeY = null;
        this.buttonPageSize = 100;

        this.appDiv = null;
        this.buttonContainer = null;

        /**
         * @type {Object.<string,Array.GridState>}
         */
        this.buttonStates = {};
        this.interval = null;
        this.eventInterval = null;
        this.observer = null;
        this.api = new Api();
        this.panelTracker = null;
        this.debug = false;
    }

    /**
     * @returns {string|null}
     */
    getUserHex() {
        return this.storage.getItem('userHex')
    }

    /**
     * 
     * @param {string} hex 
     */
    setUserHex(hex) {
        this.storage.setItem('userHex', hex);
    }

    /**
     * 
     * @param {GridState} buttonState 
     */
    async storeButtonState(buttonState) {
        const key = `${buttonState.x}_${buttonState.y}`;
        let modified = false;

        if (this.buttonStates[key]) {

            for (let i = 0; i < buttonState.buttons.length; i++) {
                if (buttonState.buttons[i].hex) {
                    this.buttonStates[key].buttons[i].hex = buttonState.buttons[i].hex;
                    modified = true;
                }
            }
        }
        else {
            this.buttonStates[key] = buttonState;
        }

        if (modified) {
            this.buttonStates[key].next = buttonState.next;
        }
    }

    /**
     * 
     * @param {number} x 
     * @param {number} y 
     * @returns {GridState}
     */
    async retrieveButtonState(x, y) {
        const key = `${x}_${y}`;

        return this.buttonStates[key];
    }

    /**
     * @returns {HTMLDialogElement}
     */
    get controlsModal() {
        return this.document.getElementById('controlDialog');
    }

    /**
     * @returns {HTMLDialogElement}
     */
    get statsModal() {
        return this.document.getElementById('statsDialog');
    }

    /**
     * @returns {HTMLDialogElement}
     */
    get mapModal() {
        return this.document.getElementById('mapDialog');
    }

    /**
     * @returns {HTMLImageElement}
     */
    get minimap() {
        return this.document.getElementById('minimap');
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

/**
 * 
 * @param {string} hash 
 * @returns {GridPoint|null}
 */
function parseHash(hash) {
    try {
        const [xStr, yStr] = hash.split(',').map(s => s.replace('#', ''));

        return { x: parseInt(xStr), y: parseInt(yStr) };
    } catch {
        return null;
    }
}

/**
 * 
 * @param {LocalState} s 
 * @param {HTMLElement} elem 
 * @returns {boolean}
 */
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

/**
 * 
 * @param {Window} w 
 * @param {LocalState} s 
 * @param {number} gridX 
 * @param {number} gridY 
 * @returns 
 */
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

        if (s.gridSizeX && s.gridSizeY) {
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

/**
 * 
 * @param {Window} w 
 * @param {LocalState} s 
 * @param {*} buttonState 
 * @returns {HTMLDivElement}
 */
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
            button.setAttribute('aria-label', `Button ${id}`);

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
        } else {
            button.classList.remove('pressed');
        }

        if (!button.parentElement) {
            div.appendChild(button);
        }
    }

    buttonState.seen = true;

    return div;
}

/**
 * 
 * @param {Window} w 
 * @param {LocalState} s 
 * @param {number} x 
 * @param {number} y 
 * @returns {Promise<HTMLDivElement|null>}
 */
async function renderGridPoint(w, s, x, y) {

    if (x < 1 || y < 1) {
        return null;
    }

    const buttonState = await s.retrieveButtonState(x, y) || await s.api.getButtons(x, y);
    await s.storeButtonState(buttonState);
    return renderButtons(w, s, buttonState);
}

/**
 * 
 * @param {Window} w 
 * @param {LocalState} s 
 * @param {HTMLButtonElement} button 
 * @returns {Promise<HTMLDivElement>}
 */
async function handleButtonClick(w, s, button) {
    const point = getGridPoint(button);
    const id = /[0-9]+/.exec(button.id)[0];

    button.classList.remove('pressed');
    button.classList.add('pressed');

    const buttonState = await s.api.pressButton(point.x, point.y, id, s.getUserHex());

    if (buttonState.success) {
        console.log('You totally clicked the button');
    } else {
        console.log('Someone else clicked the button first');
    }

    await s.storeButtonState(buttonState);
    return renderButtons(w, s, buttonState);
}

/**
 * 
 * @param {HTMLElement} element 
 * @returns {GridPoint}
 */
function getGridPoint(element) {
    const x = parseInt(element.getAttribute('data-x'));
    const y = parseInt(element.getAttribute('data-y'));
    return { x, y };
}

/**
 * 
 * @param {Window} w 
 * @param {LocalState} s 
 */
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

/**
 * 
 * @param {Window} w 
 * @param {LocalState} s 
 */
function updateDocumentCursor(w, s) {
    const colorSelect = w.document.getElementById('color-select');

    if (colorSelect && colorSelect.value) {
        s.setUserHex(colorSelect.value.substring(1));
    }

    const cursor = `url("/cursor/${s.getUserHex()}/cursor.png")`;
    const cursorProp = s.appDiv.style.getPropertyValue('--cursor-url');
    if (cursorProp !== cursor) {
        w.document.body.style.setProperty('--cursor-url', cursor);
    }
}

async function fixStates(w, s) {
    updateDocumentCursor(w, s);
}

/**
 * 
 * @param {PointerEvent} evt 
 * @param {Window} w 
 * @param {LocalState} s 
 */
function handleMinimapNavigation(evt, w, s) {
    const x = evt.pageX - s.mapModal.offsetLeft;
    const y = evt.pageY - s.mapModal.offsetTop;

    const xRatio = s.mapModal.clientWidth / s.gridMaxX;
    const yRatio = s.mapModal.clientHeight / s.gridMaxY;

    const xcoord = Math.ceil(x / xRatio);
    const ycoord = Math.ceil(y / yRatio);
    w.location.hash = `#${xcoord},${ycoord}`;
}

/**
 * 
 * @param {Window} w 
 * @param {LocalState} s 
 * @returns 
 */
async function startApplication(w, s) {
    const hash = parseHash(w.location.hash);

    if (!hash ||
        !hash.x || !hash.y ||
        hash.x < 1 || hash.y < 1 ||
        hash.x > s.gridMaxX || hash.y > gridMaxY) {
        const toX = Math.ceil((Math.random() * 1000000) % s.gridMaxX);
        const toY = Math.ceil((Math.random() * 1000000) % s.gridMaxY);

        w.location.hash = `#${toX},${toY}`;
        return;
    }

    s.gridX = hash.x;
    s.gridY = hash.y;

    s.appDiv = document.getElementById('app');
    // We reset ALL state here
    s.buttonContainer = document.getElementById('buttonBox');
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
        root: w.document.getElementById('buttonBox')
    });

    let hex = s.getUserHex()

    if (!hex) {
        s.setUserHex(generateRandomHex());
    }

    w.document.getElementById('buttonBox').addEventListener('click', async (evt) => {
        if (evt.target.tagName === 'BUTTON' && evt.target.classList.contains('button')) {
            await handleButtonClick(w, state, evt.target);
        }
    });

    w.document.getElementById('statsTrigger').addEventListener('click', async (evt) => {
        await showStats(evt, w, s);
    });

    w.document.getElementById('controlTrigger').addEventListener('click', async (evt) => {
        await showControls(evt, w, s);
    });

    w.document.getElementById('mapTrigger').addEventListener('click', async (evt) => {
        await showMap(evt, w, s);
    });

    // Modal-closing clicks
    w.document.addEventListener('click', (evt) => {
        const rect = evt.target.getBoundingClientRect();

        for (modal of [s.statsModal, s.controlsModal, s.mapModal]) {
            if (modal.open && evt.target === modal) {
                if (evt.clientX < rect.left ||
                    evt.clientX > rect.right ||
                    evt.clientY < rect.top ||
                    evt.clientY > rect.bottom) {
                    modal.close();
                }
            }
        }
    });

    s.minimap.addEventListener('click', (evt) => {
        handleMinimapNavigation(evt, w, s);
    });

    const centerDiv = await renderGridPoint(w, s, s.gridX, s.gridY);

    s.gridSizeX = centerDiv.clientWidth;
    s.gridSizeY = centerDiv.clientHeight;
    s.isScrollDirty = true;

    /* Set some control initial states */
    w.document.getElementById('color-select').value = `#${s.getUserHex()}`;

    /* Some browsers are just still terrible */
    s.interval = w.setInterval(async () => fixStates(w, s), 100);

    /* Bad way to eal with screen proportions, find something better */
    s.eventInterval = w.setInterval(async () => await eventLoop(w, s,), 1000);
}

/**
 * 
 * @param {Window} w 
 * @param {LocalState} s 
 * @param {*} data 
 */
async function onPanelStateChange(w, s, data) {
    if (data.dragging && s.appDiv) {
        s.isScrollDirty = true;
        s.appDiv.style.setProperty('--offset-x', `${-data.posX}px`);
        s.appDiv.style.setProperty('--offset-y', `${-data.posY}px`);
    }

    await eventLoop(w, s);
}

/**
 * 
 * @param {MouseEvent} evt 
 * @param {Window} w
 * @param {LocalState} s 
 */
async function showStats(evt, w, s) {
    const stats = await s.api.getStats();

    if (stats instanceof Array) {
        stats.sort((a, b) => a.order - b.order);
    }

    const dialog = s.statsModal;

    const statsDiv = dialog.getElementsByClassName('stats-container')[0];

    const fragment = w.document.createDocumentFragment();
    for (stat of stats) {
        const value = stat.val * Math.pow(10, stat.scale);

        const div = fragment.appendChild(w.document.createElement('div'));
        div.classList.add('stat-item');
        const lSpan = div.appendChild(w.document.createElement('span'));
        lSpan.classList.add('stat-label');
        lSpan.textContent = stat.stat_name;
        const vSpan = div.appendChild(w.document.createElement('span'));
        vSpan.classList.add('stat-value');
        vSpan.textContent = value.toLocaleString();
    }

    statsDiv.replaceChildren(fragment);
    dialog.showModal();
}

async function showControls(evt, w, s) {
    s.controlsModal.showModal();
}

async function showMap(evt, w, s) {
    s.minimap.src = "./minimap.png";
    s.mapModal.showModal();
}

const root = window.getComputedStyle(document.documentElement);
const gridMaxX = parseInt(root.getPropertyValue('--button-grid-count-x'));
const gridMaxY = parseInt(root.getPropertyValue('--button-grid-count-y'));

window.state = window.state || new LocalState(window, gridMaxX, gridMaxY);
window.state.panelTracker = new PanelTracker(window, (data) => onPanelStateChange(window, state, data));

function bootstrap() {
    startApplication(window, state);
    window.state.panelTracker.init();
}

window.addEventListener('load', bootstrap);
window.addEventListener('hashchange', bootstrap);

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
