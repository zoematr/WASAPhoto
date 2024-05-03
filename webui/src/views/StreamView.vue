<template>
    <div class="container">
        <div v-if="storedUsername">
            <p>Hello {{ storedUsername }}, this is your stream!</p>
            <PhotoComponent :photos="userStream" />
        </div>
        <div v-else>
            <p>Log in to see other users' photos!</p>
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
            storedUsername: localStorage.getItem("username"),
            userStream: null,
            photos: []
        }
    },
    created() {
        this.getStream();
    },
    methods: {
        async getStream() {
            try {
                console.log("getStream is being called")
                const response = await instance.get(`/users/${localStorage.getItem('username')}/mystream/`, {
                    headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                    }
            });
            this.userStream = response.data;
            console.log("this is the response data", response.data)
            } catch (error) {
            console.error(error);
            console.error('Error occurred while retrieving the stream', error);
            }
        },
    }
};
</script>

<style scoped>
.container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: 10px;
  max-height: 500px; /* max height to limit the container size */
}
</style>