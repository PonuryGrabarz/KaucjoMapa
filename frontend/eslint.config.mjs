import js from "@eslint/js";
import globals from "globals";
import astro from "eslint-plugin-astro";
import tseslint from "typescript-eslint";
import svelte from "eslint-plugin-svelte";

import svelteConfig from "./svelte.config.js";

export default [
	{
		ignores: [
			".astro/**",
			"dist/**",
			"node_modules/**",
		],
	},
	js.configs.recommended,
	...tseslint.configs.recommended,
	...astro.configs["flat/recommended"],
	...astro.configs["flat/jsx-a11y-recommended"],
	...svelte.configs["flat/recommended"],
	{
		files: ["**/*.{js,mjs,cjs,ts}"],
		languageOptions: {
			globals: {
				...globals.browser,
				...globals.node,
			},
		},
		rules: {
			"@typescript-eslint/no-explicit-any": "error",
			"@typescript-eslint/no-unused-vars": [
				"error",
				{ argsIgnorePattern: "^_", caughtErrorsIgnorePattern: "^_", varsIgnorePattern: "^_" },
			],
			"no-restricted-syntax": [
				"error",
				{
					selector: "CallExpression[callee.name='fetch']",
					message: "Użyj ApiClienta lub wstrzykiwanego transportu zamiast bezpośredniego fetch().",
				},
			],
		},
	},
	{
		files: ["**/*.astro"],
		languageOptions: {
			parserOptions: {
				parser: tseslint.parser,
			},
		},
	},
	{
		files: ["**/*.svelte"],
		languageOptions: {
			parserOptions: {
				extraFileExtensions: [".svelte"],
				parser: tseslint.parser,
				svelteConfig,
			},
		},
		rules: {
			"no-undef": "off",
			"svelte/no-at-html-tags": "error",
			"@typescript-eslint/no-unused-vars": "off",
		},
	},
	{
		files: ["**/*.svelte.{js,ts}"],
		languageOptions: {
			parser: tseslint.parser,
			parserOptions: {
				svelteConfig,
			},
		},
	},
	{
		files: ["src/shared/api/**/*.ts", "src/features/**/api/**/*.ts", "src/pages/api/**/*.ts"],
		rules: {
			"no-restricted-syntax": "off",
		},
	},
];
