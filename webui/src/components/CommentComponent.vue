<template>
  <div class="popup-container" v-if="showPopup">
    <div v-if="comments.length > 0">
      <div class="comment" v-for="comment in comments" :key="comment.commentId">
        <p><strong>{{ comment.commenter }}</strong></p>
        <p>{{ comment.content }}</p>
        <button v-if="comment.commenter === currentUser" @click="deleteComment(comment.commentId)">Delete</button>
      </div>
    </div>
    <div v-else>
      <p>No comments yet.</p>
    </div>
    <div class="add-comment">
      <textarea v-model="newComment" placeholder="Write a comment..."></textarea>
      <button @click="postComment">Post Comment</button>
    </div>
  </div>
</template>
  
<script>
import api from "@/services/axios"; 

export default {
  props: {
    photoId: {
      type: String,
      required: true
    },
    showPopup: {
      type: Boolean,
      required: true
    }
  },
  data() {
    return {
      comments: [],
      currentUser: localStorage.getItem('username'),
      newComment: '' // Data property for the new comment text
    };
  },
  mounted() {
    this.fetchComments();
  },
  methods: {
    async fetchComments() {
      try {
        const response = await api.get(`photos/${this.photoId}/comments/`,{headers: {
                        Authorization: localStorage.getItem("token")}
                    });
        this.comments = Array.isArray(response.data) ? response.data : [];
      } catch (error) {
        console.error('Error fetching comments:', error);
      }
    },
    async deleteComment(commentId) {
      try {
        await api.delete(`/comments/${commentId}`,{headers: {
                        Authorization: localStorage.getItem("token")}
                    });
        this.comments = this.comments.filter(comment => comment.commentId !== commentId);
      } catch (error) {
        console.error('Error deleting comment:', error);
      }
    },
    async postComment() {
      try {
        if (this.newComment.trim() === '') {
          alert('Please enter a comment.');
          return;
        }
        const payload = {
          content: this.newComment,
          // Include other necessary fields if required
        };
        await api.post(`photos/${this.photoId}/comments/`,payload,{headers: {
                        Authorization: localStorage.getItem("token")}
                    });
        this.newComment = ''; // Reset the text area
        this.fetchComments(); // Refresh comments list
      } catch (error) {
        console.error('Error posting comment:', error);
      }
    }
  }
};
</script>
  
  <style>

  .add-comment textarea {
  width: 100%; /* Adjust as needed */
  height: 100px; /* Adjust as needed */
}

  </style>
  