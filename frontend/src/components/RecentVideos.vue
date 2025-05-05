<template>
  <div class="recent-videos">
    <h2>Recent Videos</h2>
    <div class="videos-grid">
      <div 
        v-for="(video, index) in recentVideos" 
        :key="index" 
        class="video-item"
        @click="selectVideo(video)"
      >
        <div class="thumbnail">
          <img 
            :src="getThumbnailUrl(video)" 
            :alt="video.name"
            @error="handleImageError"
          />
        </div>
        <div class="video-info">
          <div class="video-name">{{ video.name }}</div>
          <div class="video-date">{{ formatDate(video.lastPlayed) }}</div>
        </div>
      </div>
      
      <div v-if="recentVideos.length === 0" class="no-videos">
        No recent videos found
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted, watch } from 'vue';
import { RecentVideo } from '../types';

export default defineComponent({
  name: 'RecentVideos',
  props: {
    recentVideos: {
      type: Array as () => RecentVideo[],
      required: true
    }
  },
  emits: ['select-video'],
  setup(props, { emit }) {
    const placeholderImage = ref('data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMzIwIiBoZWlnaHQ9IjE4MCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iMzIwIiBoZWlnaHQ9IjE4MCIgZmlsbD0iIzFhMWIyNiIvPjx0ZXh0IHg9IjUwJSIgeT0iNTAlIiBmb250LXNpemU9IjI0IiB0ZXh0LWFuY2hvcj0ibWlkZGxlIiBhbGlnbm1lbnQtYmFzZWxpbmU9Im1pZGRsZSIgZmlsbD0iIzg5YjRmYSI+Tm8gVGh1bWJuYWlsPC90ZXh0Pjwvc3ZnPg==');

    const preloadThumbnails = () => {
      props.recentVideos.forEach(video => {
        if (video.thumbnailPath && video.thumbnailPath.trim() !== '') {
          const img = new Image();
          img.src = `http://localhost:8080/thumbnail/${encodeURIComponent(video.thumbnailPath)}`;
        }
      });
    };
    
    const getThumbnailUrl = (video: RecentVideo) => {
      if (video.thumbnailPath && video.thumbnailPath.trim() !== '') {
        console.log("Loading thumbnail from path:", video.thumbnailPath);
        // For local files, we need to create a URL
        return `http://localhost:8080/thumbnail/${encodeURIComponent(video.thumbnailPath)}`;
      }
      console.log("No thumbnail path for:", video.name);
      return placeholderImage.value;
    };
    
    const handleImageError = (event: Event) => {
      const target = event.target as HTMLImageElement;
      target.src = placeholderImage.value;
    };
    
    const formatDate = (dateString: string) => {
      try {
        const date = new Date(dateString);
        return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
      } catch (e) {
        return 'Unknown date';
      }
    };
    
    const selectVideo = (video: RecentVideo) => {
      emit('select-video', video);
    };

    onMounted(() => {
      preloadThumbnails();
    });
    
    // Also preload when recentVideos changes
    watch(() => props.recentVideos, () => {
      preloadThumbnails();
    });
    
    return {
      getThumbnailUrl,
      handleImageError,
      formatDate,
      selectVideo
    };
  }
});
</script>

<style scoped>
.recent-videos {
  height: 100%;
  display: flex;
  flex-direction: column;
}

h2 {
  margin-bottom: 15px;
  color: #89b4fa; /* Catppuccin Mocha blue */
  font-size: 1.2rem;
}

.videos-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 15px;
  overflow-y: auto;
  padding-right: 5px;
}

.video-item {
  cursor: pointer;
  border-radius: 8px;
  overflow: hidden;
  background-color: #313244; /* Catppuccin Mocha surface0 */
  transition: transform 0.2s, box-shadow 0.2s;
}

.video-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
}

.thumbnail {
  width: 100%;
  height: 90px;
  overflow: hidden;
}

.thumbnail img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.video-info {
  padding: 8px;
}

.video-name {
  font-size: 0.9rem;
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.video-date {
  font-size: 0.7rem;
  color: #a6adc8; /* Catppuccin Mocha subtext0 */
}

.no-videos {
  grid-column: 1 / -1;
  text-align: center;
  padding: 20px;
  color: #a6adc8; /* Catppuccin Mocha subtext0 */
}

/* Mobile adjustments */
@media (max-width: 768px) {
  .videos-grid {
    grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  }
}
</style>
