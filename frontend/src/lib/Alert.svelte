<script lang="ts">
	import { X } from "@lucide/svelte";
	import { fly } from "svelte/transition";

	let open = $state(false);
	let text = $state("");
	let isError = $state(false);
	let timeoutId = $state<NodeJS.Timeout>();
	let expiresAt = $state(0);
	let hoverStart = $state(0);

	function startTimer(duration: number) {
		clearTimeout(timeoutId);
		open = true;
		timeoutId = setTimeout(() => (open = false), duration);
		expiresAt = Date.now() + duration;
	}

	export function error(err: string) {
		text = err;
		isError = true;
		startTimer(3000);
	}

	export function info(info: string) {
		text = info;
		isError = false;
		startTimer(3000);
	}
</script>

{#if open}
	<div
		transition:fly={{ duration: 100, y: -50 }}
		role="alert"
		class="absolute z-200 mt-18 flex w-full justify-center"
		onmouseenter={() => {
			hoverStart = Date.now();
			clearTimeout(timeoutId);
		}}
		onmouseleave={() => {
			if (open) {
				const now = Date.now();
				if (now - hoverStart >= 500) {
					startTimer(500);
				} else {
					startTimer(Math.max(expiresAt - now, 500));
				}
			}
		}}
	>
		<div
			class:bg-red-300={isError}
			class:dark:bg-red-500={isError}
			class="flex max-w-[60%] rounded-md bg-indigo-300 px-3 py-3 dark:bg-indigo-500"
		>
			<div class="text-center">{text}</div>
			<button
				class="cursor-pointer rounded-full pl-2 outline-offset-2 focus-visible:outline-2"
				onclick={() => (open = false)}
			>
				<X />
			</button>
		</div>
	</div>
{/if}
