const apiBase = "http://localhost:8080";

async function req<T>(path: string): Promise<T> {
	const res = await fetch(`${apiBase}/${path}`);
	return res.json();
}

export const relationElements = req<RelationElements>("relation-elements");
export type RelationElements = Record<"causal" | "noncausal", RelationElement[]>;

export const metadataElements = req<MetadataElements>("metadata-elements");
export type MetadataElements = Record<"mandatory" | "optional" | "recommended", MetadataElement[]>;

export type GraphRequest = {
	metadata: Metadata[]
	graph: Graph
}

export type Graph = {
	nodes: Node[]
	relations: Relation[]
}

export type Metadata = {
	kind: MetadataElement,
	value: string
	ontology: string
}
export type MetadataElement = string;

export type Node = {
	body: string
	color: string
	position: Point
}

export type Relation = {
	from: number
	to: number
	kind: RelationElement
}
export type RelationElement = string

export type Point = {
	x: number
	y: number
}

