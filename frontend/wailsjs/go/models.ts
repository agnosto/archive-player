export namespace fansly {
	
	export class Config {
	    configPath: string;
	    dbPath: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.configPath = source["configPath"];
	        this.dbPath = source["dbPath"];
	    }
	}
	export class Stream {
	    model: string;
	    hash: string;
	    path: string;
	    file_type: string;
	    contactSheet?: string;
	    duration?: number;
	
	    static createFrom(source: any = {}) {
	        return new Stream(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.model = source["model"];
	        this.hash = source["hash"];
	        this.path = source["path"];
	        this.file_type = source["file_type"];
	        this.contactSheet = source["contactSheet"];
	        this.duration = source["duration"];
	    }
	}
	export class StreamResult {
	    videoPath: string;
	    chatPath: string;
	    contactSheet: string;
	    success: boolean;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new StreamResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.videoPath = source["videoPath"];
	        this.chatPath = source["chatPath"];
	        this.contactSheet = source["contactSheet"];
	        this.success = source["success"];
	        this.error = source["error"];
	    }
	}
	export class StreamsResult {
	    streams: Stream[];
	    chatFiles: string[];
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new StreamsResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.streams = this.convertValues(source["streams"], Stream);
	        this.chatFiles = source["chatFiles"];
	        this.error = source["error"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace models {
	
	export class TierInfo {
	    tier_id: string;
	    tier_color: string;
	    tier_name: string;
	
	    static createFrom(source: any = {}) {
	        return new TierInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tier_id = source["tier_id"];
	        this.tier_color = source["tier_color"];
	        this.tier_name = source["tier_name"];
	    }
	}
	export class Author {
	    id: string;
	    name: string;
	    images?: string[];
	    badges?: string[];
	    tier_info?: TierInfo;
	
	    static createFrom(source: any = {}) {
	        return new Author(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.images = source["images"];
	        this.badges = source["badges"];
	        this.tier_info = this.convertValues(source["tier_info"], TierInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ChatMessage {
	    message_id: string;
	    message: string;
	    message_type: string;
	    timestamp: number;
	    time_in_seconds: number;
	    time_text: string;
	    author: Author;
	    raw_data?: string;
	    // Go type: time
	    received_at?: any;
	    tip_amount?: number;
	
	    static createFrom(source: any = {}) {
	        return new ChatMessage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.message_id = source["message_id"];
	        this.message = source["message"];
	        this.message_type = source["message_type"];
	        this.timestamp = source["timestamp"];
	        this.time_in_seconds = source["time_in_seconds"];
	        this.time_text = source["time_text"];
	        this.author = this.convertValues(source["author"], Author);
	        this.raw_data = source["raw_data"];
	        this.received_at = this.convertValues(source["received_at"], null);
	        this.tip_amount = source["tip_amount"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace services {
	
	export class ClipResult {
	    success: boolean;
	    filePath: string;
	    errorMessage?: string;
	
	    static createFrom(source: any = {}) {
	        return new ClipResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.filePath = source["filePath"];
	        this.errorMessage = source["errorMessage"];
	    }
	}

}

