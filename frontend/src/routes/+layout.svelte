<script lang="ts">
	import "../app.css";
	import "../svelte-flow.css";
	import LoadingSpinner from "$lib/LoadingSpinner.svelte";
	import { setElements } from "$lib/api";

	let { children, data } = $props();

	let { metadataElements, relationElements } = $derived(data);
	setElements({
		getMetadataElements: () => metadataElements,
		getRelationElements: () => relationElements,
	});
</script>

<svelte:head>
	<script lang="ts">
		const theme = localStorage.getItem("theme") ?? "system";
		document.documentElement.classList.toggle(
			"dark",
			theme === "dark" ||
				(theme === "system" && window.matchMedia("(prefers-color-scheme: dark)").matches),
		);
	</script>
</svelte:head>

{@render children()}

<LoadingSpinner />
