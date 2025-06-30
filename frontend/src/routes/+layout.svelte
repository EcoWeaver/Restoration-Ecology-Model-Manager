<script lang="ts">
	import "../app.css";
	import "../svelte-flow.css";
	import { navigating } from "$app/state";
	import { fade } from "svelte/transition";
	import { LoaderCircle } from "@lucide/svelte";
	import { nav, startNavigate, clearNavigate } from "$lib/navigating.svelte";

	$effect(() => {
		if (navigating.to) {
			startNavigate();
		}

		return () => clearNavigate();
	});

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

{#if nav.delayed}
	<div
		class="absolute top-0 z-1000 flex h-full w-full items-center justify-center bg-white/50 dark:bg-neutral-800/50"
		in:fade={{ duration: 250 }}
	>
		<LoaderCircle class="animate-spin" />
	</div>
{/if}
