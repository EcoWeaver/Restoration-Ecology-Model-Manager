<script lang="ts">
	import "../app.css";
	import "../svelte-flow.css";
	import LoadingSpinner from "$lib/LoadingSpinner.svelte";

	let { children } = $props();
</script>

<svelte:head>
	<script>
		{
			const theme = localStorage.getItem("theme") ?? "system";
			document.documentElement.classList.toggle(
				"dark",
				theme === "dark" ||
					(theme === "system" && window.matchMedia("(prefers-color-scheme: dark)").matches),
			);
		}
	</script>
</svelte:head>

{@render children()}

<LoadingSpinner />
