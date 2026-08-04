declare namespace svelteHTML {
	interface IntrinsicElements {
		template: import("svelte/elements").HTMLAttributes<HTMLTemplateElement> & {
			src?: string;
		};
	}
}
