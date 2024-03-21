<template>
  <div>
    <div class="info-container"></div>
    <input v-model="username" placeholder="Search Username" />
    <button @click="searchUser">Search</button>
    <div v-if="userProfile">
      <p>Followers: {{ userProfile.Followers }}</p>
      <p>Following: {{ userProfile.Following }}</p>
      <p>Posts: {{ userProfile.PhotosCount }}</p>
        <button v-if="userProfile && !isOwnProfile" @click="toggleFollow">
        {{ userProfile.IsFollowing ? 'Unfollow' : 'Follow' }}
      </button>
      <button v-if="userProfile && !isOwnProfile" @click="toggleBan">{{ userProfile.IsBanning ? 'Unban' : 'Ban' }}</button>

      
    <ImageComponent
      v-for="image in images"
      :key="image.photoId"
      :photoDetails="image"
    />
  </div>
  <div v-else-if="searched">
      <p>This user wasn't found.</p>
    </div>


    </div>
      


</template>

<script>
import ImageComponent from '@/components/ImageComponent.vue'; 
import api from "@/services/axios"; 

export default {
  components: {
    ImageComponent
  },

  data() {
    return {
      username: '', // Username to search
      userProfile: null,
      images: [],
      searched: false
    };
  },
  async created() {
    if (this.$route.params.username){
    this.username = this.$route.params.username;
    await this.fetchUserProfile();}
  },
  computed: {
    isOwnProfile() {
      return this.username === localStorage.getItem("username");
    }
  },
  methods: {
    async searchUser() {
      console.log("User searching")

    if (this.username) {
      this.$router.push({ name: 'UserProfile', params: { username: this.username } });
    }
    await this.fetchUserProfile();
  },

    async toggleBan() {
    try {
      const url = `/users/${this.username}/bans/`;
      if (this.userProfile.IsBanning) {
        await api.delete(url, { headers: { Authorization: localStorage.getItem("token") } });
      } else {
        await api.post(url, {}, { headers: { Authorization: localStorage.getItem("token") } });
      }
      this.userProfile.IsBanning = !this.userProfile.IsBanning; // Toggle the banned status
    } catch (error) {
      console.error('Error toggling ban:', error);
    }
  },


    async toggleFollow() {
      if (this.userProfile.IsFollowing) {
        await this.unfollowUser();
        
      } else {
        await this.followUser();
      }
      this.userProfile.IsFollowing= !this.userProfile.IsFollowing
    },
    async followUser() {
      // the API call to follow the user
      try {
        await api.post(`/users/${this.username}/followers/`, {}, {
          headers: { Authorization: localStorage.getItem("token") }
        });

      } catch (error) {
        console.error("Failed to follow user:", error);
      }
    },
    async unfollowUser() {
      // the API call to unfollow the user
      try {
        await api.delete(`/users/${this.username}/followers/`, {
          headers: { Authorization: localStorage.getItem("token") }
        });

      } catch (error) {
        console.error("Failed to unfollow user:", error);
      }
    }
      ,
    
  

    async fetchUserProfile() {
  try {
    const response = await api.get(`/users/${this.username}`, {
      headers: { Authorization: localStorage.getItem("token") }
    });
    if (response.data==null){
      this.searched= true;
      return;
    }
  
    // Check if the response contains 'photos' and it is an array
    this.images = response.data.Photos && Array.isArray(response.data.Photos)
      ? response.data.Photos.map(photo => ({
          ...photo,
          ImageData: this.arrayBufferToBase64(photo.ImageData)
        }))
      : [];

    this.userProfile = {
      ...response.data,
      photos: this.images
    };
  } catch (error) {
    console.error("Failed to fetch user profile:", error);
  }
},
arrayBufferToBase64(buffer) {
      let binary = '';
      let bytes = new Uint8Array(buffer);
      let len = bytes.byteLength;
      for (let i = 0; i < len; i++) {
        binary += String.fromCharCode(bytes[i]);
      }
      return window.btoa(binary);
    }
  }
}


</script>

<style>
.info-container {
  background-color: #d6b343; /* Adjust the color code to get the desired shade of yellow */
}
</style>