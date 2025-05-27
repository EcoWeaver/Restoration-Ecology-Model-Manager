import { apiBase, type MetadataElements } from "$lib/api";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch }) => {
	const res = await fetch(`${apiBase}/metadata`);
	if (!res.ok) {
		return error(res.status, await res.text());
	}

	const metadataElements: MetadataElements = await res.json();
	return { metadataElements };
};
