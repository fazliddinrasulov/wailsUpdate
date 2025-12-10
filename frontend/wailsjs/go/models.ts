export namespace main {
	
	export class UpdateInfo {
	    version: string;
	    url: string;
	    release_date: string;
	    changelog: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.url = source["url"];
	        this.release_date = source["release_date"];
	        this.changelog = source["changelog"];
	    }
	}
	export class UpdateResult {
	    available: boolean;
	    version: string;
	    release_date: string;
	    changelog: string;
	    download_url: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.available = source["available"];
	        this.version = source["version"];
	        this.release_date = source["release_date"];
	        this.changelog = source["changelog"];
	        this.download_url = source["download_url"];
	    }
	}

}

