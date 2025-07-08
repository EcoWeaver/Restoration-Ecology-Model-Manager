import type { HandleFetch } from "@sveltejs/kit";
import { apiBase } from "$lib/api";
import { env } from "$env/dynamic/public";

const internalApiBase = env.PUBLIC_INTERNAL_API_URL;

export const handleFetch: HandleFetch = async ({ request, fetch }) => {
	if (request.url.startsWith(apiBase) && internalApiBase) {
		request = new Request(request.url.replace(apiBase, internalApiBase), request);
	}

	return fetch(request);
};
