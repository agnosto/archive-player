<template>
  <div class="container" :class="{ 'theater-mode': theaterMode }">
    <div class="header">
      <!-- Hamburger menu button for mobile -->
      <div class="hamburger-menu" @click="toggleMobileMenu">
        <div class="hamburger-line"></div>
        <div class="hamburger-line"></div>
        <div class="hamburger-line"></div>
      </div>
      
      <!-- Logo and title (hidden on mobile when menu closed) -->
      <div class="logo-container" :class="{ 'mobile-hidden': !mobileMenuOpen }">
        <img 
          src="./assets/logo.svg" 
          alt="Archive Player Logo" 
          class="logo clickable" 
          @click="resetPlayer" 
          title="Reset player"
        />
        <h1 
          class="clickable" 
          @click="resetPlayer" 
          title="Reset player"
        >Archive Player</h1>
      </div>
      
      <!-- Controls (desktop only) -->
      <div class="desktop-controls">
        <button @click="openVideoFile">Open Video</button>
        <button @click="openChatFile" :disabled="!videoLoaded">Load Chat</button>
        <button @click="toggleTheaterMode" :disabled="!videoLoaded">
          {{ theaterMode ? 'Exit Theater Mode' : 'Theater Mode' }}
        </button>
        <button @click="toggleChat" :disabled="!chatLoaded">
          {{ showChat ? 'Hide Chat' : 'Show Chat' }}
        </button>
        <button @click="toggleSettings">{{ showSettings ? 'Hide Settings' : 'Show Settings' }}</button>
        <button @click="toggleRecentVideos">{{ showRecentVideos ? 'Hide Recent' : 'Recent Videos' }}</button>
        <!-- In desktop-controls div -->
        <button @click="toggleClipCreator" :disabled="!videoLoaded">
          {{ showClipCreator ? 'Hide Clip Creator' : 'Create Clip' }}
        </button>


        
        <!-- Integrations dropdown -->
        <div class="dropdown">
          <button class="dropdown-toggle" @click="toggleIntegrationsMenu">
            Integrations <span class="dropdown-arrow">▼</span>
          </button>
          <div class="dropdown-menu" v-if="showIntegrationsMenu">
            <button @click="toggleFanslyBrowser">Fansly Streams</button>
            <!-- Add future integrations here -->
            <!-- <button @click="toggleTwitchBrowser">Twitch Streams</button> -->
            <!-- <button @click="toggleYouTubeBrowser">YouTube Streams</button> -->
          </div>
        </div>
      </div>
      
      <!-- Mobile menu panel -->
      <div class="mobile-menu" :class="{ 'mobile-menu-open': mobileMenuOpen }">
        <!-- Logo and title in mobile menu -->
        <div class="mobile-menu-header">
          <img 
            src="./assets/logo.svg" 
            alt="Archive Player Logo" 
            class="logo clickable" 
            @click="resetPlayer" 
            title="Reset player"
          />
          <h1 
            class="clickable" 
            @click="resetPlayer" 
            title="Reset player"
          >Archive Player</h1>
        </div>
        
        <div class="mobile-menu-buttons">
          <button @click="openVideoFile">Open Video</button>
          <button @click="openChatFile" :disabled="!videoLoaded">Load Chat</button>
          <button @click="toggleTheaterMode" :disabled="!videoLoaded">
            {{ theaterMode ? 'Exit Theater Mode' : 'Theater Mode' }}
          </button>
          <button @click="toggleChat" :disabled="!chatLoaded">
            {{ showChat ? 'Hide Chat' : 'Show Chat' }}
          </button>
          <button @click="toggleSettings">{{ showSettings ? 'Hide Settings' : 'Show Settings' }}</button>
          <button @click="toggleRecentVideos">{{ showRecentVideos ? 'Hide Recent' : 'Recent Videos' }}</button>
          <!-- Also add to mobile-menu-buttons div -->
          <button @click="toggleClipCreator" :disabled="!videoLoaded">
            {{ showClipCreator ? 'Hide Clip Creator' : 'Create Clip' }}
          </button>
          
          <!-- Mobile integrations submenu -->
          <div class="mobile-submenu">
            <button class="mobile-submenu-toggle" @click="toggleMobileIntegrationsMenu">
              Integrations {{ showMobileIntegrationsMenu ? '▲' : '▼' }}
            </button>
            <div class="mobile-submenu-content" v-if="showMobileIntegrationsMenu">
              <button @click="toggleFanslyBrowser">Fansly Streams</button>
              <!-- Add future integrations here -->
              <!-- <button @click="toggleTwitchBrowser">Twitch Streams</button> -->
              <!-- <button @click="toggleYouTubeBrowser">YouTube Streams</button> -->
            </div>
          </div>
          
          <!-- Close button for mobile menu -->
          <button class="close-menu-btn" @click="toggleMobileMenu">Close Menu</button>
        </div>
      </div>
      
      <!-- Overlay for mobile menu -->
      <div class="mobile-menu-overlay" v-if="mobileMenuOpen" @click="toggleMobileMenu"></div>
    </div>
    
    <div class="main-content" :class="{
      'with-settings': showSettings,
      'with-recent': showRecentVideos,
      'with-integration': showFanslyBrowser || showAnyOtherIntegration,
      'with-clip-creator': showClipCreator
    }">
      <div class="video-chat-container" :class="{ 
        'chat-hidden': !showChat, 
        'chat-left': showChat && chatLoaded && theme.chatPosition === 'left',
        'chat-right': showChat && chatLoaded && theme.chatPosition === 'right'
      }">
        <div class="video-container">
          <VideoPlayer 
            :videoSrc="videoSrc"
            :theme="theme"
            :showChat="showChat && chatLoaded"
            :theaterMode="theaterMode"
            @update:currentTime="currentTime = $event"
            @update:duration="videoDuration = $event"
          />
        </div> 
        
        <div v-if="showChat && chatLoaded" class="chat-container" :style="chatContainerStyle">
          <ChatOverlay
            :messages="currentMessages"
            :theme="theme"
          />
        </div>
      </div>
      
      <SettingsPanel
        v-if="showSettings"
        :theme="theme"
        @update:theme="updateTheme"
        class="settings-container"
      />

      <ClipCreator
          v-if="showClipCreator"
          :currentTime="currentTime"
          :videoDuration="videoDuration"
          @clip-loaded="handleClipLoaded"
          class="clip-creator-container"
        />
      
      <RecentVideos
        v-if="showRecentVideos"
        :recentVideos="recentVideos"
        @select-video="loadRecentVideo"
        class="recent-videos-container"
      />
      
      <!-- Integration panels -->
      <div v-if="showFanslyBrowser || showAnyOtherIntegration" class="integration-container">
        <div class="integration-tabs">
          <button 
            @click="selectIntegration('fansly')" 
            :class="{ 'active': activeIntegration === 'fansly' }"
          >
            Fansly
          </button>
          <!-- Add tabs for future integrations -->
          <!-- <button 
            @click="selectIntegration('twitch')" 
            :class="{ 'active': activeIntegration === 'twitch' }"
          >
            Twitch
          </button> -->
        </div>
        
        <div class="integration-content">
          <FanslyBrowser
            v-if="activeIntegration === 'fansly'"
            @stream-selected="loadFanslyStream"
            class="integration-panel"
          />
          <!-- Add content for future integrations -->
          <!-- <TwitchBrowser
            v-if="activeIntegration === 'twitch'"
            @stream-selected="loadTwitchStream"
            class="integration-panel"
          /> -->
        </div>
      </div>
    </div>
    
    <div class="status-bar">
      <div v-if="videoInfo.filename">
        <span>Video: {{ videoInfo.filename }}</span>
      </div>
      <div v-if="chatLoaded">
        <span>Chat: {{ chatMessagesCount }} messages loaded</span>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted, computed, watch } from 'vue';
import VideoPlayer from './components/VideoPlayer.vue';
import SettingsPanel from './components/SettingsPanel.vue';
import ChatOverlay from './components/ChatOverlay.vue';
import RecentVideos from './components/RecentVideos.vue';
import FanslyBrowser from './components/FanslyBrowser.vue';
import ClipCreator from './components/ClipCreator.vue';
import { OpenVideoFile, OpenChatFile, GetVideoFileInfo, GetAllChatMessages, GetMessagesAtTime, LoadVideoFromPath } from '../wailsjs/go/main/App';
import { ThemeSettings, RecentVideo } from './types';

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
  name: 'App',
  components: {
    VideoPlayer,
    SettingsPanel,
    ChatOverlay,
    RecentVideos,
    FanslyBrowser,
    ClipCreator,
  },
  setup() {
    const videoSrc = ref('');
    const videoLoaded = ref(false);
    const chatLoaded = ref(false);
    const chatMessagesCount = ref(0);
    const videoInfo = ref<Record<string, string>>({ filename: '', path: '' });
    const showSettings = ref(false);
    const showRecentVideos = ref(false);
    const theaterMode = ref(false);
    const showChat = ref(true);
    const theme = ref<ThemeSettings>({ ...defaultTheme });
    const currentMessages = ref<any[]>([]);
    const currentTime = ref(0);
    const recentVideos = ref<RecentVideo[]>([]);
    const mobileMenuOpen = ref(false);
    const showFanslyBrowser = ref(false);
    const accumulatedMessages = ref<any[]>([]);
    // Inside the setup function
    const showClipCreator = ref(false);
    const videoDuration = ref(0);
    
    // Add this missing ref for chat messages
    const chatMessages = ref<any[]>([]);
    
    // New refs for integrations menu
    const showIntegrationsMenu = ref(false);
    const showMobileIntegrationsMenu = ref(false);
    const activeIntegration = ref('fansly');
    const showAnyOtherIntegration = ref(false); // For future integrations

    // Toggle clip creator panel
    const toggleClipCreator = () => {
      showClipCreator.value = !showClipCreator.value;
      
      if (showClipCreator.value) {
        showSettings.value = false;
        showRecentVideos.value = false;
        showFanslyBrowser.value = false;
        showAnyOtherIntegration.value = false;
        
        // Close mobile menu if open
        if (mobileMenuOpen.value) {
          toggleMobileMenu();
        }
      }
    };

    // Handle when a clip is loaded from the clip creator
    const handleClipLoaded = (videoPath: string) => {
      videoSrc.value = videoPath;
      videoLoaded.value = true;
      
      // Close clip creator panel
      showClipCreator.value = false;
    };
    
    // Toggle integrations dropdown menu
    const toggleIntegrationsMenu = () => {
      showIntegrationsMenu.value = !showIntegrationsMenu.value;
      
      // Close menu when clicking outside
      if (showIntegrationsMenu.value) {
        setTimeout(() => {
          const closeDropdown = (e: MouseEvent) => {
            const target = e.target as HTMLElement;
            if (!target.closest('.dropdown')) {
              showIntegrationsMenu.value = false;
              document.removeEventListener('click', closeDropdown);
            }
          };
          document.addEventListener('click', closeDropdown);
        }, 0);
      }
    };
    
    // Toggle mobile integrations submenu
    const toggleMobileIntegrationsMenu = () => {
      showMobileIntegrationsMenu.value = !showMobileIntegrationsMenu.value;
    };
    
    // Select active integration
    const selectIntegration = (integration: string) => {
      activeIntegration.value = integration;
    };
    
    const toggleFanslyBrowser = () => {
      showFanslyBrowser.value = !showFanslyBrowser.value;
      
      // Close other panels when opening Fansly browser
      if (showFanslyBrowser.value) {
        showSettings.value = false;
        showRecentVideos.value = false;
        activeIntegration.value = 'fansly';
        
        // Close mobile menu if open
        if (mobileMenuOpen.value) {
          toggleMobileMenu();
        }
        
        // Close integrations menu
        showIntegrationsMenu.value = false;
      }
    };

    // Function to toggle any integration panel
    const toggleIntegrationPanel = (integration: string) => {
      // Close other panels
      showSettings.value = false;
      showRecentVideos.value = false;
      
      // Set the active integration
      activeIntegration.value = integration;
      
      // Toggle the appropriate integration panel
      if (integration === 'fansly') {
        showFanslyBrowser.value = !showFanslyBrowser.value;
        showAnyOtherIntegration.value = false;
      } else {
        // For future integrations
        showFanslyBrowser.value = false;
        showAnyOtherIntegration.value = !showAnyOtherIntegration.value;
      }
      
      // Close mobile menu if open
      if (mobileMenuOpen.value) {
        toggleMobileMenu();
      }
      
      // Close integrations menu
      showIntegrationsMenu.value = false;
    };

    const resetVideoState = () => {
      // Reset video-related state
      videoSrc.value = '';
      videoLoaded.value = false;
      currentTime.value = 0;
      videoDuration.value = 0;
      
      // Reset chat-related state
      resetChat();
    };

    const resetPlayer = () => {
      // Reset video state
      resetVideoState();
      
      // Reset UI states
      showSettings.value = false;
      showRecentVideos.value = false;
      showFanslyBrowser.value = false;
      showAnyOtherIntegration.value = false;
      theaterMode.value = false;
      
      // Reset video info
      videoInfo.value = { filename: '', path: '' };
      
      // Close mobile menu if open
      if (mobileMenuOpen.value) {
        toggleMobileMenu();
      }
      
      console.log('Player state reset');
    };

    const loadFanslyStream = async (result: any) => {
      if (result.success) {
        try {
          // Reset current video state first
          resetVideoState();
          
          // Load the video using the same function as loadRecentVideo
          const videoPath = await LoadVideoFromPath(result.videoPath);
          videoSrc.value = videoPath;
          videoLoaded.value = true;
          
          // Extract filename for display
          const pathParts = result.videoPath.split(/[\/\\]/);
          videoInfo.value = {
            filename: pathParts[pathParts.length - 1],
            path: result.videoPath,
            thumbnailPath: result.contactSheet || '' // Use contactSheet from result
          };
          
          console.log("Loading Fansly stream with thumbnail:", result.contactSheet);
          
          // Add to recent videos with thumbnail
          addToRecentVideos(videoInfo.value);
          
          // Load chat if available
          if (result.chatPath) {
            try {
              await OpenChatFile([result.chatPath]);
              // After opening the chat file, get all messages separately
              const messages = await GetAllChatMessages();
              if (messages && messages.length > 0) {
                chatMessages.value = messages;
                chatMessagesCount.value = messages.length;
                chatLoaded.value = true;
                showChat.value = true;
                // Update current messages
                updateMessages();
              } else {
                console.error('No chat messages found in file');
                chatLoaded.value = false;
                showChat.value = false;
              }
            } catch (error) {
              console.error('Failed to load chat file:', error);
              chatLoaded.value = false;
              showChat.value = false;
            }
          } else {
            chatLoaded.value = false;
            showChat.value = false;
          }
          
          // Close Fansly browser after loading
          showFanslyBrowser.value = false;
        } catch (error) {
          console.error('Failed to load Fansly stream:', error);
          window.alert(`Failed to load stream: ${error}`);
        }
      }
    };

    // Function to toggle mobile menu
    const toggleMobileMenu = () => {
      mobileMenuOpen.value = !mobileMenuOpen.value;
      
      // Toggle active class on hamburger menu for animation
      const hamburgerMenu = document.querySelector('.hamburger-menu');
      if (hamburgerMenu) {
        if (mobileMenuOpen.value) {
          hamburgerMenu.classList.add('active');
        } else {
          hamburgerMenu.classList.remove('active');
        }
      }
      // If opening the menu, close settings and recent videos panels
      if (mobileMenuOpen.value) {
        showSettings.value = false;
        showRecentVideos.value = false;
        showFanslyBrowser.value = false;
        showAnyOtherIntegration.value = false;
      }
    };

    // Computed style for chat container based on theme and screen size
    const chatContainerStyle = computed(() => {
      const isMobile = window.innerWidth < 768;
      
      // Convert hex color to rgba
      const hexToRgba = (hex: string, opacity: number) => {
        hex = hex.replace('#', '');
        const r = parseInt(hex.substring(0, 2), 16);
        const g = parseInt(hex.substring(2, 4), 16);
        const b = parseInt(hex.substring(4, 6), 16);
        return `rgba(${r}, ${g}, ${b}, ${opacity})`;
      };

      return {
        width: isMobile ? '100%' : `${theme.value.chatWidth}px`,
        backgroundColor: hexToRgba(theme.value.chatBgColor, theme.value.chatOpacity),
        color: theme.value.chatTextColor,
        fontSize: `${theme.value.chatFontSize}px`,
      };
    });

    // Function to update messages based on current time
    const updateMessages = async () => {
      if (videoLoaded.value && chatLoaded.value) {
        try {
          // Get messages up to the current time with a 30 second window
          const messages = await GetMessagesAtTime(currentTime.value, 30);
          
          if (messages && Array.isArray(messages)) {
            // Find new messages that aren't already in accumulatedMessages
            const existingIds = new Set(accumulatedMessages.value.map(m => m.message_id));
            const newMessages = messages.filter(m => !existingIds.has(m.message_id));
            
            if (newMessages.length > 0) {
              // Add new messages to accumulated messages
              accumulatedMessages.value = [...accumulatedMessages.value, ...newMessages];
              
              // Sort by timestamp to ensure chronological order
              accumulatedMessages.value.sort((a, b) => a.timestamp - b.timestamp);
              
              console.log(`Added ${newMessages.length} new messages, total: ${accumulatedMessages.value.length}`);
            }
            
            // Update current messages for display
            currentMessages.value = [...accumulatedMessages.value];
          } else {
            console.warn('GetMessagesAtTime returned invalid data:', messages);
          }
        } catch (error) {
          console.error('Error getting messages:', error);
        }
      }
    };

    // Watch for time changes to update messages
    watch(currentTime, () => {
      updateMessages();
    });

    const openVideoFile = async () => {
      try {
        // Reset current video state first
        resetVideoState();
        
        const videoUrl = await OpenVideoFile();
        if (videoUrl) {
          videoSrc.value = videoUrl;
          videoLoaded.value = true;
          
          // Get video file info
          const info = await GetVideoFileInfo();
          videoInfo.value = info;
          
          // Add to recent videos
          addToRecentVideos(info);
        }
      } catch (error) {
        console.error('Error opening video file:', error);
        alert('Failed to open video file');
      }
    };

    const openChatFile = async () => {
      try {
        // Reset accumulated messages
        accumulatedMessages.value = [];
        
        // Pass an empty array to indicate we want to use the file dialog
        const filePath = await OpenChatFile([]);
        if (filePath) {
          // After opening the chat file, get all messages separately
          const messages = await GetAllChatMessages();
          chatMessagesCount.value = messages.length;
          chatLoaded.value = true;
          showChat.value = true;
          
          // Update current messages
          updateMessages();
        }
      } catch (error) {
        console.error('Error opening chat file:', error);
        alert('Failed to load chat file: ' + error);
      }
    };

    const loadChatFile = async (path?: string) => {
      try {
        let chatPath = "";
        
        if (!path) {
          // If no path provided, use OpenChatFile with empty array to trigger file dialog
          chatPath = await OpenChatFile([]);
        } else {
          // If path is provided, use it directly in an array
          chatPath = await OpenChatFile([path]);
        }
        
        if (!chatPath) return; // User cancelled
        
        // Get all messages to check if they loaded correctly
        const allMessages = await GetAllChatMessages();
        console.log("All chat messages:", allMessages);
        
        if (allMessages && allMessages.length > 0) {
          chatLoaded.value = true;
          chatMessagesCount.value = allMessages.length;
          showChat.value = true;
          
          // Update current messages
          updateMessages();
        } else {
          console.error("No chat messages found in file");
          window.alert("No chat messages found in the file. Please check the file format.");
          chatLoaded.value = false;
          chatMessagesCount.value = 0;
        }
      } catch (error) {
        console.error("Error loading chat file:", error);
        window.alert("Error loading chat file: " + error);
        chatLoaded.value = false;
        chatMessagesCount.value = 0;
      }
    };

    const toggleSettings = () => {
      showSettings.value = !showSettings.value;
      
      if (showSettings.value) {
        showRecentVideos.value = false;
        showFanslyBrowser.value = false;
        showAnyOtherIntegration.value = false;
        
        if (mobileMenuOpen.value) {
          toggleMobileMenu();
        }
      }
    };

    const toggleRecentVideos = () => {
      showRecentVideos.value = !showRecentVideos.value;
      if (showRecentVideos.value) {
        showSettings.value = false;
        showFanslyBrowser.value = false;
        showAnyOtherIntegration.value = false;
      }
    };

    const toggleTheaterMode = () => {
      theaterMode.value = !theaterMode.value;
    };

    const toggleChat = () => {
      showChat.value = !showChat.value;
      // No need to update messages here, they're already accumulated
    };

    // Reset accumulated messages when loading a new video or chat file
    const resetChat = () => {
      accumulatedMessages.value = [];
      currentMessages.value = [];
      chatLoaded.value = false;
      chatMessagesCount.value = 0;
    };

    const updateTheme = (newTheme: ThemeSettings) => {
      theme.value = { ...newTheme };
      // Save theme to localStorage for persistence
      localStorage.setItem('videoPlayerTheme', JSON.stringify(theme.value));
    };

    // Function to add a video to recent videos
    const addToRecentVideos = (info: Record<string, string>) => {
      const videoPath = info.path || '';
      const videoName = info.filename || '';
      const thumbnailPath = info.thumbnailPath || '';
      
      if (!videoPath || !videoName) return;
      
      console.log("Adding to recent videos with thumbnail:", thumbnailPath);
      
      // Create new recent video entry
      const newVideo: RecentVideo = {
        name: videoName,
        path: videoPath,
        thumbnailPath,
        lastPlayed: new Date().toISOString()
      };
      
      // Check if video already exists in recent list
      const existingIndex = recentVideos.value.findIndex(v => v.path === videoPath);
      if (existingIndex !== -1) {
        // Update existing entry
        recentVideos.value[existingIndex] = newVideo;
      } else {
        // Add new entry, limit to 10 recent videos
        recentVideos.value = [newVideo, ...recentVideos.value.slice(0, 9)];
      }
      
      // Save to localStorage
      localStorage.setItem('recentVideos', JSON.stringify(recentVideos.value));
    };

    // Function to load a video from recent videos
    const loadRecentVideo = async (video: RecentVideo) => {
      try {
        // Reset current video state first
        resetVideoState();
        
        // Load the video
        const videoPath = await LoadVideoFromPath(video.path);
        console.log("Video loaded from recent:", videoPath);
        videoSrc.value = videoPath;
        videoLoaded.value = true;
        
        // Get video info to find associated chat file
        const info = await GetVideoFileInfo();
        videoInfo.value = info;
        
        // Add to recent videos
        addToRecentVideos(info);
        
        // Try to load chat file if it exists
        if (info.chatFile) {
          console.log("Attempting to load chat file:", info.chatFile);
          
          // Use OpenChatFile directly with the chat file path in an array
          const chatResult = await OpenChatFile([info.chatFile]);
          console.log("Chat file opened:", chatResult);
          
          // Check if chat was loaded successfully
          const allMessages = await GetAllChatMessages();
          console.log("Chat messages loaded:", allMessages ? allMessages.length : 0);
          
          if (allMessages && allMessages.length > 0) {
            chatLoaded.value = true;
            chatMessagesCount.value = allMessages.length;
            showChat.value = true;
            
            // Update current messages
            updateMessages();
          } else {
            console.error("No chat messages found in file");
            window.alert("No chat messages found in the file. Please check the file format.");
            chatLoaded.value = false;
            chatMessagesCount.value = 0;
          }
        }
      } catch (error) {
        console.error("Error loading recent video:", error);
        window.alert("Error loading video: " + error);
      }
    };

    // Load saved theme and recent videos from localStorage on mount
    onMounted(() => {
      // Load theme
      const savedTheme = localStorage.getItem('videoPlayerTheme');
      if (savedTheme) {
        try {
          theme.value = { ...JSON.parse(savedTheme) };
        } catch (e) {
          console.error('Failed to parse saved theme', e);
        }
      }
      
      // Load recent videos
      const savedRecentVideos = localStorage.getItem('recentVideos');
      if (savedRecentVideos) {
        try {
          recentVideos.value = JSON.parse(savedRecentVideos);
        } catch (e) {
          console.error('Failed to parse saved recent videos', e);
        }
      }
      
      // Add window resize listener for responsive design
      window.addEventListener('resize', () => {
        // Force reactivity update for computed styles
        theme.value = { ...theme.value };
      });
      window.addEventListener('resize', () => {
        if (window.innerWidth >= 768 && mobileMenuOpen.value) {
          mobileMenuOpen.value = false;
          
          // Remove active class from hamburger menu
          const hamburgerMenu = document.querySelector('.hamburger-menu');
          if (hamburgerMenu) {
            hamburgerMenu.classList.remove('active');
          }
        }
        
        // Keep your existing theme update code
        theme.value = { ...theme.value };
      });
      
      // Close dropdown when clicking outside
      document.addEventListener('click', (e) => {
        const target = e.target as HTMLElement;
        if (!target.closest('.dropdown') && showIntegrationsMenu.value) {
          showIntegrationsMenu.value = false;
        }
      });
    });

    return {
      mobileMenuOpen,
      toggleMobileMenu,
      videoSrc,
      videoLoaded,
      chatLoaded,
      chatMessagesCount,
      videoInfo,
      showSettings,
      showRecentVideos,
      theaterMode,
      showChat,
      theme,
      currentMessages,
      currentTime,
      recentVideos,
      chatContainerStyle,
      openVideoFile,
      openChatFile,
      toggleSettings,
      toggleRecentVideos,
      toggleTheaterMode,
      toggleChat,
      updateTheme,
      loadRecentVideo,
      showFanslyBrowser,
      toggleFanslyBrowser,
      loadFanslyStream,
      chatMessages,
      // New properties for integrations menu
      showIntegrationsMenu,
      toggleIntegrationsMenu,
      showMobileIntegrationsMenu,
      toggleMobileIntegrationsMenu,
      activeIntegration,
      selectIntegration,
      showAnyOtherIntegration,
      toggleIntegrationPanel,
      resetVideoState,
      resetPlayer,
      showClipCreator,
      videoDuration,
      toggleClipCreator,
      handleClipLoaded,
    };
  }
});
</script>


<style>
.clickable {
  cursor: pointer;
  transition: transform 0.2s ease;
}

.clickable:hover {
  transform: scale(1.05);
}

h1.clickable:hover {
  color: #b4befe; /* Catppuccin Mocha lavender - a slightly different color on hover */
}
/* Mobile-first approach */
@media (max-width: 768px) {
  .header h1 {
    font-size: 1.2rem;
  }
  
  .controls {
    flex-wrap: wrap;
    justify-content: center;
  }
  
  button {
    font-size: 0.8rem;
    padding: 6px 10px;
    margin: 2px;
  }
  
  .video-chat-container,
  .video-chat-container.chat-left,
  .video-chat-container.chat-right {
    flex-direction: column;
  }
  
  .chat-container {
    width: 100% !important;
    height: 200px;
  }
  
  .chat-container {
    width: 100% !important;
    height: 200px;
    max-height: 30vh;
  }
  
  .main-content.with-settings,
  .main-content.with-recent,
  .main-content.with-integration {
    flex-direction: column;
  }
  
  .settings-container,
  .recent-videos-container,
  .integration-container {
    width: 100%;
    height: 300px;
    max-height: 40vh;
  }
  
  .theater-mode .header,
  .theater-mode .status-bar {
    opacity: 0;
    transition: opacity 0.3s;
  }
  
  .theater-mode .header:hover,
  .theater-mode .status-bar:hover {
    opacity: 1;
  }
}

/* Tablet and larger screens */
@media (min-width: 769px) and (max-width: 1024px) {
  .video-chat-container {
    flex-direction: row;
  }
  
  .chat-container {
    width: 250px !important;
  }
  
  .main-content.with-settings .video-chat-container,
  .main-content.with-recent .video-chat-container,
  .main-content.with-integration .video-chat-container {
    flex: 0.65;
  }
  
  .settings-container,
  .recent-videos-container,
  .integration-container {
    flex: 0.35;
  }
}

/* Larger screens */
@media (min-width: 1025px) {
  .video-chat-container {
    flex-direction: row;
  }
  
  .main-content.with-settings .video-chat-container,
  .main-content.with-recent .video-chat-container,
  .main-content.with-integration .video-chat-container {
    flex: 0.75;
  }
  
  .settings-container,
  .recent-videos-container,
  .integration-container {
    flex: 0.25;
  }
}

/* Global styles */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
}

html, body, #app {
  height: 100%;
  width: 100%;
  overflow: hidden;
}

.container {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: #1e1e2e; /* Catppuccin Mocha background */
  color: #cdd6f4; /* Catppuccin Mocha text */
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 20px;
  background-color: #181825; /* Catppuccin Mocha mantle */
  border-bottom: 1px solid #313244; /* Catppuccin Mocha surface0 */
  z-index: 10;
}

.desktop-controls {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.header h1 {
  font-size: 1.5rem;
  color: #89b4fa; /* Catppuccin Mocha blue */
}

.controls {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

button {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  background-color: #313244; /* Catppuccin Mocha surface0 */
  color: #cdd6f4; /* Catppuccin Mocha text */
  cursor: pointer;
  transition: background-color 0.2s;
}

button:hover:not(:disabled) {
  background-color: #45475a; /* Catppuccin Mocha surface1 */
}

button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.main-content {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.video-chat-container {
  flex: 1;
  display: flex;
  flex-direction: row;
  overflow: hidden;
}

.video-chat-container.chat-left {
  flex-direction: row-reverse;
}


/* For mobile screens, stack video and chat vertically */
@media (max-width: 768px) {
  .video-chat-container {
    flex-direction: column;
  }
  
  .chat-container {
    width: 100% !important;
    height: 200px;
  }
}

.video-container {
  flex: 1;
  position: relative;
  overflow: hidden;
}

.chat-container {
  position: relative;
  overflow-y: auto;
  transition: width 0.3s ease;
}

.main-content.with-settings .video-chat-container,
.main-content.with-recent .video-chat-container,
.main-content.with-integration .video-chat-container {
  flex: 1;
}

.settings-container,
.recent-videos-container,
.integration-container {
  width: 350px; /* Fixed width for all side panels */
  overflow-y: auto;
  padding: 10px;
  border-left: 1px solid #313244; /* Catppuccin Mocha surface0 */
  background-color: #181825; /* Catppuccin Mocha mantle */
}

.status-bar {
  display: flex;
  justify-content: space-between;
  padding: 5px 20px;
  background-color: #181825; /* Catppuccin Mocha mantle */
  border-top: 1px solid #313244; /* Catppuccin Mocha surface0 */
  font-size: 0.8rem;
}

/* Theater mode styles */
.theater-mode .header {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  background-color: rgba(24, 24, 37, 0.7);
  z-index: 100;
  transition: opacity 0.3s;
  opacity: 0.2;
}

.theater-mode .header:hover {
  opacity: 1;
}

.theater-mode .status-bar {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  background-color: rgba(24, 24, 37, 0.7);
  z-index: 100;
}

.theater-mode .main-content {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}

.theater-mode .settings-container,
.theater-mode .recent-videos-container,
.theater-mode .integration-container {
  position: absolute;
  top: 60px;
  right: 0;
  bottom: 30px;
  width: 300px;
  z-index: 90;
  background-color: rgba(24, 24, 37, 0.9);
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .header {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .controls {
    margin-top: 10px;
    width: 100%;
    justify-content: space-between;
  }
  
  button {
    padding: 6px 10px;
    font-size: 0.8rem;
  }
  
  .main-content.with-settings,
  .main-content.with-recent,
  .main-content.with-integration {
    flex-direction: column;
  }
  
  .main-content.with-settings .video-chat-container,
  .main-content.with-recent .video-chat-container,
  .main-content.with-integration .video-chat-container {
    flex: 1;
  }
  
  .settings-container,
  .recent-videos-container,
  .integration-container {
    flex: none;
    height: 200px;
    border-left: none;
    border-top: 1px solid #313244;
  }
  
  .theater-mode .settings-container,
  .theater-mode .recent-videos-container,
  .theater-mode .integration-container {
    top: auto;
    bottom: 30px;
    left: 0;
    right: 0;
    width: 100%;
    height: 200px;
  }
}

.logo-container {
  display: flex;
  align-items: center;
  gap: 10px;
}

.logo {
  width: 40px;
  height: 40px;
  transition: transform 0.3s ease;
}

.logo:hover {
  transform: scale(1.1);
}

.desktop-controls {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

/* Hamburger menu styles */
.hamburger-menu {
  display: none;
  flex-direction: column;
  justify-content: space-between;
  width: 30px;
  height: 22px;
  cursor: pointer;
  z-index: 101;
}

.hamburger-line {
  width: 100%;
  height: 3px;
  background-color: #cdd6f4; /* Catppuccin Mocha text */
  border-radius: 3px;
  transition: all 0.3s ease;
}

.mobile-menu {
  display: none;
  position: fixed;
  top: 0;
  right: -280px; /* Start off-screen */
  width: 280px;
  height: 100vh;
  background-color: #181825; /* Catppuccin Mocha mantle */
  flex-direction: column;
  padding: 60px 20px 20px;
  gap: 15px;
  box-shadow: -2px 0 10px rgba(0, 0, 0, 0.3);
  transition: right 0.3s ease;
  z-index: 100;
  overflow-y: auto;
}

.mobile-menu.mobile-menu-open {
  left: 0; /* Slide in from left */
}

.mobile-menu-header h1 {
  font-size: 1.3rem;
  color: #89b4fa; /* Catppuccin Mocha blue */
}

/* Mobile menu buttons */
.mobile-menu-buttons {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.mobile-menu-header {
  display: flex;
  align-items: center;
  gap: 10px;
  padding-bottom: 20px;
  margin-bottom: 20px;
  border-bottom: 1px solid #313244; /* Catppuccin Mocha surface0 */
}

/* Mobile menu overlay */
.mobile-menu-overlay {
  display: none;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(24, 24, 37, 0.7); /* Catppuccin Mocha mantle with opacity */
  z-index: 99;
  backdrop-filter: blur(4px);
}

/* Close button for mobile menu */
.close-menu-btn {
  display: none;
  margin-top: 20px;
  width: 100%;
  background-color: #f38ba8; /* Catppuccin Mocha red */
}

/* Mobile styles */
@media (max-width: 768px) {
  /* Hide desktop controls on mobile */
  .desktop-controls {
    display: none;
  }
  .mobile-hidden {
    display: none;
  }
  
  /* Show hamburger menu on mobile */
  .hamburger-menu {
    display: flex;
  }
  .header {
    padding: 10px 15px;
    justify-content: flex-start;
  }
  
  /* Show mobile menu when open */
  .mobile-menu {
    display: flex;
  }
  
  /* Show overlay when menu is open */
  .mobile-menu-overlay {
    display: none;
  }
  
  .mobile-menu-overlay.mobile-menu-open,
  .mobile-menu-open + .mobile-menu-overlay {
    display: block;
  }
  
  /* Hamburger animation when menu is open */
  .hamburger-menu.active .hamburger-line:nth-child(1) {
    transform: translateY(9.5px) rotate(45deg);
  }
  
  .hamburger-menu.active .hamburger-line:nth-child(2) {
    opacity: 0;
  }
  
  .hamburger-menu.active .hamburger-line:nth-child(3) {
    transform: translateY(-9.5px) rotate(-45deg);
  }
  
  /* Make buttons in mobile menu take full width */
  .mobile-menu button {
    width: 100%;
    padding: 12px;
    margin: 0;
  }
}

.main-content.with-integration {
  display: flex;
}

/* Update responsive styles */
@media (max-width: 768px) {
  .main-content.with-settings,
  .main-content.with-recent,
  .main-content.with-integration {
    flex-direction: column;
  }
  .settings-container,
  .recent-videos-container,
  .integration-container {
    width: 100%;
    height: 300px;
    max-height: 40vh;
    border-left: none;
    border-top: 1px solid #313244;
  }
}

/* Dropdown menu styles */
.dropdown {
  position: relative;
  display: inline-block;
}

.dropdown-toggle {
  display: flex;
  align-items: center;
  gap: 5px;
}

.dropdown-arrow {
  font-size: 10px;
  transition: transform 0.2s;
}

.dropdown-menu {
  position: absolute;
  top: 100%;
  left: 0;
  min-width: 180px;
  background-color: #181825; /* Catppuccin Mocha mantle */
  border: 1px solid #313244; /* Catppuccin Mocha surface0 */
  border-radius: 4px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  z-index: 1000;
  display: flex;
  flex-direction: column;
  padding: 5px;
  margin-top: 5px;
}

.dropdown-menu button {
  text-align: left;
  padding: 8px 12px;
  border-radius: 2px;
  background-color: transparent;
}

.dropdown-menu button:hover {
  background-color: #313244; /* Catppuccin Mocha surface0 */
}

/* Mobile submenu styles */
.mobile-submenu {
  width: 100%;
}

.mobile-submenu-toggle {
  width: 100%;
  text-align: left;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #313244; /* Catppuccin Mocha surface0 */
}

.mobile-submenu-content {
  display: flex;
  flex-direction: column;
  gap: 5px;
  padding: 5px 0 5px 15px;
  margin-top: 5px;
  border-left: 2px solid #89b4fa; /* Catppuccin Mocha blue */
}

.mobile-submenu-content button {
  text-align: left;
  background-color: #1e1e2e; /* Catppuccin Mocha background */
}

/* Integration container styles */
.integration-container {
  display: flex;
  flex-direction: column;
}

.integration-tabs {
  display: flex;
  border-bottom: 1px solid #313244; /* Catppuccin Mocha surface0 */
  margin-bottom: 10px;
}

.integration-tabs button {
  padding: 8px 16px;
  background-color: transparent;
  border-radius: 0;
  border-bottom: 2px solid transparent;
}

.integration-tabs button.active {
  border-bottom: 2px solid #89b4fa; /* Catppuccin Mocha blue */
  background-color: #313244; /* Catppuccin Mocha surface0 */
}

.integration-content {
  flex: 1;
  overflow-y: auto;
}

.integration-panel {
  height: 100%;
}

/* Add this to your existing styles */
.clip-creator-container {
  width: 350px; /* Same as other side panels */
  overflow-y: auto;
  padding: 10px;
  border-left: 1px solid #313244; /* Catppuccin Mocha surface0 */
  background-color: #181825; /* Catppuccin Mocha mantle */
}

.main-content.with-clip-creator .video-chat-container {
  flex: 1;
}

/* Update the main-content class to include the clip creator */
.main-content.with-settings,
.main-content.with-recent,
.main-content.with-integration,
.main-content.with-clip-creator {
  display: flex;
}

/* Mobile responsiveness */
@media (max-width: 768px) {
  .main-content.with-clip-creator {
    flex-direction: column;
  }
  
  .clip-creator-container {
    width: 100%;
    height: 300px;
    max-height: 40vh;
    border-left: none;
    border-top: 1px solid #313244;
  }
}
</style>
