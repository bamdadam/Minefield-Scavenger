// game_logic.js

let playerBetAmount = 10; // Default bet amount
document.addEventListener('DOMContentLoaded', function () {
    fetchUserData();
    const resultMessage = document.getElementById('resultMessage');

    // Hide placeholders initially
    document.querySelectorAll('.choice-placeholder').forEach(function (placeholder) {
        placeholder.style.backgroundImage = '';
    });

    document.querySelectorAll('.choiceButton').forEach(button => {
        button.addEventListener('click', function () {
            setButtonsClickable(false);
            const playerChoice = this.getAttribute('data-choice');
            // const playerBet = 10; // Update as needed or get from user input
            const gameVersion = document.getElementById('gameVersionToggle').checked ? 1 : 2;
            // const choiceImages = ['images/rock.png', 'images/paper.png', 'images/scissor.png'];
            setImageChoice('playerChoiceImg', playerChoice);
            // Set player's choice image
            // document.getElementById('playerChoice').style.backgroundImage = `url(${choiceImages[playerChoice]})`;

            // Prepare request data
            const requestData = {
                player_choice: parseInt(playerChoice, 10),
                player_bet: parseInt(playerBetAmount, 10),
                game_version: gameVersion
            };
            console.log(requestData);

            // Post request to the server
            fetch('https://minefield.onrender.com/player/play/rps/', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': 'Bearer ' + localStorage.getItem('jwtToken') // Include the JWT token
                },
                body: JSON.stringify(requestData)
            })
                .then(response => response.json())
                .then(data => {
                    // Update enemy's choice image
                    setImageChoice('enemyChoiceImg', data.enemy_choice);
                    // Update points and username
                    updateUserData(data)

                    if (data.user_choice === data.enemy_choice) {
                        resultMessage.textContent = "It's a draw!";
                    } else {
                        resultMessage.textContent = data.has_won ? "You won!" : "You lost!";
                    }

                    // Show result message
                    // resultMessage.textContent = data.has_won ? "You won!" : "You lost!";
                    resultMessage.classList.remove('fade-out');
                    resultMessage.style.opacity = 1; // Reset opacity if it was faded out before

                    // Fade out result message
                    setTimeout(() => {
                        setButtonsClickable(true);
                        resultMessage.classList.add('fade-out');
                        clearImageChoice('playerChoiceImg');
                        clearImageChoice('enemyChoiceImg');
                    }, 2000);

                    // Clear message after fade out
                    setTimeout(() => {
                        resultMessage.textContent = '';
                    }, 2000);
                })
                .catch(error => {
                    console.error('Error:', error);
                    setButtonsClickable(true);
                });
        });
    });
});

function setImageChoice(choiceElementId, choice) {
    const choiceImages = {
        '0': 'images/rock.png',
        '1': 'images/paper.png',
        '2': 'images/scissors.png'
    };

    const imgElement = document.getElementById(choiceElementId);
    imgElement.src = choiceImages[choice];
    imgElement.parentElement.classList.add('image-set'); // Add class to hide placeholder text
}

function clearImageChoice(choiceElementId) {
    const imgElement = document.getElementById(choiceElementId);
    imgElement.src = '';
    imgElement.parentElement.classList.add('image-set'); // Add class to hide placeholder text
}


function updateGameModeText() {
    const gameModeText = document.getElementById('gameModeText');
    const gameVersionToggle = document.getElementById('gameVersionToggle');
    gameModeText.textContent = 'Game Mode: ' + (gameVersionToggle.checked ? '1' : '2');
}

// Set initial game mode text when the page loads
updateGameModeText();

// Event listener for the game version toggle switch
document.getElementById('gameVersionToggle').addEventListener('change', function () {
    updateGameModeText();
});

function setButtonsClickable(state) {
    document.querySelectorAll('.choiceButton').forEach(button => {
        button.disabled = !state; // If state is true, buttons are clickable (not disabled)
    });
}

document.querySelectorAll('.bet-switch input[type="radio"]').forEach(radio => {
    radio.addEventListener('change', function() {
        // const indicator = document.querySelector('.bet-indicator');
        const betOption = document.querySelector(`label[for="${this.id}"]`);
        const betOptions = document.querySelectorAll('.bet-option');
        const betBackground = document.querySelector('.bet-switch');

        // Update the color of the switch based on the bet amount
        let color;
        switch(this.value) {
            case '10': color = '#4CAF50'; break;
            case '20': color = '#8BC34A'; break;
            case '30': color = '#CDDC39'; break;
            case '40': color = '#FFEB3B'; break;
            case '50': color = '#FFC107'; break;
            default: color = '#4CAF50';
        }

        // Move the indicator under the selected label
        // indicator.style.left = `${betOption.offsetLeft}px`;
        // indicator.style.width = `${betOption.offsetWidth}px`;
        betBackground.style.backgroundColor = color;

        // Reset other options to default
        betOptions.forEach(option => {
            if (option.htmlFor !== radio.id) {
                option.style.backgroundColor = '#fff';
                option.style.color = 'black';
            }
            option.addEventListener('change', updateBetAmount);
        });
    });
});

document.getElementById('bet10').checked = true;
document.getElementById('bet10').dispatchEvent(new Event('change'));
document.addEventListener('DOMContentLoaded', updateBetAmount);

function updateBetAmount() {
    const betOptions = document.querySelectorAll('.bet-switch input[type="radio"]');
    betOptions.forEach(option => {
        if (option.checked) {
            playerBetAmount = parseInt(option.value, 10);
            console.log("Player's bet amount is now:", playerBetAmount);
        }
    });
}

const betOptions = document.querySelectorAll('.bet-switch input[type="radio"]');
betOptions.forEach(option => {
    option.addEventListener('change', updateBetAmount);
});

function updateUserData(userData) {
    document.getElementById('username').textContent = userData.username;
    document.getElementById('pointsLeft').textContent = userData.points_left;
}

// Function to query the backend for user data
function fetchUserData() {
    fetch('https://minefield.onrender.com/player/data/rps/', {
        method: 'GET', // or 'POST' if required
        headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('jwtToken'),
            'Content-Type': 'application/json'
        }
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok ' + response.statusText);
        }
        return response.json();
    })
    .then(userData => {
        updateUserData(userData); // Update the frontend with the received data
    })
    .catch(error => {
        console.error('There has been a problem with your fetch operation:', error);
    });
}