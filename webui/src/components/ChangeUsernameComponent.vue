<template>
  <div>
    <h1>Change Username here</h1>
    <form @submit.prevent="changeUsername">
      <div class="form-group">
        <label for="newUsername">New Username</label>
        <input v-model="newUsername" type="text" class="form-control" id="newUsername" placeholder="New username" required>
      </div>
      <button type="submit" class="btn btn-primary">Change YOUR Username</button>
    </form>
  </div>
</template>

<script>
</script>

<style scoped>
</style>