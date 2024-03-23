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
    <input type="file" ref="fileInput" hidden @change="uploadImage" accept="image/*" />
    <div class="container-fluid">
      <div class="row">
        <nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
          <div class="position-sticky pt-3 sidebar-sticky">
            <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
              <span>General</span>
            </h6>
            <ul class="nav flex-column">
              <li class="nav-item">
                <RouterLink to="/" class="nav-link">
                  <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#home"/></svg>
                  Home
                </RouterLink>
              </li>
              <li class="nav-item">
                <RouterLink :to="'/user/' + usernameComputed" class="nav-link">
                  <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#home"/></svg>
                  My Profile
                </RouterLink>
              </li>
              <li class="nav-item">
                <RouterLink to="/stream" class="nav-link">
                  <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#layout"/></svg>
                  My stream
                </RouterLink>
              </li>
              <li class="nav-item">
                <RouterLink to="/user/" class="nav-link">
                  <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#key"/></svg>
                  A users profile
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
              <li class="nav-item">
                <a class="nav-link" href="#" @click="changeUsername">
                  <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#file-text"/></svg>
                  Change Username
                </a>
              </li>
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

<script setup>
import { ref, computed } from 'vue';
import { RouterLink, RouterView } from 'vue-router';
import axios from 'axios'; // Import Axios

// Computed property to get the username from local storage
const usernameComputed = computed(() => localStorage.getItem('username') || '');
// Ref for file input
const fileInput = ref(null);

// Method to trigger file input
function triggerFileInput() {
  fileInput.value.click();
}

const newUsername = ref('');

async function changeUsername() {
  try {
    const response = await axios.patch(`/users/${localStorage.getItem('username')}`, { newusername: newUsername.value }, {
      headers: {
        'Authorization': localStorage.getItem('token'),
        'Content-Type': 'application/json'
      }
    });
    // Handle success response
    console.log(response.data);
  } catch (error) {
    // Handle error response
    console.error('Error changing username:', error.response.data);
  }
}

async function logmeout() {
  localStorage.setItem("username","");
  localStorage.setItem("token","");
  alert('You have logged out successfully! Now you can only search people, not interact with profiles.'); 
  location.reload();
} 

// Async method to handle image upload
/*
async function uploadImage(event) {
  const file = event.target.files[0];
  if (file) {
    const formData = new FormData();
    formData.append('image', file);

    try {
      const response = await axios.post(`/users/${usernameComputed.value}/photos/`, formData, {
        headers: {
          'Authorization': localStorage.getItem('token'),
          'Content-Type': 'multipart/form-data' // Set content type for FormData
        }
      });
      alert('Image uploaded successfully!');
      location.reload(); // Reload the page after successful upload
    } catch (error) {
      console.error('Error uploading image:', error);
      alert('Error uploading image: ' + error.message); // Display error message
    }
  }
}*/
// Async method to handle image upload
async function uploadImage(event) {
  const file = event.target.files[0];
  if (file) {
    const formData = new FormData();
    formData.append('image', file);

    try {
      const username = localStorage.getItem('username');
      const response = await axios.post(`/users/${username}/photos/`, formData, {  
        headers: { Authorization: localStorage.getItem("token") }
      });
      alert('Image uploaded successfully!'); 
      location.reload();
      // Handle the response, e.g., showing a success message
    } catch (error) {
      console.error('Error uploading image:', error);
      // Handle the error, e.g., showing an error message
    }
  }
}

</script>

<style>
</style>
