<template>
  <div class="comment-component">
    <div class="comment-input">
      <input v-model="commentText" type="text" placeholder="Enter your comment">
      <button class="comment-button" @click="postComment">Comment</button>
    </div>
    <div class="comment-list">
      <div v-for="(comment, index) in comments" :key="index" class="comment-item">
        {{ comment.UserAuthorOfComment }}: {{ comment.CommentContent }}
      </div>
    </div>
  </div>
</template>

<script>
import instance from '../services/axios.js';
export default {
  data() {
    return {
      commentText: ''
    };
  },
  methods: {
    async postComment(photo) {
      if (this.isValidComment(this.commentText)) {
        // call api function
        this.commentPhoto(photo, this.commentText);
        // clear input field
        this.commentText = '';
      } else {
        alert('Oops! The comment has to be between 1 and 400 characters long.');
      }
    },
    async isValidComment(comment) {
      return comment.length > 0 && comment.length <= 400;
    },
    async commentPhoto(photo, commentcontent) {
      try {
        await instance.post(`/users/${photo.Username}/photos/${photo.PhotoId}/comments/`, { commentcontent },  {
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        });
        alert('You commented the photo!');
        location.reload();
      } catch (error) {
          console.error('Error commenting the photo:', error);
      }
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