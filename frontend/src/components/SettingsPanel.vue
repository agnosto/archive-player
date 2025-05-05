<template>
  <div class="settings-panel">
    <h2>Settings</h2>
    
    <!-- Integrations Section -->
    <div class="settings-section">
      <h3>Integrations</h3>
      
      <!-- Fansly Integration -->
      <div class="setting-item">
        <div class="setting-label">
          <span>Fansly Scraper</span>
          <span class="status-indicator" :class="{ 'status-configured': isFanslyConfigured }">
            {{ isFanslyConfigured ? 'Configured' : 'Not Configured' }}
          </span>
        </div>
        <div class="setting-controls">
          <button @click="configureFanslyIntegration">
            {{ isFanslyConfigured ? 'Reconfigure' : 'Configure' }}
          </button>
        </div>
      </div>
      
      <!-- Placeholder for future integrations -->
      <div class="setting-item disabled">
        <div class="setting-label">
          <span>Other Integrations</span>
          <span class="status-indicator status-coming-soon">Coming Soon</span>
        </div>
        <div class="setting-controls">
          <button disabled>Configure</button>
        </div>
      </div>
    </div>
    
    <!-- Appearance Settings Section -->
    <div class="settings-section">
      <h3>Appearance</h3>
      
      <div class="setting-group">
        <h4>Chat Position</h4>
        <div class="radio-group">
          <label>
            <input 
              type="radio"
              name="chatPosition"
              value="left"
              :checked="localTheme.chatPosition === 'left'"
              @change="updateSetting('chatPosition', 'left')"
            />
            Left
          </label>
          <label>
            <input 
              type="radio"
              name="chatPosition"
              value="right"
              :checked="localTheme.chatPosition === 'right'"
              @change="updateSetting('chatPosition', 'right')"
            />
            Right
          </label>
        </div>
      </div>
      
      <div class="setting-group">
        <h4>Chat Width</h4>
        <div class="slider-container">
          <input 
            type="range"
            min="200"
            max="500"
            step="10"
            v-model.number="localTheme.chatWidth"
            @change="emitThemeUpdate"
          />
          <span>{{ localTheme.chatWidth }}px</span>
        </div>
      </div>
      
      <div class="setting-group">
        <h4>Chat Background</h4>
        <div class="color-picker">
          <input 
            type="color"
            v-model="localTheme.chatBgColor"
            @change="emitThemeUpdate"
          />
          <span>{{ localTheme.chatBgColor }}</span>
        </div>
        
        <div class="slider-container">
          <label>Opacity</label>
          <input 
            type="range"
            min="0"
            max="1"
            step="0.05"
            v-model.number="localTheme.chatOpacity"
            @change="emitThemeUpdate"
          />
          <span>{{ Math.round(localTheme.chatOpacity * 100) }}%</span>
        </div>
      </div>
      
      <div class="setting-group">
        <h4>Text Settings</h4>
        <div class="color-picker">
          <label>Text Color</label>
          <input 
            type="color"
            v-model="localTheme.chatTextColor"
            @change="emitThemeUpdate"
          />
          <span>{{ localTheme.chatTextColor }}</span>
        </div>
        
        <div class="color-picker">
          <label>Author Name Color</label>
          <input 
            type="color"
            v-model="localTheme.authorNameColor"
            @change="emitThemeUpdate"
          />
          <span>{{ localTheme.authorNameColor }}</span>
        </div>
        
        <div class="slider-container">
          <label>Font Size</label>
          <input 
            type="range"
            min="10"
            max="24"
            step="1"
            v-model.number="localTheme.chatFontSize"
            @change="emitThemeUpdate"
          />
          <span>{{ localTheme.chatFontSize }}px</span>
        </div>
        
        <div class="slider-container">
          <label>Message Spacing</label>
          <input 
            type="range"
            min="4"
            max="20"
            step="1"
            v-model.number="localTheme.messageSpacing"
            @change="emitThemeUpdate"
          />
          <span>{{ localTheme.messageSpacing }}px</span>
        </div>
      </div>
      
      <div class="setting-group">
        <button class="reset-button" @click="resetToDefaults">Reset to Defaults</button>
      </div>
    </div>
  </div>

  <!-- Fansly Configuration Modal -->
  <div class="modal" v-if="showFanslyConfigModal">
    <div class="modal-content">
      <h3>Fansly Integration Setup</h3>
      
      <div class="config-field">
        <label for="configPath">Config Path:</label>
        <div class="path-input-group">
          <input type="text" id="configPath" v-model="fanslyConfigPath" readonly />
          <button @click="browseFanslyConfigPath">Browse</button>
        </div>
        <small>Default locations: ~/.config/fansly-scraper/config.toml (Mac/Linux) or %APPDATA%\fansly-scraper\config.toml (Windows)</small>
      </div>
      
      <div class="config-field">
        <label for="dbPath">Database Path:</label>
        <div class="path-input-group">
          <input type="text" id="dbPath" v-model="fanslyDbPath" readonly />
          <button @click="browseFanslyDbPath">Browse</button>
        </div>
        <small>Usually located in the save_location specified in your config</small>
      </div>
      
      <div class="modal-actions">
        <button @click="saveFanslyConfig" class="save-btn">Save</button>
        <button @click="cancelFanslyConfig" class="cancel-btn">Cancel</button>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, watch, onMounted } from 'vue';
import { ThemeSettings } from '../types';
import { GetFanslyConfig, SaveFanslyConfig, BrowseForFile, BrowseForFolder } from '../../wailsjs/go/main/App';

// Default theme settings
const defaultTheme: ThemeSettings = {
  chatBgColor: '#1e1e2e', // Catppuccin Mocha background
  chatTextColor: '#cdd6f4', // Catppuccin Mocha text
  chatOpacity: 0.8,
  chatFontSize: 14,
  authorNameColor: '#89b4fa', // Catppuccin Mocha blue
  messageSpacing: 8,
  chatWidth: 300,
  chatPosition: 'right'
};

export default defineComponent({
  name: 'SettingsPanel',
  props: {
    theme: {
      type: Object as () => ThemeSettings,
      required: true
    }
  },
  emits: ['update:theme'],
  setup(props, { emit }) {
    // Create a local copy of the theme to work with
    const localTheme = ref<ThemeSettings>({ ...props.theme });
    const showFanslyConfigModal = ref(false);
    const fanslyConfigPath = ref('');
    const fanslyDbPath = ref('');
    const isFanslyConfigured = ref(false);

    const configureFanslyIntegration = async () => {
      try {
        const config = await GetFanslyConfig();
        fanslyConfigPath.value = config.configPath || '';
        fanslyDbPath.value = config.dbPath || '';
        showFanslyConfigModal.value = true;
      } catch (err) {
        console.error('Failed to get Fansly config:', err);
        fanslyConfigPath.value = '';
        fanslyDbPath.value = '';
        showFanslyConfigModal.value = true;
      }
    };
    
    const browseFanslyConfigPath = async () => {
      try {
        const path = await BrowseForFile('Select Fansly Scraper Config File', 'TOML files (*.toml)|*.toml');
        if (path) {
          fanslyConfigPath.value = path;
        }
      } catch (err) {
        console.error('Failed to browse for config file:', err);
      }
    };
    
    const browseFanslyDbPath = async () => {
      try {
        const path = await BrowseForFolder('Select Folder Containing downloads.db');
        if (path) {
          fanslyDbPath.value = path;
        }
      } catch (err) {
        console.error('Failed to browse for database folder:', err);
      }
    };
    
    const saveFanslyConfig = async () => {
      try {
        await SaveFanslyConfig({
          configPath: fanslyConfigPath.value,
          dbPath: fanslyDbPath.value
        });
        showFanslyConfigModal.value = false;
        checkFanslyConfig(); // Update status after saving
      } catch (err) {
        console.error('Failed to save Fansly config:', err);
        window.alert(`Failed to save configuration: ${err}`);
      }
    };
    
    const cancelFanslyConfig = () => {
      showFanslyConfigModal.value = false;
    };
    
    // Check if Fansly is configured
    const checkFanslyConfig = async () => {
      try {
        const config = await GetFanslyConfig();
        isFanslyConfigured.value = !!(config.configPath && config.dbPath);
      } catch (err) {
        console.error('Failed to check Fansly config:', err);
        isFanslyConfigured.value = false;
      }
    };

    // Watch for changes in the props
    watch(() => props.theme, (newTheme) => {
      localTheme.value = { ...newTheme };
    });

    // Update a specific setting
    const updateSetting = (key: keyof ThemeSettings, value: any) => {
      (localTheme.value as any)[key] = value;
      emitThemeUpdate();
    };

    // Emit the theme update event
    const emitThemeUpdate = () => {
      emit('update:theme', { ...localTheme.value });
    };

    // Reset to default theme
    const resetToDefaults = () => {
      localTheme.value = { ...defaultTheme };
      emitThemeUpdate();
    };
    
    // Check Fansly config on mount
    onMounted(() => {
      checkFanslyConfig();
    });

    return {
      localTheme,
      updateSetting,
      emitThemeUpdate,
      resetToDefaults,
      showFanslyConfigModal,
      fanslyConfigPath,
      fanslyDbPath,
      isFanslyConfigured,
      configureFanslyIntegration,
      browseFanslyConfigPath,
      browseFanslyDbPath,
      saveFanslyConfig,
      cancelFanslyConfig
    };
  }
});
</script>

<style scoped>
.settings-panel {
  padding: 15px;
  height: 100%;
  overflow-y: auto;
}

h2 {
  margin-bottom: 20px;
  color: #89b4fa; /* Catppuccin Mocha blue */
  font-size: 1.2rem;
}

h3 {
  margin-bottom: 15px;
  color: #89b4fa; /* Catppuccin Mocha blue */
  font-size: 1.1rem;
  border-bottom: 1px solid #313244; /* Catppuccin Mocha surface0 */
  padding-bottom: 5px;
}

h4 {
  margin-bottom: 10px;
  color: #cdd6f4; /* Catppuccin Mocha text */
  font-size: 1rem;
}

.settings-section {
  margin-bottom: 30px;
}

.setting-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  background-color: #181825; /* Catppuccin Mocha mantle */
  border-radius: 6px;
  margin-bottom: 10px;
}

.setting-item.disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.setting-label {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.status-indicator {
  font-size: 0.8rem;
  padding: 2px 8px;
  border-radius: 10px;
  background-color: #313244; /* Catppuccin Mocha surface0 */
  display: inline-block;
}

.status-configured {
  background-color: #a6e3a1; /* Catppuccin Mocha green */
  color: #1e1e2e; /* Catppuccin Mocha background */
}

.status-coming-soon {
  background-color: #f9e2af; /* Catppuccin Mocha yellow */
  color: #1e1e2e; /* Catppuccin Mocha background */
}

.setting-controls button {
  padding: 6px 12px;
  background-color: #89b4fa; /* Catppuccin Mocha blue */
  color: #1e1e2e; /* Catppuccin Mocha background */
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: bold;
  transition: background-color 0.2s;
}

.setting-controls button:hover {
  background-color: #b4befe; /* Catppuccin Mocha lavender */
}

.setting-controls button:disabled {
  background-color: #45475a; /* Catppuccin Mocha surface1 */
  cursor: not-allowed;
}

.setting-group {
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid #313244; /* Catppuccin Mocha surface0 */
}

.setting-group:last-child {
  border-bottom: none;
}

.radio-group {
  display: flex;
  gap: 15px;
}

.radio-group label {
  display: flex;
  align-items: center;
  gap: 5px;
  cursor: pointer;
}
.slider-container {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 10px;
}

.slider-container label {
  min-width: 100px;
}

.slider-container input[type="range"] {
  flex: 1;
}

.color-picker {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 10px;
}

.color-picker label {
  min-width: 100px;
}

input[type="color"] {
  width: 40px;
  height: 25px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.reset-button {
  padding: 8px 16px;
  background-color: #f38ba8; /* Catppuccin Mocha red */
  color: #1e1e2e; /* Catppuccin Mocha background */
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: bold;
  transition: background-color 0.2s;
}

.reset-button:hover {
  background-color: #f5c2e7; /* Catppuccin Mocha pink */
}

/* Scrollbar styling */
.settings-panel::-webkit-scrollbar {
  width: 6px;
}

.settings-panel::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.1);
  border-radius: 3px;
}

.settings-panel::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.3);
  border-radius: 3px;
}

.settings-panel::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.5);
}

.modal {
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

.modal-content {
  background-color: #1e1e2e; /* Catppuccin Mocha background */
  border-radius: 8px;
  padding: 20px;
  width: 90%;
  max-width: 600px;
  box-shadow: 0 5px 20px rgba(0, 0, 0, 0.3);
}

.modal-content h3 {
  margin-top: 0;
  margin-bottom: 20px;
  color: #89b4fa; /* Catppuccin Mocha blue */
  text-align: center;
}

.config-field {
  margin-bottom: 15px;
}

.config-field label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
}

.config-field small {
  display: block;
  margin-top: 5px;
  color: #a6adc8; /* Catppuccin Mocha subtext1 */
  font-size: 0.8rem;
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
  background-color: #313244; /* Catppuccin Mocha surface0 */
  color: #cdd6f4; /* Catppuccin Mocha text */
  cursor: pointer;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 20px;
}

.save-btn {
  padding: 8px 16px;
  border-radius: 4px;
  border: none;
  background-color: #a6e3a1; /* Catppuccin Mocha green */
  color: #1e1e2e; /* Catppuccin Mocha background */
  cursor: pointer;
  font-weight: bold;
}

.cancel-btn {
  padding: 8px 16px;
  border-radius: 4px;
  border: none;
  background-color: #313244; /* Catppuccin Mocha surface0 */
  color: #cdd6f4; /* Catppuccin Mocha text */
  cursor: pointer;
}

.save-btn:hover {
  background-color: #94e2d5; /* Catppuccin Mocha teal */
}

.cancel-btn:hover {
  background-color: #45475a; /* Catppuccin Mocha surface1 */
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .path-input-group {
    flex-direction: column;
  }
  
  .modal-content {
    width: 95%;
    padding: 15px;
  }
  
  .setting-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  
  .setting-controls {
    width: 100%;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
