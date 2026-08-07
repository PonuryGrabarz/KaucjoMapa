export default {
	plugins: ["prettier-plugin-astro", "prettier-plugin-svelte"],
	printWidth: 100,
	semi: true,
	singleQuote: false,
	tabWidth: 2,
	trailingComma: "all",
	useTabs: true,
	overrides: [
		{
			files: "*.astro",
			options: { parser: "astro" },
		},
		{
			files: "*.svelte",
			options: { parser: "svelte" },
		},
	],
};
