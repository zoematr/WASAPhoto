<template>
  <div>
    <div class="banner">
      <p v-if="storedUsername">Welcome lol, {{ storedUsername }}</p>
      <p v-else>Welcome lol, Guest</p>
    </div>
    <div>
      <input v-model="username" placeholder="Enter Username" />
      <button @click="login">Login</button>
    </div>
  </div>
</template>

<script>
</script>

<style scoped>
</style>