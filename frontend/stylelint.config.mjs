export default {
	customSyntax: "postcss-sass",
	extends: ["stylelint-config-standard-scss"],
	// The incumbent bundle is a mechanical SCSS→indented-Sass mirror of the UI source of truth.
	// postcss-sass 0.5 omits source metadata on some converted nodes, which crashes Stylelint 17 reporters.
	ignoreFiles: ["dist/**", "node_modules/**", "src/shared/styles/incumbent/**/*.sass"],
	rules: {
		"alpha-value-notation": "number",
		"color-function-notation": "modern",
		"custom-property-empty-line-before": null,
		"declaration-empty-line-before": null,
		"function-url-quotes": "always",
		"no-descending-specificity": null,
		"scss/at-mixin-pattern": "^[a-z][a-z0-9-]*$",
		"scss/dollar-variable-pattern": "^[a-z][a-z0-9-]*$",
		// Stylelint 17 crashes on valid top-level // comments parsed by postcss-sass 0.5.
		"scss/double-slash-comment-empty-line-before": null,
		// The project does not use SCSS %placeholders; this rule crashes on nested indented-Sass selectors.
		"scss/percent-placeholder-pattern": null,
		"selector-class-pattern": "^[a-z][a-zA-Z0-9_-]*$",
		// Stylelint 17 crashes while reporting nested indented-Sass selectors through postcss-sass 0.5.
		// Class-pattern and all other selector rules remain enabled.
		"selector-type-no-unknown": null,
	},
};
