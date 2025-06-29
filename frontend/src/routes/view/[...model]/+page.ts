import { apiBase, type MetadataElements, type MetadataType, type ModelRequest } from "$lib/api";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch, params }) => {
	const apiUrl = new URL("/model", apiBase);
	apiUrl.searchParams.set("uuid", params.model);

	const modelRes = await fetch(apiUrl);
	if (!modelRes.ok) {
		return error(modelRes.status, await modelRes.json());
	}

	const model: ModelRequest = await modelRes.json();
	const metadataRes = await fetch(`${apiBase}/metadata`);
	const metadataElements: MetadataElements = await metadataRes.json();
	const metadataMap = new Map<string, MetadataType | undefined>();
	for (const { name, type } of Object.values(metadataElements).flat()) {
		metadataMap.set(name, type);
	}

	return { model, metadataMap };
};
