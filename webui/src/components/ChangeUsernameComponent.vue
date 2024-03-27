<!-- ChangeUsername.vue -->
<template>
  <div>
    <h1>Change Username</h1>
    <form @submit.prevent="changeUsername">
      <div class="form-group">
        <label for="newUsername">New Username</label>
        <input v-model="newUsername" type="text" class="form-control" id="newUsername" placeholder="New username" required>
      </div>
      <button type="submit" class="btn btn-primary">Change Username</button>
    </form>
  </div>
</template>

<script>
import { ref, computed } from 'vue';
import instance from '../services/axios.js';

export default {
  setup() {
    const newUsername = ref('');
    const usernameComputed = computed(() => localStorage.getItem('username') || '');

    async function changeUsername() {
      try {
        const response = await instance.patch(`/users/${usernameComputed}`, { newusername: newUsername.value }, {
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

    return {
      newUsername,
      changeUsername,
    };
  },
};
</script>

<style scoped>
</style>
