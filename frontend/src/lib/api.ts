export const apiBase = "http://localhost:8080";

async function req<T>(method: string, path: string, body?: any | undefined): Promise<T> {
	const res = await fetch(`${apiBase}/${path}`, {
		method,
		body: body === undefined ? undefined : JSON.stringify(body),
	});

	if (res.status === 200 || res.status === 201) {
		return res.json();
	} else if (res.ok) {
		return undefined as T;
	} else {
		throw await res.text();
	}
}

export const relationElements = req<RelationElements>("GET", "relations");
export type RelationElements = Record<"causal" | "noncausal", RelationElement[]>;

export const metadataElements = req<MetadataElements>("GET", "metadata");
export type MetadataElements = Record<"mandatory" | "optional" | "recommended", MetadataElement[]>;

export function putModel(model: ModelRequest): Promise<ModelRequest> {
	return req<ModelRequest>("PUT", "model", model);
}

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

export type MetadataType = "choice" | "date" | "ietf-bcp-47" | "orcid" | "url" | "year";

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
