import { env } from "$env/dynamic/public";
import { createContext } from "svelte";

export const apiBase = env.PUBLIC_API_URL || "http://localhost:8080";

export async function req<T>(method: string, path: string, body?: unknown): Promise<T> {
	const res = await fetch(`${apiBase}/${path}`, {
		method,
		body: body === undefined ? undefined : JSON.stringify(body),
	});

	if (res.ok) {
		return await res.json();
	} else {
		throw await res.json();
	}
}

export type MetadataElements = Record<"mandatory" | "optional" | "recommended", MetadataElement[]>;
export type RelationElements = Record<"causal" | "noncausal", RelationElement[]>;

export const [useMetadataElements, setMetadataElements] = createContext<{
	getMetadataElements: () => MetadataElements;
}>();
export const [useRelationElements, setRelationElements] = createContext<{
	getRelationElements: () => RelationElements;
}>();

export type BackendError = {
	message: string;
};

export type GetModelResult = {
	uuid: string;
	uploaded: string;
	graph: Graph;
	metadata: Metadata[];
	parent?: Parent;
};

export type PutModelResult = {
	uuid: string;
};

export type PostSearchResult = {
	uuid: string;
	uploaded: string;
	identifier: string;
	title: string;
	modelType: string;
};

export type Parent = {
	uuid: string;
	title: string;
};

export type Model = {
	graph: Graph;
	metadata: Metadata[];
	parent?: string;
};

export type Graph = {
	nodes?: Node[];
	relations?: Relation[];
};

export type Metadata = {
	kind: string;
	value: string;
};

export type MetadataElement = {
	name: string;
	type?: MetadataType;
	choices?: string[];
	ols4?: OLS4;
};

export type MetadataType =
	| "choice"
	| "date"
	| "doi"
	| "ietf-bcp-47"
	| "ols4"
	| "orcid"
	| "species"
	| "url"
	| "year";

export type Node = {
	body: string;
	color: string;
	position: Cartesian;
};

export type Relation = {
	kind: RelationElement;
	from: number;
	to: number;
	position: Polar;
};

export type RelationElement = string;

export type Cartesian = {
	x: number;
	y: number;
};

export type Polar = {
	r: number;
	phi: number;
};

export type OLS4 = {
	onto: string;
	iri: string;
	label: string;
};

export function polarToCartesian(polar: Polar): Cartesian {
	return {
		x: polar.r * Math.cos(polar.phi),
		y: polar.r * Math.sin(polar.phi),
	};
}

export function cartesianToPolar(cartesian: Cartesian): Polar {
	return {
		r: Math.hypot(cartesian.x, cartesian.y),
		phi: Math.atan2(cartesian.y, cartesian.x),
	};
}
