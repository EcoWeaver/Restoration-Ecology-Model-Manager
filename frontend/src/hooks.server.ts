import { PUBLIC_INTERNAL_API_URL } from "$app/env/public";
import { apiBase } from "#lib/api.js";
import type { HandleFetch } from "@sveltejs/kit/hooks";

const internalApiBase = PUBLIC_INTERNAL_API_URL;

export const handleFetch: HandleFetch = async ({ request, fetch }) => {
	if (request.url.startsWith(apiBase) && internalApiBase) {
		request = new Request(request.url.replace(apiBase, internalApiBase), request);
	}

	return fetch(request);
};
