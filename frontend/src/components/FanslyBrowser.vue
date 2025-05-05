<template>
  <div class="fansly-browser">
    <div class="browser-header">
      <h2>Fansly Streams</h2>
      <div class="filter-controls">
        <input 
          type="text"
          v-model="searchQuery"
          placeholder="Search by model name..."
          class="search-input"
        />
        <select v-model="selectedModel" class="model-filter">
          <option value="">All Models</option>
          <option v-for="model in uniqueModels" :key="model" :value="model">
            {{ model }}
          </option>
        </select>
        <button @click="refreshContent" class="refresh-btn">
          Refresh
        </button>
      </div>
    </div>
    
    <div class="streams-container">
      <div v-if="loading" class="loading-indicator">
        <div class="spinner"></div>
        <span>Loading content...</span>
      </div>
      <div v-else-if="error" class="error-message">
        <p>{{ error }}</p>
        <button @click="setupFanslyIntegration" class="setup-btn">
          Setup Fansly Integration
        </button>
      </div>
      <div v-else-if="filteredStreams.length === 0" class="no-content">
        <p>No streams found. Try a different search or model filter.</p>
      </div>
      <div v-else class="streams-grid">
        <div 
          v-for="stream in filteredStreams"
          :key="stream.hash"
          class="stream-card"
          @click="selectStream(stream)"
        >
          <div class="thumbnail-container">
            <img 
              v-if="stream.contactSheet"
              :src="getThumbnailUrl(stream.contactSheet)"
              :alt="`${stream.model} stream thumbnail`"
              class="stream-thumbnail"
              @error="handleImageError"
            />
            <div v-else class="no-thumbnail">
              <span>No Preview</span>
            </div>
            <div class="stream-info-overlay">
              <span class="stream-duration" v-if="stream.duration">{{ formatDuration(stream.duration) }}</span>
              <span class="chat-badge" v-if="stream.hasChat">Chat Available</span>
            </div>
          </div>
          <div class="stream-info">
            <h3 class="stream-title">{{ formatStreamTitle(stream.path) }}</h3>
            <div class="stream-details">
              <span class="stream-model">{{ stream.model }}</span>
              <span class="stream-date">{{ formatDate(extractDateFromPath(stream.path)) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="config-panel" v-if="showConfigPanel">
      <div class="config-panel-content">
        <h3>Fansly Integration Setup</h3>
        
        <div class="config-field">
          <label for="configPath">Config Path:</label>
          <div class="path-input-group">
            <input type="text" id="configPath" v-model="configPath" readonly />
            <button @click="browseConfigPath">Browse</button>
          </div>
          <small>Default locations: ~/.config/fansly-scraper/config.toml (Mac/Linux) or %APPDATA%\fansly-scraper\config.toml (Windows)</small>
        </div>
        
        <div class="config-field">
          <label for="dbPath">Database Path:</label>
          <div class="path-input-group">
            <input type="text" id="dbPath" v-model="dbPath" readonly />
            <button @click="browseDbPath">Browse</button>
          </div>
          <small>Usually located in the save_location specified in your config</small>
        </div>
        
        <div class="config-actions">
          <button @click="saveConfig" class="save-btn">Save Configuration</button>
          <button @click="cancelConfig" class="cancel-btn">Cancel</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, computed, onMounted, watch } from 'vue';
import { 
  GetFanslyStreams, 
  LoadFanslyStream, 
  GetFanslyConfig, 
  SaveFanslyConfig,
  BrowseForFile,
  BrowseForFolder
} from '../../wailsjs/go/main/App';

interface FanslyStream {
  model: string;
  hash: string;
  path: string;
  file_type: string;
  contactSheet?: string;
  hasChat: boolean;
  duration?: number;
}

export default defineComponent({
  name: 'FanslyBrowser',
  emits: ['stream-selected'],
  
  setup(props, { emit }) {
    const streams = ref<FanslyStream[]>([]);
    const loading = ref(true);
    const error = ref('');
    const searchQuery = ref('');
    const selectedModel = ref('');
    const showConfigPanel = ref(false);
    const configPath = ref('');
    const dbPath = ref('');

    const placeholderImage = ref('data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMzIwIiBoZWlnaHQ9IjE4MCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iMzIwIiBoZWlnaHQ9IjE4MCIgZmlsbD0iIzFhMWIyNiIvPjx0ZXh0IHg9IjUwJSIgeT0iNTAlIiBmb250LXNpemU9IjI0IiB0ZXh0LWFuY2hvcj0ibWlkZGxlIiBhbGlnbm1lbnQtYmFzZWxpbmU9Im1pZGRsZSIgZmlsbD0iIzg5YjRmYSI+Tm8gVGh1bWJuYWlsPC90ZXh0Pjwvc3ZnPg==');

    const getThumbnailUrl = (path: string) => {
      if (path) {
        return `http://localhost:8080/thumbnail/${encodeURIComponent(path)}`;
      }
      return placeholderImage.value;
    };

    const handleImageError = (event: Event) => {
      const target = event.target as HTMLImageElement;
      target.src = placeholderImage.value;
    };
    
    // Computed properties
    const uniqueModels = computed(() => {
      const models = new Set<string>();
      streams.value.forEach(stream => {
        if (stream.model) models.add(stream.model);
      });
      return Array.from(models).sort();
    });
    
    const filteredStreams = computed(() => {
      return streams.value.filter(stream => {
        // Only include livestream type files
        if (stream.file_type !== 'livestream') {
          return false;
        }
        
        // Filter by model if selected
        if (selectedModel.value && stream.model !== selectedModel.value) {
          return false;
        }
        
        // Filter by search query
        if (searchQuery.value) {
          const query = searchQuery.value.toLowerCase();
          return (
            stream.model.toLowerCase().includes(query) ||
            stream.path.toLowerCase().includes(query)
          );
        }
        
        return true;
      });
    });
    
    // Methods
    const loadStreams = async () => {
      loading.value = true;
      error.value = '';
      
      try {
        const result = await GetFanslyStreams();
        
        if (result.error) {
          error.value = result.error;
          streams.value = [];
          return;
        }
        
        // Process streams to match contact sheets with videos
        const processedStreams: FanslyStream[] = [];
        const contactSheets = new Map<string, string>();
        
        // First pass - collect all contact sheets
        result.streams.forEach((stream: any) => {
          if (stream.file_type === 'contact_sheet') {
            // Extract base path to match with video
            const basePath = stream.path.replace('_contact_sheet.jpg', '');
            contactSheets.set(basePath, stream.path);
          }
        });
        
        // Second pass - process livestream files and match with contact sheets
        result.streams.forEach((stream: any) => {
          if (stream.file_type === 'livestream') {
            // Check for chat file
            const chatPath = stream.path.replace(/\.(mp4|ts)$/, '_chat.json');
            const hasChat = result.chatFiles.includes(chatPath);
            
            // Match contact sheet
            const basePath = stream.path.replace(/\.(mp4|ts)$/, '');
            const contactSheet = contactSheets.get(basePath);
            
            // Create processed stream object
            processedStreams.push({
              model: stream.model || '',
              hash: stream.hash || '',
              path: stream.path,
              file_type: stream.file_type,
              hasChat: hasChat,
              contactSheet: contactSheet,
              duration: stream.duration || 0
            });
          }
        });
        
        // Sort streams by date (newest first) if possible
        processedStreams.sort((a, b) => {
          const dateA = extractDateFromPath(a.path);
          const dateB = extractDateFromPath(b.path);
          if (dateA && dateB) {
            return dateB.localeCompare(dateA); // Newest first
          }
          return 0;
        });
        
        streams.value = processedStreams;

      } catch (err) {
        console.error('Failed to load Fansly streams:', err);
        error.value = 'Failed to load streams. Please check your configuration.';
        streams.value = [];
      } finally {
        loading.value = false;
      }
    };
    
    const refreshContent = () => {
      loadStreams();
    };
    
    const selectStream = async (stream: FanslyStream) => {
      try {
        console.log("Selecting stream with contact sheet:", stream.contactSheet);
        const result = await LoadFanslyStream(stream.path);
        if (result.success) {
          // If the result doesn't have a contact sheet but we know it from the stream, add it
          if (!result.contactSheet && stream.contactSheet) {
            result.contactSheet = stream.contactSheet;
          }
          emit('stream-selected', result);
        } else {
          console.error('Failed to load stream:', result.error);
          window.alert(`Failed to load stream: ${result.error}`);
        }
      } catch (err) {
        console.error('Error loading stream:', err);
        window.alert(`Error loading stream: ${err}`);
      }
    };
    
    const setupFanslyIntegration = async () => {
      try {
        const config = await GetFanslyConfig();
        configPath.value = config.configPath || '';
        dbPath.value = config.dbPath || '';
        showConfigPanel.value = true;
      } catch (err) {
        console.error('Failed to get Fansly config:', err);
        // Set default paths based on OS
        const isWindows = navigator.platform.indexOf('Win') > -1;
        if (isWindows) {
          configPath.value = '%APPDATA%\\fansly-scraper\\config.toml';
        } else {
          configPath.value = '~/.config/fansly-scraper/config.toml';
        }
        dbPath.value = '';
        showConfigPanel.value = true;
      }
    };
    
    const browseConfigPath = async () => {
      try {
        const path = await BrowseForFile('Select Fansly Scraper Config File', 'TOML files (*.toml)|*.toml');
        if (path) {
          configPath.value = path;
        }
      } catch (err) {
        console.error('Failed to browse for config file:', err);
      }
    };
    
    const browseDbPath = async () => {
      try {
        const path = await BrowseForFolder('Select Folder Containing downloads.db');
        if (path) {
          dbPath.value = path;
        }
      } catch (err) {
        console.error('Failed to browse for database folder:', err);
      }
    };
    
    const saveConfig = async () => {
      try {
        await SaveFanslyConfig({
          configPath: configPath.value,
          dbPath: dbPath.value
        });
        showConfigPanel.value = false;
        loadStreams();
      } catch (err) {
        console.error('Failed to save Fansly config:', err);
        window.alert(`Failed to save configuration: ${err}`);
      }
    };
    
    const cancelConfig = () => {
      showConfigPanel.value = false;
    };
    
    // Helper functions
    const formatStreamTitle = (path: string) => {
      // Extract just the filename from the path
      const filename = path.split('/').pop() || path.split('\\').pop() || path;
      // Remove file extension
      return filename.replace(/\.(mp4|ts)$/, '');
    };
    
    const extractDateFromPath = (path: string) => {
      // Extract date from filename pattern like "model_20250503_191759_id_v7202.mp4"
      const filename = path.split('/').pop() || path.split('\\').pop() || '';
      const dateMatch = filename.match(/_(\d{8}_\d{6})_/);
      if (dateMatch && dateMatch[1]) {
        return dateMatch[1];
      }
      return '';
    };
    
    const formatDate = (dateStr: string) => {
      if (!dateStr) return '';
      
      try {
        // Parse date string in format "20250503_191759"
        const year = dateStr.substring(0, 4);
        const month = dateStr.substring(4, 6);
        const day = dateStr.substring(6, 8);
        const hour = dateStr.substring(9, 11);
        const minute = dateStr.substring(11, 13);
        const second = dateStr.substring(13, 15);
        
        const date = new Date(`${year}-${month}-${day}T${hour}:${minute}:${second}`);
        return date.toLocaleString();
      } catch (e) {
        return dateStr;
      }
    };
    
    const formatDuration = (seconds: number) => {
      const hours = Math.floor(seconds / 3600);
      const minutes = Math.floor((seconds % 3600) / 60);
      const secs = Math.floor(seconds % 60);
      
      if (hours > 0) {
        return `${hours}:${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
      } else {
        return `${minutes}:${secs.toString().padStart(2, '0')}`;
      }
    };
    
    // Lifecycle hooks
    onMounted(() => {
      loadStreams();
    });
    
    return {
      streams,
      loading,
      error,
      searchQuery,
      selectedModel,
      showConfigPanel,
      configPath,
      dbPath,
      uniqueModels,
      filteredStreams,
      refreshContent,
      selectStream,
      setupFanslyIntegration,
      browseConfigPath,
      browseDbPath,
      saveConfig,
      cancelConfig,
      formatStreamTitle,
      extractDateFromPath,
      formatDate,
      getThumbnailUrl,
      handleImageError,
      formatDuration
    };
  }
});
</script>

<style scoped>
.fansly-browser {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: #1e1e2e; /* Catppuccin Mocha background */
  color: #cdd6f4; /* Catppuccin Mocha text */
  position: relative;
}

.browser-header {
  padding: 15px;
  border-bottom: 1px solid #313244; /* Catppuccin Mocha surface0 */
  display: flex;
  flex-direction: column;
  gap: 10px;
  flex-shrink: 0; /* Prevent header from shrinking */
}

.browser-header h2 {
  margin: 0;
  color: #89b4fa; /* Catppuccin Mocha blue */
}

.filter-controls {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.search-input, .model-filter {
  padding: 8px 12px;
  border-radius: 4px;
  border: 1px solid #313244; /* Catppuccin Mocha surface0 */
  background-color: #181825; /* Catppuccin Mocha mantle */
  color: #cdd6f4; /* Catppuccin Mocha text */
  flex: 1;
  min-width: 150px;
}

.refresh-btn {
  padding: 8px 16px;
  border-radius: 4px;
  border: none;
  background-color: #89b4fa; /* Catppuccin Mocha blue */
  color: #1e1e2e; /* Catppuccin Mocha background */
  cursor: pointer;
  font-weight: bold;
}

.refresh-btn:hover {
  background-color: #b4befe; /* Catppuccin Mocha lavender */
}

/* Add a container for the streams content with proper scrolling */
.streams-container {
  flex: 1;
  overflow: auto;
  display: flex;
  flex-direction: column;
  position: relative;
}

.loading-indicator, .error-message, .no-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  gap: 15px;
  padding: 20px;
  text-align: center;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid rgba(137, 180, 250, 0.3); /* Catppuccin Mocha blue with opacity */
  border-radius: 50%;
  border-top-color: #89b4fa; /* Catppuccin Mocha blue */
  animation: spin 1s ease-in-out infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.setup-btn {
  padding: 10px 20px;
  border-radius: 4px;
  border: none;
  background-color: #f38ba8; /* Catppuccin Mocha red */
  color: #1e1e2e; /* Catppuccin Mocha background */
  cursor: pointer;
  font-weight: bold;
}

.setup-btn:hover {
  background-color: #f5c2e7; /* Catppuccin Mocha pink */
}

/* Fix the streams grid to maintain structure and scrollability */
.streams-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 15px;
  padding: 15px;
  align-content: start; /* Start from the top */
}

.stream-card {
  background-color: #181825; /* Catppuccin Mocha mantle */
  border-radius: 8px;
  overflow: hidden;
  transition: transform 0.2s, box-shadow 0.2s;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  height: auto; /* Let the card size naturally based on content */
  position: relative; /* Ensure proper stacking context */
}

.stream-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.3);
}

.thumbnail-container {
  position: relative;
  width: 100%;
  padding-top: 56.25%; /* 16:9 aspect ratio */
  background-color: #11111b; /* Catppuccin Mocha crust */
  flex-shrink: 0; /* Don't allow shrinking */
  overflow: hidden; /* Ensure content doesn't overflow */
}

.stream-thumbnail {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.no-thumbnail {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #313244; /* Catppuccin Mocha surface0 */
  color: #6c7086; /* Catppuccin Mocha subtext0 */
}

.stream-info-overlay {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  justify-content: space-between;
  padding: 5px 10px;
  background: linear-gradient(transparent, rgba(24, 24, 37, 0.8));
  z-index: 2;
}

.stream-duration {
  background-color: rgba(24, 24, 37, 0.7);
  padding: 3px 6px;
  border-radius: 3px;
  font-size: 0.8rem;
}

.chat-badge {
  background-color: #89b4fa; /* Catppuccin Mocha blue */
  color: #1e1e2e; /* Catppuccin Mocha background */
  padding: 3px 6px;
  border-radius: 3px;
  font-size: 0.8rem;
  font-weight: bold;
}

.stream-info {
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex: 1; /* Take up remaining space */
  z-index: 1;
}

.stream-title {
  margin: 0;
  font-size: 1rem;
  line-height: 1.3;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.stream-details {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.85rem;
  color: #a6adc8; /* Catppuccin Mocha subtext1 */
  margin-top: auto; /* Push to bottom of container */
}

.stream-model {
  font-weight: bold;
  color: #89b4fa; /* Catppuccin Mocha blue */
  max-width: 50%; /* Limit width to prevent overflow */
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.stream-date {
  max-width: 50%; /* Limit width to prevent overflow */
  text-align: right;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* Config panel styles */
.config-panel {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(17, 17, 27, 0.8); /* Catppuccin Mocha crust with opacity */
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  backdrop-filter: blur(5px);
}

.config-panel-content {
  background-color: #1e1e2e; /* Catppuccin Mocha background */
  border-radius: 8px;
  padding: 20px;
  width: 80%;
  max-width: 600px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.5);
}

.config-panel-content h3 {
  margin-top: 0;
  color: #89b4fa; /* Catppuccin Mocha blue */
  border-bottom: 1px solid #313244; /* Catppuccin Mocha surface0 */
  padding-bottom: 10px;
}

.config-field {
  margin-bottom: 15px;
}

.config-field label {
  display: block;
  margin-bottom: 5px;
  color: #cdd6f4; /* Catppuccin Mocha text */
}

.path-input-group {
  display: flex;
  gap: 10px;
}

.path-input-group input {
  flex: 1;
  padding: 8px 12px;
  border-radius: 4px;
  border: 1px solid #313244; /* Catppuccin Mocha surface0 */
  background-color: #181825; /* Catppuccin Mocha mantle */
  color: #cdd6f4; /* Catppuccin Mocha text */
}

.path-input-group button {
  padding: 8px 16px;
  border-radius: 4px;
  border: none;
  background-color: #89b4fa; /* Catppuccin Mocha blue */
  color: #1e1e2e; /* Catppuccin Mocha background */
  cursor: pointer;
}

.config-field small {
  display: block;
  margin-top: 5px;
  color: #a6adc8; /* Catppuccin Mocha subtext1 */
  font-size: 0.8rem;
}

.config-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 20px;
}

.save-btn, .cancel-btn {
  padding: 10px 20px;
  border-radius: 4px;
  border: none;
  cursor: pointer;
  font-weight: bold;
}

.save-btn {
  background-color: #a6e3a1; /* Catppuccin Mocha green */
  color: #1e1e2e; /* Catppuccin Mocha background */
}

.cancel-btn {
  background-color: #f38ba8; /* Catppuccin Mocha red */
  color: #1e1e2e; /* Catppuccin Mocha background */
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .filter-controls {
    flex-direction: column;
  }
  
  .streams-grid {
    grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  }
  
  .config-panel-content {
    width: 95%;
    padding: 15px;
  }
  
  .path-input-group {
    flex-direction: column;
  }
}
</style>
