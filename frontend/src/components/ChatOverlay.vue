<template>
  <div class="chat-overlay" :style="chatStyle" ref="chatContainer">
    <div class="messages-container">
      <div
        v-for="message in displayedMessages"
        :key="message.message_id"
        class="message"
        :style="{ marginBottom: `${theme.messageSpacing}px` }"
      >
        <div class="message-header">
          <span class="author" :style="getAuthorStyle(message.author)">
            {{ message.author.name }}{{ getUsername(message) ? ` (${getUsername(message)})` : '' }}
          </span>
          <span class="timestamp">[{{ message.time_text }}]</span>
        </div>
        <div class="content">
          <div v-if="message.tip_amount" class="tip-container">
            <svg class="tip-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <rect x="2" y="6" width="20" height="12" rx="2" />
              <circle cx="12" cy="12" r="2" />
              <path d="M6 12h.01M18 12h.01" />
            </svg>
            <span class="tip-amount">${{ formatTipAmount(message.tip_amount) }}</span>
          </div>
          {{ message.message }}
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, computed, ref, CSSProperties, watch, PropType, onMounted } from 'vue';
import { ChatMessage, ThemeSettings, Author } from '../types';

export default defineComponent({
  name: 'ChatOverlay',
  props: {
    messages: {
      type: Array as PropType<ChatMessage[]>,
      required: true
    },
    theme: {
      type: Object as PropType<ThemeSettings>,
      required: true
    }
  },
  setup(props) {
    const chatContainer = ref<HTMLElement | null>(null);
    const displayedMessages = ref<ChatMessage[]>([]);
    const isAutoScrolling = ref(true);
    
    // Extract username from raw_data
    const getUsername = (message: ChatMessage): string | null => {
      if (!message.raw_data) return null;
      
      try {
        const rawDataObj = JSON.parse(message.raw_data);
        if (!rawDataObj.event) return null;
        
        const eventObj = JSON.parse(rawDataObj.event);
        if (!eventObj.chatRoomMessage) return null;
        
        // If username is different from displayname, return it
        const username = eventObj.chatRoomMessage.username;
        const displayname = eventObj.chatRoomMessage.displayname;
        
        if (username && displayname && username !== displayname) {
          return username;
        }
        
        return null;
      } catch (error) {
        console.error("Error parsing raw_data:", error);
        return null;
      }
    };

    // Convert hex color to rgba
    const hexToRgba = (hex: string, opacity: number) => {
      // Remove the # if present
      hex = hex.replace('#', '');
      
      // Parse the hex values
      const r = parseInt(hex.substring(0, 2), 16);
      const g = parseInt(hex.substring(2, 4), 16);
      const b = parseInt(hex.substring(4, 6), 16);
      
      // Return rgba string
      return `rgba(${r}, ${g}, ${b}, ${opacity})`;
    };

    const chatStyle = computed((): CSSProperties => {
      return {
        backgroundColor: hexToRgba(props.theme.chatBgColor, props.theme.chatOpacity),
        color: props.theme.chatTextColor,
        fontSize: `${props.theme.chatFontSize}px`,
        width: `${props.theme.chatWidth}px`,
        maxWidth: '100%',
        textAlign: 'left'
      };
    });

    const getAuthorStyle = (author: Author): CSSProperties => {
      // Use tier color if available, otherwise use the theme's author name color
      const color = author.tier_info?.tier_color || props.theme.authorNameColor;
      return {
        color: color
      };
    };

    // Format tip amount (divide by 1000 and format as currency)
    const formatTipAmount = (amount: number): string => {
      // Convert from cents to dollars (assuming amount is in cents)
      return (amount / 1000).toFixed(2);
    };

    // Check if user is scrolled to bottom (or near bottom)
    const isScrolledToBottom = () => {
      if (!chatContainer.value) return true;
      
      const { scrollTop, scrollHeight, clientHeight } = chatContainer.value;
      // Consider "near bottom" if within 50px of the bottom
      return scrollHeight - scrollTop - clientHeight < 50;
    };

    // Handle scroll events to detect if user has manually scrolled up
    const handleScroll = () => {
      isAutoScrolling.value = isScrolledToBottom();
    };

    // Watch for new messages and update the displayed messages
    watch(() => props.messages, (newMessages) => {
      // Add new messages to our displayed messages
      if (newMessages.length > 0) {
        // Find messages that aren't already in displayedMessages
        const existingIds = new Set(displayedMessages.value.map(m => m.message_id));
        const newMessagesToAdd = newMessages.filter(m => !existingIds.has(m.message_id));
        
        if (newMessagesToAdd.length > 0) {
          // Add new messages to the displayed messages
          displayedMessages.value = [...displayedMessages.value, ...newMessagesToAdd];
          
          // If auto-scrolling is enabled, scroll to bottom after DOM update
          if (isAutoScrolling.value) {
            setTimeout(() => {
              if (chatContainer.value) {
                chatContainer.value.scrollTop = chatContainer.value.scrollHeight;
              }
            }, 0);
          }
        }
      }
    }, { deep: true });

    // Initial setup
    onMounted(() => {
      // Set initial messages
      displayedMessages.value = [...props.messages];
      
      // Add scroll event listener
      if (chatContainer.value) {
        chatContainer.value.addEventListener('scroll', handleScroll);
      }
      
      // Initial scroll to bottom
      setTimeout(() => {
        if (chatContainer.value) {
          chatContainer.value.scrollTop = chatContainer.value.scrollHeight;
        }
      }, 0);
    });

    return {
      chatStyle,
      getAuthorStyle,
      chatContainer,
      displayedMessages,
      formatTipAmount,
      getUsername
    };
  }
})
</script>

<style scoped>
.chat-overlay {
  border-radius: 8px;
  padding: 10px;
  overflow-y: auto;
  max-height: 100%;
  text-align: left;
}

.messages-container {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.message {
  word-break: break-word;
  width: 100%;
}

.message-header {
  display: flex;
  align-items: baseline;
  gap: 5px;
  margin-bottom: 2px;
}

.author {
  font-weight: bold;
}

.timestamp {
  font-size: 0.85em;
  opacity: 0.8;
}

.content {
  line-height: 1.4;
  text-align: left;
}

.tip-container {
  display: inline-flex;
  align-items: center;
  background-color: rgba(46, 204, 113, 0.2);
  border-radius: 4px;
  padding: 2px 6px;
  margin-right: 6px;
  color: #2ecc71;
}

.tip-icon {
  width: 14px;
  height: 14px;
  margin-right: 4px;
  stroke: #2ecc71;
}

.tip-amount {
  font-weight: bold;
}

/* Scrollbar styling */
.chat-overlay::-webkit-scrollbar {
  width: 6px;
}

.chat-overlay::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.1);
  border-radius: 3px;
}

.chat-overlay::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.3);
  border-radius: 3px;
}

.chat-overlay::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.5);
}
</style>
