<template>
  <div>
    <div class="banner">
      <p v-if="storedUsername">Welcome, {{ storedUsername }}</p>
      <p v-else>Welcome, Guest</p>
    </div>
    <div>
      <input v-model="username" placeholder="Enter Username" />
      <button @click="login">Login</button>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
import api from "@/services/axios"; 
import router from '@/router';

export default {
  data() {
    return {
      username: '', // Tracks login status
      storedUsername: localStorage.getItem("username") || '' 
    }
  },
  methods: {
    async login() {
      try {
        const inputUsername = `"${this.username}"`;
        const response = await this.$axios.post('/session', inputUsername, {
  headers: {
    'Content-Type': 'text/plain', // Set content type to text/plain
  },
});
        // Check if the response status is 200 or 201
        if (response.status === 200 || response.status === 201) {
          // Store token in local storage
          localStorage.setItem("token", response.data);
          // Store username in local storage
          localStorage.setItem("username", this.username);
          // Set token in axios defaults for future requests
          axios.defaults.headers.common['Authorization'] = `Bearer ${response.data}`;
          // Redirect to home after successful login
          router.push('/');
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
  }
}
</script>
