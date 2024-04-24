<template>
  <div class="comment-component">
    <div class="comment-list">
      <div v-for="(comment, index) in comments" :key="index" class="comment-item">
        {{ comment.UserAuthorOfComment }}: {{ comment.CommentContent }}
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    comments: {
      type: Array,
      required: true
    }
  }
}
</script>

<style scoped>
.comment-component {
  margin-top: 10px;
}

.comment-input {
  display: flex;
}

.comment-input input {
  flex: 1;
}

.comment-button {
  margin-left: 0px;
}

.comment-list {
  margin-top: 10px;
  max-height: 200px; /* Set a maximum height to limit the container size */
  overflow-y: auto; /* Enable vertical scrolling */
}

.comment-item {
  margin-bottom: 5px;
}
</style>