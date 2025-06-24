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
	<div class="absolute top-0 z-10000 h-full w-full" in:fade={{ duration: 250 }}>
		<div class="absolute z-10 h-full w-full bg-white opacity-50 dark:bg-neutral-800"></div>
		<div class="absolute z-20 flex h-full w-full animate-spin items-center justify-center">
			<LoaderCircle />
		</div>
	</div>
{/if}
