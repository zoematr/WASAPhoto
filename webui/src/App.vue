<template>
  <div>
    <header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
      <div :key="componentKey">
        <a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">WASAPhoto</a>
      </div>
      <button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
        <span class="navbar-toggler-icon"></span>
      </button>
    </header>
    <input type="file" ref="fileInput" hidden @change="uploadPhoto" accept="image/*" />
    <div class="container-fluid">
      <div class="row">
        <nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
          <div class="position-sticky pt-3 sidebar-sticky">
            <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
            </h6>
            <ul class="nav flex-column">
              <div class="banner">
                <p v-if="storedUsername">Welcome, {{ usernameComputed }}</p>
                <p v-else>Welcome, Guest</p>
              </div>
              <form @submit.prevent="login">
                <div class="form-group">
                  <input type="text" v-model="inputUsername" class="form-control" id="inputUsername" placeholder="Enter your username" required>
                </div>
                <button type="submit">Login</button>
              </form>
              <li class="nav-item">
                <RouterLink to="/stream" class="nav-link">
                  <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#layout"/></svg>
                  My stream
                </RouterLink>
              </li>
              <li class="nav-item">
                <RouterLink to="/users/" class="nav-link">
                  <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#key"/></svg>
                  Search users
                </RouterLink>
              </li>
            </ul>
            <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
              <span>Actions menu</span>
            </h6>
            <ul class="nav flex-column">
              <li class="nav-item">
                <a class="nav-link" href="#" @click="triggerFileInput">
                  <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#file-text"/></svg>
                  Upload Photo
                </a>
              </li>
              <li class="nav-item">
                <a class="nav-link" href="#" @click="logmeout">
                  <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#file-text"/></svg>
                  Log out
                </a>
              </li>
              <h1 class="small">Change Username</h1>
              <form @submit.prevent="changeUsername">
                <div class="form-group">
                  <input type="text" v-model="newUsername" class="form-control" id="newUsername" placeholder="Enter new username" required>
                </div>
                <button type="submit">Change Username</button>
              </form>
            </ul>
          </div>
        </nav>
        <main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
          <RouterView />
        </main>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed } from 'vue';
import { RouterLink, RouterView } from 'vue-router';
import instance from './services/axios.js';

export default {
  data() {
    return {
      storedUsername: localStorage.getItem("username") || '' 
    }
  },
  setup() {
    // Computed property to get the username from local storage
    const usernameComputed = computed(() => {
      const username = localStorage.getItem("username") || "";
      console.log("Username:", username);
      return username;
    });
    const fileInput = ref(null);
    const newUsername = ref('');

    // Method to trigger file input
    function triggerFileInput() {
      fileInput.value.click();
    }
    
    async function uploadPhoto(event) {
      const file = event.target.files[0];
      
      // Create a new FileReader to read the file data
      const reader = new FileReader();
      
      // Define an onload function for the reader
      reader.onload = async function () {
        // Convert the file data to base64
        const fileBase64 = reader.result.split(',')[1];
        console.log("Base64 encoded photo:", fileBase64);
        
        if (fileBase64) {
          try {
            // Send the base64-encoded photo data to the backend
            const username = localStorage.getItem('username');
            const response = await instance.post(`/users/${username}/photos/`, { photofile: fileBase64 }, {
              headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('token')}`
              }
            });
            alert('You posted your photo!');
            location.reload();
          } catch (error) {
            console.error('Error uploading image:', error);
          }
        }
      };
      
      // Define an onerror function for the reader
      reader.onerror = function (error) {
        console.log('Error reading file:', error);
      };
      
      // Read the file as a data URL (base64)
      reader.readAsDataURL(file);
    }

    async function changeUsername() {
      console.log("Change username function called");
      console.log("Username computed value:", usernameComputed.value);
      try {
        
        const response = await instance.patch(`/users/${usernameComputed.value}`, JSON.stringify({ newusername: newUsername.value }), {
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        });
        localStorage.setItem("username", newUsername.value)
        alert('You successfully changed your username!');
        window.location.reload();
      } catch (error) {
          console.error('An error occurred:', error);
          alert("There was an error changing the username :( Try another one!");
      }
    }

    async function login() {
      try {
        console.log("this is the inputUsername");
        console.log(inputUsername.value);
        const username = `"${inputUsername.value}"`;

        const response = await instance.post('/session', username, {
          headers: {
            'Content-Type': 'text/plain'
          }
        });

        if (response.status === 200 || response.status === 201) {
          // Store token in local storage
          const token = response.data;
          console.log(token)
          localStorage.setItem("token",token);
          console.log("login: this is the token", token)
          // Store username in local storage
          localStorage.setItem("username", inputUsername.value);
          console.log("this is usernameComputed.value")
          console.log(usernameComputed.value)
          // Set token in axios defaults for future requests
          instance.defaults.headers.common['Authorization'] = `Bearer ${token}`;
          // Redirect to home after successful login
          // router.push('/stream');
          // router.push('/');
          // Reload the page to ensure the token is properly set
          location.reload();
        } else {
          // Handle other status codes if needed
          console.error("Login failed:", response.statusText);
        }
      } catch (error) {
        // Handle error responses
        console.error("Login failed:", error);
      }
    }

    async function logmeout() {
      localStorage.setItem("username","");
      localStorage.setItem("token","");
      alert('You have logged out! Goodbye'); 
      location.reload();
    }


    return {
      usernameComputed,
      fileInput,
      newUsername,
      triggerFileInput,
      changeUsername,
      logmeout,
      uploadPhoto,
      login,
    };
  },
};

</script>

<style>
  .small {
    font-size: 1.0rem;
  }

  button {
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 0 5px 5px 0;
  font-size: 16px;
  cursor: pointer;
  }
</style>