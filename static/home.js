
import { initializeApp } from "https://www.gstatic.com/firebasejs/10.10.0/firebase-app.js";
import { getAuth, signOut, onAuthStateChanged } from "https://www.gstatic.com/firebasejs/10.10.0/firebase-auth.js";

const firebaseConfig = {
    apiKey: "AIzaSyAOPZjFuX_9WWuVwso8RtEnY6W4ntt4ZBg",
    authDomain: "blissful-bites-9985b.firebaseapp.com",
    projectId: "blissful-bites-9985b",
    storageBucket: "blissful-bites-9985b.appspot.com",
    messagingSenderId: "1070336763456",
    appId: "1:1070336763456:web:d4eefe70dd2d4c53fd61a5",
    measurementId: "G-LWJW5JXMEL"
};

// Initialize Firebase
const app = initializeApp(firebaseConfig);
const auth = getAuth(app);

auth.onAuthStateChanged(console.log)
var email;
onAuthStateChanged(auth, (user) => {
    if (user) {
        // User is signed in
        email = user.email
        console.log("User is signed in:", user);
        console.log(email);
        const url = `/userDetails?email=${user.email}`;

        fetch(url)
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.json();
            })
            .then(data => {
                console.log("received user's data")
                var weightInKg = data.weight;
                var heightInCm = data.height;
                var heightInM = heightInCm / 100;
                var bmi = weightInKg / (heightInM * heightInM);

                var userData = `
                        <p>Name: ${data.name}</p>
                        <p>Heathscore: ${data.healthscore}</p>
                        <p>BMI: ${bmi}</p>
                        <p>Your diet plan provided by us: ${data.diet_plan}</p>
                     `;
                var trackData = data.track;
                var weightMap = extractWeightValues(trackData);
                console.log(weightMap);
                createWeightChart(weightMap);

                // createTrackDataContainer("trackDataContainer", trackData)
                document.getElementById('name').innerHTML = data.name;
                document.getElementById('hs').innerHTML = data.healthscore;
                document.getElementById('bmi').innerHTML = bmi;
                document.getElementById('diet_plan').innerHTML = data.diet_plan;
            })
            .catch(error => {
                console.error('There was a problem with the fetch operation:', error);
            });
    } else {
        // No user is signed in
        console.log("No user is signed in");
        // You can handle this case, for example, redirecting the user to the login page
    }
});

function createTrackDataContainer(elementId, trackData) {
    // Get a reference to the track data container element
    var trackDataContainer = document.getElementById(elementId);

    // Loop through each track in the track data
    if (trackData && trackData.length > 0) {
        trackData.forEach(function (track) {
            // Create a new div for the track
            var trackDiv = document.createElement("div");
            trackDiv.classList.add("track");

            // Create paragraphs for each property of the track
            for (var key in track) {
                if (track.hasOwnProperty(key)) {
                    var trackProperty = document.createElement("p");
                    trackProperty.textContent = key + ": " + track[key];
                    trackDiv.appendChild(trackProperty);
                }
            }

            // Append the track div to the track data container
            trackDataContainer.appendChild(trackDiv);
        });
    } else {
        console.log("trackData is null or empty");
    }
}

function extractWeightValues(trackData) {
    var weightMap = {}; // Object to store weight values mapped with date

    // Loop through each track in the track data
    if (trackData && trackData.length > 0) {
        trackData.forEach(function (track) {
            if (track.hasOwnProperty("weight") && track.hasOwnProperty("date")) {
                if (weightMap.hasOwnProperty(track.date)) {
                    weightMap[track.date].push(track.weight);
                } else {
                    weightMap[track.date] = [track.weight];
                }
            }
        });
    } else {
        console.log("trackData is null or empty");
    }

    return weightMap;
}


const logout = document.getElementById("logoutButton");
logout.addEventListener('click', () => {
    signOut(auth).then(() => {
        // Sign-out successful
        console.log("User signed out successfully");
        window.location.href = "/login";

    }).catch((error) => {
        console.error("Error signing out:", error);
    });
});

const logout2 = document.getElementById("logoutButton2");
logout2.addEventListener('click', () => {
    signOut(auth).then(() => {
        // Sign-out successful
        console.log("User signed out successfully");
        window.location.href = "/login";

    }).catch((error) => {
        console.error("Error signing out:", error);
    });
});