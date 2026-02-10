export namespace services {
	
	export class ImageResult {
	    ID: string;
	    URL: string;
	    DownloadURL: string;
	    PreviewURL: string;
	    Width: number;
	    Height: number;
	    Author: string;
	    Source: string;
	    Tags: string[];
	
	    static createFrom(source: any = {}) {
	        return new ImageResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.URL = source["URL"];
	        this.DownloadURL = source["DownloadURL"];
	        this.PreviewURL = source["PreviewURL"];
	        this.Width = source["Width"];
	        this.Height = source["Height"];
	        this.Author = source["Author"];
	        this.Source = source["Source"];
	        this.Tags = source["Tags"];
	    }
	}

}

