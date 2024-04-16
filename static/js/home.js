
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
let userData2 = null;

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
                userData2 = data
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

                var calorieMap = extractCalorieValues(trackData);
                console.log(calorieMap);

                createWeightChart(weightMap, "weightChart", "weight");
                createWeightChart(calorieMap, "calorieChart", "Total Calories");

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

    document.getElementById('gendietplan').addEventListener('click', function () {
        // Call the function to fetch data from the API
        fetchDataFromAPI(userData2);
    });
});

function fetchDataFromAPI(userData) {
    // Make a GET request using fetch API
    // const url = `/genDietPlan?email=${email}`;
    console.log(userData)
    fetch('/genDietPlan', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(userData)
    })
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            // Handle the data returned from the API
            var dietPlan = data.diet_plan;

            // Check if 'dietPlan' is a string or an array/object
            if (typeof dietPlan === 'string') {
                // If 'dietPlan' is a string, set it as the inner HTML directly
                document.getElementById('diet_plan').innerHTML = dietPlan;
            } else if (Array.isArray(dietPlan) || typeof dietPlan === 'object') {
                // If 'dietPlan' is an array or object, you need to convert it to a string first
                var jsonString = JSON.stringify(dietPlan);
                document.getElementById('diet_plan').innerHTML = jsonString;
            } else {
                // Handle other data types if needed
                console.error('Unexpected data type for diet_plan:', typeof dietPlan);
            }
        })
        .catch(error => {
            // Handle errors
            console.error('There was a problem with the fetch operation:', error);
            alert(error)
        });
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
function extractCalorieValues(trackData) {
    var weightMap = {}; // Object to store total calorie values mapped with date
  
    // Loop through each track in the track data
    if (trackData && trackData.length > 0) {
      trackData.forEach(function (track) {
        if (track.hasOwnProperty("date")) {
          var totalCalories = 0;
  
          // Check if breakfast, lunch, and dinner objects exist
          if (track.hasOwnProperty("breakfast") && track.breakfast.hasOwnProperty("Total calories")) {
            totalCalories += track.breakfast["Total calories"];
          }
          if (track.hasOwnProperty("lunch") && track.lunch.hasOwnProperty("Total calories")) {
            totalCalories += track.lunch["Total calories"];
          }
          if (track.hasOwnProperty("dinner") && track.dinner.hasOwnProperty("Total calories")) {
            totalCalories += track.dinner["Total calories"];
          }
  
          weightMap[track.date] = totalCalories;
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