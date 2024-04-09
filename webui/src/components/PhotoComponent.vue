<template>
  <div class="photo-container">
    <div v-for="photo in photos" :key="photo.PhotoId" class="photo-item">
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
  display: flex;
  flex-wrap: wrap;
}

.photo-item {
  margin: 10px;
}

.photo {
  width: 200px;
  height: 200px;
  object-fit: cover;
}

.photo-info {
  background-color: rgba(50, 48, 50, 0.5);
  color: white;
  padding: 5px;
}
</style>
