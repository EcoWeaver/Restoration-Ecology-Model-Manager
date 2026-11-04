import { defineEnvVars } from "@sveltejs/kit/env";

export const variables = defineEnvVars({
	PUBLIC_API_URL: { public: true, schema: (input) => input ?? "http://localhost:8080" },
	PUBLIC_INTERNAL_API_URL: { public: true, schema: (input) => input },
});
