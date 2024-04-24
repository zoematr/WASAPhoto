<template>
  <div class="photo-container">
    <div v-for="photo in sortedPhotos" :key="photo.PhotoId" class="photo-item">
      <img :src="getPhotoDataURL(photo.PhotoFile)" alt="Photo" class="photo">
      <div class="photo-info">
        <div class="info-line">
          <p class="username">{{ photo.Username }}</p>
          <p class="date">{{ formatDate(photo.Date) }}</p>
        </div>
        <div class="comment-input">
          <input v-model="commentText" type="text" placeholder="Enter your comment">
          <button class="comment-button" @click="commentPhoto(photo, commentText)">Comment</button>
        </div>
        <comment-component :comments="photo.Comments" ref="commentComponent"></comment-component>
      </div>
    </div>
  </div>
</template>

<script>
import instance from '../services/axios.js';
import CommentComponent from './CommentComponent.vue';

export default {
  components: {
    CommentComponent
  },
  props: {
    photos: {
      type: Array,
      required: true
    }
  },
  data() {
    return {
      commentText: '',
      storedUsername: localStorage.getItem("username") 
    };
  },
  computed: {
    sortedPhotos() {
      // sort photos in reverse chronological order
      return this.photos.slice().sort((a, b) => new Date(b.Date) - new Date(a.Date));
    }
  },
  methods: {
    getPhotoDataURL(photoFile) {
      try {
        if (!photoFile) {
          console.error('Photo file is null or undefined:', photoFile);
          return '';
        }

        const binaryData = atob(photoFile);
        const byteArray = new Uint8Array(binaryData.length);
        for (let i = 0; i < binaryData.length; i++) {
          byteArray[i] = binaryData.charCodeAt(i);
        }
        const blob = new Blob([byteArray], { type: 'image/jpeg' });

        // Create blob URL for the image
        const imageUrl = URL.createObjectURL(blob);

        // Check if imageUrl is not null or empty
        if (!imageUrl) {
          console.error('Generated image URL is null or empty:', imageUrl);
          return ''; // Return empty string if imageUrl is not valid
        }

        return imageUrl;
      } catch (error) {
        console.error('Error converting photo data to URL:', error);
        return ''; // Return empty string if an error occurs
      }
    },
    formatDate(dateString) {
      const date = new Date(dateString);
      const options = {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        hour12: false 
      };
      return date.toLocaleString(undefined, options);
    },

    async isValidComment(commenttext) {
      return commenttext.length > 1 && commenttext.length <= 400;
    },

    async commentPhoto(photo, commentText) {
      if (this.isValidComment(commentText)) {
        console.log(this.commentText)
        console.log(photo.Username)
        console.log(photo.PhotoId)
        console.log("this is the token", localStorage.getItem('token'))
        try {
          const response = await instance.post(`/users/${photo.Username}/photos/${photo.PhotoId}/comments/`, commentText ,  {
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${localStorage.getItem('token')}`,
            }
          });
          // const newComment = response.data;
          // photo.Comments.push(newComment);
          alert('You commented the photo!');
          // location.reload();
          this.commentText = '';
        } catch (error) {
          console.error('Error commenting the photo:', error);
        }
      } else {
        alert('Oops! The comment has to be between 1 and 400 characters long.');
      }
    }
  }
}
</script>

<style scoped>
.photo-container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); /* Adjust the minmax width to your preference */
  gap: 10px;
  overflow-y: auto; /* Enable vertical scrolling */
  max-height: 500px; /* Set a maximum height to limit the container size */
}

.photo-item {
  display: flex;
  flex-direction: column;
}

.photo {
  width: 100%;
  height: auto;
  object-fit: cover;
}

.photo-info {
  background-color: rgba(42, 41, 41, 0.5);
  color: white;
  padding: 5px;
}

.info-line {
  display: flex;
  justify-content: space-between;
}

.username,
.date {
  margin: 0;
}
.comment-section {
  margin-top: 10px;
}
</style>

