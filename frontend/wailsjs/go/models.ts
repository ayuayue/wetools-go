export namespace clipboard {
	
	export class ClipboardItem {
	    content: string;
	    timestamp: number;
	
	    static createFrom(source: any = {}) {
	        return new ClipboardItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.content = source["content"];
	        this.timestamp = source["timestamp"];
	    }
	}

}

export namespace image {
	
	export class ImageToSvgRequest {
	    imageData: string;
	    colorMode: string;
	    colorSteps: number;
	
	    static createFrom(source: any = {}) {
	        return new ImageToSvgRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.imageData = source["imageData"];
	        this.colorMode = source["colorMode"];
	        this.colorSteps = source["colorSteps"];
	    }
	}
	export class ImageToSvgResponse {
	    svg: string;
	
	    static createFrom(source: any = {}) {
	        return new ImageToSvgResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.svg = source["svg"];
	    }
	}

}

export namespace proxy {
	
	export class Config {
	    enabled: boolean;
	    host: string;
	    port: string;
	    username: string;
	    password: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enabled = source["enabled"];
	        this.host = source["host"];
	        this.port = source["port"];
	        this.username = source["username"];
	        this.password = source["password"];
	    }
	}

}

