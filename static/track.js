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
        document.getElementById("email").value = email;
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
                
                var trackData = data.track;
                
                createTrackDataContainer("trackDataContainer", trackData)
            })
            .catch(error => {
                console.error('There was a problem with the fetch operation:', error);
            });
    } else {
        console.log("No user is signed in");
    }
});

document.getElementById('date').value = new Date().toISOString().slice(0, 10);
document.getElementById('trackForm').addEventListener('submit', function (event) {
    event.preventDefault();

    var formData = new FormData(this);

    fetch('/trackMeal', {
        method: 'POST',
        body: formData
    })
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            // Clear the form input
            document.getElementById('breakfast').value = '';
            document.getElementById('lunch').value = '';
            document.getElementById('dinner').value = '';
            document.getElementById('weight').value = '';

            alert('Your meals have been tracked successfully!');
        })
        .catch(error => {
            alert('Oops! Something went wrong. Please try again later.');
        });
});

function createTrackDataContainer(elementId, trackData) {
    // Get a reference to the track data container element
    var trackDataContainer = document.getElementById(elementId);

    // Loop through each track in the track data
    if (trackData && trackData.length > 0) {
        trackData.forEach(function (track) {
            // Create a new container div for the track
            var trackContainer = document.createElement("div");
            trackContainer.classList.add("u-clearfix", "u-sheet", "u-sheet-1");

            // Create the inner container div
            var innerContainer = document.createElement("div");
            innerContainer.classList.add("u-container-style", "u-custom-color-2", "u-expanded-width-lg", "u-expanded-width-md", "u-expanded-width-sm", "u-expanded-width-xs", "u-group", "u-radius", "u-shape-round", "u-group-1");

            // Create the layout div
            var layoutDiv = document.createElement("div");
            layoutDiv.classList.add("u-container-layout", "u-container-layout-1");

            // Loop through each property in the track object
            for (var key in track) {
                if (track.hasOwnProperty(key)) {
                    // Create a new paragraph element for each property
                    var paragraph = document.createElement("b");
                    paragraph.classList.add("u-custom-font", "u-text", "u-text-default");
                    if (key == "weight"){
                        paragraph.textContent = "Calculated Weight: " + track[key];
                    } else {
                        paragraph.textContent = key+":   "+ track[key];
                    }
                    
                    // Append the paragraph to the layout div
                    layoutDiv.appendChild(paragraph);
                }
            }

            // Append the inner container to the track container
            innerContainer.appendChild(layoutDiv);

            // Append the inner container to the track container
            trackContainer.appendChild(innerContainer);

            // Append the track container to the track data container
            trackDataContainer.appendChild(trackContainer);
        });
    } else {
        console.log("trackData is null or empty");
    }
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