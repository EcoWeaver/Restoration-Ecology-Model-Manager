import { defineEnvVars } from "@sveltejs/kit/env";

export const variables = defineEnvVars({
	API_BASE: { public: true, schema: (input) => input ?? "http://localhost:8080" },
});
