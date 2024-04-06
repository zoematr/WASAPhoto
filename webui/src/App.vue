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
                  <label for="newUsername">New Username</label>
                  <input type="text" v-model="newUsername" class="form-control" id="newUsername" placeholder="Enter new username" required>
                </div>
                <button type="submit" class="btn btn-primary">Change Username</button>
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
  setup() {
    // Computed property to get the username from local storage
    const usernameComputed = computed(() => {
      const username = localStorage.getItem("username") || "";
      console.log("Username:", username);
      return username;
    });
    // Ref for file input
    const fileInput = ref(null);
    // ref 
    const newUsername = ref('');

    // Method to trigger file input
    function triggerFileInput() {
      fileInput.value.click();
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
        window.location.reload();
      } catch (error) {
    // Handle error response
        if (error.response && error.response.data && error.response.data.message) {
          console.error('Error changing username:', error.response.data.message);
        } else {
          console.error('Error changing username:', error);
        }
      }
    }


    async function logmeout() {
      localStorage.setItem("username","");
      localStorage.setItem("token","");
      alert('You have logged out! Goodbye'); 
      location.reload();
    }


    async function uploadPhoto(event) {
      const file = event.target.files[0];
      if (file) {
        const formData = new FormData();
        formData.append('image', file);

        try {
          const username = localStorage.getItem('username');
          const response = await instance.post(`/users/${username}/photos/`, formData, {  
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
          });
          alert('You posted your photo!'); 
          location.reload();
          // Handle the response, e.g., showing a success message
        } catch (error) {
          console.error('Error uploading image:', error);
          // Handle the error, e.g., showing an error message
        }
      }
    }


    return {
      usernameComputed,
      fileInput,
      newUsername,
      triggerFileInput,
      changeUsername,
      logmeout,
      uploadPhoto,
    };
  },
};

</script>

<style>
  .small {
    font-size: 1.2rem;
  }
</style>