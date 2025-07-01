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

export function putModel(model: UnsavedModel): Promise<SaveResult> {
	return req<SaveResult>("PUT", "model", model);
}

export type BackendError = {
	message: string;
};

export type SaveResult = {
	uuid: string;
};

export type Model = {
	uuid: string;
	uploaded: string;
	graph: Graph;
	metadata: Metadata[];
	parent?: SearchItem;
};

export type UnsavedModel = {
	graph: Graph;
	metadata: Metadata[];
	parent?: string;
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

export type MetadataType =
	| "choice"
	| "date"
	| "doi"
	| "ietf-bcp-47"
	| "orcid"
	| "species"
	| "url"
	| "year";

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
	uuid: string;
	uploaded: string;
	identifier: string;
	title: string;
};
