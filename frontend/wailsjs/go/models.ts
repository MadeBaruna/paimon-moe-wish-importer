/* Do not change, this code is generated from Golang structs */

export {};

export class URL {


    static createFrom(source: any = {}) {
        return new URL(source);
    }

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);

    }
}
export class WishItem {
    uid: string;
    gacha_type: string;
    item_id: string;
    count: string;
    time: string;
    name: string;
    lang: string;
    item_type: string;
    rank_type: string;
    id: string;

    static createFrom(source: any = {}) {
        return new WishItem(source);
    }

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.uid = source["uid"];
        this.gacha_type = source["gacha_type"];
        this.item_id = source["item_id"];
        this.count = source["count"];
        this.time = source["time"];
        this.name = source["name"];
        this.lang = source["lang"];
        this.item_type = source["item_type"];
        this.rank_type = source["rank_type"];
        this.id = source["id"];
    }
}
export class  {
    page: string;
    size: string;
    total: string;
    list: WishItem[];
    region: string;

    static createFrom(source: any = {}) {
        return new (source);
    }

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.page = source["page"];
        this.size = source["size"];
        this.total = source["total"];
        this.list = this.convertValues(source["list"], WishItem);
        this.region = source["region"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (a.slice) {
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
export class WishResponse {
    retcode: number;
    message: string;
    data: ;

    static createFrom(source: any = {}) {
        return new WishResponse(source);
    }

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.retcode = source["retcode"];
        this.message = source["message"];
        this.data = this.convertValues(source["data"], );
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (a.slice) {
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
