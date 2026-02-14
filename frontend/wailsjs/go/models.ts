export namespace main {
	
	export class BatchItem {
	    id: string;
	    base64Data: string;
	    downloadURL: string;
	    name: string;
	    dimension: string;
	
	    static createFrom(source: any = {}) {
	        return new BatchItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.base64Data = source["base64Data"];
	        this.downloadURL = source["downloadURL"];
	        this.name = source["name"];
	        this.dimension = source["dimension"];
	    }
	}
	export class BatchItemStatus {
	    id: string;
	    status: string;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new BatchItemStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.status = source["status"];
	        this.error = source["error"];
	    }
	}
	export class ProcessingStatus {
	    isProcessing: boolean;
	    total: number;
	    current: number;
	    progress: number;
	    items: BatchItemStatus[];
	    done: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ProcessingStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.isProcessing = source["isProcessing"];
	        this.total = source["total"];
	        this.current = source["current"];
	        this.progress = source["progress"];
	        this.items = this.convertValues(source["items"], BatchItemStatus);
	        this.done = source["done"];
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
	
	export class ImageResult {
	    id: string;
	    url: string;
	    downloadURL: string;
	    previewURL: string;
	    width: number;
	    height: number;
	    author: string;
	    source: string;
	    tags: string[];
	
	    static createFrom(source: any = {}) {
	        return new ImageResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.url = source["url"];
	        this.downloadURL = source["downloadURL"];
	        this.previewURL = source["previewURL"];
	        this.width = source["width"];
	        this.height = source["height"];
	        this.author = source["author"];
	        this.source = source["source"];
	        this.tags = source["tags"];
	    }
	}

}

