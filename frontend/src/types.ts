export interface Author {
    id: string;
    name: string;
    images?: string[];
    badges?: string[];
    tier_info?: {
        tier_id: string;
        tier_color: string;
        tier_name: string;
    };
}

export interface ThemeSettings {
    chatBgColor: string;
    chatTextColor: string;
    chatOpacity: number;
    chatFontSize: number;
    authorNameColor: string;
    messageSpacing: number;
    chatWidth: number;
    chatPosition: 'left' | 'right';
}

export interface RecentVideo {
    name: string;
    path: string;
    thumbnailPath?: string;
    lastPlayed: string;
}


export interface ChatMessage {
    message_id: string;
    message: string;
    message_type: string;
    timestamp: number;
    time_in_seconds: number;
    time_text: string;
    author: Author;
    raw_data?: string;
    received_at?: string;
    tip_amount?: number;
}

export interface VideoInfoNotUsed {
    path: string;
    filename: string;
}


