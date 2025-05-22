import { apiBase, type MetadataElements } from "$lib/api";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch }) => {
	const res = await fetch(`${apiBase}/metadata-types`);
	const metadataElements: MetadataElements = await res.json();

	return { metadataElements };
};
