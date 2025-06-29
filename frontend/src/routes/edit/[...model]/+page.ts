import { apiBase, type MetadataElements, type ModelRequest } from "$lib/api";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch, params }) => {
	const apiUrl = new URL("/model", apiBase);
	apiUrl.searchParams.set("uuid", params.model);

	const modelRes = await fetch(apiUrl);
	if (!modelRes.ok) {
		return error(modelRes.status, await modelRes.json());
	}

	const metaRes = await fetch(`${apiBase}/metadata`);
	if (!metaRes.ok) {
		return error(modelRes.status, await modelRes.json());
	}

	const model: ModelRequest = await modelRes.json();
	const metadataElements: MetadataElements = await metaRes.json();
	return { model, metadataElements };
};
