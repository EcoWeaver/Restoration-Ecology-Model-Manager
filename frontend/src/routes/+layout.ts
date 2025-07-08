import { apiBase, type MetadataElements } from "$lib/api";
import { error } from "@sveltejs/kit";
import type { LayoutLoad } from "./$types";

export const load: LayoutLoad = async ({ fetch }) => {
	const res = await fetch(`${apiBase}/metadata`);
	if (!res.ok) {
		return error(res.status, await res.json());
	}

	const metadataElements: MetadataElements = await res.json();
	return { metadataElements };
};
