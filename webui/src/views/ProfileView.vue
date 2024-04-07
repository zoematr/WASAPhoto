<template>
  <div class="container">
    <div class="search-bar-container">
      <div class="search-bar">
        <input v-model="searchedUsername" placeholder="Search Username" @keyup.enter="searchUser">
        <button @click="searchUser">Search</button>
      </div>
      <div v-if="userProfile" class="user-info">
        <p class="username">{{ userProfile.Username }}</p>
        <div v-if="!userProfile.OwnProfile" class="button-row">
          <button v-if="userProfile && !isOwnProfile" @click="toggleFollow">
            {{ userProfile.AlreadyFollowed ? 'Unfollow' : 'Follow' }}
          </button>
          <button v-if="userProfile && !isOwnProfile" @click="toggleBan">
            {{ userProfile.AlreadyBanned ? 'Unban' : 'Ban' }}
          </button>
        </div>
      </div>
    </div>
    <div v-if="userProfile">
      <!-- Display user's photos -->
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
      searchedUsername: '', // dont confuse with username (of the logged in user)
      username: localStorage.getItem('username'),
      userProfile: null,
      searched: false
    }
  },
  methods: {
    async searchUser() {
      this.searched = true;
      try {
        const response = await instance.get(`/users/${this.searchedUsername}`, {
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
        // Optionally, you can reset userProfile to null here
        this.userProfile = null;
        // Preserve the searchedUsername when an error occurs
        // This prevents the input field from being cleared
        // You may also want to display a relevant error message to the user
        // For simplicity, I'm logging the error here
        console.error('Error occurred while searching for user:', error);
      }
    },

    async toggleFollow() {
      try {
        if (this.userProfile.AlreadyFollowed) {
          await this.unfollowUser();
        } else {
          await this.followUser();
        }
        // Update AlreadyFollowed based on the action performed
        this.userProfile.AlreadyFollowed = !this.userProfile.AlreadyFollowed;
      } catch (error) {
        console.error(error);
      }
    },
    async followUser() {
      try {
        console.log(this.searchedUsername)
        const response = await instance.post(`/users/${this.username}/following/`, JSON.stringify({ username: this.searchedUsername }), {
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        });
        this.userProfile = response.data;
        this.userProfile.AlreadyFollowed = true;
      } catch (error) {
        console.error(error);
      }
    },

    async unfollowUser() {
      try {
        await instance.delete(`/users/${this.username}/following/${this.searchedUsername}`, {
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        });
        // Update AlreadyFollowed directly
        this.userProfile.AlreadyFollowed = false;
      } catch (error) {
        console.error(error);
      }
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

.button-row {
  display: flex;
  gap: 10px;
}
</style>