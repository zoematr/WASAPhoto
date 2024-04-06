<template>
  <div>
    <div class="search-bar">
      <input v-model="username" placeholder="Search Username" @keyup.enter="searchUser">
      <button @click="searchUser">Search</button>
    </div>
    <div v-if="userProfile">
      <h2>{{ userProfile.Username }}</h2>
      <p>Followers: {{ userProfile.Followers.length }}</p>
      <p>Following: {{ userProfile.Following.length }}</p>
      <div v-for="photo in userProfile.Photos" :key="photo.PhotoId">
        <img :src="'data:image/jpeg;base64,' + photo.PhotoFile" alt="Photo">
        <p>{{ photo.Date }}</p>
        <p>Likes: {{ photo.Likes.length }}</p>
        <p>Comments: {{ photo.Comments.length }}</p>
      </div>
    </div>
    <div v-else-if="searched && !userProfile">
      <p>User not found</p>
    </div>
  </div>
</template>

<script>
import instance from '../services/axios.js';

export default {
  data() {
    return {
      username: '',
      userProfile: null,
      searched: false
    }
  },
  methods: {
    async searchUser() {
      this.searched = true
      try {
        const response = await instance.get(`/users/${this.username}`, {
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        });
        this.userProfile = response.data
      } catch (error) {
        console.error(error)
      }
    }
  }
}
</script>

<style scoped>
.search-bar {
  display: flex;
  justify-content: center;
  align-items: center;
  margin-bottom: 20px;
}

input[type='text'] {
  width: 300px;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 5px 0 0 5px;
  font-size: 16px;
}

button {
  padding: 10px 20px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 0 5px 5px 0;
  font-size: 16px;
  cursor: pointer;
}

button:hover {
  background-color: #3e8e41;
}

img {
  max-width: 100%;
  height: auto;
}
</style>
