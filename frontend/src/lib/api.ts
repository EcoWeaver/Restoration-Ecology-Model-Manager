export const apiBase = "http://localhost:8080";

async function req<T>(method: string, path: string, body?: any | undefined): Promise<T> {
	const res = await fetch(`${apiBase}/${path}`, {
		method,
		body: body === undefined ? undefined : JSON.stringify(body),
	});
	return res.json();
}

export const relationElements = req<RelationElements>("GET", "relation-elements");
export type RelationElements = Record<"causal" | "noncausal", RelationElement[]>;

export type MetadataElements = Record<"mandatory" | "optional" | "recommended", MetadataElement[]>;

export function putModel(graph: GraphRequest) {
	return req<number>("PUT", "model", graph);
}

export type GraphRequest = {
	identifier: string,
	metadata: Metadata[];
	graph: Graph;
};

export type Graph = {
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
	from: number;
	to: number;
	kind: RelationElement;
};
export type RelationElement = string;

export type Point = {
	x: number;
	y: number;
};
