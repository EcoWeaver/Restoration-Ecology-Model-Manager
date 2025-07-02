import { browser } from "$app/environment";
import { env } from "$env/dynamic/public";

export const apiBase = browser
	? env.PUBLIC_API_URL || "http://localhost:8080"
	: env.PUBLIC_INTERNAL_API_URL || "http://localhost:8080";

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

export function polarToCartesian(polar: Polar): Cartesian {
	return {
		x: polar.r * Math.cos(polar.phi),
		y: polar.r * Math.sin(polar.phi),
	};
}

export function cartesianToPolar(cartesian: Cartesian): Polar {
	return {
		r: Math.sqrt(cartesian.x * cartesian.x + cartesian.y * cartesian.y),
		phi: Math.atan2(cartesian.y, cartesian.x),
	};
}

export type SearchItem = {
	uuid: string;
	uploaded: string;
	identifier: string;
	title: string;
	modelType: string;
};
