<template>
  <div class="clip-creator">
    <h3>Create Clip</h3>
    
    <div class="clip-form">
      <div class="form-group">
        <label for="clipTitle">Clip Title</label>
        <input 
          type="text"
          id="clipTitle"
          v-model="clipTitle"
          placeholder="Enter a title for your clip"
        />
      </div>
      
      <div class="form-group">
        <label>Clip Duration (max 5 minutes)</label>
        <div class="duration-controls">
          <input
            type="range"
            v-model.number="clipDuration"
            min="1"
            max="300"
            step="1"
            @input="updateDuration"
          />
          <span>{{ formatDuration(clipDuration) }}</span>
        </div>
      </div>
      
      <div class="form-group">
        <label>Start Time</label>
        <div class="time-controls">
          <button @click="setStartTimeToCurrentTime">Set to Current Time</button>
          <span>{{ formatDuration(startTime) }}</span>
        </div>
      </div>
      
      <div class="form-group">
        <label>Save Location</label>
        <select v-model="storageOption" @change="updateStorageOption">
          <option value="videos_dir">Videos/fansly-clips</option>
          <option value="source_video_dir">Same folder as video</option>
          <option value="custom_dir">Custom location</option>
        </select>
        <div v-if="storageOption === 'custom_dir'" class="custom-dir-controls">
          <button @click="selectCustomDir">Select Folder</button>
        </div>
        <div class="current-dir">
          <small>Clips will be saved to: {{ currentClipsDir }}</small>
        </div>
      </div>
      
      <div class="clip-preview">
        <div class="time-range">
          <div class="start-time">{{ formatDuration(startTime) }}</div>
          <div class="end-time">{{ formatDuration(startTime + clipDuration) }}</div>
        </div>
        <div class="preview-bar">
          <div class="preview-indicator" :style="previewStyle"></div>
        </div>
      </div>
      
      <div class="clip-actions">
        <button 
          @click="createClip"
          :disabled="isCreatingClip"
          class="create-clip-btn"
        >
          {{ isCreatingClip ? 'Creating...' : 'Create Clip' }}
        </button>
      </div>
    </div>
    
    <div v-if="clipResult" class="clip-result">
      <div v-if="clipResult.success" class="success-message">
        <p>Clip created successfully!</p>
        <button @click="openClipsFolder">Open Clips Folder</button>
      </div>
      <div v-else class="error-message">
        <p>Failed to create clip: {{ clipResult.errorMessage }}</p>
      </div>
    </div>
    
    <div class="saved-clips" v-if="savedClips && savedClips.length > 0">
      <h4>Saved Clips</h4>
      <ul class="clips-list">
        <li v-for="(clip, index) in savedClips" :key="index" class="clip-item">
          <span>{{ getClipName(clip) }}</span>
          <button @click="loadClip(clip)">Play</button>
        </li>
      </ul>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, computed, onMounted } from 'vue';
import { 
  CreateClip, 
  GetClips, 
  OpenClipsFolder, 
  LoadVideoFromPath, 
  GetCurrentClipsDir,
  SetClipStorageOption,
  BrowseForFolder
} from '../../wailsjs/go/main/App';

// Define the ClipResult interface to match the Go struct
interface ClipResult {
  success: boolean;
  filePath: string;
  errorMessage: string;
}

export default defineComponent({
  name: 'ClipCreator',
  props: {
    currentTime: {
      type: Number,
      default: 0
    },
    videoDuration: {
      type: Number,
      default: 0
    }
  },
  emits: ['clip-loaded'],
  setup(props, { emit }) {
    const clipTitle = ref('');
    const startTime = ref(0);
    const clipDuration = ref(30); // Default to 30 seconds
    const isCreatingClip = ref(false);
    const clipResult = ref<ClipResult | null>(null);
    const savedClips = ref<string[]>([]);
    const storageOption = ref<string>('videos_dir'); // Default to Videos directory
    const currentClipsDir = ref<string>('');
    const customDir = ref<string>('');
    
    // Set start time to current video position
    const setStartTimeToCurrentTime = () => {
      startTime.value = props.currentTime;
      
      // Adjust clip duration if it would exceed video length
      if (props.videoDuration > 0 && startTime.value + clipDuration.value > props.videoDuration) {
        clipDuration.value = Math.floor(props.videoDuration - startTime.value);
      }
    };
    
    // Format seconds to MM:SS format
    const formatDuration = (seconds: number) => {
      const mins = Math.floor(seconds / 60);
      const secs = Math.floor(seconds % 60);
      return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
    };
    
    // Explicitly update duration (for reactivity)
    const updateDuration = (event: Event) => {
      const target = event.target as HTMLInputElement;
      clipDuration.value = parseInt(target.value, 10);
      // Force UI update
      setTimeout(() => {}, 0);
    };
    
    // Update the storage option
    const updateStorageOption = async () => {
      try {
        await SetClipStorageOption(storageOption.value, customDir.value);
        await refreshCurrentClipsDir();
      } catch (error) {
        console.error('Failed to update storage option:', error);
      }
    };
    
    // Select a custom directory for saving clips
    const selectCustomDir = async () => {
      try {
        // Use the Go method instead of runtime directly
        const selectedDir = await BrowseForFolder("Select Directory for Clips");
        if (selectedDir) {
          customDir.value = selectedDir;
          await SetClipStorageOption('custom_dir', selectedDir);
          await refreshCurrentClipsDir();
        }
      } catch (error) {
        console.error('Failed to select directory:', error);
      }
    };
    
    // Refresh the current clips directory display
    const refreshCurrentClipsDir = async () => {
      try {
        currentClipsDir.value = await GetCurrentClipsDir();
      } catch (error) {
        console.error('Failed to get current clips directory:', error);
        currentClipsDir.value = 'Unknown';
      }
    };
    
    // Computed style for preview indicator
    const previewStyle = computed(() => {
      if (props.videoDuration <= 0) return { width: '0%' };
      
      const startPercent = (startTime.value / props.videoDuration) * 100;
      const durationPercent = (clipDuration.value / props.videoDuration) * 100;
      
      return {
        left: `${startPercent}%`,
        width: `${durationPercent}%`
      };
    });
    
    // Create the clip
    const createClip = async () => {
      if (isCreatingClip.value) return;
      
      isCreatingClip.value = true;
      clipResult.value = null;
      
      try {
        const result = await CreateClip(startTime.value, clipDuration.value, clipTitle.value);
        clipResult.value = result as ClipResult;
        
        if (result.success) {
          // Reset form and refresh clips list
          clipTitle.value = '';
          await loadSavedClips();
        }
      } catch (error: unknown) {
        clipResult.value = {
          success: false,
          filePath: '',
          errorMessage: `Error: ${error instanceof Error ? error.message : String(error)}`
        };
      } finally {
        isCreatingClip.value = false;
      }
    };
    
    // Load saved clips
    const loadSavedClips = async () => {
      try {
        const clips = await GetClips();
        savedClips.value = clips || [];
      } catch (error) {
        console.error('Failed to load saved clips:', error);
        savedClips.value = [];
      }
    };
    
    // Open clips folder
    const openClipsFolder = async () => {
      try {
        await OpenClipsFolder();
      } catch (error) {
        console.error('Failed to open clips folder:', error);
      }
    };
    
    // Get clip name from path
    const getClipName = (clipPath: string) => {
      if (!clipPath) return 'Unknown';
      const parts = clipPath.split(/[\/\\]/);
      return parts[parts.length - 1];
    };
    
    // Load a clip for playback
    const loadClip = async (clipPath: string) => {
      try {
        const videoPath = await LoadVideoFromPath(clipPath);
        if (videoPath) {
          emit('clip-loaded', videoPath);
        }
      } catch (error) {
        console.error('Failed to load clip:', error);
      }
    };
    
    // Load saved clips and get current clips directory on component mount
    onMounted(async () => {
      await loadSavedClips();
      await refreshCurrentClipsDir();
    });
    
    return {
      clipTitle,
      startTime,
      clipDuration,
      isCreatingClip,
      clipResult,
      savedClips,
      storageOption,
      currentClipsDir,
      setStartTimeToCurrentTime,
      formatDuration,
      updateDuration,
      previewStyle,
      createClip,
      openClipsFolder,
      getClipName,
      loadClip,
      updateStorageOption,
      selectCustomDir
    };
  }
});
</script>

<style scoped>
.clip-creator {
  padding: 15px;
  color: #cdd6f4; /* Catppuccin Mocha text */
}

h3 {
  margin-bottom: 15px;
  color: #89b4fa; /* Catppuccin Mocha blue */
}

.clip-form {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

label {
  font-size: 0.9rem;
  color: #bac2de; /* Catppuccin Mocha subtext1 */
}

input[type="text"], select {
  padding: 8px;
  border-radius: 4px;
  border: 1px solid #313244; /* Catppuccin Mocha surface0 */
  background-color: #1e1e2e; /* Catppuccin Mocha background */
  color: #cdd6f4; /* Catppuccin Mocha text */
}

input[type="range"] {
  width: 100%;
  accent-color: #89b4fa; /* Catppuccin Mocha blue */
}

.duration-controls, .time-controls, .custom-dir-controls {
  display: flex;
  align-items: center;
  gap: 10px;
}

.current-dir {
  margin-top: 5px;
  font-size: 0.8rem;
  color: #a6adc8; /* Catppuccin Mocha overlay0 */
  word-break: break-all;
}

button {
  padding: 5px 10px;
  font-size: 0.8rem;
  background-color: #313244; /* Catppuccin Mocha surface0 */
  color: #cdd6f4; /* Catppuccin Mocha text */
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

button:hover:not(:disabled) {
  background-color: #45475a; /* Catppuccin Mocha surface1 */
}

.clip-preview {
  margin: 15px 0;
  padding: 10px;
  background-color: #181825; /* Catppuccin Mocha mantle */
  border-radius: 4px;
}

.time-range {
  display: flex;
  justify-content: space-between;
  margin-bottom: 5px;
  font-size: 0.8rem;
}

.preview-bar {
  position: relative;
  height: 10px;
  background-color: #313244; /* Catppuccin Mocha surface0 */
  border-radius: 5px;
  overflow: hidden;
}

.preview-indicator {
  position: absolute;
  height: 100%;
  background-color: #89b4fa; /* Catppuccin Mocha blue */
  border-radius: 5px;
}

.clip-actions {
  display: flex;
  justify-content: center;
  margin-top: 10px;
}

.create-clip-btn {
  padding: 10px 20px;
  background-color: #89b4fa; /* Catppuccin Mocha blue */
  color: #1e1e2e; /* Catppuccin Mocha background */
  font-weight: bold;
  border-radius: 4px;
  border: none;
  cursor: pointer;
  transition: background-color 0.2s;
}

.create-clip-btn:hover:not(:disabled) {
  background-color: #b4befe; /* Catppuccin Mocha lavender */
}

.create-clip-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.clip-result {
  margin-top: 15px;
  padding: 10px;
  border-radius: 4px;
}

.success-message {
  background-color: rgba(166, 227, 161, 0.2); /* Catppuccin Mocha green with opacity */
  color: #a6e3a1; /* Catppuccin Mocha green */
  padding: 10px;
  border-radius: 4px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.error-message {
  background-color: rgba(243, 139, 168, 0.2); /* Catppuccin Mocha red with opacity */
  color: #f38ba8; /* Catppuccin Mocha red */
  padding: 10px;
  border-radius: 4px;
}

.saved-clips {
  margin-top: 20px;
}

h4 {
  margin-bottom: 10px;
  color: #89b4fa; /* Catppuccin Mocha blue */
}

.clips-list {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 5px;
  max-height: 200px;
  overflow-y: auto;
}

.clip-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px;
  background-color: #181825; /* Catppuccin Mocha mantle */
  border-radius: 4px;
}

.clip-item span {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 70%;
}

.clip-item button {
  padding: 4px 8px;
  font-size: 0.8rem;
  background-color: #89b4fa; /* Catppuccin Mocha blue */
  color: #1e1e2e; /* Catppuccin Mocha background */
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.clip-item button:hover {
  background-color: #b4befe; /* Catppuccin Mocha lavender */
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .clip-creator {
    padding: 10px;
  }
  
  .clip-actions {
    flex-direction: column;
    gap: 10px;
  }
  
  .create-clip-btn {
    width: 100%;
  }
}
</style>

