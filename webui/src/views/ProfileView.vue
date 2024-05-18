<template>
  <div class="container">
    <div class="search-bar-container">
      <div class="search-bar">
        <input v-model="searchedUsername" placeholder="Search Username" @keyup.enter="searchUser">
        <button @click="searchUser">Search</button>
      </div>
      <div v-if="userProfile" class="user-info">
        <p class="username">{{ userProfile.Username }}</p>
        <p class="follower-count">Followers: {{ userProfile.Followers ? userProfile.Followers.length : 0 }}</p>
        <p class="following-count">Following: {{ userProfile.Following ? userProfile.Following.length : 0 }}</p>
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
      <PhotoComponent :photos="userProfile.Photos" />
    </div>
    <div v-else-if="searched && !userProfile">
      <p>User not found</p>
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
        this.userProfile = null
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
        // Call searchUser again to update userProfile with the latest data
        await this.searchUser();
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
        this.userProfile.AlreadyFollowed = false;
      } catch (error) {
        console.error(error);
      }
    },

    async toggleBan() {
      try {
        if (this.userProfile.AlreadyBanned) {
          await this.unbanUser();
        } else {
          await this.banUser();
        }
        // Call searchUser again to update userProfile with the latest data
        await this.searchUser();
      } catch (error) {
        console.error(error);
      }
    },
    async banUser() {
      try {
        console.log(this.searchedUsername)
        const response = await instance.post(`/users/${this.username}/banned/`, JSON.stringify({ username: this.searchedUsername }), {
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        });
        this.userProfile = response.data;
        this.userProfile.AlreadyBanned = true;
      } catch (error) {
        console.error(error);
      }
    },

    async unbanUser() {
      try {
        await instance.delete(`/users/${this.username}/banned/${this.searchedUsername}`, {
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        });
        this.userProfile.AlreadyBanned = false;
      } catch (error) {
        console.error(error);
      }
    }


  },
  created() {
    if (this.searchedUsername) {
      this.searchUser();
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