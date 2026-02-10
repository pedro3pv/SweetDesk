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

