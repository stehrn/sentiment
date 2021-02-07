// Import stylesheets
import "./style.css";
// Firebase App (the core Firebase SDK) is always required and must be listed first
import * as firebase from "firebase/app";

// Add the Firebase products that you want to use
import "firebase/auth";
import "firebase/firestore";
import "firebase/storage";

import * as firebaseui from "firebaseui";

// Document elements
const startLoginButton = document.getElementById("login");
const photoContainer = document.getElementById("photo-container");

const form = document.getElementById("upload-photo");
const input = document.getElementById("message");
const photos = document.getElementById("photos");

const dropArea = document.getElementById("drop-area");

var photoListener = null;

async function main() {
  // Add Firebase project configuration object here
  var firebaseConfig = {
    apiKey: "AIzaSyD2KDL5GYEWRXU5j1q7ELQEHKlW7gH0xuU",
    authDomain: "sentiment-302320.firebaseapp.com",
    projectId: "sentiment-302320",
    storageBucket: "sentiment-302320.appspot.com",
    messagingSenderId: "708629737434",
    appId: "1:708629737434:web:d0c7672cf8a559500a1ae3",
    measurementId: "G-JQ7L4VXGMC"
  };

  firebase.initializeApp(firebaseConfig);

  // FirebaseUI config
  const uiConfig = {
    credentialHelper: firebaseui.auth.CredentialHelper.NONE,
    signInOptions: [
      // Email / Password Provider.
      firebase.auth.EmailAuthProvider.PROVIDER_ID
    ],
    callbacks: {
      signInSuccessWithAuthResult: function(authResult, redirectUrl) {
        // Handle sign-in.
        // Return false to avoid redirect.
        return false;
      }
    }
  };

  const ui = new firebaseui.auth.AuthUI(firebase.auth());

  // Log in/out
  startLoginButton.addEventListener("click", () => {
    if (firebase.auth().currentUser) {
      // User is signed in; allows user to sign out
      firebase.auth().signOut();
    } else {
      // No user is signed in; allows user to sign in
      ui.start("#firebaseui-auth-container", uiConfig);
    }
  });

  // Change state depending on if user logged in/out
  firebase.auth().onAuthStateChanged(user => {
    if (user) {
      startLoginButton.textContent = "LOGOUT";
      // Show guestbook to logged-in users
      photoContainer.style.display = "block";

      // Subscribe to the guestbook collection
      subscribePhotos();
    } else {
      startLoginButton.textContent = "Log In";
      // Hide guestbook for non-logged-in users
      photoContainer.style.display = "none";

      // Unsubscribe from the guestbook collection
      unsubscribePhotos();
    }
  });

  // Listen to photo updates
  function subscribePhotos() {
    console.log("subscribing to photos");
    // Create query for photos
    photoListener = firebase
      .firestore()
      .collection("users")
      .doc(firebase.auth().currentUser.uid)
      .collection("photos")
      .orderBy("timestamp", "desc")
      .onSnapshot(snaps => {
        // Reset page
        photos.innerHTML = "";
        // Loop through documents in database
        snaps.forEach(doc => {
          // Create an HTML entry for each photoContainer
          const text = document.createElement("p");
          text.textContent = doc.data().text;
          photos.appendChild(text);

          const sentiment = document.createElement("p");
          var sentimentDoc = doc.data().sentiment;
          text.textContent =
            "Joy: " +
            sentimentDoc.JoyLikelihood +
            ", anger: " +
            sentimentDoc.AngerLikelihood;
          photos.appendChild(sentiment);

          const thumb = document.createElement("img");
          thumb.setAttribute("src", doc.data().imageUrl);
          photos.appendChild(thumb);

          console.log("photo update: %s", doc.data().imageUrl);
        });
      });
  }

  // Unsubscribe from photo updates
  function unsubscribePhotos() {
    console.log("unsubscribing from photos");
    if (photoListener != null) {
      photoListener();
      photoListener = null;
    }
  }

  dropArea.addEventListener("dragover", event => {
    event.stopPropagation();
    event.preventDefault();
    // Style the drag-and-drop as a "copy file" operation.
    event.dataTransfer.dropEffect = "copy";
  });

  dropArea.addEventListener("drop", event => {
    event.stopPropagation();
    event.preventDefault();
    const fileList = event.dataTransfer.files;
    console.log(fileList);
  });

  var files = [];
  document.getElementById("file").addEventListener("change", function(e) {
    files = e.target.files;
  });

  form.addEventListener("submit", e => {
    // Prevent the default form redirect
    e.preventDefault();

    //checks if files are selected
    if (files.length == 0) {
      alert("No file chosen");
      return;
    }

    // Handle photo upload
    addPhoto(files[0]);

    // Return false to avoid redirect
    return false;
  });

  // 1. add or update user doc
  // 2. add new doc to photo collection of user
  // 3. add photo to cloud storage
  function addPhoto(file) {
    console.log(
      "Adding photo: %s, for user: %s",
      file.name,
      firebase.auth().currentUser.displayName
    );

    var userDocRef = firebase
      .firestore()
      .collection("users")
      .doc(firebase.auth().currentUser.uid);

    // 1. add or update user doc
    userDocRef
      .set({
        name: firebase.auth().currentUser.displayName,
        timestamp: firebase.firestore.FieldValue.serverTimestamp()
      })
      .then(function() {
        // 2. add new doc to photo collection of user
        console.log("User docoument set, adding photo");
        return userDocRef
          .collection("photos")
          .add({
            text: input.value,
            timestamp: firebase.firestore.FieldValue.serverTimestamp()
          })
          .then(function(photoDocRef) {
            console.log("Photo added, doc ID: ", photoDocRef.id);
            return photoDocRef;
          })
          .catch(function(error) {
            console.error("Error writing photo to database", error);
          });
      })
      .then(function(photoDocRef) {
        // 3. add photo to cloud storage
        var metadata = {
          customMetadata: {
            photoDocId: photoDocRef.id // used by cloud storage function to udate photo with sentiment
          }
        };

        var extention = file.type.split("/").pop();
        var filePath =
          firebase.auth().currentUser.uid +
          "/" +
          photoDocRef.id +
          "." +
          extention;
        return firebase
          .storage()
          .ref(filePath)
          .put(file, metadata)
          .then(function(fileSnapshot) {
            console.log("Photo uploaded to cloud storage");
            // generate a public URL for the file
            return fileSnapshot.ref.getDownloadURL().then(url => {
              // update the photo doc with the image's URL
              console.log("Updating photo doc: ", photoDocRef.id);
              return photoDocRef.update({
                imageUrl: url,
                storageUri: fileSnapshot.metadata.fullPath
              });
            });
          })
          .catch(function(error) {
            console.error("Error uploading photo", error);
          });
      })
      .catch(function(error) {
        console.error("Error writing user document to database", error);
      })
      .then(function() {
        input.value = "";
        input.file = "";
      });
  } //addPhoto
} // main

main();
