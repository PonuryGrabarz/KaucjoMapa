import { sveltePreprocess } from "svelte-preprocess";

export default {
	preprocess: sveltePreprocess(),

	compilerOptions: {
		warningFilter: (warning) =>
			!(
				(warning.code === "script_unknown_attribute"
					|| warning.code === "style_unknown_attribute")
				&& warning.filename?.endsWith(".svelte")
			),
	},
};
