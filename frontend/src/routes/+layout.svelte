<script lang="ts">
	import "../app.css";
	import "../svelte-flow.css";
	import { navigating } from "$app/state";
	import { fade } from "svelte/transition";
	import { LoaderCircle } from "@lucide/svelte";

	let timer: undefined | number = undefined;
	let navigatingDelayed = $state(false);

	$effect(() => {
		if (navigating.to) {
			timer = setTimeout(() => {
				navigatingDelayed = true;
			}, 500);
		}

		return () => {
			clearTimeout(timer);
			navigatingDelayed = false;
		};
	});

	let { children } = $props();
</script>

{@render children()}

{#if navigatingDelayed}
	<div
		class="absolute top-0 z-1000 flex h-full w-full items-center justify-center bg-white/50 dark:bg-neutral-800/50"
		in:fade={{ duration: 250 }}
	>
		<LoaderCircle class="animate-spin" />
	</div>
{/if}
