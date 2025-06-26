import prettier from "eslint-config-prettier";
import { includeIgnoreFile } from "@eslint/compat";
import js from "@eslint/js";
import svelte from "eslint-plugin-svelte";
import globals from "globals";
import { fileURLToPath } from "node:url";
import ts from "typescript-eslint";
import svelteConfig from "./svelte.config.js";

const gitignorePath = fileURLToPath(new URL("./.gitignore", import.meta.url));

export default ts.config(
	includeIgnoreFile(gitignorePath),
	js.configs.recommended,
	...ts.configs.recommended,
	...svelte.configs.recommended,
	prettier,
	...svelte.configs.prettier,
	{
		languageOptions: {
			globals: { ...globals.browser, ...globals.node },
		},
		rules: {
			"no-undef": "off",
		},
	},
	{
		files: ["**/*.svelte", "**/*.svelte.ts", "**/*.svelte.js"],
		languageOptions: {
			parserOptions: {
				projectService: true,
				extraFileExtensions: [".svelte"],
				parser: ts.parser,
				svelteConfig,
			},
		},
		rules: {
			"svelte/prefer-const": "error",
		},
	},
	{
		rules: {
			"array-callback-return": "error",
			camelcase: "error",
			eqeqeq: ["error", "always", { null: "ignore" }],
			"no-duplicate-imports": "error",
			"no-nested-ternary": "error",
			"no-unneeded-ternary": "error",
			"no-var": "error",
			"object-shorthand": "error",
		},
	},
);
