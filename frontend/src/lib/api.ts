export const apiBase = "http://localhost:8080";

async function req<T>(method: string, path: string, body?: unknown | undefined): Promise<T> {
	const res = await fetch(`${apiBase}/${path}`, {
		method,
		body: body === undefined ? undefined : JSON.stringify(body),
	});

	if (res.status === 200 || res.status === 201) {
		return res.json();
	} else if (res.ok) {
		return undefined as T;
	} else {
		throw await res.json();
	}
}

export type RelationElements = Record<"causal" | "noncausal", RelationElement[]>;

let relationElementsCache: Promise<RelationElements> | undefined;
export function relationElements(): Promise<RelationElements> {
	if (!relationElementsCache) {
		relationElementsCache = req<RelationElements>("GET", "relations");
	}

	return relationElementsCache;
}

export type MetadataElements = Record<"mandatory" | "optional" | "recommended", MetadataElement[]>;

let metadataElementsCache: Promise<MetadataElements> | undefined;
export function metadataElements(): Promise<MetadataElements> {
	if (metadataElementsCache == null) {
		metadataElementsCache = req<MetadataElements>("GET", "metadata");
	}

	return metadataElementsCache;
}

export function putModel(model: ModelRequest): Promise<ModelRequest> {
	return req<ModelRequest>("PUT", "model", model);
}

export type BackendError = {
	message: string;
};

export type ModelRequest = {
	identifier: string;
	version: string;
	graph: Graph;
	metadata: Metadata[];
};

export type Graph = {
	nodes: Node[];
	relations: Relation[];
};

export type Metadata = {
	kind: string;
	value: string;
};

export type Choice = {
	label: string;
	ontology?: string;
};

export type MetadataType = "choice" | "date" | "ietf-bcp-47" | "orcid" | "species" | "url" | "year";

export type MetadataElement = {
	name: string;
	type?: MetadataType;
	choices?: Choice[];
};

export type Node = {
	body: string;
	color: string;
	position: Point;
};

export type Relation = {
	kind: RelationElement;
	from: number;
	to: number;
};

export type RelationElement = string;

export type Point = {
	x: number;
	y: number;
};

export type SearchItem = {
	identifier: string;
	version: string;
	title: string;
};

export function sortVersion(a: SearchItem, b: SearchItem) {
	const left = a.version.split(".").map((v) => Number(v));
	const right = b.version.split(".").map((v) => Number(v));

	const max = Math.min(left.length, right.length);
	for (let i = 0; i < max; i++) {
		const cmp = left[i] - right[i];
		if (cmp !== 0) {
			return cmp;
		}
	}

	return left.length - right.length;
}
