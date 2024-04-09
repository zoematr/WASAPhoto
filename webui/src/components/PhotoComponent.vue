<template>
  <div class="photo-container">
    <div v-for="photo in sortedPhotos" :key="photo.PhotoId" class="photo-item">
      <img :src="getPhotoDataURL(photo.PhotoFile)" alt="Photo" class="photo">
      <div class="photo-info">
        <p>{{ photo.Username }}</p>
        <p>{{ formatDate(photo.Date) }}</p>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    photos: {
      type: Array,
      required: true
    }
  },
  computed: {
    sortedPhotos() {
      // Sort photos by Date in reverse chronological order (so i have and instagram kind of thing)
      return this.photos.slice().sort((a, b) => new Date(b.Date) - new Date(a.Date));
    }
  },
  methods: {
    getPhotoDataURL(photoFile) {
      try {
        // Check if photoFile is not null or undefined
        if (!photoFile) {
          console.error('Photo file is null or undefined:', photoFile);
          return ''; // Return empty string if photoFile is not valid
        }

        // Convert Base64 string to binary data
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
        hour12: false // Use 24-hour format
      };
      return date.toLocaleString(undefined, options);
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
  background-color: rgba(54, 51, 53, 0.5);
  color: white;
  padding: 5px;
}
</style>