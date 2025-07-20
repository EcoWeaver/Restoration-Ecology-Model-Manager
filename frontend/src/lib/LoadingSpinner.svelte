<script lang="ts" module>
	import { LoaderCircle } from "@lucide/svelte";
	import { navigating } from "$app/state";
	import { fade } from "svelte/transition";

	let timer: undefined | number = undefined;
	let spinner: "idle" | "waiting" | "spinning" = $state("idle");

	export function startNavigate() {
		if (spinner === "idle") {
			spinner = "waiting";
			timer = setTimeout(() => (spinner = "spinning"), 500);
		}
	}

	export function clearNavigate() {
		clearTimeout(timer);
		spinner = "idle";
	}
</script>

<script>
	$effect(() => {
		if (navigating.to) {
			startNavigate();
		} else {
			clearNavigate();
		}
	});
</script>

{#if spinner === "spinning"}
	<div
		class="absolute top-0 z-1000 flex h-full w-full items-center justify-center bg-white/50 dark:bg-neutral-800/50"
		in:fade={{ duration: 250 }}
	>
		<LoaderCircle class="animate-spin" />
	</div>
{/if}
