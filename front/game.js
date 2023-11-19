// game.js

const cellImages = {
    not_seen: 'images/not_seen.png',
    bomb: 'images/bomb.png',
    empty: 'images/empty.png',
    key: 'images/key.png'
};

document.addEventListener('DOMContentLoaded', function () {
    // Initialize the minefield grid
    fetchGameState();
});

function fetchGameState() {
    fetch('https://minefield.onrender.com/player/data', {
        method: 'GET',
        headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('jwtToken') // Use the stored JWT token
        }
    })
        .then(response => response.json())
        .then(data => {
            displayUserInfo(data); // Display user info
            initializeGrid(data.game_state); // Initialize the game grid
        })
        .catch(error => {
            console.error('Error:', error);
        });
}

// function initializeGrid(size) {
//     const minefield = document.getElementById('minefield');
//     const notSeenIconUrl = 'images/not_seen.png'; // Replace with the path to your not_seen icon

//     minefield.innerHTML = ''; // Clear previous cells if any

//     for (let i = 0; i < size * size; i++) {
//         const cell = document.createElement('div');
//         cell.className = 'cell not_seen'; // Initially, all cells are not seen
//         cell.dataset.status = 'not_seen'; // Data attribute to track status
//         cell.style.backgroundImage = `url('${notSeenIconUrl}')`; // Set not_seen icon as background
//         cell.dataset.index = i; // Store the 1D index
//         cell.addEventListener('click', cellClicked);
//         minefield.appendChild(cell);
//     }
// }

function initializeGrid(gameState) {
    const minefield = document.getElementById('minefield');
    minefield.innerHTML = ''; // Clear the grid first
    const cellSize = `calc(700px / ${gameState.field_len})`; // Calculate cell size
    minefield.style.gridTemplateColumns = `repeat(${gameState.field_len}, ${cellSize})`;
    minefield.style.gridTemplateRows = `repeat(${gameState.field_len}, ${cellSize})`;

    // Create a cell for each position in the grid
    for (let i = 0; i < gameState.field_len * gameState.field_len; i++) {
        const cell = document.createElement('div');
        cell.className = 'cell';
        cell.dataset.index = i; // Store the index for later use in click handling
        cell.dataset.fieldLen = gameState.field_len;
        // Determine the type of the cell and add appropriate class if they are not null
        if (gameState.bombs && gameState.bombs.includes(i)) {
            cell.classList.add('bomb');
        } else if (gameState.keys && gameState.keys.includes(i)) {
            cell.classList.add('key');
        } else if (gameState.empty && gameState.empty.includes(i)) {
            cell.classList.add('empty');
        } else {
            // If none are specified or if the cell index is not included in the arrays, 
            // it's either not_seen or empty depending on the is_default_not_seen flag
            if (gameState.is_default_not_seen ||
                (gameState.not_seen && gameState.not_seen.includes(i))) {
                cell.classList.add('not_seen');
                cell.addEventListener('click', cellClicked); // Only not_seen cells are clickable
            } else {
                cell.classList.add('empty');
            }
        }
        // Set cell type as a class and background image
        // cell.classList.add(cellType);
        // console.log(cell.classList);
        cell.style.backgroundImage = `url('${cellImages[cell.classList[1]]}')`;
        minefield.appendChild(cell);
    }
}

function initializeGridFromState(gameState) {
    const minefield = document.getElementById('minefield');
    minefield.innerHTML = ''; // Clear the grid first

    // Use the gameState to create the grid
    gameState.forEach((row, rowIndex) => {
        row.forEach((cellState, colIndex) => {
            const cell = createCellElement(rowIndex, colIndex, cellState);
            minefield.appendChild(cell);
        });
    });
}

function displayUserInfo(data) {
    // Update the UI with the fetched data
    document.getElementById('username').textContent = data.username;
    document.getElementById('points-left').textContent = data.points_left;
    document.getElementById('num-of-keys').textContent = data.number_of_keys;
    document.getElementById('next-move-cost').textContent = data.next_move_cost;
    document.getElementById('normal-move-cost').textContent = data.normal_move_cost;
    document.getElementById('bomb-move-cost').textContent = data.bomb_move_cost;
    // Update other elements for NextMoveCost, NormalMoveCost, and BombMoveCost if necessary
}

function cellClicked(event) {
    const cell = event.target;
    const index = parseInt(cell.dataset.index, 10);
    const fieldLen = parseInt(cell.dataset.fieldLen, 10); // Get the fieldLen from the parent element    console.log(fieldLen);
    console.log(fieldLen);
    const y = index % fieldLen;
    const x = Math.floor(index / fieldLen);

    // Check if the cell is clickable (not_seen)
    if (!cell.classList.contains('not_seen')) {
        return; // If the cell is not clickable, do nothing
    }

    // Construct the request payload
    const payload = {
        x: x,
        y: y
    };
    // Confirm with the player
    if (confirm('Do you want to pay to open this cell?')) {
        // Send POST request to server
        fetch('https://minefield.onrender.com/player/play', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + localStorage.getItem('jwtToken') // Include the JWT token
            },
            body: JSON.stringify(payload),
        })
            .then(response => {
                if (!response.ok) {
                    if (response.status === 400) {
                        // Handle the 400 status without consuming the body here
                        response.text().then(text => {
                            const data = JSON.parse(text);
                            if (data === "can't make any more moves, player has won") {
                                displayWinMessage();
                                disableAllCells();
                            }
                        });
                    }
                    // Throw an error to skip the next .then() because we've already handled the error case
                    throw new Error('Bad response from server');
                }
                return response.json();
            })
            .then(data => {
                if (data.next_move_cost === data.bomb_move_cost) {
                    // User clicked on a bomb, ask if they want to continue
                    const userWantsToContinue = confirm("You've clicked on a bomb! Do you want to pay " + data.bomb_move_cost + " points to continue?");
                    if (userWantsToContinue) {
                        // User chooses to continue, proceed as normal
                        displayUserInfo(data); // Display user info
                        updateGameState(data.game_state);
                    } else {
                        // User chooses not to continue, send a request to lose the game
                        loseGame();
                    }
                } else {
                    // Handle updating the cell based on the response
                    displayUserInfo(data); // Display user info
                    updateGameState(data.game_state);
                    if (data.number_of_keys === 5) {
                        displayWinMessage();
                        disableAllCells();
                    }
                }
            })
            .catch(error => {
                console.error('Error:', error);
            });
    }
}

function loseGame() {
    fetch('https://minefield.onrender.com/player/lose', {
        method: 'POST',
        headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('jwtToken'),
            'Content-Type': 'application/json'
        }
        // Include any additional required body data for your POST request
    })
        .then(response => {
            if (response.ok) {
                fetchGameState(); // Fetch the game state again to reset the game
            } else {
                console.error('Error: Failed to update game loss status');
            }
        })
        .catch(error => {
            console.error('Error:', error);
        });
}

function updateGameState(gameState) {
    // Call your existing initializeGrid function to re-create the grid with the new game state
    initializeGrid(gameState);
}

function updateGrid(data) {
    const size = data.field_len;
    const allCells = document.getElementsByClassName('cell');
    const defaultClass = data.is_default_not_seen ? 'not_seen' : 'empty';

    // Initially set all cells to the default state based on is_default_not_seen
    for (let i = 0; i < allCells.length; i++) {
        allCells[i].className = `cell ${defaultClass}`;
    }

    // Update cells based on the lists provided in the response
    // Each list represents cells of a different type
    ['bombs', 'keys', 'empty', 'not_seen'].forEach(type => {
        if (Array.isArray(data[type])) {
            data[type].forEach(index => {
                console.log(index, type);
                updateCell(index, type, size, allCells);
            });
        }
    });

    // Assuming 'key_shards' might be displayed elsewhere on the page
    // Update key shards display if needed
    // document.getElementById('keyShardsDisplay').textContent = 'Key Shards: ' + data.key_shards;
}

function updateCell(index, type, size, allCells) {
    const x = index % size;
    const y = Math.floor(index / size);
    const cellIndex = y * size + x;
    const cellElement = allCells[cellIndex];
    cellElement.className = `cell ${type}`;
    cellElement.dataset.status = type; // Update status
    if (type !== 'not_seen') {
        cellElement.removeEventListener('click', cellClicked);
    } else {
        // Ensure that the event listener is in place for 'not_seen' cells
        cellElement.addEventListener('click', cellClicked);
    }
    switch (type) {
        case 'not_seen':
            cellElement.style.backgroundImage = "url('images/not_seen.png')";
            break;
        case 'bombs':
            cellElement.style.backgroundImage = "url('images/bomb.png')";
            break;
        case 'empty':
            cellElement.style.backgroundImage = "url('images/empty.png')";
            break;
        case 'keys':
            cellElement.style.backgroundImage = "url('images/key.png')";
            break;
    }
}

document.getElementById('restart-button').addEventListener('click', function () {
    const topUpPoints = parseInt(document.getElementById('top-up-points').value, 10);
    const openingCost = parseInt(document.getElementById('opening-cost').value, 10);
    const bombOpeningCost = parseInt(document.getElementById('bomb-opening-cost').value, 10);
    const fieldLen = parseInt(document.getElementById('field-len').value, 10);
    const bombPercent = parseInt(document.getElementById('bomb-percentage').value, 10);
    const numOfKeys = parseInt(document.getElementById('number-of-keys').value, 10);

    const restartData = {
        top_up_points: topUpPoints,
        opening_cost: openingCost,
        bomb_opening_cost: bombOpeningCost,
        field_len: fieldLen,
        bomb_percentage: bombPercent,
        number_of_keys: numOfKeys
    };

    // Replace with the actual URL of your backend endpoint
    fetch('https://minefield.onrender.com/player/restart', {
        method: 'POST',
        headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('jwtToken'),
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(restartData)
    })
        .then(response => {
            if (response.ok) {
                // Restart was successful, fetch the game state again
                clearWinMessage();
                fetchGameState();
            } else {
                // Handle errors, maybe the user didn't have enough points to restart
                console.error('Restart failed');
            }
        })
        .catch(error => {
            console.error('Error:', error);
        });
});


function displayWinMessage() {
    const winMessageBar = document.createElement('div');
    winMessageBar.setAttribute('id', 'win-message-bar');
    winMessageBar.textContent = "Congratulations! You Win!";
    winMessageBar.classList.add('win-message-bar'); // Use the class for styling
    document.body.appendChild(winMessageBar);
}

function disableAllCells() {
    const cells = document.querySelectorAll('#minefield .cell');
    cells.forEach(cell => {
        cell.removeEventListener('click', cellClicked);
        cell.dataset.clickable = 'false'; // Mark the cell as non-clickable
    });
}


function clearWinMessage() {
    const winMessageBar = document.getElementById('win-message-bar');
    if (winMessageBar) {
        winMessageBar.remove();
    }
}