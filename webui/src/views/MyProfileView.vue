<template>
  <div class="container">
    <div class="search-bar-container">
      <div v-if="userProfile" class="user-info">
        <p class="username">{{ userProfile.Username }}</p>
        <p class="follower-count">Followers: {{ userProfile.Followers ? userProfile.Followers.length : 0 }}</p>
        <p class="following-count">Following: {{ userProfile.Following ? userProfile.Following.length : 0 }}</p>
      </div>
    </div>
    <div v-if="userProfile">
      <PhotoComponent :photos="userProfile.Photos" />
    </div>
    <div v-else-if="!localStorage.getItem('username')">
      <p>Log in to see your profile!</p>
    </div>
  </div>
</template>


<script>
import instance from '../services/axios.js';
import PhotoComponent from '../components/PhotoComponent.vue';

export default {
  components: {
    PhotoComponent
  },
  data() {
    return {
      username: localStorage.getItem('username'),
      userProfile: null,
    }
  },
  methods: {
    async getMyUserProfile() {
      try {
        const response = await instance.get(`/users/${localStorage.getItem('username')}`, {
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        });
        console.log("this is the response data")
        console.log(response.data)
        this.userProfile = response.data;
      } catch (error) {
        console.error(error);
        this.userProfile = null
        console.error('Error occurred while searching for user:', error);
      }
    }
  },
  created() {
    if (localStorage.getItem('username')) {
      this.getMyUserProfile();
    }
  }
}
</script>

<style scoped>
.container {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  margin-left: 20px;
}

.search-bar-container {
  margin-bottom: 20px;
}

.search-bar {
  display: flex;
  align-items: center;
}

input[type='text'] {
  width: 300px;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 5px 0 0 5px;
  font-size: 16px;
}

button {
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

.username {
  font-size: 14px;
  margin-bottom: 10px;
}

.follower-count,
.following-count {
  text-align: left;
  margin-bottom: 5px;
}

.button-row {
  display: flex;
  gap: 10px;
}
</style>