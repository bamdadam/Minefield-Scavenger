function login() {
    var username = document.getElementById('username').value;
    var points = document.getElementById('points').value;
    var openingCost = document.getElementById('openingCost').value;
    var bombOpeningCost = document.getElementById('bombOpeningCost').value;
    var FieldLen = document.getElementById('FieldLen').value;
    var BombPercent = document.getElementById('BombPercent').value;
    var NumOfKeys = document.getElementById('NumOfKeys').value;
    console.log(JSON.stringify({
        username: username,
        points: parseInt(points, 10),
        opening_cost: parseInt(openingCost, 10),
        bomb_opening_cost: parseInt(bombOpeningCost, 10),
        field_len: parseInt(FieldLen, 10),
        bomb_percentage: parseInt(BombPercent, 10),
        number_of_keys: parseInt(NumOfKeys, 10)
    }));

    // The rest of your fetch logic remains the same
    fetch('http://127.0.0.1:8080/player/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            username: username,
            points: parseInt(points, 10),
            opening_cost: parseInt(openingCost, 10),
            bomb_opening_cost: parseInt(bombOpeningCost, 10),
            field_len: parseInt(FieldLen, 10),
            bomb_percentage: parseInt(BombPercent, 10),
            number_of_keys: parseInt(NumOfKeys, 10)
        }),
    })
    .then(response => response.json())
    .then(data => {
        // Your success logic
        localStorage.setItem('jwtToken', data.token);
        localStorage.setItem('expirationTime', data.expirationTime);

        // Redirect to the game page
        window.location.href = 'game.html';
    })
    .catch(error => {
        console.error('Error:', error);
        // Your error logic
    });
}