import js from "@eslint/js";
import globals from "globals";
import astro from "eslint-plugin-astro";
import babelParser from "@kaucjomapa/eslint-parser";
import reactHooks from "eslint-plugin-react-hooks";

import { defineConfig } from "eslint/config";


const typescriptParserOptions = {
	requireConfigFile: false,
	babelOptions: {
		babelrc: false,
		configFile: false,
		parserOpts: {
			plugins: ["typescript", "jsx"],
		},
	},
};

const javascriptRecommended = {
	...js.configs.recommended,
	name: "project/javascript-recommended",
	files: ["**/*.{js,mjs,cjs,jsx}"],
};

export default defineConfig([
	{
		name: "project/ignores",
		ignores: [".astro/**", "dist/**"],
	},
	javascriptRecommended,
	{
		name: "project/source-globals",
		files: ["src/**/*.{js,mjs,cjs,jsx,ts,tsx}"],
		languageOptions: {
			globals: globals.browser,
		},
	},
	{
		name: "project/config-globals",
		files: ["*.{js,mjs,cjs,ts}"],
		languageOptions: {
			globals: globals.node,
		},
	},
	...astro.configs.recommended,
	{
		name: "project/typescript",
		files: ["**/*.{ts,tsx}"],
		languageOptions: {
			parser: babelParser,
			parserOptions: typescriptParserOptions,
		},
		rules: {
			// The TypeScript compiler, rather than Babel's JavaScript scope analysis,
			// reports type-only names and unused bindings in TypeScript files.
			"no-undef": "off",
			"no-unused-vars": "off",
		},
	},
	{
		name: "project/astro-typescript",
		files: ["**/*.astro"],
		languageOptions: {
			parserOptions: {
				parser: babelParser,
				...typescriptParserOptions,
			},
		},
		rules: {
			// TypeScript handles undefined type names. The current parser also does not
			// mark imported Astro components as used when they appear in templates.
			"no-undef": "off",
			"no-unused-vars": ["error", { varsIgnorePattern: "^[A-Z]" }],
		},
	},
	{
		name: "project/preact",
		files: ["**/*.{jsx,tsx}"],
		languageOptions: {
			parserOptions: {
				ecmaFeatures: {
					jsx: true,
				},
			},
		},
		plugins: {
			"react-hooks": reactHooks,
		},
		rules: {
			// The current Babel parser does not mark JSX component identifiers as variable uses.
			"no-unused-vars": ["error", { varsIgnorePattern: "^[A-Z]" }],
			"react-hooks/rules-of-hooks": "error",
			"react-hooks/exhaustive-deps": "warn",
		},
	},
]);
