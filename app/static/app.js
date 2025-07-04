/* GLOBAL BOOTSTRAPPING */

window.CSS.registerProperty({
    name: '--cursor-url',
    inherits: true,
    initialValue: 'url("/cursor/fefefe/cursor.png"',
});

function Api() {
    async function getButtonsAtCoordinates(x, y) {
        const resp = await fetch(`/api/${x}/${y}`);
        if (resp.status === 200) {
            const state = resp.json();
            return state;
        }

        return null;
    };

    async function pressButton(x, y, id, hex) {
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
        })
    };

    this.gets = {};
    this.posts = {};
    this.getButtons = (x, y) => {
        const key = `${x}_${y}`;

        const cache = this.gets[key];
        if (cache) {
            return cache;
        }

        this.gets[key] = getButtonsAtCoordinates(x, y)
            .catch(e => {
                this.gets[key] = null;
                return e;
            }).then(r => {
                this.gets[key] = null;
                return r;
            });

        return this.gets[key];
    };
    this.pressButton = (x, y, id, hex) => {
        const key = `${x}_${y}`;

        const cache = this.posts[key];
        if (cache) {
            return cache;
        }

        this.posts[key] = pressButton(x, y, id, hex)
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



function renderButtons(w, s, buttonState) {
    const buttonCount = buttonState.buttons.length;
    const rowLength = Math.sqrt(buttonCount);
    const gridX = buttonState.x;
    const gridY = buttonState.y;

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
        div.setAttribute('data-x', gridX.toString());
        div.setAttribute('data-y', gridY.toString());
        div.style['grid-template-columns'] = `repeat(${rowLength}, auto)`;
    }

    for (let i = 0; i < buttonCount; i++) {
        const id = `b${buttonState.buttons[i].id}`;
        let button = w.document.getElementById(id);

        if (!button) {
            button = w.document.createElement('button');

            const x = i % rowLength;
            const y = Math.floor(i / rowLength);
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

    if (!div.parentElement) {
        s.buttonContainer.appendChild(div);
        s.observer.observe(div);
    }

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

async function renderLoop(w, s) {

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
        const checkBounds = [
            { coord: [1, 1], vec: [-1, -1] }, // Top left 
            { coord: [1, w.innerHeight / 2], vec: [-1, 0] }, // Middle Left
            { coord: [1, w.innerHeight - 1], vec: [-1, 1] }, // Bottom Left
            { coord: [w.innerWidth / 2, 1], vec: [0, -1] }, // Top Middle
            { coord: [w.innerWidth - 1, 1], vec: [1, -1] }, // Top Right
            { coord: [w.innerWidth - 1, w.innerHeight / 2], vec: [1, 0] }, // Middle Right
            { coord: [w.innerWidth - 1, w.innerHeight - 1], vec: [1, 1] }, // Bottom Right
            { coord: [w.innerWidth / 2, w.innerHeight - 1], vec: [0, 1] }, // Bottom Middle
        ];

        const center = w.document.elementsFromPoint(w.innerWidth / 2, w.innerHeight / 2)
            .find((d) => d.classList.contains('grid-container'));

        if (!center) {
            console.log('fuck...');
        } else {

            const gridPoint = getGridPoint(center);
            const gridX = gridPoint.x;
            const gridY = gridPoint.y;

            for (let j = 0; j < checkBounds.length; j++) {
                const [x, y] = checkBounds[j].coord;
                const elems = w.document.elementsFromPoint(x, y);
                const elem = elems.find((d) => d.classList.contains('grid-container'));

                if (!elem) {

                    if (!gridX || !gridY || isNaN(gridX) || isNaN(gridY) || gridX < 1 || gridY < 1) {
                        continue;
                    }
                    const [dX, dY] = checkBounds[j].vec;
                    const renderGridX = gridX + dX;
                    const renderGridY = gridY + dY;

                    await renderGridPoint(w, s, renderGridX, renderGridY);
                }
            }

            s.isScrollDirty = false;
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
        w.location.hash = '#1,1';
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
}

function startDrag(evt, s) {
    // We're targeting a button so ignore this.
    if (evt.target.localName === 'button') {
        return;
    }

    let check = evt.target;

    while (check) {
        if (check.classList.contains('control-box')) {
            return;
        }

        check = check.parentElement;
    }

    s.dragging = true;
}

function stopDrag(evt, s) {
    s.dragging = false;
}

async function windowMove(evt, s) {
    evt.preventDefault();

    if (!s.dragging) {
        return;
    }

    // Keep posX, posY convenient and invert in css.
    s.posX += -evt.movementX;
    s.posY += -evt.movementY;

    if (s.appDiv) {
        s.appDiv.style.setProperty('--offset-x', `${-s.posX}px`);
        s.appDiv.style.setProperty('--offset-y', `${-s.posY}px`);
    }

    s.isScrollDirty = true;

    await eventLoop(window, s);
}

window.state = window.state || {
    dragging: false,
    posX: 0,
    posY: 0,
    gridX: null,
    gridY: null,
    gridSizeX: null,
    gridSizeY: null,
    buttonPageSize: 100,
    appDiv: null,
    buttonContainer: null,

    buttonStates: {},

    interval: null,
    observer: null,
    api: new Api(),

    // User
    hexCode: generateRandomHex(),
};

window.addEventListener('load', function () { startApplication(window, state); });
window.addEventListener('hashchange', function () {
    state.posX = 0;
    state.posY = 0;

    if (state.appDiv) {
        state.appDiv.style.setProperty('--offset-x', `${-state.posX}px`);
        state.appDiv.style.setProperty('--offset-y', `${-state.posY}px`);
    }
    startApplication(window, state);
});

window.document.addEventListener('mousemove', async (evt) => await windowMove(evt, state), false);
window.document.addEventListener('mousedown', (evt) => startDrag(evt, state), false);
window.document.addEventListener('mouseleave', (evt) => stopDrag(evt, state), false);
window.document.addEventListener('mouseup', (evt) => stopDrag(evt, state), false);

window.keys = window.keys || { ctrl: false };

window.document.addEventListener('keydown', (evt) => {
    if (evt.ctrlKey && evt.key.toLowerCase() === 'b') {
        const appDiv = window.document.getElementById('app');
        appDiv.classList.toggle('app-debug');

        const debugDiv = window.document.getElementById('debug');

        debugDiv.style.display = debugDiv.style.display === 'none'
            ? 'block'
            : 'none';
    }
});

window.document.addEventListener('click', async (evt) => {
    if (evt.target.tagName === 'BUTTON' && evt.target.classList.contains('button')) {
        await handleButtonClick(window, state, evt.target);
    }

});