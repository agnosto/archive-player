<template>
  <div class="video-player-container" :class="{ 'theater-mode': theaterMode }">
    <div class="video-wrapper">
      <video
        ref="videoRef"
        controls
        @timeupdate="onTimeUpdate"
        @loadedmetadata="onVideoLoaded"
        v-if="videoSrc"
      >
        <source :src="videoSrc" :type="videoType" />
        Your browser does not support the video tag.
      </video>
      <div v-else class="no-video">
        <p>No video loaded. Click "Open Video" to select a video file.</p>
      </div>
    </div>
    
    <div v-if="showChat && currentMessages.length > 0" 
         class="chat-overlay"
         :class="{ 'chat-left': theme.chatPosition === 'left', 'chat-right': theme.chatPosition === 'right' }">
      <ChatOverlay
        :messages="currentMessages"
        :theme="theme"
      />
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted, watch, computed } from 'vue';
import { GetMessagesAtTime } from '../../wailsjs/go/main/App';
import { ChatMessage, ThemeSettings } from '../types';
import ChatOverlay from './ChatOverlay.vue';

export default defineComponent({
  name: 'VideoPlayer',
  components: {
    ChatOverlay
  },
  props: {
    videoSrc: {
      type: String,
      default: ''
    },
    theme: {
      type: Object as () => ThemeSettings,
      required: true
    },
    showChat: {
      type: Boolean,
      default: true
    },
    theaterMode: {
      type: Boolean,
      default: false
    }
  },
  emits: ['update:currentTime'],
  setup(props, { emit }) {
    const videoRef = ref<HTMLVideoElement | null>(null);
    const currentTime = ref(0);
    const currentMessages = ref<ChatMessage[]>([]);
    const chatWindowSize = ref(30); // Show messages for 30 seconds
    
    // Compute video type based on file extension
    const videoType = computed(() => {
      if (!props.videoSrc) return 'video/mp4';
      
      const extension = props.videoSrc.split('.').pop()?.toLowerCase();
      switch (extension) {
        case 'mp4': return 'video/mp4';
        case 'webm': return 'video/webm';
        case 'mkv': return 'video/x-matroska';
        case 'avi': return 'video/x-msvideo';
        default: return 'video/mp4';
      }
    });

    const onTimeUpdate = async () => {
      if (videoRef.value) {
        currentTime.value = videoRef.value.currentTime;
        emit('update:currentTime', currentTime.value);
        
        try {
          const messages = await GetMessagesAtTime(currentTime.value, chatWindowSize.value);
          currentMessages.value = messages || [];
        } catch (error) {
          console.error('Error getting messages:', error);
        }
      }
    };

    const onVideoLoaded = () => {
      if (videoRef.value) {
        console.log('Video loaded, duration:', videoRef.value.duration);
      }
    };

    watch(() => props.videoSrc, (newSrc) => {
      if (newSrc && videoRef.value) {
        // Reset current time and messages when video changes
        currentTime.value = 0;
        currentMessages.value = [];
        console.log('Video source changed to:', newSrc);
      }
    });

    return {
      videoRef,
      currentTime,
      currentMessages,
      videoType,
      onTimeUpdate,
      onVideoLoaded
    };
  }
});
</script>

<style scoped>
.video-player-container {
  width: 100%;
  height: 100%;
  position: relative;
  background-color: #000;
  overflow: hidden;
}

.video-wrapper {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

video {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.no-video {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
  color: #cdd6f4;
  text-align: center;
  padding: 20px;
}

.chat-overlay {
  position: absolute;
  bottom: 10px;
  max-height: 80%;
  overflow-y: auto;
  z-index: 10;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.chat-left {
  left: 10px;
}

.chat-right {
  right: 10px;
}

/* Theater mode adjustments */
.theater-mode video {
  object-fit: cover;
}

/* Mobile adjustments */
@media (max-width: 768px) {
  .chat-overlay {
    max-height: 50%;
    width: 90% !important;
    left: 5% !important;
    right: 5% !important;
  }
}
</style>
