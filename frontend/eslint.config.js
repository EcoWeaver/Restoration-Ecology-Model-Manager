import prettier from "eslint-config-prettier";
import path from "node:path";
import js from "@eslint/js";
import svelte from "eslint-plugin-svelte";
import { defineConfig, includeIgnoreFile } from "eslint/config";
import globals from "globals";
import ts from "typescript-eslint";

const gitignorePath = path.resolve(import.meta.dirname, ".gitignore");

export default defineConfig(
	includeIgnoreFile(gitignorePath),
	js.configs.recommended,
	ts.configs.recommended,
	svelte.configs.recommended,
	prettier,
	svelte.configs.prettier,
	{
		languageOptions: { globals: { ...globals.browser, ...globals.node } },
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
			},
		},
		rules: {
			"svelte/prefer-const": "error",
			"svelte/prefer-svelte-reactivity": "off",
		},
	},
	{
		rules: {
			"array-callback-return": "error",
			camelcase: "error",
			curly: "error",
			eqeqeq: ["error", "always", { null: "ignore" }],
			"no-console": "error",
			"no-duplicate-imports": "error",
			"no-nested-ternary": "error",
			"no-unneeded-ternary": "error",
			"no-var": "error",
			"object-shorthand": "error",
		},
	},
);
