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
import api from "@/services/axios"; 

export default {
  components: {
    ImageComponent
  },
  data() {
    return {
      images: [] // This array will hold the processed photo objects
    }
  },
  mounted() {
    this.fetchImages();
  },
  methods: {
    async fetchImages() {
      try {
        
        const response = await api.get('/stream',{headers: {
                        Authorization: localStorage.getItem("token")}
                    }); // Replace with the full API URL if necessary
                    this.images = response.data && Array.isArray(response.data)
        ? response.data.map(photo => ({
            ...photo,
            ImageData: `data:image/jpeg;base64,${photo.ImageData}`
          }))
        : [];
    } catch (error) {
      console.error('Error fetching images:', error);
    }
  },
    arrayBufferToBase64(buffer) {
      let binary = '';
      let bytes = new Uint8Array(buffer);
      let len = bytes.byteLength;
      for (let i = 0; i < len; i++) {
        binary += String.fromCharCode(bytes[i]);
      }
      return window.btoa(binary);
    }
  }
}
</script>