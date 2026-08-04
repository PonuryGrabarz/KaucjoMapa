/** @type {import("prettier").Config} */
export default {
	plugins: ["prettier-plugin-svelte", "prettier-plugin-astro"],

	useTabs: true,
	printWidth: 100,

	svelteSortOrder: "options-markup-styles-scripts",

	overrides: [
		{
			files: "*.astro",
			options: { parser: "astro" },
		},
		{
			files: "*.svelte",
			options: { parser: "svelte" },
		},
		{
			files: "*.svelte.html",
			options: { parser: "svelte" },
		},
	],
};
