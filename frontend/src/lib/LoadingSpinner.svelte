<script lang="ts" module>
	import { LoaderCircle } from "@lucide/svelte";
	import { navigating } from "$app/state";
	import { fade } from "svelte/transition";

	let timer: undefined | number = undefined;
	let navDelayed = $state(false);

	export function startNavigate() {
		timer = setTimeout(() => (navDelayed = true), 500);
	}

	export function clearNavigate() {
		clearTimeout(timer);
		navDelayed = false;
	}
</script>

<script>
	$effect(() => {
		if (navigating.to) {
			startNavigate();
		}

		return () => clearNavigate();
	});
</script>

{#if navDelayed}
	<div
		class="absolute top-0 z-1000 flex h-full w-full items-center justify-center bg-white/50 dark:bg-neutral-800/50"
		in:fade={{ duration: 250 }}
	>
		<LoaderCircle class="animate-spin" />
	</div>
{/if}
