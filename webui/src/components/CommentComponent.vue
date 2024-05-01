<template>
  <div class="comment-component">
    <div class="comment-list">
      <div v-for="(comment, index) in comments" :key="index" class="comment-item">
        {{ comment.Username }}: {{ comment.CommentContent }}
        <button
          v-if="comment.Username === getLocalStorageUsername()"
          @click="deleteComment(comment)"
          class="delete-button"
        >
          Delete
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import instance from '../services/axios.js';
export default {
  props: {
    comments: {
      type: Array,
      required: true,
    },
  },
  methods: {
    getLocalStorageUsername() {
      return localStorage.getItem("username");
    },
    deleteComment(comment) {
      this.$emit("delete-comment", comment);
    },
  },
};
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

.delete-button {
  background-color: #cd1414;
  color: white;
  border: none;
  border-radius: 0 5px 5px 0;
  font-size: 11px;
  cursor: pointer;
  }
</style>