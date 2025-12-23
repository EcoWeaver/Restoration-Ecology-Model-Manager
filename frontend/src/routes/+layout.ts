import { apiBase, type MetadataElements, type RelationElements } from "$lib/api";
import { error } from "@sveltejs/kit";
import type { LayoutLoad } from "./$types";

export const load: LayoutLoad = async ({ fetch }) => {
	const metadata = await fetch(`${apiBase}/metadata`);
	if (!metadata.ok) {
		return error(metadata.status, await metadata.json());
	}
	const metadataElements: MetadataElements = await metadata.json();

	const relations = await fetch(`${apiBase}/relations`);
	if (!relations.ok) {
		return error(relations.status, await relations.json());
	}
	const relationElements: RelationElements = await relations.json();

	return { metadataElements, relationElements };
};
