<template>
  <div v-if="images.length">
    <ImageComponent
      v-for="image in images"
      :key="image.photoId"
      :photoDetails="image"
    />
  </div>
</template>

<script>
import ImageComponent from '@/components/ImageComponent.vue'; 
import axios from 'axios';

export default {
  components: {
    ImageComponent
  },
  data() {
    return {
      images: []
    }
  },
  mounted() {
    this.fetchImages();
  },
  methods: {
    async fetchImages() {
      try {
        // Modify the URL to match your backend API endpoint
        const username = 'replace_with_username'; // Replace with the actual username
        const response = await axios.get(`/users/${username}/mystream`, {
          headers: {
            Authorization: localStorage.getItem("token")
          },
          // Ensure responseType is set to 'arraybuffer' to receive binary data
          responseType: 'arraybuffer'
        });

        // Process binary image data into a format suitable for displaying
        this.images = response.data.map(photo => ({
          ...photo,
          ImageData: this.arrayBufferToBase64(photo.ImageData)
        }));
      } catch (error) {
        console.error('Error fetching images:', error);
      }
    },
    arrayBufferToBase64(buffer) {
      // Convert binary data to base64
      let binary = '';
      let bytes = new Uint8Array(buffer);
      let len = bytes.byteLength;
      for (let i = 0; i < len; i++) {
        binary += String.fromCharCode(bytes[i]);
      }
      return 'data:image/jpeg;base64,' + window.btoa(binary);
    }
  }
}
</script>
