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
import instance from '../services/axios.js';
import api from "@/services/axios"; 
import router from '@/router';

export default {
  data() {
    return {
      username: '', // Tracks login status
      storedUsername: localStorage.getItem("username") || '' 
    }
  },
  mounted(){
    storedUsername: localStorage.getItem("username");

  },
  methods: {
    async login() {
      try {
        const inputUsername = `"${this.username}"`;
        const response = await instance.post('/session', inputUsername, {
  headers: {
    'Content-Type': 'text/plain', // Set content type to text/plain
  },
});
        // Check if the response status is 200 or 201
        console.log("this resp");
        console.log(response.data);


        if (response.status === 200 || response.status === 201) {
          // Store token in local storage
          const token = response.data;
          console.log(token)
          /* const authorizationHeader = response.headers['Authorization'];
          console.log("this is the aut "+authorizationHeader);
          if (authorizationHeader) {
            const token = authorizationHeader.split(" ")[1];
          console.log("this is the token frontend")
          console.log(token) */
          localStorage.setItem("token",token);
          console.log("login: this is the token", token)
          // Store username in local storage
          localStorage.setItem("username", this.username);
          // Set token in axios defaults for future requests
          instance.defaults.headers.common['Authorization'] = `Bearer ${token}`;
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
