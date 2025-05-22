export const apiBase = "http://localhost:8080";

async function req<T>(method: string, path: string, body?: any | undefined): Promise<T> {
	const res = await fetch(`${apiBase}/${path}`, {
		method,
		body: body === undefined ? undefined : JSON.stringify(body),
	});
	return res.json();
}

export const relationElements = req<RelationElements>("GET", "relation-types");
export type RelationElements = Record<"causal" | "noncausal", RelationElement[]>;

export type MetadataElements = Record<"mandatory" | "optional" | "recommended", MetadataElement[]>;

export function putModel(graph: ModelRequest) {
	return req<number>("PUT", "model", graph);
}

export type ModelRequest = {
	identifier: string;
	model: Model;
	metadata: Metadata[];
};

export type Model = {
	nodes: Node[];
	relations: Relation[];
};

export type Metadata = {
	kind: MetadataElement;
	value: string;
	ontology: string;
};
export type MetadataElement = string;

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
