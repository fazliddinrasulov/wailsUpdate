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

}

